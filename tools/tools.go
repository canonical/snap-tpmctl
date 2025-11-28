// TiCS: disabled // This is to pin the tools versions that we use.

//go:build tools

package tools

import (
	_ "github.com/AlekSi/gocov-xml"
	_ "github.com/axw/gocov/gocov"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
)
