//go:build !windows
// +build !windows

package main

func ConsoleCP(*bool)                  {}
func IsCygwinTerminal(fd uintptr) bool { return false }
func consoleMode() (reset func()) {
	return func() {}
}

func prepareConsole(uintptr) (reset func() error, err error) {
	reset = func() error { return nil }
	return
}
