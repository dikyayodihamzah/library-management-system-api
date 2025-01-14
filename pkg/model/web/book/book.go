package book

import (
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
)

type BookRequest struct {
	Title           string `json:"title,omitempty" validate:"required"`
	Author          string `json:"author,omitempty" validate:"required"`
	Genre           string `json:"genre,omitempty" validate:"required"`
	Rating          int    `json:"rating,omitempty" validate:"required"`
	CoverURL        string `json:"cover_url,omitempty"`
	CoverColor      string `json:"cover_color,omitempty"`
	Description     string `json:"description,omitempty" validate:"required"`
	TotalCopies     int    `json:"total_copies,omitempty" validate:"required"`
	AvailableCopies int    `json:"available_copies,omitempty"`
	VideoURL        string `json:"video_url,omitempty"`
	Summary         string `json:"summary,omitempty" validate:"required"`
	Price           int    `json:"price,omitempty" validate:"required"`
}

type Book struct {
	model.Base
	BookRequest
}
