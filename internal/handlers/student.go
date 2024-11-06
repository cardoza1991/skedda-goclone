// internal/handlers/student.go
package handlers

import (
	"encoding/json"
	"net/http"
	"skedda-goclone/internal/models"

	"gorm.io/gorm"
)

type StudentHandler struct {
	DB *gorm.DB
}

// AddStudent allows a teacher to add a new student
func (h *StudentHandler) AddStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use GORM to create a student record
	if err := h.DB.Create(&student).Error; err != nil {
		http.Error(w, "Error saving student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Student added successfully"})
}

// ListStudents fetches all students for display in a dropdown
func (h *StudentHandler) ListStudents(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	if err := h.DB.Find(&students).Error; err != nil {
		http.Error(w, "Error fetching students", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(students)
}
