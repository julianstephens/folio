package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/julianstephens/folio/internal/utils"
)

var (
	config *pdfmodel.Configuration
	once   sync.Once
)

type pdfSource interface {
	Size() (int64, error)
	ReadAt(p []byte, off int64) (n int, err error)
}

type filePDFSource struct {
	path string
}

func (f *filePDFSource) Size() (int64, error) {
	info, err := os.Stat(f.path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (f *filePDFSource) ReadAt(p []byte, off int64) (n int, err error) {
	file, err := os.Open(f.path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return file.ReadAt(p, off)
}

type readerPDFSource struct {
	rs io.ReadSeeker
}

func (r *readerPDFSource) Size() (int64, error) {
	seeker, ok := r.rs.(io.Seeker)
	if !ok {
		return 0, errors.New("reader is not seekable")
	}
	currentPos, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, fmt.Errorf("could not get current position: %w", err)
	}
	endPos, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, fmt.Errorf("could not seek to end: %w", err)
	}
	_, err = seeker.Seek(currentPos, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("could not restore position: %w", err)
	}
	return endPos, nil
}

func (r *readerPDFSource) ReadAt(p []byte, off int64) (n int, err error) {
	seeker, ok := r.rs.(io.Seeker)
	if !ok {
		return 0, errors.New("reader is not seekable")
	}
	_, err = seeker.Seek(off, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("could not seek to offset: %w", err)
	}
	n, err = r.rs.Read(p)
	if err != nil {
		return n, err
	}
	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		return n, fmt.Errorf("could not restore position: %w", err)
	}
	return n, nil
}

func detectEncryptionInTail(tail []byte) bool {
	trailerIdx := bytes.LastIndex(tail, []byte("trailer"))
	if trailerIdx == -1 {
		return false
	}

	dictStart := bytes.Index(tail[trailerIdx:], []byte("<<"))
	if dictStart == -1 {
		return false
	}
	dictStart += trailerIdx

	dictEnd := bytes.Index(tail[dictStart:], []byte(">>"))
	if dictEnd == -1 {
		return false
	}
	dictEnd += dictStart

	dictBytes := tail[dictStart:dictEnd]
	return bytes.Contains(dictBytes, []byte("/Encrypt"))
}

func validatePDFContent(filePath *string, reader *io.ReadSeeker) error {
	if filePath == nil && reader == nil {
		return utils.NewErr("either filePath or reader must be provided", ErrInvalidPDF)
	}

	var source pdfSource
	var fp string

	if filePath != nil {
		fp = *filePath
		source = &filePDFSource{path: fp}
	} else {
		fp = "<reader>"
		source = &readerPDFSource{rs: *reader}
	}

	fileSize, err := source.Size()
	if err != nil {
		return utils.WrapErr(fmt.Sprintf("could not determine size of %s", fp), ErrInvalidPDF, err)
	}

	if fileSize == 0 {
		return utils.NewErr(fmt.Sprintf("file %q is empty", fp), ErrInvalidPDF)
	}

	header := make([]byte, 8)
	if _, err := source.ReadAt(header, 0); err != nil {
		return utils.WrapErr(fmt.Sprintf("could not read header of %s", fp), ErrInvalidPDF, err)
	}

	if !strings.HasPrefix(string(header), "%PDF-") {
		return utils.NewErr(fmt.Sprintf("file %q does not appear to be a valid PDF", fp), ErrInvalidPDF)
	}

	bufSize := int64(1024)
	if fileSize < bufSize {
		bufSize = fileSize
	}

	trailer := make([]byte, bufSize)
	if _, err := source.ReadAt(trailer, fileSize-bufSize); err != nil {
		return utils.WrapErr(fmt.Sprintf("could not read trailer of %s", fp), ErrInvalidPDF, err)
	}

	if !strings.Contains(string(trailer), "%%EOF") {
		return utils.NewErr(fmt.Sprintf("file %q does not appear to be a valid PDF (missing EOF marker)", fp), ErrInvalidPDF)
	}

	if detectEncryptionInTail(trailer) {
		return utils.NewErr(fmt.Sprintf("file %q appears to be an encrypted PDF", fp), ErrInvalidPDF)
	}

	return nil
}

// IsPDFFile checks if the file at filePath is a valid, non-empty PDF file.
func IsPDFFile(filePath string) error {
	if strings.ToLower(filepath.Ext(filePath)) != ".pdf" {
		return utils.NewErr(fmt.Sprintf("file %q is not a PDF", filePath), ErrInvalidPDF)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return utils.WrapErr(fmt.Sprintf("could not stat file %q", filePath), ErrInvalidPDF, err)
	}

	if fileInfo.Size() == 0 {
		return utils.NewErr(fmt.Sprintf("file %q is empty", filePath), ErrInvalidPDF)
	}

	if err = validatePDFContent(&filePath, nil); err != nil {
		return err
	}

	return nil
}

// GetPDFConfig returns a singleton PDF configuration instance.
func GetPDFConfig() *pdfmodel.Configuration {
	once.Do(func() {
		config = pdfmodel.NewDefaultConfiguration()
	})
	return config
}

// GetPDFReader reads the PDF file located at filePath and returns an io.ReadSeeker for it.
func GetPDFReader(filePath string) (reader io.ReadSeeker, err error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		err = utils.WrapErr(fmt.Sprintf("could not read file %q", filePath), ErrInvalidPDF, err)
		return
	}
	reader = bytes.NewReader(fileContent)
	if err = validatePDFContent(nil, &reader); err != nil {
		return
	}
	return
}

// GetPDFInfo retrieves information about the PDF file located at filePath.
func GetPDFInfo(reader io.ReadSeeker, filePath string) (info *pdfcpu.PDFInfo, err error) {
	info, err = pdfapi.PDFInfo(reader, filePath, nil, true, GetPDFConfig())
	if err != nil {
		if strings.Contains(err.Error(), "correct password") {
			err = utils.NewErr("encrypted PDFs are not supported", ErrInvalidPDF)
		} else {
			err = utils.WrapErr("could not get PDF info", ErrInvalidPDF, err)
		}
		return
	}

	if info.Encrypted {
		err = utils.NewErr("encrypted PDFs are not supported", ErrInvalidPDF)
		return
	}

	if info.PageCount == 0 || info.PageBoundaries == nil {
		err = utils.NewErr("PDF contains no pages", ErrInvalidPDF)
		return
	}

	return
}

// ExtractPageBox extracts the PageBox information from the given PageBoundaries.
func ExtractPageBox(pageBoundary pdfmodel.PageBoundaries) (*PageBox, error) {
	mediaBox := pageBoundary.MediaBox()
	if mediaBox == nil {
		return nil, utils.NewErr("media box not found", ErrInvalidPDF)
	}

	width := mediaBox.Width()
	height := mediaBox.Height()

	return &PageBox{
		Width:       width,
		Height:      height,
		Size:        DetectKnownSize(width, height),
		Orientation: DetectOrientation(width, height),
	}, nil
}
