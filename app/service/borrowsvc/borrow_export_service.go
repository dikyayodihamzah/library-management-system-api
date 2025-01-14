package borrowsvc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
)

var (
	excelTemplate = "borrow-template.xlsx"
	sheetName     = "Borrow List"
)

func (s *borrowService) writeTemplate() error {
	// create template
	s.Logger.Info("Creating template")
	f, err := os.Create(excelTemplate)
	if err != nil {
		s.Logger.Errorw("Failed to create template", "error", err)
		return err
	}
	defer f.Close()

	// write template
	file := excelize.NewFile()

	// set sheet name
	file.SetSheetName("Sheet1", sheetName)

	// set header
	headers := []string{
		"No",
		"Book ID",
		"Book Title",
		"Borrower",
		"Due Date",
		"Status",
		"Total Price (IDR)",
		"Created At",
	}

	for i, header := range headers {
		cell := string(rune(65+i)) + "1"
		file.SetCellValue(sheetName, cell, header)
	}

	// set width
	file.SetColWidth(sheetName, "A", "A", 5)
	file.SetColWidth(sheetName, "B", "B", 30)
	file.SetColWidth(sheetName, "C", "C", 30)
	file.SetColWidth(sheetName, "D", "D", 30)
	file.SetColWidth(sheetName, "E", "E", 30)
	file.SetColWidth(sheetName, "F", "F", 10)
	file.SetColWidth(sheetName, "G", "G", 15)
	file.SetColWidth(sheetName, "H", "H", 30)

	// set style
	headerStyle, err := file.NewStyle(`
	{
		"alignment":
			{
				"horizontal":"center",
				"vertical":"center"
			},
		"font":{"bold":true},
		"fill":
			{
				"type":"pattern",
				"color":["#f4b084"],
				"pattern":1
			}
	}`)
	if err != nil {
		s.Logger.Errorw("Failed to create header style", "error", err)
	}

	// set style to body
	bodyStyle, err := file.NewStyle(`{"alignment":{"wrap_text":true}}`)
	if err != nil {
		s.Logger.Errorw("Failed to create body style", "error", err)
	}

	// set style to header
	for i := 1; i <= len(headers); i++ {
		cell := string(rune(65+i-1)) + "1"
		file.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// set style to body
	for i := 2; i <= 1000; i++ {
		for j := 1; j <= len(headers); j++ {
			cell := string(rune(65+j-1)) + fmt.Sprintf("%d", i)
			file.SetCellStyle(sheetName, cell, cell, bodyStyle)
		}
	}

	// save file
	if err := file.SaveAs(excelTemplate); err != nil {
		s.Logger.Errorw("Failed to save template", "error", err)
		return err
	}

	return nil
}

func (s *borrowService) GenerateExcel(c context.Context, filter *book.BorrowQuery, timezone int) (*bytes.Buffer, error) {
	if _, err := os.Stat(excelTemplate); os.IsNotExist(err) {
		if err := s.writeTemplate(); err != nil {
			return nil, exception.ErrorInternal("Failed to write template")
		}
	}

	borrows, _, err := s.FindAll(c, filter)
	if err != nil {
		return nil, exception.ErrorInternal("Failed to get borrows")
	}

	// get template
	file, err := os.Open(excelTemplate)
	if err != nil {
		return nil, exception.ErrorInternal("Failed to open template")
	}

	// read template
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, exception.ErrorInternal("Failed to read template")
	}

	additionalDuration := time.Duration(timezone) * time.Hour

	// Write data rows
	for i, borrow := range borrows {
		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), i+1)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), borrow.Book.ID)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), borrow.Book.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), borrow.User.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), borrow.DueDate.Add(additionalDuration).Format("2006-01-02 15:04:05"))
		xlsx.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), borrow.Status)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), borrow.TotalPrice)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("H%d", i+2), borrow.CreatedAt.Add(additionalDuration).Format("2006-01-02 15:04:05"))
	}

	// Save file to buffer
	var buf bytes.Buffer
	if err := xlsx.Write(&buf); err != nil {
		s.Logger.Errorw("Failed to write Excel file to buffer", "error", err)
		return nil, err
	}

	return &buf, nil
}
