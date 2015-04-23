package toolkit

import "errors"

// Errors
var (
	ErrCompile      = errors.New("gsweb app compile error")
	ErrApp          = errors.New("gsweb app error")
	ErrAppDuplicate = errors.New("gsweb app error")
	ErrAppNotFound  = errors.New("gsweb app error")
)
