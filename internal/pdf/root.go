package pdf

import "errors"

var ErrInvalidPDF = errors.New("invalid PDF file")

type WriteOptions struct {
	SheetSize string // e.g., "letter", "a4", or ""
	Pad       string // "auto" or "none"
}
