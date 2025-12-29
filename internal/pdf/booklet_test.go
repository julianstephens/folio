package pdf_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/julianstephens/folio/internal/pdf"
	testhelpers "github.com/julianstephens/folio/internal/tests"
)

func TestCreateBooklet_4Pages(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_4p.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}

	// Verify output has correct page count (4 input pages -> 2 output pages)
	reader, err := pdf.GetPDFReader(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	info, err := pdf.GetPDFInfo(reader, outputPath)
	if err != nil {
		t.Fatalf("Failed to get PDF info: %v", err)
	}

	if info.PageCount != 2 {
		t.Errorf("Expected 2 pages in output, got %d", info.PageCount)
	}
}

func TestCreateBooklet_8Pages(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "8p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_8p.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() failed: %v", err)
	}

	// Verify output has correct page count (8 input pages -> 4 output pages)
	reader, err := pdf.GetPDFReader(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	info, err := pdf.GetPDFInfo(reader, outputPath)
	if err != nil {
		t.Fatalf("Failed to get PDF info: %v", err)
	}

	if info.PageCount != 4 {
		t.Errorf("Expected 4 pages in output, got %d", info.PageCount)
	}
}

func TestCreateBooklet_5Pages_AutoPad(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "5p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_5p_autopad.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() with auto-pad failed: %v", err)
	}

	// Verify output was created and padded to 8 pages -> 4 output pages
	reader, err := pdf.GetPDFReader(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	info, err := pdf.GetPDFInfo(reader, outputPath)
	if err != nil {
		t.Fatalf("Failed to get PDF info: %v", err)
	}

	if info.PageCount != 4 {
		t.Errorf("Expected 4 pages in output (padded from 5), got %d", info.PageCount)
	}
}

func TestCreateBooklet_5Pages_NoPad(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "5p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_5p_nopad.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "none",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateBooklet() with no-pad should have failed for 5-page input")
	}

	if !errors.Is(err, pdf.ErrInvalidPageCount) {
		t.Errorf("Expected ErrInvalidPageCount, got: %v", err)
	}
}

func TestCreateBooklet_A4(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_a4.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_4p_a4.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() with A4 failed: %v", err)
	}

	// Verify output is A4 size
	reader, err := pdf.GetPDFReader(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	info, err := pdf.GetPDFInfo(reader, outputPath)
	if err != nil {
		t.Fatalf("Failed to get PDF info: %v", err)
	}

	if info.PageCount != 2 {
		t.Errorf("Expected 2 pages in output, got %d", info.PageCount)
	}

	firstBox, err := pdf.ExtractPageBox(info.PageBoundaries[0])
	if err != nil {
		t.Fatalf("Failed to extract page box: %v", err)
	}

	if firstBox.Size != pdf.KnownSizes.A4 {
		t.Errorf("Expected A4 size, got: %s", firstBox.Size)
	}
}

func TestCreateBooklet_ShortEdgeBinding(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_short_edge.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "short",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() with short-edge binding failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
}

func TestCreateBooklet_ExplicitSheetSize(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_explicit_size.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "letter",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() with explicit sheet size failed: %v", err)
	}

	// Verify output has correct size
	reader, err := pdf.GetPDFReader(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	info, err := pdf.GetPDFInfo(reader, outputPath)
	if err != nil {
		t.Fatalf("Failed to get PDF info: %v", err)
	}

	firstBox, err := pdf.ExtractPageBox(info.PageBoundaries[0])
	if err != nil {
		t.Fatalf("Failed to extract page box: %v", err)
	}

	if firstBox.Size != pdf.KnownSizes.Letter {
		t.Errorf("Expected Letter size, got: %s", firstBox.Size)
	}
}

func TestCreateBooklet_MixedSizes(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "mixed_sizes.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_mixed.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateBooklet() with mixed page sizes should have failed")
	}

	if !errors.Is(err, pdf.ErrMixedPageSizes) {
		t.Errorf("Expected ErrMixedPageSizes, got: %v", err)
	}
}

func TestCreateBooklet_NonWritableOutput(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := "/nonexistent/directory/booklet.pdf"

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateBooklet() with non-writable output should have failed")
	}

	if !errors.Is(err, pdf.ErrOutputNotWritable) {
		t.Errorf("Expected ErrOutputNotWritable, got: %v", err)
	}
}

func TestCreateBooklet_EmptyPDF(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "empty.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_invalid.pdf")

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateBooklet() with invalid input should have failed")
	}

	if !errors.Is(err, pdf.ErrInvalidPDF) {
		t.Errorf("Expected ErrInvalidPDF, got: %v", err)
	}
}

func TestCreateBooklet_CustomSize_ExplicitSheetSize(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "Booklet_organized.pdf")
	outputPath := filepath.Join(t.TempDir(), "booklet_custom.pdf")

	// Check if file exists first
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Skip("Booklet_organized.pdf not found")
	}

	opts := pdf.BookletOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "letter",
			Pad:       "auto",
		},
		BindingEdge: "long",
	}

	err := pdf.CreateBooklet(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateBooklet() with custom size and explicit sheet size failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
}
