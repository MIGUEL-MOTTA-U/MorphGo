package parsers

import (
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
		t.Fatal("expected workbook sheets to be detected")
	}
	if len(summary.Columns) == 0 {
		t.Fatal("expected sheet names in summary")
	}
	if summary.Rows == 0 {
		t.Fatal("expected row count greater than zero")
	}
}

func TestInspectExcel_FileNotFound(t *testing.T) {
	_, err := InspectExcel(filepath.Join(t.TempDir(), "missing.xlsx"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
