/*
 * Example using SUBLEQHA to echo keypresses from STDIN back to STDOUT
 *
 * Copyright (C) 2023 Lawrence Woodman
 *
 * Licensed under a BSD 0-Clause licence. Please see 0BSD_LICENCE.md for details.
 */

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/lawrencewoodman/go-subleqha"
)

func createInputHandler() (func(int64) (int64, error), io.Closer, error) {
	rt, err := NewRawTerm()
	if err != nil {
		return nil, nil, err
	}

	handler := func(operandA int64) (int64, error) {
		const stdin = 1
		key := make([]byte, 1)

		if operandA == stdin {
			n, err := rt.Read(key)
			if err != nil {
				return -1, err
			}
			if n == 1 {
				// Exit on CTRL-\ from keyboard
				if key[0] == 0x1C {
					fmt.Println("Quit")
					os.Exit(0)
					// TODO: use a flag to exit nicely
				}
				if key[0] == 'x' {
					return -1, nil
				}
				return int64(key[0]), nil
			}
		}
		return 0, nil

	}
	return handler, rt, nil
}

func testOutputHandler(valA, operandB int64) (bool, error) {
	// Location in memory of hltVal
	// If this is used as a destination location then a HLT is executed
	const hltLoc = 0
	const stdout = 1

	switch operandB {
	case hltLoc:
		return true, nil
	case stdout:
		//fmt.Printf("%c\n", -valA)
		os.Stdout.Write([]byte{byte(-valA)})
	default:
		return true, fmt.Errorf("unknown IO location for B: %d", operandB)
	}
	return false, nil
}

func main() {
	var dataSize int64 = 31000 // The size of the data area including the io area
	var ioSize int64 = 1000    // The size of the io area before the true data area

	testInputHandler, rt, err := createInputHandler()
	if err != nil {
		panic(err)
	}
	defer rt.Close()

	code, data, codeSymbols, dataSymbols, err := subleqha.Asm("echo.asm", ioSize)
	if err != nil {
		panic(fmt.Sprintf("asm() err: %v", err))
	}
	v := subleqha.New(ioSize, dataSize, testInputHandler, testOutputHandler)
	v.LoadRoutine(code, data, codeSymbols, dataSymbols)
	if err := v.Run(); err != nil {
		panic(fmt.Sprintf("Run() err: %v", err))
	}
}
