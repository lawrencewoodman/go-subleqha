/*
 * SUBLEQHA Virtual Machine
 *
 * A SUBLEQHA VM using a Harvard Architecture to keep code and data in
 * separate spaces.  Negative numbers passed as operands are indirect
 * addresses.
 *
 * Copyright (C) 2023 Lawrence Woodman
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package subleqha

import (
	"fmt"
)

// Location in memory of hltVal
// If this is used as a destination location then a HLT is executed
const hltLoc = 0

type SUBLEQHA struct {
	ioSize        int64
	dataSize      int64
	codeSize      int64            // The size of the code / program
	code          []int64          // Code / Program
	data          []int64          // Data
	pc            int64            // Program Counter
	hltVal        int64            // A value returned by HLT
	codeSymbols   map[string]int64 // The code symbols table from the assembler - to aid debugging
	dataSymbols   map[string]int64 // The data symbols table from the assembler - to aid debugging
	inputHandler  InputHandler
	outputHandler OutputHandler
}

type InputHandler func(operandA int64) (int64, error)
type OutputHandler func(valA, operandB int64) (bool, error)

func New(ioSize, dataSize int64, inputHandler InputHandler, outputHandler OutputHandler) *SUBLEQHA {
	return &SUBLEQHA{
		data:          make([]int64, dataSize),
		ioSize:        ioSize,
		dataSize:      dataSize,
		codeSize:      0,
		inputHandler:  inputHandler,
		outputHandler: outputHandler,
	}
}

func (v *SUBLEQHA) Run() error {
	var hlt bool = false
	var valA int64 = 0
	var err error

	for !hlt {
		if v.pc+2 >= v.codeSize {
			return fmt.Errorf("PC: %d, outside code range", v.pc)
		}
		operandA := v.code[v.pc]
		operandB := v.code[v.pc+1]
		operandC := v.code[v.pc+2]

		if operandA < 0 {
			operandA = v.data[-operandA]
			if operandA < 0 {
				return fmt.Errorf("PC: %d, double indirect not supported", v.pc)
			}
			if operandA >= v.dataSize {
				return fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandA)
			}
		}

		if operandB < 0 {
			operandB = v.data[-operandB]
			if operandB < 0 {
				return fmt.Errorf("PC: %d, double indirect not supported", v.pc)
			}
			if operandB >= v.dataSize {
				return fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandB)
			}
		}
		if operandC < 0 {
			operandC = v.data[-operandC]
			if operandC < 0 {
				return fmt.Errorf("PC: %d, double indirect not supported", v.pc)
			}
		}

		//fmt.Printf("PC: %7s    SUBLEQHA %s, %s, %s\n", v.addr2symbol(v.pc, true), v.addr2symbol(operandA), v.addr2symbol(operandB), v.addr2symbol(operandC, true))
		//fmt.Printf("                      %d - %d = ", v.data[operandB], v.data[operandA])

		// If an IO operation
		if operandA < v.ioSize {
			valA, err = v.inputHandler(operandA)
			if err != nil {
				// TODO: Add description to error
				return err
			}
		} else {
			valA = v.data[operandA]
		}

		// If an IO operation
		if operandB < v.ioSize {
			hlt, err = v.outputHandler(valA, operandB)
			if err != nil {
				// TODO: Add description to error
				return err
			}
			if hlt {
				v.hltVal = -valA
			}
			v.pc += 3
		} else {
			v.data[operandB] -= valA
			if v.data[operandB] <= 0 {
				v.pc = operandC
			} else {
				v.pc += 3
			}
			//fmt.Printf("%d\n", v.data[operandB])
		}

	}
	return nil
}

func (v *SUBLEQHA) LoadRoutine(code []int64, data []int64, codeSymbols map[string]int64, dataSymbols map[string]int64) {
	v.codeSize = int64(len(code))
	v.code = make([]int64, v.codeSize)
	copy(v.code, code)
	copy(v.data[v.ioSize:], data)
	v.codeSymbols = codeSymbols
	v.dataSymbols = dataSymbols
}

func (v *SUBLEQHA) addr2symbol(addr int64, onlyCode ...bool) string {
	if len(onlyCode) == 0 {
		for k, v := range v.dataSymbols {
			if v == addr {
				return k
			}
		}
	}

	for k, v := range v.codeSymbols {
		if v == addr {
			return k
		}
	}
	return fmt.Sprintf("%d", addr)
}
