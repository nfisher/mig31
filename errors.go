package main

import (
	"fmt"
	"os"
)

type configError struct {
	message  string
	exitCode int
}

func newError(m string, rc int) (err *configError) {
	err = &configError{message: m, exitCode: rc}
	return
}

func (e *configError) Error() string {
	return e.message
}

func (e *configError) ExitCode() int {
	return e.exitCode
}

// ExitWithError prints the error and exits using the supplied exit code.
func ExitWithError(err error, exitCode int) {
	ExitWithMessage(err.Error(), exitCode)
}

// ExitWithMessage prints the message and exits using the supplied exit code.
func ExitWithMessage(message string, exitCode int) {
	fmt.Println(message)
	os.Exit(exitCode)
}
