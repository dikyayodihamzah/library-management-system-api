package book

import (
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/google/uuid"
)

type BorrowRequest struct {
	BookIDs []uuid.UUID `json:"book_ids,omitempty" validate:"required"`
	UserID  uuid.UUID   `json:"user_id,omitempty" validate:"required"`
	DueDate time.Time   `json:"due_date,omitempty"`
}

type BorrowRecord struct {
	model.Base
	BorrowRequest
	BookID       uuid.UUID  `json:"book_id,omitempty"`
	BorrowDate   time.Time  `json:"borrow_date,omitempty" validate:"required"`
	ReturnedDate *time.Time `json:"returned_date,omitempty"`
	Status       string     `json:"status,omitempty"`
	TotalPrice   int        `json:"total_price,omitempty"`
}

type BorrowDTO struct {
	BorrowRecord
	UserName  string
	BookTitle string
}

type BorrowResponse struct {
	BorrowRecord
	TotalPrice int                  `json:"total_price"`
	User       model.SimpleResponse `json:"user"`
	Book       model.SimpleResponse `json:"book"`
}

type BorrowUserResponse struct {
	BorrowDate time.Time              `json:"borrow_date"`
	DueDate    time.Time              `json:"due_date"`
	TotalPrice int                    `json:"total_price"`
	Books      []model.SimpleResponse `json:"books"`
}

type BorrowQuery struct {
	model.QueryParam
	UserID    uuid.UUID `query:"user_id,omitempty"`
	BookID    uuid.UUID `query:"book_id,omitempty"`
	Status    string    `query:"status,omitempty"`
	StartDate time.Time `query:"start_date,omitempty"`
	EndDate   time.Time `query:"end_date,omitempty"`
}

func (req *BorrowRequest) ToBorrowRecord() []BorrowRecord {
	records := make([]BorrowRecord, 0)
	for _, id := range req.BookIDs {
		record := BorrowRecord{
			BorrowRequest: *req,
			BookID:        id,
			BorrowDate:    time.Now(),
			Status:        "BORROWED",
		}
		record.BookIDs = nil
		record.ID = uuid.New()
		record.CreatedAt = lib.Pointer(time.Now())
		records = append(records, record)
	}

	return records
}

func (dto *BorrowDTO) ToResponse() *BorrowResponse {
	user := model.SimpleResponse{
		ID:   dto.UserID,
		Name: dto.UserName,
	}

	book := model.SimpleResponse{
		ID:   dto.BookID,
		Name: dto.BookTitle,
	}

	dto.BookID = uuid.Nil
	dto.UserID = uuid.Nil

	return &BorrowResponse{
		BorrowRecord: dto.BorrowRecord,
		Book:         book,
		User:         user,
		TotalPrice:   dto.TotalPrice,
	}
}

func DtoToUserResponse(dtos []BorrowDTO) *BorrowUserResponse {
	books := make([]model.SimpleResponse, 0)
	totalPrice := 0
	for _, dto := range dtos {
		books = append(books, model.SimpleResponse{
			ID:   dto.BookID,
			Name: dto.BookTitle,
		})
		totalPrice += 10000
	}

	durationDays := dtos[0].DueDate.Sub(dtos[0].BorrowDate).Hours() / 24
	totalPrice = totalPrice * int(durationDays)

	return &BorrowUserResponse{
		BorrowDate: dtos[0].BorrowDate,
		DueDate:    dtos[0].DueDate,
		TotalPrice: totalPrice,
		Books:      books,
	}
}

func ModelToUserResponse(borrows []BorrowRecord) *BorrowUserResponse {
	dtos := make([]BorrowDTO, 0)
	for _, borrow := range borrows {
		dtos = append(dtos, BorrowDTO{
			BorrowRecord: borrow,
			UserName:     "User Name",
			BookTitle:    "Book Title",
		})
	}

	return DtoToUserResponse(dtos)
}
