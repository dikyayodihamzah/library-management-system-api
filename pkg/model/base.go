package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID  `json:"id" swaggerignore:"true"`
	CreatedAt *time.Time `json:"created_at,omitempty" swaggerignore:"true" format:"date-time"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" swaggerignore:"true" format:"date-time"`
	DeletedAt *time.Time `json:"-" swaggerignore:"true" format:"date-time"`
}

type BaseSerial struct {
	ID        int        `json:"id" swaggerignore:"true"`
	CreatedAt *time.Time `json:"created_at,omitempty" swaggerignore:"true" format:"date-time"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" swaggerignore:"true" format:"date-time"`
	DeletedAt *time.Time `json:"-" swaggerignore:"true" format:"date-time"`
}

type BaseString struct {
	ID        string     `json:"id" swaggerignore:"true"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at" swaggerignore:"true" format:"date-time"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at" swaggerignore:"true" format:"date-time"`
	DeletedAt *time.Time `json:"-" db:"deleted_at" swaggerignore:"true" format:"date-time"`
}

type SimpleResponse struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

type Table interface {
	Table() TableData
}

type TableData struct {
	Name       string
	PrimaryKey string
	Alias      string
}
