package utils

import (
	"errors"
	"testing"
)

var ErrMyPackage = errors.New("mypackage error")

func TestWrapErr(t *testing.T) {
	originalErr := errors.New("database timeout")
	err := WrapErr("failed to connect", ErrMyPackage, originalErr)

	if !errors.Is(err, ErrMyPackage) {
		t.Error("WrapErr should wrap ErrMyPackage")
	}

	if !errors.Is(err, originalErr) {
		t.Error("WrapErr should preserve original error")
	}

	t.Logf("WrapErr result: %v", err)
}

func TestMultipleLevels(t *testing.T) {
	dbErr := errors.New("connection refused")
	err1 := WrapErr("query failed", ErrMyPackage, dbErr)
	err2 := WrapErr("fetch user failed", ErrMyPackage, err1)

	if !errors.Is(err2, ErrMyPackage) {
		t.Error("multiple wraps should preserve ErrMyPackage")
	}

	if !errors.Is(err2, dbErr) {
		t.Error("original error should be preserved through multiple wraps")
	}

	t.Logf("Multiple wraps result: %v", err2)
}

func TestNewErr(t *testing.T) {
	err := NewErr("failed to %s", ErrMyPackage, "process")

	if !errors.Is(err, ErrMyPackage) {
		t.Error("NewErr should wrap ErrMyPackage")
	}

	expectedMsg := "mypackage error: failed to process"
	if err.Error() != expectedMsg {
		t.Errorf("NewErr message mismatch. got: %q, want: %q", err.Error(), expectedMsg)
	}

	t.Logf("NewErr result: %v", err)
}
