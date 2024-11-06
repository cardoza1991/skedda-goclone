// internal/models/models.go
package models

type Student struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Subject struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StudentSubject struct {
	StudentID int64 `json:"student_id" gorm:"primaryKey"`
	SubjectID int64 `json:"subject_id" gorm:"primaryKey"`
}
