package user

import (
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	model.Base
	SignUpRequest
	Role             string    `json:"role" db:"role"`
	LastActivityDate time.Time `json:"last_activity_date" db:"last_activity_date"`
}

func (u *User) GenerateID() {
	u.ID = uuid.New()
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}
