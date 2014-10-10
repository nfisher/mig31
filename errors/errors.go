package errors


type FatalError struct {
	message  string
	exitCode int
}

func New(m string, rc int) (err *FatalError) {
	err = &FatalError{message: m, exitCode: rc}
	return
}

func (e *FatalError) Error() string {
	return e.message
}

func (e *FatalError) ExitCode() int {
	return e.exitCode
}

