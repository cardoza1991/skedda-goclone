// internal/handlers/teacher.go
package handlers

import (
	"encoding/json"
	"net/http"
	"skedda-goclone/internal/models"

	"gorm.io/gorm"
)

type TeacherHandler struct {
	DB *gorm.DB
}

// RegisterTeacher handles teacher registration
func (h *TeacherHandler) RegisterTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Use GORM to create teacher
	if err := h.DB.Create(&teacher).Error; err != nil {
		http.Error(w, "Could not register teacher", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Teacher registered successfully")
}

// LoginTeacher handles teacher login
func (h *TeacherHandler) LoginTeacher(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Retrieve teacher from the database using GORM
	var teacher models.Teacher
	if err := h.DB.Where("email = ?", loginRequest.Email).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Teacher not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	if !teacher.CheckPassword(loginRequest.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Successful login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Teacher logged in successfully")
}
