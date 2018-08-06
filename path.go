// Copyright 2013, 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

// prefixSize is used internally to trim the user specific path from the
// front of the returned filenames from the runtime call stack.
var prefixSize int

// goPath is the deduced path based on the location of this file as compiled.
var goPath string
var srcDir string

func init() {
	goPath = getGOPATH()
	srcDir = filepath.Join(goPath, "src")
}

func getGOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func trimGoPath(filename string) string {
	return strings.Replace(filename, fmt.Sprintf("%s/", srcDir), "", 1)
}
