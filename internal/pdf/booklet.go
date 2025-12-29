package pdf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/julianstephens/folio/internal/utils"
)

var (
	ErrMixedPageSizes      = errors.New("mixed page sizes")
	ErrCannotInferSize     = errors.New("cannot infer sheet size")
	ErrInvalidPageCount    = errors.New("invalid page count")
	ErrOutputNotWritable   = errors.New("output path not writable")
)

// BookletOptions contains configuration for booklet creation.
type BookletOptions struct {
	SheetSize   string // "letter" or "a4" or empty to infer
	BindingEdge string // "long" or "short"
	Duplex      string // "long", "short", or "none"
	Pad         string // "auto" or "none"
}

// CreateBooklet creates a booklet-imposed PDF from an input PDF.
func CreateBooklet(inputPath, outputPath string, opts BookletOptions) error {
	// Validate output path is writable
	if err := validateOutputPath(outputPath); err != nil {
		return err
	}

	// Read and validate input PDF
	reader, err := GetPDFReader(inputPath)
	if err != nil {
		return err
	}

	info, err := GetPDFInfo(reader, inputPath)
	if err != nil {
		return err
	}

	// Validate uniform page sizes
	firstBox, err := ExtractPageBox(info.PageBoundaries[0])
	if err != nil {
		return err
	}

	for i := 1; i < len(info.PageBoundaries); i++ {
		box, err := ExtractPageBox(info.PageBoundaries[i])
		if err != nil {
			return err
		}
		if !approxEqual(box.Width, firstBox.Width, tolerance) ||
			!approxEqual(box.Height, firstBox.Height, tolerance) {
			return utils.NewErr("input PDF has mixed page sizes", ErrMixedPageSizes)
		}
	}

	// Determine sheet size
	sheetSize, err := determineSheetSize(opts.SheetSize, firstBox)
	if err != nil {
		return err
	}

	// Validate or pad page count
	pageCount := info.PageCount
	if pageCount%4 != 0 {
		if opts.Pad == "none" {
			return utils.NewErr(
				fmt.Sprintf("page count (%d) is not a multiple of 4 and padding is disabled", pageCount),
				ErrInvalidPageCount,
			)
		}
		// With auto padding, pdfcpu handles this automatically
	}

	// Configure booklet
	conf := GetPDFConfig()
	nup, err := createBookletConfig(sheetSize, opts.BindingEdge, conf)
	if err != nil {
		return err
	}

	// Create booklet
	err = pdfapi.BookletFile([]string{inputPath}, outputPath, nil, nup, conf)
	if err != nil {
		return utils.WrapErr("failed to create booklet", ErrInvalidPDF, err)
	}

	return nil
}

// validateOutputPath checks if the output path is writable.
func validateOutputPath(outputPath string) error {
	dir := filepath.Dir(outputPath)
	
	// Check if directory exists
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.NewErr(fmt.Sprintf("output directory %q does not exist", dir), ErrOutputNotWritable)
		}
		return utils.WrapErr(fmt.Sprintf("cannot access output directory %q", dir), ErrOutputNotWritable, err)
	}
	
	if !info.IsDir() {
		return utils.NewErr(fmt.Sprintf("output directory %q is not a directory", dir), ErrOutputNotWritable)
	}

	// Try to check write permissions by attempting to create a temp file
	testFile := filepath.Join(dir, ".folio_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return utils.NewErr(fmt.Sprintf("output directory %q is not writable", dir), ErrOutputNotWritable)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}

// determineSheetSize infers or validates the sheet size.
func determineSheetSize(requestedSize string, inputBox *PageBox) (string, error) {
	// If explicitly provided, validate and use it
	if requestedSize != "" {
		requestedSize = strings.ToLower(requestedSize)
		if requestedSize != "letter" && requestedSize != "a4" {
			return "", utils.NewErr(
				fmt.Sprintf("invalid sheet size %q, must be 'letter' or 'a4'", requestedSize),
				ErrCannotInferSize,
			)
		}
		// Validate that the requested size is compatible with input
		// For booklet, input should be half the sheet size
		sheetWidth, sheetHeight := getSheetDimensions(requestedSize)
		expectedWidth := sheetWidth / 2
		expectedHeight := sheetHeight
		
		// Check if input is portrait half-sheet
		if inputBox.Orientation == Orientations.Portrait &&
			approxEqual(inputBox.Width, expectedWidth, tolerance*2) &&
			approxEqual(inputBox.Height, expectedHeight, tolerance*2) {
			return requestedSize, nil
		}
		
		// Check if input is landscape half-sheet  
		if inputBox.Orientation == Orientations.Landscape &&
			approxEqual(inputBox.Width, expectedHeight, tolerance*2) &&
			approxEqual(inputBox.Height, expectedWidth, tolerance*2) {
			return requestedSize, nil
		}
		
		// Input matches the full sheet size, which is valid for booklet
		if (approxEqual(inputBox.Width, sheetWidth, tolerance*2) && 
			approxEqual(inputBox.Height, sheetHeight, tolerance*2)) ||
		   (approxEqual(inputBox.Width, sheetHeight, tolerance*2) && 
			approxEqual(inputBox.Height, sheetWidth, tolerance*2)) {
			return requestedSize, nil
		}
		
		// Allow if it's a known size that could work
		if inputBox.Size == KnownSizes.Letter && requestedSize == "letter" {
			return requestedSize, nil
		}
		if inputBox.Size == KnownSizes.A4 && requestedSize == "a4" {
			return requestedSize, nil
		}
	}

	// Try to infer from input
	if inputBox.Size == KnownSizes.Letter {
		return "letter", nil
	}
	if inputBox.Size == KnownSizes.A4 {
		return "a4", nil
	}

	// Cannot infer
	return "", utils.NewErr(
		fmt.Sprintf("cannot infer sheet size from input page size (%.0fx%.0f pt), please specify --sheet-size", 
			inputBox.Width, inputBox.Height),
		ErrCannotInferSize,
	)
}

// getSheetDimensions returns the width and height for a sheet size.
func getSheetDimensions(size string) (width, height float64) {
	switch size {
	case "letter":
		return LetterWidth, LetterHeight
	case "a4":
		return A4Width, A4Height
	default:
		return 0, 0
	}
}

// createBookletConfig creates a pdfcpu NUp configuration for booklet.
func createBookletConfig(sheetSize, bindingEdge string, conf *pdfmodel.Configuration) (*pdfmodel.NUp, error) {
	// Map sheet size to pdfcpu format
	pdfcpuSize := strings.Title(sheetSize) // "Letter" or "A4"
	desc := fmt.Sprintf("formsize:%s", pdfcpuSize)

	// Create base booklet configuration
	nup, err := pdfapi.PDFBookletConfig(2, desc, conf)
	if err != nil {
		return nil, utils.WrapErr("failed to create booklet config", ErrInvalidPDF, err)
	}

	// Set binding edge
	if bindingEdge == "short" {
		nup.BookletBinding = pdfmodel.ShortEdge
	} else {
		nup.BookletBinding = pdfmodel.LongEdge
	}

	return nup, nil
}
