package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestInspectExcel_ValidWorkbook(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.xlsx")

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	if err := f.SetCellValue(sheet, "A1", "name"); err != nil {
		t.Fatalf("set cell: %v", err)
	}
	if err := f.SetCellValue(sheet, "A2", "ana"); err != nil {
		t.Fatalf("set cell: %v", err)
	}
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("save workbook: %v", err)
	}
	_ = f.Close()

	summary, err := InspectExcel(path)
	if err != nil {
		t.Fatalf("InspectExcel returned error: %v", err)
	}
	if summary.Format != "excel" {
		t.Fatalf("expected format excel, got %q", summary.Format)
	}
	if !summary.HasHeader {
		t.Fatal("expected workbook structure to be detected")
	}
	if summary.Rows == 0 {
		t.Fatal("expected row count greater than zero")
	}
}

func TestInspectExcel_EmptySheet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty-sheet.xlsx")

	f := excelize.NewFile()
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("save workbook: %v", err)
	}
	_ = f.Close()

	summary, err := InspectExcel(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Rows != 0 {
		t.Fatalf("expected 0 rows for empty sheet, got %d", summary.Rows)
	}
}

func TestInspectExcel_FileNotFound(t *testing.T) {
	_, err := InspectExcel(filepath.Join(t.TempDir(), "missing.xlsx"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInspectExcel_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.xlsx")
	if err := os.WriteFile(path, []byte("not an excel file"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := InspectExcel(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
