/*
 * A simple assembler for SUBLEQ2HA virtual machine
 *
 * Copyright (C) 2023 Lawrence Woodman
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package subleq2ha

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func readFile(filename string) ([]string, error) {
	lines := make([]string, 0)

	f, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return []string{}, err
	}
	return lines, nil
}

// Regular expressions for parts of a line
var reLabel = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-z]*):`)
var reInstr2 = regexp.MustCompile(`^\s*([\[\]0-9a-zA-Z\-\+]+)\s+([\[\]0-9a-zA-Z\-\+]+)`)
var reInstr3 = regexp.MustCompile(`^\s*([\[\]0-9a-zA-Z\-\+]+)\s+([\[\]0-9a-zA-Z\-\+]+)\s+([\[\]0-9a-zA-Z\-\+]+)`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+).*`)
var reExpr = regexp.MustCompile(`^\s*([0-9a-zA-z]+)([\-\+])([0-9a-zA-Z]+)`)
var reDirective = regexp.MustCompile(`^\.([a-zA-Z]+)$`)
var reIndirect = regexp.MustCompile(`^\s*(\[([0-9a-zA-z\-\+]+)\])`)

// Build symbol tables
func pass1(srcLines []string, ioSize int64) (map[string]int64, map[string]int64) {
	symbolType := "c"
	var dataPos int64 = ioSize
	var codePos int64 = 0
	codeSymbols := make(map[string]int64, 0)
	dataSymbols := make(map[string]int64, 0)
	for _, line := range srcLines {
		// If there is a directive
		if reDirective.MatchString(line) {
			directive := reDirective.FindStringSubmatch(line)[1]
			if directive == "data" {
				symbolType = "d"
			} else {
				panic(fmt.Sprintf("unknown directive: .%s", directive))
			}
		}
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			if symbolType == "c" {
				codeSymbols[label] = codePos
			} else {
				dataSymbols[label] = dataPos
			}
			line = line[matchIndices[1]:]
		}
		if reInstr3.MatchString(line) || reInstr2.MatchString(line) {
			// If there is an instruction
			codePos += 3
		} else if reExpr.MatchString(line) {
			// If there is an expression
			dataPos++
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			dataPos++
		}

	}
	return codeSymbols, dataSymbols
}

func pass2(srcLines []string, codeSymbols, dataSymbols map[string]int64) ([]int64, []int64) {
	outputType := "c"
	code := make([]int64, 0)
	data := make([]int64, ioSize)
	codePos := 0
	lineNum := 0

	for _, line := range srcLines {
		// fmt.Printf("%s\n", line)
		lineNum++
		// If there is a directive
		if reDirective.MatchString(line) {
			directive := reDirective.FindStringSubmatch(line)[1]
			if directive == "data" {
				outputType = "d"
			} else {
				panic(fmt.Sprintf("unknown directive: .%s", directive))
			}
		}
		// If there is a label
		if reLabel.MatchString(line) {
			// Remove from line
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			line = line[matchIndices[1]:]
		}
		if reInstr3.MatchString(line) {
			// If there is a 3 operand instruction
			operandA := reInstr3.FindStringSubmatch(line)[1]
			operandB := reInstr3.FindStringSubmatch(line)[2]
			operandC := reInstr3.FindStringSubmatch(line)[3]
			matchIndices := reInstr3.FindStringSubmatchIndex(line)
			line = line[matchIndices[7]:]
			if len(line) > 0 {
				panic(fmt.Sprintf("remaining line: %s", line))
			}
			code = append(code, asmInstr(codeSymbols, dataSymbols, operandA, operandB, operandC)...)
			codePos += 3
		} else if reInstr2.MatchString(line) {
			// If there is a 2 operand instruction
			operandA := reInstr2.FindStringSubmatch(line)[1]
			operandB := reInstr2.FindStringSubmatch(line)[2]
			operandC := strconv.Itoa(codePos + 3)
			matchIndices := reInstr2.FindStringSubmatchIndex(line)
			line = line[matchIndices[5]:]
			if len(line) > 0 {
				panic(fmt.Sprintf("line number: %d, remaining line: %s", lineNum, line))
			}
			code = append(code, asmInstr(codeSymbols, dataSymbols, operandA, operandB, operandC)...)
			codePos += 3
		} else if reExpr.MatchString(line) {
			// If there is an expression
			a := reExpr.FindStringSubmatch(line)[1]
			op := reExpr.FindStringSubmatch(line)[2]
			b := reExpr.FindStringSubmatch(line)[3]
			if outputType == "d" {
				data = append(data, resolveExpr(codeSymbols, dataSymbols, a, op, b))
			} else {
				panic("expression shouldn't appear here")
			}

		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			num, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				panic(err)
			}
			if outputType == "d" {
				data = append(data, num)
			} else {
				panic("literal shouldn't appear here")
			}
		}

	}

	// Add an infinite loop at the end
	// TODO: consider an instruction which will raise an error / exception
	// TODO: do we need this guard, look at alternative
	code = append(code, asmInstr(codeSymbols, dataSymbols, "0", "0", fmt.Sprintf("%d", codePos))...)
	return code, data
}

func resolveOperand(codeSymbols, dataSymbols map[string]int64, operand string) int64 {
	if reIndirect.MatchString(operand) {
		s := reIndirect.FindStringSubmatch(operand)[2]
		return 0 - resolveOperand(codeSymbols, dataSymbols, s)
	} else if reExpr.MatchString(operand) {
		// If operand is an expression
		a := reExpr.FindStringSubmatch(operand)[1]
		op := reExpr.FindStringSubmatch(operand)[2]
		b := reExpr.FindStringSubmatch(operand)[3]
		return resolveExpr(codeSymbols, dataSymbols, a, op, b)
	} else if reLiteral.MatchString(operand) {
		// If operand is a literal value
		i64, err := strconv.ParseInt(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		return i64
	}
	addr, err := resolveSymbol(codeSymbols, dataSymbols, operand)
	if err != nil {
		panic(err)
	}
	return addr
}

func resolveSymbol(codeSymbols, dataSymbols map[string]int64, sym string) (int64, error) {
	v, ok := codeSymbols[sym]
	if !ok {
		v, ok = dataSymbols[sym]
		if !ok {
			return v, fmt.Errorf("unknown symbol: %s", sym)
		}
	}
	return v, nil

}

func resolveExpr(codeSymbols, dataSymbols map[string]int64, a, op, b string) int64 {
	aVal := resolveOperand(codeSymbols, dataSymbols, a)
	bVal := resolveOperand(codeSymbols, dataSymbols, b)
	switch op {
	case "+":
		return aVal + bVal
	case "-":
		return aVal - bVal
	default:
		panic(fmt.Sprintf("unknown op: %s", op))
	}
}

func asmInstr(codeSymbols, dataSymbols map[string]int64, operandA string, operandB string, operandC string) []int64 {
	code := []int64{
		resolveOperand(codeSymbols, dataSymbols, operandA),
		resolveOperand(codeSymbols, dataSymbols, operandB),
		resolveOperand(codeSymbols, dataSymbols, operandC),
	}
	return code
}

func checkJumpsInRange(code []int64) error {
	for i := 0; i < len(code); i += 3 {
		c := code[i+2]
		if c > 0 && c >= int64(len(code)) {
			return fmt.Errorf("C operand outside code: %d", c)
		}
	}
	return nil
}

func checkMemInRange(code []int64, data []int64) error {
	for i := 0; i < len(code); i += 3 {
		a := code[i]
		b := code[i+1]
		if a > 0 && a >= int64(len(data)) {
			return fmt.Errorf("A operand outside data: %d", a)
		}
		if b > 0 && b >= int64(len(data)) {
			return fmt.Errorf("B operand outside data: %d", b)
		}

	}
	return nil
}

func printSymbols(symbols map[string]int64) {
	labels := make([]string, 0, len(symbols))
	for l := range symbols {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	fmt.Printf("Symbols\n=======\n")
	for _, l := range labels {
		fmt.Printf("%s: %d\n", l, symbols[l])
	}
	fmt.Printf("\n")
}

func asm(filename string, ioSize int64) ([]int64, []int64, map[string]int64, map[string]int64, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []int64{}, []int64{}, map[string]int64{}, map[string]int64{}, err
	}
	codeSymbols, dataSymbols := pass1(srcLines, ioSize)
	code, data := pass2(srcLines, codeSymbols, dataSymbols)
	if err := checkJumpsInRange(code); err != nil {
		return []int64{}, []int64{}, map[string]int64{}, map[string]int64{}, err
	}
	if err := checkMemInRange(code, data); err != nil {
		return []int64{}, []int64{}, map[string]int64{}, map[string]int64{}, err
	}
	// printSymbols(symbols)
	// fmt.Printf("%v\n", code)
	return code, data[ioSize:], codeSymbols, dataSymbols, nil
}
