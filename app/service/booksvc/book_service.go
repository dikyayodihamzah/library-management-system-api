package booksvc

import (
	"context"
	"time"

	"github.com/dikyayodihamzah/library-management-api/app/repository/bookrepo"
	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/dikyayodihamzah/library-management-api/pkg/transaction"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BookService interface {
	Create(c context.Context, req *book.BookRequest) (*book.Book, error)
	FindAll(c context.Context, filter *model.QueryParam) ([]book.Book, int, error)
	FindByID(c context.Context, id uuid.UUID) (*book.Book, error)
	Update(c context.Context, id uuid.UUID, req *book.BookRequest) (*book.Book, error)
	Delete(c context.Context, id uuid.UUID) error
}

type bookService struct {
	Validate  *validator.Validate
	TxManager transaction.Manager
	BookRepo  bookrepo.BookRepository
}

func New(
	validate *validator.Validate,
	txManager transaction.Manager,
	bookRepo bookrepo.BookRepository,
) BookService {
	return &bookService{
		Validate:  validate,
		TxManager: txManager,
		BookRepo:  bookRepo,
	}
}

func (s *bookService) Create(c context.Context, req *book.BookRequest) (*book.Book, error) {
	// validate request
	if err := s.Validate.Struct(req); err != nil {
		return nil, exception.ErrorBadRequest(err.Error())
	}

	arr := map[string]string{
		"title":       req.Title,
		"author":      req.Author,
		"genre":       req.Genre,
		"description": req.Description,
	}

	for key, value := range arr {
		if len(value) > 255 {
			return nil, exception.ErrorBadRequest(key + " must be less than 255 characters")
		}
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, exception.ErrorBadRequest("rating must be between 1 and 5")
	}

	if req.TotalCopies < 1 {
		return nil, exception.ErrorBadRequest("total_copies must be greater than 0")
	}
	req.AvailableCopies = req.TotalCopies

	if req.Price < 0 {
		return nil, exception.ErrorBadRequest("price must be greater than or equal to 0")
	}

	// create new book data
	b := book.Book{
		BookRequest: *req,
	}
	b.ID = uuid.New()
	b.CreatedAt = lib.Pointer(time.Now())

	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		return s.BookRepo.Create(c, tx, &b)
	}); err != nil {
		return nil, exception.ErrorInternal("Failed to create book")
	}

	return &b, nil
}

func (s *bookService) FindAll(c context.Context, filter *model.QueryParam) ([]book.Book, int, error) {
	// validate filter
	if filter.Sort == "" {
		filter.Sort = "-created_at"
	}

	if _, _, err := query.ValidateSort(filter.Sort, bookrepo.SortBookMap); err != nil {
		return nil, 0, exception.ErrorBadRequest(err.Error())
	}

	// get all books data
	books, err := s.BookRepo.FindAll(c, filter)
	if err != nil {
		return nil, 0, exception.ErrorInternal("Failed to get books")
	}

	// get total books data
	total, err := s.BookRepo.Count(c, filter)
	if err != nil {
		return nil, 0, exception.ErrorInternal("Failed to get total books")
	}

	return books, total, nil
}

func (s *bookService) FindByID(c context.Context, id uuid.UUID) (*book.Book, error) {
	// get book data
	b, err := s.BookRepo.FindByID(c, id)
	if err != nil {
		return nil, exception.ErrorNotFound("Book not found")
	}

	return b, nil
}

func (s *bookService) Update(c context.Context, id uuid.UUID, req *book.BookRequest) (*book.Book, error) {
	// get book data
	b, err := s.BookRepo.FindByID(c, id)
	if err != nil {
		return nil, exception.ErrorNotFound("Book not found")
	}

	// validate request
	if err := s.Validate.Struct(req); err != nil {
		return nil, exception.ErrorBadRequest(err.Error())
	}

	arr := map[string]string{
		"title":       req.Title,
		"author":      req.Author,
		"genre":       req.Genre,
		"description": req.Description,
	}

	for key, value := range arr {
		if len(value) > 255 {
			return nil, exception.ErrorBadRequest(key + " must be less than 255 characters")
		}
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, exception.ErrorBadRequest("rating must be between 1 and 5")
	}

	if req.TotalCopies < 1 {
		return nil, exception.ErrorBadRequest("total_copies must be greater than 0")
	}

	if req.AvailableCopies > req.TotalCopies {
		return nil, exception.ErrorBadRequest("available_copies must be less than or equal to total_copies")
	}

	if req.Price < 0 {
		return nil, exception.ErrorBadRequest("price must be greater than or equal to 0")
	}

	// update book data
	b.Title = req.Title
	b.Author = req.Author
	b.Genre = req.Genre
	b.Rating = req.Rating
	b.Description = req.Description
	b.TotalCopies = req.TotalCopies
	b.AvailableCopies = req.AvailableCopies
	b.VideoURL = req.VideoURL
	b.Summary = req.Summary
	b.Price = req.Price

	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		return s.BookRepo.Update(c, tx, b)
	}); err != nil {
		return nil, exception.ErrorInternal("Failed to update book")
	}

	return b, nil
}

func (s *bookService) Delete(c context.Context, id uuid.UUID) error {
	// get book data
	b, err := s.BookRepo.FindByID(c, id)
	if err != nil {
		return exception.ErrorNotFound("Book not found")
	}

	// check if book has been borrowed
	if b.TotalCopies != b.AvailableCopies {
		return exception.ErrorBadRequest("Book has been borrowed")
	}

	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		return s.BookRepo.Delete(c, tx, b.ID)
	}); err != nil {
		return exception.ErrorInternal("Failed to delete book")
	}

	return nil
}
