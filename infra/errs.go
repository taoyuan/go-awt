package infra

import "errors"

var ErrExecLookPathFailed = errors.New("ErrExecLookPathFailed")
var ErrExecStartFailed = errors.New("ErrExecStartFailed")
var ErrExecStdoutPipeFailed = errors.New("ErrExecStdoutPipeFailed")
