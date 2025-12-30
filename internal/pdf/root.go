package pdf

import "errors"

var (
	ErrInvalidPDF        = errors.New("invalid PDF file")
	ErrMixedPageSizes    = errors.New("mixed page sizes")
	ErrCannotInferSize   = errors.New("cannot infer sheet size")
	ErrBookletConfig     = errors.New("invalid booklet configuration")
	ErrNupConfig         = errors.New("invalid nup configuration")
	ErrInvalidPageCount  = errors.New("invalid page count")
	ErrOutputNotWritable = errors.New("output path not writable")
)

type WriteOptions struct {
	SheetSize string // e.g., "letter", "a4", or ""
	Pad       string // "auto" or "none"
}
