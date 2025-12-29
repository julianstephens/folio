package pdf_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/julianstephens/folio/internal/pdf"
	testhelpers "github.com/julianstephens/folio/internal/tests"
)

func TestIsPDFFile(t *testing.T) {
	var tests []struct {
		name     string
		filePath string
		err      error
	}
	for _, fname := range testhelpers.VALID_PDFS {
		tests = append(tests, struct {
			name     string
			filePath string
			err      error
		}{
			name:     "valid: " + fname,
			filePath: filepath.Join(testhelpers.TESTDATA_DIR, fname),
			err:      nil,
		})
	}
	tests = append(tests, struct {
		name     string
		filePath string
		err      error
	}{
		name:     "invalid: " + "empty.pdf",
		filePath: filepath.Join(testhelpers.TESTDATA_DIR, "empty.pdf"),
		err:      pdf.ErrInvalidPDF,
	})
	tests = append(tests, struct {
		name     string
		filePath string
		err      error
	}{
		name:     "invalid: " + "truncated.pdf",
		filePath: filepath.Join(testhelpers.TESTDATA_DIR, "truncated.pdf"),
		err:      pdf.ErrInvalidPDF,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pdf.IsPDFFile(tt.filePath)
			if !errors.Is(err, tt.err) {
				t.Errorf("IsPDFFile(%q) returned error: %v, want error: %v", tt.filePath, err, tt.err)
			}
		})
	}
}

func TestGetPDFReader(t *testing.T) {
	var tests []struct {
		name     string
		filePath string
		err      error
	}
	for _, fname := range testhelpers.VALID_PDFS {
		tests = append(tests, struct {
			name     string
			filePath string
			err      error
		}{
			name:     "valid: " + fname,
			filePath: filepath.Join(testhelpers.TESTDATA_DIR, fname),
			err:      nil,
		})
	}
	tests = append(tests, struct {
		name     string
		filePath string
		err      error
	}{
		name:     "invalid: " + "not_a_pdf.pdf",
		filePath: filepath.Join(testhelpers.TESTDATA_DIR, "not_a_pdf.pdf"),
		err:      pdf.ErrInvalidPDF,
	})
	tests = append(tests, struct {
		name     string
		filePath string
		err      error
	}{
		name:     "invalid: " + "encrypted.pdf",
		filePath: filepath.Join(testhelpers.TESTDATA_DIR, "encrypted.pdf"),
		err:      pdf.ErrInvalidPDF,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pdf.GetPDFReader(tt.filePath)
			if !errors.Is(err, tt.err) {
				t.Errorf("GetPDFReader(%q) returned error: %v, want error: %v", tt.filePath, err, tt.err)
			}
		})
	}
}

func TestGetPDFInfo(t *testing.T) {
	for _, fname := range testhelpers.VALID_PDFS {
		t.Run("info: "+fname, func(t *testing.T) {
			filePath := filepath.Join(testhelpers.TESTDATA_DIR, fname)
			reader, err := pdf.GetPDFReader(filePath)
			if err != nil {
				t.Fatalf("GetPDFReader(%q) returned error: %v", filePath, err)
			}
			info, err := pdf.GetPDFInfo(reader, filePath)
			if err != nil {
				t.Fatalf("GetPDFInfo(%q) returned error: %v", filePath, err)
			}
			if info.PageCount == 0 {
				t.Errorf("GetPDFInfo(%q) returned PageCount=0, want >0", filePath)
			}
			if len(info.PageDimensions) == 0 {
				t.Errorf("GetPDFInfo(%q) returned no PageDimensions, want at least one", filePath)
			}
		})
	}

	skipped := map[string]bool{
		"mixed_sizes.pdf": true, // mixed page sizes are valid
	}

	for _, fname := range testhelpers.INVALID_PDFS {
		if skipped[fname] {
			continue
		}

		t.Run("invalid info: "+fname, func(t *testing.T) {
			filePath := filepath.Join(testhelpers.TESTDATA_DIR, fname)
			reader, err := pdf.GetPDFReader(filePath)
			if err == nil {
				t.Fatalf("GetPDFReader(%q) expected error, got nil", filePath)
			}
			if reader != nil {
				_, err = pdf.GetPDFInfo(reader, filePath)
				if err == nil {
					t.Fatalf("GetPDFInfo(%q) expected error, got nil", filePath)
				}
			}
		})
	}
}
