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
)

var (
	config *pdfmodel.Configuration
	once   sync.Once
)

/* IsPDFFile checks if the file at filePath is a valid, non-empty PDF file. */
func IsPDFFile(filePath string) error {
	if strings.ToLower(filepath.Ext(filePath)) != ".pdf" {
		return fmt.Errorf("file %q is not a PDF", filePath)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("could not stat file %q: %w", filePath, err)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("file %q is empty", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read file %q: %w", filePath, err)
	}

	if !strings.HasPrefix(string(content), "%PDF-1.") {
		return fmt.Errorf("file %q does not appear to be a valid PDF", filePath)
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
		err = fmt.Errorf("could not read file %q: %w", filePath, err)
		return
	}
	reader = bytes.NewReader(fileContent)
	return
}

/* GetPDFInfo retrieves information about the PDF file located at filePath. */
func GetPDFInfo(reader io.ReadSeeker, filePath string) (info *pdfcpu.PDFInfo, err error) {
	info, err = pdfapi.PDFInfo(reader, filePath, nil, true, GetPDFConfig())
	if err != nil {
		err = fmt.Errorf("could not get PDF info: %w", err)
		return
	}

	if info.Encrypted {
		err = errors.New("encrypted PDFs are not supported")
		return
	}

	if info.PageCount == 0 || info.PageBoundaries == nil {
		err = errors.New("PDF contains no pages")
		return
	}

	return
}
