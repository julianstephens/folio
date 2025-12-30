package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/julianstephens/folio/internal/utils"
)

// validateOutputPath verifies that the directory for outputPath exists, is a directory,
// and is writable by attempting to create and remove a temporary file within it.
// It returns an error wrapping ErrOutputNotWritable via utils.NewErr or
// utils.WrapErr if the directory does not exist, cannot be accessed, is not
// a directory, or is not writable. It returns nil if all checks pass.
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

	if err := f.Close(); err != nil {
		return utils.WrapErr(fmt.Sprintf("failed to close test file in output directory %q", dir), ErrOutputNotWritable, err)
	}

	if err := os.Remove(testFile); err != nil {
		return utils.WrapErr(fmt.Sprintf("output directory %q is not writable", dir), ErrOutputNotWritable, err)
	}
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
