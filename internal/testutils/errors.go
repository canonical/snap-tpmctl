package testutils

import "github.com/matryer/is"

// CheckError is a helper function to check if an error occurred based on the expectation of wantErr.
// It only returns false if no error was expected and none occurred.
func CheckError(is *is.I, err error, wantErr bool) (hadError bool) {
	is.Helper()
	hadError = true // Default to true even if is is doing assertion to mark the test as failed and panic

	if wantErr {
		is.True(err != nil)
		return hadError
	}

	is.NoErr(err)
	return false
}
