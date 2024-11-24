//go:build windows
// +build windows

package main

import (
	"log"
	"os"

	"github.com/abakum/cancelreader"
	"github.com/mattn/go-isatty"
	"github.com/xlab/closer"
	"golang.org/x/sys/windows"
)

func ConsoleCP(once *bool) {
	if *once {
		return
	}
	*once = false
	const CP_UTF8 uint32 = 65001
	var kernel32 = windows.NewLazyDLL("kernel32.dll")

	getConsoleCP := func() uint32 {
		result, _, _ := kernel32.NewProc("GetConsoleCP").Call()
		return uint32(result)
	}

	getConsoleOutputCP := func() uint32 {
		result, _, _ := kernel32.NewProc("GetConsoleOutputCP").Call()
		return uint32(result)
	}

	setConsoleCP := func(cp uint32) {
		kernel32.NewProc("SetConsoleCP").Call(uintptr(cp))
	}

	setConsoleOutputCP := func(cp uint32) {
		kernel32.NewProc("SetConsoleOutputCP").Call(uintptr(cp))
	}

	inCP := getConsoleCP()
	outCP := getConsoleOutputCP()
	setConsoleCP(CP_UTF8)
	setConsoleOutputCP(CP_UTF8)
	closer.Bind(func() { setConsoleCP(inCP) })
	closer.Bind(func() { setConsoleOutputCP(outCP) })
}

func IsCygwinTerminal(fd uintptr) bool {
	return isatty.IsCygwinTerminal(fd)
}

func consoleMode() (reset func()) {
	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		settings, err := sttySettings()
		if err != nil {
			return func() {}
		}
		log.Println("sttySettings", settings)
		return func() {
			log.Println("sttyReset", settings)
			sttyReset(settings)
		}
	}
	input := windows.Handle(os.Stdin.Fd())
	var originalMode uint32

	err := windows.GetConsoleMode(input, &originalMode)
	if err != nil {
		return func() {}
	}
	log.Printf("GetConsoleMode %x\r\n", originalMode)
	return func() {
		log.Printf("SetConsoleMode %x\r\n", originalMode)
		windows.SetConsoleMode(input, originalMode)
	}
}

func prepareConsole(uinput uintptr) (reset func() error, err error) {
	log.Println("prepareConsole", isatty.IsCygwinTerminal(os.Stdin.Fd()))
	reset = func() error { return nil }
	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		s := ""
		s, err = sttySettings()
		log.Println("sttySettings", err, s)
		if err != nil {
			return
		}
		err = sttyMakeRaw()
		log.Println("sttyMakeRaw", err)
		if err != nil {
			return
		}
		log.Println(sttySettings())

		return func() error {
			sttyReset(s)
			return nil
		}, nil
	}
	return cancelreader.PrepareConsole(uinput)
}
