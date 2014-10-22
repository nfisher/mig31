package runtime

import (
	"fmt"
	"os"
)

const (
	ExitIncorrectFlag = iota + 1
	ExitErrorReadingEnvConfig
	ExitNoEnvironmentDefined
	ExitErrorReadingMigrations
	ExitUnableToCreateSchema
	ExitErrorGeneratingOfflineSchema
	ExitMigrationMismatch
)

// ExitWithError prints the error and exits using the supplied exit code.
func ExitWithError(err error, exitCode int) {
	ExitWithMessage(err.Error(), exitCode)
}

// ExitWithMessage prints the message and exits using the supplied exit code.
func ExitWithMessage(message string, exitCode int) {
	fmt.Println(message)
	os.Exit(exitCode)
}
