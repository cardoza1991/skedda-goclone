// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"skedda-goclone/internal/database"
	"skedda-goclone/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Run migrations for all models
	db.Migrate()

	// Initialize router
	router := mux.NewRouter()

	// Register handlers
	teacherHandler := handlers.TeacherHandler{DB: db.DB}
	studentHandler := handlers.StudentHandler{DB: db.DB}
	subjectHandler := handlers.SubjectHandler{DB: db.DB}

	// Define API endpoints
	router.HandleFunc("/api/teachers/register", teacherHandler.RegisterTeacher).Methods("POST")
	router.HandleFunc("/api/students", studentHandler.AddStudent).Methods("POST")
	router.HandleFunc("/api/students", studentHandler.ListStudents).Methods("GET")
	router.HandleFunc("/api/subjects", subjectHandler.CreateSubject).Methods("POST")
	router.HandleFunc("/api/subjects/assign", subjectHandler.AssignSubjectToStudent).Methods("POST")

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
