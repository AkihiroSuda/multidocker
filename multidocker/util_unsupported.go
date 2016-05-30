// +build !linux

package multidocker

import (
	"fmt"
	"runtime"
)

var unsupportedErr = fmt.Errorf("unsupported: %s(%s) ", runtime.GOOS, runtime.GOARCH)
