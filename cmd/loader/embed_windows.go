//go:build windows
// +build windows

package main

import (
	_ "embed"
)

//go:embed core.exe
var embeddedMainApp []byte
