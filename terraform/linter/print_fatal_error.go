package linter

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// PrintFatalError prints the error based on the type and
// if stack traces should be shown then exits.
func PrintFatalError(err error, hideStackTraces bool) {
	errorVerbosityFormat := "%+v"
	if hideStackTraces {
		errorVerbosityFormat = "%s"
	}

	white := color.New(color.FgWhite, color.Bold)
	formattedError := white.Sprintf(errorVerbosityFormat, err)

	if linterErr, ok := err.(Error); ok {
		red := color.New(color.FgRed, color.Bold).SprintfFunc()
		formattedError = fmt.Sprintf(
			"%s - %s",
			red(linterErr.Resource),
			white.Sprintf(errorVerbosityFormat, linterErr.Cause()),
		)
	}

	fmt.Printf("\n%s\n", formattedError)
	os.Exit(1)
}
