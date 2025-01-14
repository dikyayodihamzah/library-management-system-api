package borrowsvc

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/dikyayodihamzah/library-management-api/app/repository/bookrepo"
	"github.com/dikyayodihamzah/library-management-api/app/repository/borrowrepo"
	"github.com/dikyayodihamzah/library-management-api/app/repository/userrepo"
	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/dikyayodihamzah/library-management-api/pkg/transaction"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type BorrowService interface {
	Borrow(c context.Context, req *book.BorrowRequest) (*book.BorrowUserResponse, error)
	Return(c context.Context, req *book.BorrowRequest) error
	FindAll(c context.Context, filter *book.BorrowQuery) ([]book.BorrowResponse, int, error)

	GenerateExcel(c context.Context, filter *book.BorrowQuery, timezone int) (*bytes.Buffer, error)
}

type borrowService struct {
	Logger     *zap.SugaredLogger
	Validate   *validator.Validate
	TxManager  transaction.Manager
	UserRepo   userrepo.UserRepository
	BookRepo   bookrepo.BookRepository
	BorrowRepo borrowrepo.BorrowRepository
}

func New(
	logger *zap.SugaredLogger,
	validate *validator.Validate,
	txManager transaction.Manager,
	userRepo userrepo.UserRepository,
	bookRepo bookrepo.BookRepository,
	borrowRepo borrowrepo.BorrowRepository,
) BorrowService {
	return &borrowService{
		Logger:     logger,
		Validate:   validate,
		TxManager:  txManager,
		UserRepo:   userRepo,
		BookRepo:   bookRepo,
		BorrowRepo: borrowRepo,
	}
}

func (s *borrowService) Borrow(c context.Context, req *book.BorrowRequest) (*book.BorrowUserResponse, error) {
	// validate request
	if err := s.Validate.Struct(req); err != nil {
		return nil, exception.ErrorBadRequest(err.Error())
	}

	// get user data
	if _, err := s.UserRepo.FindByColumn(c, "id", req.UserID); err != nil {
		return nil, exception.ErrorNotFound("User not found")
	}

	uniqueBookIDs := []uuid.UUID{}
	bookReqMap := map[uuid.UUID]bool{}
	for _, id := range req.BookIDs {
		if _, ok := bookReqMap[id]; !ok {
			uniqueBookIDs = append(uniqueBookIDs, id)
			bookReqMap[id] = true
		}
	}

	if len(uniqueBookIDs) > 5 {
		return nil, exception.ErrorBadRequest("Cannot borrow more than 5 books")
	}

	// get all borrowed books data
	filter := &book.BorrowQuery{UserID: req.UserID}
	dtos, err := s.BorrowRepo.FindAll(c, filter)
	if err != nil {
		return nil, exception.ErrorInternal("Failed to get borrowed books")
	}
	borrowMap := make(map[uuid.UUID]string)
	for _, dto := range dtos {
		if dto.BorrowRecord.Status == "BORROWED" {
			borrowMap[dto.BookID] = dto.BookTitle
		}
	}

	// get book data
	borrowedBooks := make([]book.Book, 0)
	for _, id := range req.BookIDs {
		b, err := s.BookRepo.FindByID(c, id)
		if err != nil {
			return nil, exception.ErrorNotFound(fmt.Sprintf("Book with ID %s not found", id))
		}

		// check if book is available
		if b.AvailableCopies == 0 {
			return nil, exception.ErrorBadRequest(fmt.Sprintf("Book with ID %s is not available", id))
		}

		// check if user has borrowed the book
		if _, ok := borrowMap[id]; ok {
			return nil, exception.ErrorBadRequest(fmt.Sprintf("User has borrowed the book with ID %s", id))
		}

		b.AvailableCopies--
		borrowedBooks = append(borrowedBooks, *b)
	}

	borrowedBookMap := make(map[uuid.UUID]book.Book)
	for _, b := range borrowedBooks {
		borrowedBookMap[b.ID] = b
	}

	// validate due date to at least 1 day from now
	if req.DueDate.IsZero() {
		return nil, exception.ErrorBadRequest("Due date is required")
	}

	if time.Until(req.DueDate).Hours() < 24 {
		return nil, exception.ErrorBadRequest("Due date must be at least 1 day from now")
	}

	// create new borrow record data
	borrowRecords := req.ToBorrowRecord()

	// set price
	durationDays := time.Until(req.DueDate).Hours() / 24
	for i := range borrowRecords {
		borrowRecords[i].TotalPrice = borrowedBookMap[borrowRecords[i].BookID].Price * int(durationDays)
	}

	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		for _, b := range borrowedBooks {
			if err := s.BookRepo.Update(c, tx, &b); err != nil {
				return err
			}
		}

		return s.BorrowRepo.Add(c, tx, borrowRecords...)
	}); err != nil {
		return nil, exception.ErrorInternal("Failed to borrow book")
	}

	// construct response
	res := &book.BorrowUserResponse{
		BorrowDate: time.Now(),
		DueDate:    req.DueDate,
	}

	totalPrice := 0
	books := make([]model.SimpleResponse, 0)
	for key, value := range borrowedBookMap {
		books = append(books, model.SimpleResponse{
			ID:   key,
			Name: value.Title,
		})
		totalPrice += value.Price
	}

	totalPrice = totalPrice * int(durationDays)

	res.TotalPrice = totalPrice
	res.Books = books

	return res, nil
}

func (s *borrowService) Return(c context.Context, req *book.BorrowRequest) error {
	// validate request
	if err := s.Validate.Struct(req); err != nil {
		return exception.ErrorBadRequest(err.Error())
	}

	// get user data
	if _, err := s.UserRepo.FindByColumn(c, "id", req.UserID); err != nil {
		return exception.ErrorNotFound("User not found")
	}

	uniqueBookIDs := []uuid.UUID{}
	bookReqMap := map[uuid.UUID]bool{}
	for _, id := range req.BookIDs {
		if _, ok := bookReqMap[id]; !ok {
			uniqueBookIDs = append(uniqueBookIDs, id)
			bookReqMap[id] = true
		}
	}

	// get all borrowed books data
	filter := &book.BorrowQuery{UserID: req.UserID}
	dtos, err := s.BorrowRepo.FindAll(c, filter)
	if err != nil {
		return exception.ErrorInternal("Failed to get borrowed books")
	}
	borrowMap := make(map[uuid.UUID]book.BorrowRecord)
	for _, dto := range dtos {
		if dto.BorrowRecord.Status == "BORROWED" {
			borrowMap[dto.BookID] = dto.BorrowRecord
		}
	}

	// get book data
	returnedBooks := make([]book.Book, 0)
	updatedBorrowRecords := make([]book.BorrowRecord, 0)
	for _, id := range uniqueBookIDs {
		b, err := s.BookRepo.FindByID(c, id)
		if err != nil {
			return exception.ErrorNotFound(fmt.Sprintf("Book with ID %s not found", id))
		}

		// check if user has borrowed the book
		if _, ok := borrowMap[id]; !ok {
			return exception.ErrorBadRequest(fmt.Sprintf("User has not borrowed the book with ID %s", id))
		}

		b.AvailableCopies++
		returnedBooks = append(returnedBooks, *b)

		borrowRecord := borrowMap[id]
		borrowRecord.ReturnedDate = lib.TimeNowPtr()
		borrowRecord.Status = "RETURNED"
		updatedBorrowRecords = append(updatedBorrowRecords, borrowRecord)
	}

	// update book data
	if err := s.TxManager.WithTx(c, func(tx pgx.Tx) error {
		for i, b := range returnedBooks {
			if err := s.BookRepo.Update(c, tx, &b); err != nil {
				return err
			}

			if err := s.BorrowRepo.Update(c, tx, &updatedBorrowRecords[i]); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return exception.ErrorInternal("Failed to return book")
	}

	return nil
}

func (s *borrowService) FindAll(c context.Context, filter *book.BorrowQuery) ([]book.BorrowResponse, int, error) {
	// validate filter
	if filter.Sort == "" {
		filter.Sort = "-borrow_date"
	}

	if _, _, err := query.ValidateSort(filter.Sort, borrowrepo.SortBorrowMap); err != nil {
		return nil, 0, exception.ErrorBadRequest(err.Error())
	}

	// get all borrow records data
	dtos, err := s.BorrowRepo.FindAll(c, filter)
	if err != nil {
		return nil, 0, exception.ErrorInternal("Failed to get borrow records")
	}

	res := make([]book.BorrowResponse, 0)
	for _, dto := range dtos {
		res = append(res, *dto.ToResponse())
	}

	// get total borrow records data
	total, err := s.BorrowRepo.Count(c, filter)
	if err != nil {
		return nil, 0, exception.ErrorInternal("Failed to get total borrow records")
	}

	return res, total, nil
}
