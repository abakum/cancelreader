//go:build !windows
// +build !windows

package main

// можно импортировать любые модули

const (
	arg0 = "bash"
	arg1 = "-c"
	arg2 = "echo Press any key to continue . . .;read -rn1"
)

func ConsoleCP(*bool) {}
