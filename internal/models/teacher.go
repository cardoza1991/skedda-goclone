// internal/models/teacher.go
package models

import "golang.org/x/crypto/bcrypt"

type Teacher struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// SetPassword hashes and sets the teacher's password
func (t *Teacher) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies if the password matches the hashed password
func (t *Teacher) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(t.PasswordHash), []byte(password)) == nil
}
