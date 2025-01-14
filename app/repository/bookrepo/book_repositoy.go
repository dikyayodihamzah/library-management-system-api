package bookrepo

import (
	"context"

	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BookRepository interface {
	Create(c context.Context, tx pgx.Tx, b *book.Book) error

	FindAll(c context.Context, filter *model.QueryParam) ([]book.Book, error)
	Count(c context.Context, filter *model.QueryParam) (int, error)
	FindByID(c context.Context, id uuid.UUID) (*book.Book, error)

	Update(c context.Context, tx pgx.Tx, b *book.Book) error

	Delete(c context.Context, tx pgx.Tx, id uuid.UUID) error
}

type bookRepository struct {
	Logger *zap.SugaredLogger
	DB     *pgxpool.Pool
}

func New(
	logger *zap.SugaredLogger,
	db *pgxpool.Pool,
) BookRepository {
	return &bookRepository{
		Logger: logger,
		DB:     db,
	}
}

func (r *bookRepository) Create(c context.Context, tx pgx.Tx, b *book.Book) error {
	queryStr := `
	INSERT INTO books (
		id,
		title,
		author,
		genre,
		rating,
		cover_url,
		cover_color,
		description,
		total_copies,
		available_copies,
		video_url,
		summary,
		price_idr,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	if _, err := tx.Exec(c, queryStr,
		b.ID,
		b.Title,
		b.Author,
		b.Genre,
		b.Rating,
		b.CoverURL,
		b.CoverColor,
		b.Description,
		b.TotalCopies,
		b.AvailableCopies,
		b.VideoURL,
		b.Summary,
		b.Price,
		b.CreatedAt,
	); err != nil {
		r.Logger.Errorw("failed to create book", "error", err)
		return err
	}

	r.Logger.Infow("book created", "id", b.ID)
	return nil
}

func (r *bookRepository) FindAll(c context.Context, filter *model.QueryParam) ([]book.Book, error) {
	queryStr := `
	SELECT
		id,
		title,
		author,
		genre,
		rating,
		cover_url,
		cover_color,
		description,
		total_copies,
		available_copies,
		video_url,
		summary,
		price_idr,
		created_at
	FROM books`

	queryStr, args := filterBooks(queryStr, filter)

	// sort
	queryStr, err := query.Sort(queryStr, filter.Sort, SortBookMap)
	if err != nil {
		r.Logger.Errorw("failed to sort query", "error", err)
		return nil, err
	}

	// pagination
	queryStr = query.Paginate(queryStr, filter.Page, filter.Limit)

	rows, err := r.DB.Query(c, queryStr, args...)
	if err != nil {
		r.Logger.Errorw("failed to get books", "error", err)
		return nil, err
	}
	defer rows.Close()

	books := make([]book.Book, 0)
	for rows.Next() {
		var b book.Book
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.Genre,
			&b.Rating,
			&b.CoverURL,
			&b.CoverColor,
			&b.Description,
			&b.TotalCopies,
			&b.AvailableCopies,
			&b.VideoURL,
			&b.Summary,
			&b.Price,
			&b.CreatedAt,
		)
		if err != nil {
			r.Logger.Errorw("failed to scan books", "error", err)
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (r *bookRepository) Count(c context.Context, filter *model.QueryParam) (int, error) {
	queryStr := `
	SELECT 
		COUNT(b.id)
	FROM books b`

	queryStr, args := filterBooks(queryStr, filter)

	var count int
	if err := r.DB.QueryRow(c, queryStr, args...).Scan(&count); err != nil {
		r.Logger.Errorw("failed to count books", "error", err)
		return 0, err
	}

	return count, nil
}

func (r *bookRepository) FindByID(c context.Context, id uuid.UUID) (*book.Book, error) {
	queryStr := `
	SELECT
		id,
		title,
		author,
		genre,
		rating,
		cover_url,
		cover_color,
		description,
		total_copies,
		available_copies,
		video_url,
		summary,
		price_idr,
		created_at
	FROM books
	WHERE id = $1`

	var b book.Book
	err := r.DB.QueryRow(c, queryStr, id).Scan(
		&b.ID,
		&b.Title,
		&b.Author,
		&b.Genre,
		&b.Rating,
		&b.CoverURL,
		&b.CoverColor,
		&b.Description,
		&b.TotalCopies,
		&b.AvailableCopies,
		&b.VideoURL,
		&b.Summary,
		&b.Price,
		&b.CreatedAt,
	)
	if err != nil {
		r.Logger.Errorw("failed to get book", "error", err)
		return nil, err
	}

	return &b, nil
}

func (r *bookRepository) Update(c context.Context, tx pgx.Tx, b *book.Book) error {
	queryStr := `
	UPDATE books
	SET
		title = $1,
		author = $2,
		genre = $3,
		rating = $4,
		cover_url = $5,
		cover_color = $6,
		description = $7,
		total_copies = $8,
		available_copies = $9,
		video_url = $10,
		summary = $11,
		price_idr = $12
	WHERE id = $13`

	if _, err := tx.Exec(c, queryStr,
		b.Title,
		b.Author,
		b.Genre,
		b.Rating,
		b.CoverURL,
		b.CoverColor,
		b.Description,
		b.TotalCopies,
		b.AvailableCopies,
		b.VideoURL,
		b.Summary,
		b.Price,
		b.ID,
	); err != nil {
		r.Logger.Errorw("failed to update book", "error", err)
		return err
	}

	r.Logger.Infow("book updated", "id", b.ID)
	return nil
}

func (r *bookRepository) Delete(c context.Context, tx pgx.Tx, id uuid.UUID) error {
	queryStr := `
	DELETE FROM books
	WHERE id = $1`

	if _, err := tx.Exec(c, queryStr, id); err != nil {
		r.Logger.Errorw("failed to delete book", "error", err)
		return err
	}

	r.Logger.Infow("book deleted", "id", id)
	return nil
}
