package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

type Booking struct {
	ID        int64
	Space     string
	StartTime time.Time
	EndTime   time.Time
	User      string
	Notes     string
	Status    string // Added status field
}

type BookingSystem struct {
	spaces   []string
	bookings []Booking
	window   fyne.Window
	db       *sql.DB
}

func NewBookingSystem() *BookingSystem {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", "./bookings.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS spaces (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER PRIMARY KEY,
			space_id INTEGER,
			start_time DATETIME,
			end_time DATETIME,
			user TEXT,
			notes TEXT,
			status TEXT,
			FOREIGN KEY(space_id) REFERENCES spaces(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	bs := &BookingSystem{
		spaces:   []string{"Room A", "Room B", "Room C", "Conference Hall"},
		bookings: make([]Booking, 0),
		db:       db,
	}

	// Load existing bookings
	bs.loadBookings()
	return bs
}

func (bs *BookingSystem) loadBookings() {
	rows, err := bs.db.Query(`
		SELECT id, space_id, start_time, end_time, user, notes, status 
		FROM bookings 
		WHERE end_time >= datetime('now')
		ORDER BY start_time
	`)
	if err != nil {
		log.Printf("Error loading bookings: %v", err)
		return
	}
	defer rows.Close()

	bs.bookings = nil
	for rows.Next() {
		var b Booking
		var spaceID int
		err := rows.Scan(&b.ID, &spaceID, &b.StartTime, &b.EndTime, &b.User, &b.Notes, &b.Status)
		if err != nil {
			log.Printf("Error scanning booking: %v", err)
			continue
		}
		b.Space = bs.spaces[spaceID] // Convert space_id to space name
		bs.bookings = append(bs.bookings, b)
	}
}

func (bs *BookingSystem) createMainUI() fyne.CanvasObject {
	// Add a toolbar with common actions
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			bs.showBookingDialog(time.Now())
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			bs.loadBookings()
			bs.window.Content().Refresh()
		}),
	)

	// Create main content with tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Calendar", bs.createCalendarView()),
		container.NewTabItem("Spaces", bs.createSpacesView()),
		container.NewTabItem("Bookings", bs.createBookingsView()),
	)

	// Create status bar
	statusBar := widget.NewLabel("")
	go func() {
		for {
			time.Sleep(time.Minute)
			statusBar.SetText(fmt.Sprintf("Last updated: %s", time.Now().Format("15:04")))
		}
	}()

	return container.NewBorder(
		toolbar,
		statusBar,
		nil, nil,
		tabs,
	)
}

func (bs *BookingSystem) createCalendarView() fyne.CanvasObject {
	// Add navigation buttons for months
	prevMonth := widget.NewButton("←", nil)
	nextMonth := widget.NewButton("→", nil)
	currentDate := time.Now()

	monthLabel := widget.NewLabel(currentDate.Format("January 2006"))
	
	navigation := container.NewBorder(
		nil, nil,
		prevMonth, nextMonth,
		monthLabel,
	)

	// Create a grid for the calendar
	grid := container.NewGridWithColumns(7)
	
	updateCalendar := func(date time.Time) {
		grid.Objects = nil // Clear existing grid
		monthLabel.SetText(date.Format("January 2006"))

		// Add calendar header
		days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
		for _, day := range days {
			grid.Add(widget.NewLabel(day))
		}
		
		firstDay := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		lastDay := firstDay.AddDate(0, 1, -1)
		
		// Fill in empty cells
		for i := 0; i < int(firstDay.Weekday()); i++ {
			grid.Add(widget.NewLabel(""))
		}
		
		// Add days with booking indicators
		for day := 1; day <= lastDay.Day(); day++ {
			currentDay := day
			hasBookings := bs.hasBookingsOnDate(date.Year(), int(date.Month()), day)
			
			btn := widget.NewButton(fmt.Sprintf("%d", day), func() {
				clickDate := time.Date(date.Year(), date.Month(), currentDay, 0, 0, 0, 0, date.Location())
				bs.showBookingDialog(clickDate)
			})
			
			if hasBookings {
				btn.Importance = widget.HighImportance
			}
			
			grid.Add(btn)
		}
	}

	// Update navigation button handlers
	prevMonth.OnTapped = func() {
		currentDate = currentDate.AddDate(0, -1, 0)
		updateCalendar(currentDate)
	}
	
	nextMonth.OnTapped = func() {
		currentDate = currentDate.AddDate(0, 1, 0)
		updateCalendar(currentDate)
	}

	// Initial calendar setup
	updateCalendar(currentDate)

	return container.NewVBox(
		navigation,
		grid,
	)
}

func (bs *BookingSystem) createSpacesView() fyne.CanvasObject {
    // Create a list to display spaces
    list := widget.NewList(
        func() int { return len(bs.spaces) },
        func() fyne.CanvasObject {
            return container.NewHBox(
                widget.NewIcon(theme.HomeIcon()),
                widget.NewLabel("Template Space"),
                widget.NewLabel("(0 bookings)"),
            )
        },
        func(id widget.ListItemID, item fyne.CanvasObject) {
            box := item.(*fyne.Container)
            spaceLabel := box.Objects[1].(*widget.Label)
            bookingsLabel := box.Objects[2].(*widget.Label)
            
            spaceName := bs.spaces[id]
            bookingCount := 0
            for _, booking := range bs.bookings {
                if booking.Space == spaceName {
                    bookingCount++
                }
            }
            
            spaceLabel.SetText(spaceName)
            bookingsLabel.SetText(fmt.Sprintf("(%d bookings)", bookingCount))
        },
    )

    // Add button to manage spaces
    addButton := widget.NewButton("Add Space", func() {
        entry := widget.NewEntry()
        dialog.ShowForm("Add Space", "Add", "Cancel",
            []*widget.FormItem{
                {Text: "Space Name", Widget: entry},
            },
            func(submitted bool) {
                if submitted && entry.Text != "" {
                    // Save to database
                    _, err := bs.db.Exec("INSERT INTO spaces (name) VALUES (?)", entry.Text)
                    if err != nil {
                        dialog.ShowError(err, bs.window)
                        return
                    }
                    
                    bs.spaces = append(bs.spaces, entry.Text)
                    list.Refresh()
                }
            },
            bs.window,
        )
    })

    return container.NewBorder(
        widget.NewLabel("Available Spaces"),
        addButton,
        nil, nil,
        list,
    )
}

func (bs *BookingSystem) createBookingsView() fyne.CanvasObject {
    // Create table for bookings
    table := widget.NewTable(
        func() (int, int) { return len(bs.bookings), 5 }, // Added column for status
        func() fyne.CanvasObject { 
            return widget.NewLabel("") 
        },
        func(id widget.TableCellID, cell fyne.CanvasObject) {
            label := cell.(*widget.Label)
            if id.Row >= len(bs.bookings) {
                label.SetText("")
                return
            }
            
            booking := bs.bookings[id.Row]
            
            switch id.Col {
            case 0:
                label.SetText(booking.Space)
            case 1:
                label.SetText(booking.StartTime.Format("2006-01-02 15:04"))
            case 2:
                label.SetText(booking.EndTime.Format("15:04"))
            case 3:
                label.SetText(booking.User)
            case 4:
                label.SetText(booking.Status)
            }
        },
    )

    // Set column headers
    headers := widget.NewTable(
        func() (int, int) { return 1, 5 },
        func() fyne.CanvasObject {
            return widget.NewLabel("")
        },
        func(id widget.TableCellID, cell fyne.CanvasObject) {
            label := cell.(*widget.Label)
            headers := []string{"Space", "Start Time", "End Time", "User", "Status"}
            if id.Row == 0 && id.Col < len(headers) {
                label.SetText(headers[id.Col])
                label.TextStyle = fyne.TextStyle{Bold: true}
            }
        },
    )

    // Set column widths
    columnWidths := []float32{150, 150, 100, 120, 100}
    for i, width := range columnWidths {
        table.SetColumnWidth(i, width)
        headers.SetColumnWidth(i, width)
    }

    // Add context menu for booking management
    table.OnSelected = func(id widget.TableCellID) {
        if id.Row >= len(bs.bookings) {
            return
        }
        
        booking := bs.bookings[id.Row]
        menu := fyne.NewMenu("Booking",
            fyne.NewMenuItem("Cancel Booking", func() {
                dialog.ShowConfirm("Cancel Booking",
                    "Are you sure you want to cancel this booking?",
                    func(yes bool) {
                        if yes {
                            // Update status in database
                            _, err := bs.db.Exec(
                                "UPDATE bookings SET status = 'Cancelled' WHERE id = ?",
                                booking.ID,
                            )
                            if err != nil {
                                dialog.ShowError(err, bs.window)
                                return
                            }
                            
                            // Update in memory
                            bs.bookings[id.Row].Status = "Cancelled"
                            table.Refresh()
                        }
                    },
                    bs.window,
                )
            }),
            fyne.NewMenuItem("Edit Notes", func() {
                notes := widget.NewMultiLineEntry()
                notes.SetText(booking.Notes)
                
                dialog.ShowForm("Edit Notes", "Save", "Cancel",
                    []*widget.FormItem{
                        {Text: "Notes", Widget: notes},
                    },
                    func(submitted bool) {
                        if submitted {
                            // Update in database
                            _, err := bs.db.Exec(
                                "UPDATE bookings SET notes = ? WHERE id = ?",
                                notes.Text, booking.ID,
                            )
                            if err != nil {
                                dialog.ShowError(err, bs.window)
                                return
                            }
                            
                            // Update in memory
                            bs.bookings[id.Row].Notes = notes.Text
                            table.Refresh()
                        }
                    },
                    bs.window,
                )
            }),
        )
        
        popup := widget.NewPopUpMenu(menu, bs.window.Canvas())
        popup.Show()
    }

    // Search/filter functionality
    search := widget.NewEntry()
    search.SetPlaceHolder("Search bookings...")
    search.OnChanged = func(text string) {
        // Reload bookings with filter
        rows, err := bs.db.Query(`
            SELECT id, space_id, start_time, end_time, user, notes, status 
            FROM bookings 
            WHERE (user LIKE ? OR notes LIKE ?)
            AND end_time >= datetime('now')
            ORDER BY start_time
        `, "%"+text+"%", "%"+text+"%")
        
        if err != nil {
            log.Printf("Error searching bookings: %v", err)
            return
        }
        defer rows.Close()

        bs.bookings = nil
        for rows.Next() {
            var b Booking
            var spaceID int
            err := rows.Scan(&b.ID, &spaceID, &b.StartTime, &b.EndTime, &b.User, &b.Notes, &b.Status)
            if err != nil {
                log.Printf("Error scanning booking: %v", err)
                continue
            }
            b.Space = bs.spaces[spaceID]
            bs.bookings = append(bs.bookings, b)
        }
        
        table.Refresh()
    }

    return container.NewBorder(
        container.NewVBox(
            search,
            headers,
        ),
        nil, nil, nil,
        table,
    )
}

func (bs *BookingSystem) hasBookingsOnDate(year, month, day int) bool {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	nextDay := date.AddDate(0, 0, 1)
	
	for _, booking := range bs.bookings {
		if booking.StartTime.After(date) && booking.StartTime.Before(nextDay) {
			return true
		}
	}
	return false
}

func (bs *BookingSystem) showBookingDialog(date time.Time) {
	spaceSelect := widget.NewSelect(bs.spaces, nil)
	startTime := widget.NewEntry()
	startTime.SetText("09:00")
	endTime := widget.NewEntry()
	endTime.SetText("10:00")
	user := widget.NewEntry()
	notes := widget.NewMultiLineEntry()
	
	items := []*widget.FormItem{
		{Text: "Space", Widget: spaceSelect},
		{Text: "Start Time (HH:MM)", Widget: startTime},
		{Text: "End Time (HH:MM)", Widget: endTime},
		{Text: "User", Widget: user},
		{Text: "Notes", Widget: notes},
	}

	dialog.ShowForm("New Booking", "Book", "Cancel", items, func(submitted bool) {
		if submitted {
			startTimeStr := fmt.Sprintf("%s %s", date.Format("2006-01-02"), startTime.Text)
			endTimeStr := fmt.Sprintf("%s %s", date.Format("2006-01-02"), endTime.Text)
			
			start, err := time.Parse("2006-01-02 15:04", startTimeStr)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid start time format"), bs.window)
				return
			}
			
			end, err := time.Parse("2006-01-02 15:04", endTimeStr)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid end time format"), bs.window)
				return
			}

			// Validate booking
			if end.Before(start) {
				dialog.ShowError(fmt.Errorf("end time must be after start time"), bs.window)
				return
			}

			if bs.hasConflictingBooking(spaceSelect.Selected, start, end) {
				dialog.ShowError(fmt.Errorf("booking conflicts with existing reservation"), bs.window)
				return
			}

			// Save to database
			result, err := bs.db.Exec(`
				INSERT INTO bookings (space_id, start_time, end_time, user, notes, status)
				VALUES (?, ?, ?, ?, ?, ?)
			`, bs.getSpaceID(spaceSelect.Selected), start, end, user.Text, notes.Text, "Confirmed")
			
			if err != nil {
				dialog.ShowError(err, bs.window)
				return
			}

			id, _ := result.LastInsertId()
			bs.bookings = append(bs.bookings, Booking{
				ID:        id,
				Space:     spaceSelect.Selected,
				StartTime: start,
				EndTime:   end,
				User:     user.Text,
				Notes:    notes.Text,
				Status:   "Confirmed",
			})

			bs.window.Content().Refresh()
		}
	}, bs.window)
}

func (bs *BookingSystem) hasConflictingBooking(space string, start, end time.Time) bool {
	for _, booking := range bs.bookings {
		if booking.Space == space &&
			((start.After(booking.StartTime) && start.Before(booking.EndTime)) ||
			(end.After(booking.StartTime) && end.Before(booking.EndTime)) ||
			(start.Before(booking.StartTime) && end.After(booking.EndTime))) {
			return true
		}
	}
	return false
}

func (bs *BookingSystem) getSpaceID(spaceName string) int {
	for i, space := range bs.spaces {
		if space == spaceName {
			return i
		}
	}
	return -1
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Booking System")
	
	bs := NewBookingSystem()
	bs.window = window
	
	window.SetContent(bs.createMainUI())
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
