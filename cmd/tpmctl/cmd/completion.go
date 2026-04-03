package cmd

import (
	"io"
	"strings"

	"github.com/urfave/cli/v3"
)

// bashCompletionWriter is a wrapper around an io.Writer that removes the "-o bashdefault -o default" options from the completion script output.
type bashCompletionWriter struct {
	io.Writer
}

func (w bashCompletionWriter) Write(p []byte) (int, error) {
	patched := strings.ReplaceAll(string(p), "-o bashdefault -o default", "")

	_, err := w.Writer.Write([]byte(patched))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// alreadyCompleted checks if the command has already been completed by comparing
// the number of arguments provided with the number of expected arguments.
func alreadyCompleted(cmd *cli.Command) bool {
	return cmd.Args().Len() >= len(cmd.Arguments)
}
