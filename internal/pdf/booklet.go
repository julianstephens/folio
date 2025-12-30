package pdf

import (
	"errors"
	"fmt"

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
//
// Padding behavior is controlled by opts.Pad:
//   - If opts.Pad == "none" and the input page count is not a multiple of 4,
//     CreateBooklet returns ErrInvalidPageCount without modifying the file.
//   - For any other value (including "auto"), the underlying pdfcpu library
//     automatically pads the document as needed, for example by inserting
//     blank pages, so that the page count is suitable for booklet imposition.
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
		return utils.NewErr("failed to create booklet", err)
	}

	return nil
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
