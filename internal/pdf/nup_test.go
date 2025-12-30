package pdf_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/julianstephens/folio/internal/pdf"
	testhelpers "github.com/julianstephens/folio/internal/tests"
)

func TestCreateNup_2up_4Pages(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_2up_4p.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() failed: %v", err)
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

func TestCreateNup_4up_8Pages(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "8p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_4up_8p.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    4,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() failed: %v", err)
	}

	// Verify output has correct page count (8 input pages -> 2 output pages)
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

func TestCreateNup_2up_5Pages_AutoPad(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "5p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_2up_5p_autopad.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() with auto-pad failed: %v", err)
	}

	// Verify output was created and padded to 6 pages -> 3 output pages
	reader, err := pdf.GetPDFReader(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	info, err := pdf.GetPDFInfo(reader, outputPath)
	if err != nil {
		t.Fatalf("Failed to get PDF info: %v", err)
	}

	if info.PageCount != 3 {
		t.Errorf("Expected 3 pages in output (padded from 5), got %d", info.PageCount)
	}
}

func TestCreateNup_2up_5Pages_NoPad(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "5p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_2up_5p_nopad.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "none",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateNup() with no-pad should have failed for 5-page input with 2-up")
	}

	if !errors.Is(err, pdf.ErrInvalidPageCount) {
		t.Errorf("Expected ErrInvalidPageCount, got: %v", err)
	}
}

func TestCreateNup_4up_5Pages_NoPad(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "5p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_4up_5p_nopad.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "none",
		},
		PerSheet:    4,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateNup() with no-pad should have failed for 5-page input with 4-up")
	}

	if !errors.Is(err, pdf.ErrInvalidPageCount) {
		t.Errorf("Expected ErrInvalidPageCount, got: %v", err)
	}
}

func TestCreateNup_A4(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_a4.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_2up_4p_a4.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() with A4 failed: %v", err)
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

func TestCreateNup_ExplicitSheetSize(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_explicit_size.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "letter",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() with explicit sheet size failed: %v", err)
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

func TestCreateNup_PortraitOrientation(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_portrait.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    4,
		Orientation: "portrait",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() with portrait orientation failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
}

func TestCreateNup_LandscapeOrientation(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_landscape.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "landscape",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err != nil {
		t.Fatalf("CreateNup() with landscape orientation failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created")
	}
}

func TestCreateNup_MixedSizes(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "mixed_sizes.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_mixed.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateNup() with mixed page sizes should have failed")
	}

	if !errors.Is(err, pdf.ErrMixedPageSizes) {
		t.Errorf("Expected ErrMixedPageSizes, got: %v", err)
	}
}

func TestCreateNup_NonWritableOutput(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "4p_letter.pdf")
	outputPath := "/nonexistent/directory/nup.pdf"

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateNup() with non-writable output should have failed")
	}

	if !errors.Is(err, pdf.ErrOutputNotWritable) {
		t.Errorf("Expected ErrOutputNotWritable, got: %v", err)
	}
}

func TestCreateNup_EmptyPDF(t *testing.T) {
	inputPath := filepath.Join(testhelpers.TESTDATA_DIR, "empty.pdf")
	outputPath := filepath.Join(t.TempDir(), "nup_invalid.pdf")

	opts := pdf.NupOptions{
		WriteOptions: pdf.WriteOptions{
			SheetSize: "",
			Pad:       "auto",
		},
		PerSheet:    2,
		Orientation: "",
	}

	err := pdf.CreateNup(inputPath, outputPath, opts)
	if err == nil {
		t.Fatalf("CreateNup() with invalid input should have failed")
	}

	if !errors.Is(err, pdf.ErrInvalidPDF) {
		t.Errorf("Expected ErrInvalidPDF, got: %v", err)
	}
}
