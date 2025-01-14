package borrowrepo

import (
	"context"
	"fmt"

	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BorrowRepository interface {
	Add(c context.Context, tx pgx.Tx, borrow ...book.BorrowRecord) error
	FindAll(c context.Context, filter *book.BorrowQuery) ([]book.BorrowDTO, error)
	Count(c context.Context, filter *book.BorrowQuery) (int, error)
	Update(c context.Context, tx pgx.Tx, borrow *book.BorrowRecord) error
}

type borrowRepository struct {
	Logger *zap.SugaredLogger
	DB     *pgxpool.Pool
}

func New(
	logger *zap.SugaredLogger,
	db *pgxpool.Pool,
) BorrowRepository {
	return &borrowRepository{
		Logger: logger,
		DB:     db,
	}
}

func (r borrowRepository) Add(c context.Context, tx pgx.Tx, borrow ...book.BorrowRecord) error {
	queryStr := `
	INSERT INTO borrow_records (
		book_id, 
		user_id, 
		borrow_date, 
		due_date, 
		status,
		created_at,
		total_price
	) VALUES `

	args := make([]interface{}, 0)
	for i, b := range borrow {
		n := i * 7
		queryStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5, n+6, n+7)
		if i < len(borrow)-1 {
			queryStr += ", "
		}

		args = append(args,
			b.BookID,
			b.UserID,
			b.BorrowDate,
			b.DueDate,
			b.Status,
			b.CreatedAt,
			b.TotalPrice,
		)
	}

	if _, err := tx.Exec(c, queryStr, args...); err != nil {
		r.Logger.Errorw("failed to add borrow", "error", err)
		return err
	}

	return nil
}

func (r *borrowRepository) FindAll(c context.Context, filter *book.BorrowQuery) ([]book.BorrowDTO, error) {
	queryStr := `
	SELECT
		br.id,
		br.book_id,
		br.user_id,
		br.borrow_date,
		br.due_date,
		br.return_date,
		br.status,
		br.total_price,
		br.created_at,
		u.full_name,
		b.title
	FROM borrow_records br
	INNER JOIN users u ON br.user_id = u.id
	INNER JOIN books b ON br.book_id = b.id`

	queryStr, args := filterBorrows(queryStr, filter)

	// sort
	queryStr, err := query.Sort(queryStr, filter.Sort, SortBorrowMap)
	if err != nil {
		r.Logger.Errorw("failed to sort query", "error", err)
		return nil, err
	}

	// pagination
	queryStr = query.Paginate(queryStr, filter.Page, filter.Limit)

	rows, err := r.DB.Query(c, queryStr, args...)
	if err != nil {
		r.Logger.Errorw("failed to get borrows", "error", err)
		return nil, err
	}
	defer rows.Close()

	borrows := make([]book.BorrowDTO, 0)
	for rows.Next() {
		var b book.BorrowDTO
		err := rows.Scan(
			&b.ID,
			&b.BookID,
			&b.UserID,
			&b.BorrowDate,
			&b.DueDate,
			&b.ReturnedDate,
			&b.Status,
			&b.TotalPrice,
			&b.CreatedAt,
			&b.UserName,
			&b.BookTitle,
		)
		if err != nil {
			r.Logger.Errorw("failed to scan borrows", "error", err)
			return nil, err
		}

		borrows = append(borrows, b)
	}

	return borrows, nil
}

func (r *borrowRepository) Count(c context.Context, filter *book.BorrowQuery) (int, error) {
	queryStr := `
	SELECT 
		COUNT(br.id)
	FROM borrow_records br
	INNER JOIN users u ON br.user_id = u.id
	INNER JOIN books b ON br.book_id = b.id`

	queryStr, args := filterBorrows(queryStr, filter)

	var count int
	if err := r.DB.QueryRow(c, queryStr, args...).Scan(&count); err != nil {
		r.Logger.Errorw("error on count borrow", "error", err)
		return 0, err
	}

	return count, nil
}

func (r *borrowRepository) Update(c context.Context, tx pgx.Tx, borrow *book.BorrowRecord) error {
	queryStr := `
	UPDATE borrow_records
	SET
		return_date = $1,
		status = $2
	WHERE id = $3`

	if _, err := tx.Exec(c, queryStr,
		borrow.ReturnedDate,
		borrow.Status,
		borrow.ID,
	); err != nil {
		r.Logger.Errorw("failed to update borrow", "error", err)
		return err
	}

	return nil
}
