package pdf

import (
	"bytes"
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

	var fp string
	if filePath != nil {
		fp = *filePath
	} else {
		fp = "<reader>"
	}

	var fileSize int64
	var err error

	if filePath != nil {
		fileInfo, err := os.Stat(fp)
		if err != nil {
			return utils.WrapErr(fmt.Sprintf("could not stat file %q", fp), ErrInvalidPDF, err)
		}
		fileSize = fileInfo.Size()
	} else {
		if seeker, ok := (*reader).(io.Seeker); ok {
			currentPos, err := seeker.Seek(0, io.SeekCurrent)
			if err != nil {
				return utils.WrapErr("could not get current position of reader", ErrInvalidPDF, err)
			}
			endPos, err := seeker.Seek(0, io.SeekEnd)
			if err != nil {
				return utils.WrapErr("could not seek to end of reader", ErrInvalidPDF, err)
			}
			_, err = seeker.Seek(currentPos, io.SeekStart)
			if err != nil {
				return utils.WrapErr("could not restore position of reader", ErrInvalidPDF, err)
			}
			fileSize = endPos
		} else {
			return utils.NewErr("reader is not seekable", ErrInvalidPDF)
		}
	}

	if fileSize == 0 {
		return utils.NewErr(fmt.Sprintf("file %q is empty", fp), ErrInvalidPDF)
	}

	var header []byte = make([]byte, 8)
	if filePath != nil {
		file, err := os.Open(fp)
		if err != nil {
			return utils.WrapErr(fmt.Sprintf("could not open file %q", fp), ErrInvalidPDF, err)
		}
		defer file.Close()

		_, err = file.Read(header)
		if err != nil {
			return utils.WrapErr(fmt.Sprintf("could not read file %q", fp), ErrInvalidPDF, err)
		}
	} else {
		_, err = (*reader).Read(header)
		if err != nil {
			return utils.WrapErr("could not read from reader", ErrInvalidPDF, err)
		}
		_, err = (*reader).Seek(0, io.SeekStart)
		if err != nil {
			return utils.WrapErr("could not restore position of reader", ErrInvalidPDF, err)
		}
	}

	if !strings.HasPrefix(string(header), "%PDF-1.") {
		return utils.NewErr(fmt.Sprintf("file %q does not appear to be a valid PDF", fp), ErrInvalidPDF)
	}

	var bufSize int64 = 1024
	if fileSize < bufSize {
		bufSize = fileSize
	}

	var trailer []byte = make([]byte, bufSize)
	if filePath != nil {
		file, err := os.Open(fp)
		if err != nil {
			return utils.WrapErr(fmt.Sprintf("could not open file %q", fp), ErrInvalidPDF, err)
		}
		defer file.Close()

		_, err = file.ReadAt(trailer, fileSize-bufSize)
		if err != nil {
			return utils.WrapErr(fmt.Sprintf("could not read end of file %q", fp), ErrInvalidPDF, err)
		}
	} else {
		seeker, ok := (*reader).(io.Seeker)
		if !ok {
			return utils.NewErr("reader is not seekable", ErrInvalidPDF)
		}
		_, err = seeker.Seek(fileSize-bufSize, io.SeekStart)
		if err != nil {
			return utils.WrapErr("could not seek to end of reader", ErrInvalidPDF, err)
		}
		_, err = (*reader).Read(trailer)
		if err != nil {
			return utils.WrapErr("could not read from reader", ErrInvalidPDF, err)
		}
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return utils.WrapErr("could not restore position of reader", ErrInvalidPDF, err)
		}
	}

	if !strings.Contains(string(trailer), "%%EOF") {
		return utils.NewErr(fmt.Sprintf("file %q does not appear to be a valid PDF (missing EOF marker)", fp), ErrInvalidPDF)
	}

	if detectEncryptionInTail(trailer) {
		return utils.NewErr(fmt.Sprintf("file %q appears to be an encrypted PDF", fp), ErrInvalidPDF)
	}

	return nil
}

/* IsPDFFile checks if the file at filePath is a valid, non-empty PDF file. */
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

/* GetPDFConfig returns a singleton PDF configuration instance. */
func GetPDFConfig() *pdfmodel.Configuration {
	once.Do(func() {
		config = pdfmodel.NewDefaultConfiguration()
		config.SetUnit("in")
	})
	return config
}

/* GetPDFReader reads the PDF file located at filePath and returns an io.ReadSeeker for it. */
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

/* GetPDFInfo retrieves information about the PDF file located at filePath. */
func GetPDFInfo(reader io.ReadSeeker, filePath string) (info *pdfcpu.PDFInfo, err error) {
	info, err = pdfapi.PDFInfo(reader, filePath, nil, true, GetPDFConfig())
	if err != nil {
		if strings.Contains(err.Error(), "correct password") {
			err = utils.NewErr("encrypted PDFs are not supported", ErrInvalidPDF)
		}
		err = utils.WrapErr("could not get PDF info", ErrInvalidPDF, err)
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
