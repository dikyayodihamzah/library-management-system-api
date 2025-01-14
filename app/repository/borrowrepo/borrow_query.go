package borrowrepo

import (
	"fmt"

	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
	"github.com/google/uuid"
)

var SortBorrowMap = map[string]string{
	"borrow_date": "br.borrow_date",
	"due_date":    "br.due_date",
	"return_date": "br.return_date",
	"status":      "br.status",
	"user_name":   "u.full_name",
	"book_title":  "b.title",
	"created_at":  "br.created_at",
	"price":       "br.total_price",
}

func filterBorrows(queryStr string, filter *book.BorrowQuery) (string, []interface{}) {
	if filter == nil {
		return queryStr, make([]interface{}, 0)
	}

	var args []interface{}

	if !filter.StartDate.IsZero() {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("br.borrow_date >= $%d", len(args)+1)
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("br.borrow_date <= $%d", len(args)+1)
		args = append(args, filter.EndDate)
	}

	if filter.Status != "" {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("br.status = $%d", len(args)+1)
		args = append(args, filter.Status)
	}

	if filter.BookID != uuid.Nil {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("br.book_id = $%d", len(args)+1)
		args = append(args, filter.BookID)
	}

	if filter.UserID != uuid.Nil {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("br.user_id = $%d", len(args)+1)
		args = append(args, filter.UserID)
	}

	return queryStr, args
}
