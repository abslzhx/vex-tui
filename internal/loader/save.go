package loader

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/CodeOne45/vex-tui/pkg/models"
	"github.com/xuri/excelize/v2"
)

// SaveExcel saves sheets to an Excel file
func SaveExcel(sheets []models.Sheet, filename string) error {
	var f *excelize.File
	var err error

	// Try to open existing file to preserve formatting
	if _, statErr := os.Stat(filename); statErr == nil {
		f, err = excelize.OpenFile(filename)
	}

	// If file doesn't exist or can't be opened, create new
	if f == nil || err != nil {
		f = excelize.NewFile()
		// Only rename Sheet1 for new files
		if len(sheets) > 0 {
			f.SetSheetName("Sheet1", sheets[0].Name)
		}
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", err)
		}
	}()

	for _, sheet := range sheets {
		sheetName := sheet.Name

		// Check if sheet exists
		idx, err := f.GetSheetIndex(sheetName)
		if err != nil {
			// If error occurs (e.g. invalid name), assume we might need to create it or skip
			idx = -1
		}

		if idx == -1 {
			_, err := f.NewSheet(sheetName)
			if err != nil {
				return fmt.Errorf("failed to create sheet %s: %w", sheetName, err)
			}
		}

		for rowIdx, row := range sheet.Rows {
			for colIdx, cell := range row {
				cellRef, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
				if err != nil {
					continue
				}

				// Try to parse value as number if it looks like one, but preserve leading zeros for IDs etc.
				var val interface{} = cell.Value
				if num, err := strconv.ParseFloat(cell.Value, 64); err == nil {
					// Check for leading zeros (e.g. "0123") which should stay as text
					// Exception: "0", "0.xxx"
					isLeadingZero := len(cell.Value) > 1 && cell.Value[0] == '0' && cell.Value[1] != '.'
					if !isLeadingZero {
						val = num
					}
				}

				if cell.Formula != "" {
					if err := f.SetCellFormula(sheetName, cellRef, cell.Formula); err != nil {
						if err := f.SetCellValue(sheetName, cellRef, val); err != nil {
							continue
						}
					}
				} else {
					if err := f.SetCellValue(sheetName, cellRef, val); err != nil {
						continue
					}
				}
			}
		}
	}

	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// SaveCSV saves a sheet to CSV format
func SaveCSV(sheet models.Sheet, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", closeErr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range sheet.Rows {
		record := make([]string, 0, len(row))
		for _, cell := range row {
			value := cell.Value
			if cell.Formula != "" {
				value = "=" + cell.Formula
			}
			record = append(record, value)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	return nil
}
