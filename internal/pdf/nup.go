package pdf

import (
	"fmt"
	"strings"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/julianstephens/folio/internal/utils"
)

// NupOptions contains configuration for n-up layout.
type NupOptions struct {
	WriteOptions
	PerSheet    int    // 2 or 4 pages per sheet
	Orientation string // "portrait", "landscape", or "" (auto)
}

// CreateNup creates an n-up PDF from an input PDF.
//
// The function arranges pages in reading order on each sheet without special reordering.
// For 2-up: pages 1-2 appear on the first sheet, 3-4 on the next, etc.
// For 4-up: pages 1-4 appear on the first sheet, 5-8 on the next, etc.
//
// Padding behavior is controlled by opts.Pad:
//   - If opts.Pad == "none" and the input page count is not a multiple of opts.PerSheet,
//     CreateNup returns ErrInvalidPageCount without modifying the file.
//   - For any other value (including "auto"), blank pages are added to make the
//     page count a multiple of opts.PerSheet.
func CreateNup(inputPath, outputPath string, opts NupOptions) error {
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
	if pageCount%opts.PerSheet != 0 {
		if opts.Pad == "none" {
			return utils.NewErr(
				fmt.Sprintf("page count (%d) is not a multiple of %d and padding is disabled", pageCount, opts.PerSheet),
				ErrInvalidPageCount,
			)
		}
	}

	orientation := opts.Orientation
	if orientation == "" {
		// Auto-detect orientation based on sheet size and n-up layout
		orientation = determineNupOrientation(sheetSize, opts.PerSheet)
	}

	conf := GetPDFConfig()
	nup, err := createNupConfig(sheetSize, orientation, opts.PerSheet, conf)
	if err != nil {
		return err
	}

	err = pdfapi.NUpFile([]string{inputPath}, outputPath, nil, nup, conf)
	if err != nil {
		return utils.NewErr("failed to create n-up layout", err)
	}

	return nil
}

// determineNupOrientation infers the best orientation for n-up layout.
// For 2-up: typically landscape (side-by-side)
// For 4-up: typically portrait (2x2 grid)
func determineNupOrientation(sheetSize string, perSheet int) string {
	// For 2-up, landscape is standard (pages side-by-side)
	if perSheet == 2 {
		return "landscape"
	}
	// For 4-up, portrait is standard (2x2 grid)
	return "portrait"
}

// createNupConfig creates a pdfcpu NUp configuration for n-up layout.
func createNupConfig(sheetSize, orientation string, perSheet int, conf *pdfmodel.Configuration) (*pdfmodel.NUp, error) {
	var pdfcpuSize string
	switch sheetSize {
	case "letter":
		pdfcpuSize = "Letter"
	case "a4":
		pdfcpuSize = "A4"
	default:
		return nil, utils.NewErr(fmt.Sprintf("unsupported sheet size: %s", sheetSize), ErrCannotInferSize)
	}

	// Add orientation suffix if specified
	// L = Landscape, P = Portrait
	orientLower := strings.ToLower(orientation)
	if orientLower == "landscape" {
		pdfcpuSize += "L"
	} else if orientLower == "portrait" {
		pdfcpuSize += "P"
	}

	// Build description string for pdfcpu
	// Format: "formsize:Letter" or "formsize:LetterL" etc
	desc := fmt.Sprintf("formsize:%s", pdfcpuSize)

	nup, err := pdfapi.PDFNUpConfig(perSheet, desc, conf)
	if err != nil {
		return nil, utils.WrapErr("failed to create nup config", ErrNupConfig, err)
	}

	return nup, nil
}
