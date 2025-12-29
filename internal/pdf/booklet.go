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
	ErrMixedPageSizes    = errors.New("mixed page sizes")
	ErrCannotInferSize   = errors.New("cannot infer sheet size")
	ErrBookletConfig     = errors.New("invalid booklet configuration")
	ErrInvalidPageCount  = errors.New("invalid page count")
	ErrOutputNotWritable = errors.New("output path not writable")
)

// BookletOptions contains configuration for booklet creation.
type BookletOptions struct {
	WriteOptions
	BindingEdge string // "long" or "short"
}

// CreateBooklet creates a booklet-imposed PDF from an input PDF.
func CreateBooklet(inputPath, outputPath string, opts BookletOptions) error {
	// Validate output path is writable
	if err := validateOutputPath(outputPath); err != nil {
		return err
	}

	reader, err := GetPDFReader(inputPath)
	if err != nil {
		return err
	}

	info, err := GetPDFInfo(reader, inputPath)
	if err != nil {
		return err
	}

	if len(info.PageBoundaries) == 0 {
		return utils.NewErr("PDF contains no pages", ErrInvalidPDF)
	}

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

	sheetSize, err := determineSheetSize(opts.SheetSize, firstBox)
	if err != nil {
		return err
	}

	pageCount := info.PageCount
	if pageCount%4 != 0 {
		if opts.Pad == "none" {
			return utils.NewErr(
				fmt.Sprintf("page count (%d) is not a multiple of 4 and padding is disabled", pageCount),
				ErrInvalidPageCount,
			)
		}
	}

	conf := GetPDFConfig()
	nup, err := createBookletConfig(sheetSize, opts.BindingEdge, conf)
	if err != nil {
		return err
	}

	err = pdfapi.BookletFile([]string{inputPath}, outputPath, nil, nup, conf)
	if err != nil {
		return utils.WrapErr("failed to create booklet", ErrInvalidPDF, err)
	}

	return nil
}

// validateOutputPath checks if the output path is writable.
func validateOutputPath(outputPath string) error {
	dir := filepath.Dir(outputPath)

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

	testFile := filepath.Join(dir, ".folio_write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return utils.NewErr(fmt.Sprintf("output directory %q is not writable", dir), ErrOutputNotWritable)
	}
	defer f.Close()
	defer os.Remove(testFile)

	return nil
}

// determineSheetSize infers or validates the sheet size.
func determineSheetSize(requestedSize string, inputBox *PageBox) (string, error) {
	if requestedSize != "" {
		requestedSize = strings.ToLower(requestedSize)
		if requestedSize != "letter" && requestedSize != "a4" {
			return "", utils.NewErr(
				fmt.Sprintf("invalid sheet size %q, must be 'letter' or 'a4'", requestedSize),
				ErrCannotInferSize,
			)
		}
		return requestedSize, nil
	}

	if inputBox.Size == KnownSizes.Letter {
		return "letter", nil
	}
	if inputBox.Size == KnownSizes.A4 {
		return "a4", nil
	}

	return "", utils.NewErr(
		fmt.Sprintf("cannot infer sheet size from input page size (%.0fx%.0f pt), please specify --sheet-size",
			inputBox.Width, inputBox.Height),
		ErrCannotInferSize,
	)
}

// createBookletConfig creates a pdfcpu NUp configuration for booklet.
func createBookletConfig(sheetSize, bindingEdge string, conf *pdfmodel.Configuration) (*pdfmodel.NUp, error) {
	var pdfcpuSize string
	switch sheetSize {
	case "letter":
		pdfcpuSize = "Letter"
	case "a4":
		pdfcpuSize = "A4"
	default:
		return nil, utils.NewErr(fmt.Sprintf("unsupported sheet size: %s", sheetSize), ErrCannotInferSize)
	}
	desc := fmt.Sprintf("formsize:%s", pdfcpuSize)

	nup, err := pdfapi.PDFBookletConfig(2, desc, conf)
	if err != nil {
		return nil, utils.WrapErr("failed to create booklet config", ErrBookletConfig, err)
	}

	if bindingEdge == "short" {
		nup.BookletBinding = pdfmodel.ShortEdge
	} else {
		nup.BookletBinding = pdfmodel.LongEdge
	}

	return nup, nil
}
