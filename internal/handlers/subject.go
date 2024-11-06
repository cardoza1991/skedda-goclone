// internal/handlers/subject.go
package handlers

import (
	"encoding/json"
	"net/http"
	"skedda-goclone/internal/models"

	"gorm.io/gorm"
)

type SubjectHandler struct {
	DB *gorm.DB
}

// CreateSubject allows a teacher to create a subject
func (h *SubjectHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var subject models.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use GORM to create a subject record
	if err := h.DB.Create(&subject).Error; err != nil {
		http.Error(w, "Error saving subject", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Subject created successfully"})
}

// AssignSubjectToStudent assigns a subject to a student
func (h *SubjectHandler) AssignSubjectToStudent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StudentID int64 `json:"student_id"`
		SubjectID int64 `json:"subject_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use GORM to assign the subject to the student
	studentSubject := models.StudentSubject{
		StudentID: input.StudentID,
		SubjectID: input.SubjectID,
	}

	if err := h.DB.Create(&studentSubject).Error; err != nil {
		http.Error(w, "Error assigning subject to student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Subject assigned to student successfully"})
}
