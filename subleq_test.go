package subleqha

import (
	"fmt"
	"path/filepath"
	"testing"
)

var tests = []struct {
	filename string
	want     map[int64]int64 // [memloc]value
}{

	{"add12_v1.asm", map[int64]int64{1002: 4}},
	{"and_v1.asm", map[int64]int64{1012: 4499}},
	{"and_v2.asm", map[int64]int64{1011: 4499}},
	{"isz_v1.asm", map[int64]int64{1009: 9, 1010: 24}},
	{"jsr_v1.asm", map[int64]int64{1004: 50}},
	{"loopuntil_v1.asm", map[int64]int64{1002: 5000}},
	{"subleq_v1.asm", map[int64]int64{1027: 5000}},
	{"subleq_v2.asm", map[int64]int64{1025: 5000}},
	// TODO: reimplement switch_v1?
	//{"switch_v1.asm", map[int64]int64{3: 2255}},
	{"switch_v2.asm", map[int64]int64{1000: 2255}},
	{"switch_v3.asm", map[int64]int64{1000: 2255}},
	{"tad_v1.asm", map[int64]int64{1008: 32}},
}

func testInputHandler(operandA int64) (int64, error) {
	return 0, nil
}

func testOutputHandler(valA, operandB int64) (bool, error) {
	// Location in memory of hltVal
	// If this is used as a destination location then a HLT is executed
	const hltLoc = 0

	if operandB == hltLoc {
		return true, nil
	}
	return false, nil
}

func TestRun(t *testing.T) {
	var dataSize int64 = 31000 // The size of the data area including the io area
	var ioSize int64 = 1000    // The size of the io area before the true data area

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			code, data, codeSymbols, dataSymbols, err := Asm(filepath.Join("fixtures", test.filename), ioSize)
			if err != nil {
				t.Fatalf("asm() err: %v", err)
			}
			v := New(ioSize, dataSize, testInputHandler, testOutputHandler)
			v.LoadRoutine(code, data, codeSymbols, dataSymbols)
			if err := v.Run(); err != nil {
				t.Fatalf("Run() err: %v", err)
			}
			for dataLoc, wantValue := range test.want {
				if v.data[dataLoc] != wantValue {
					t.Errorf("mem[%d] got: %d, want: %d", dataLoc, v.data[dataLoc], wantValue)
				}
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	var dataSize int64 = 31000 // The size of the data area including the io area
	var ioSize int64 = 1000    // The size of the io area before the true data area

	for _, test := range tests {
		code, data, codeSymbols, dataSymbols, err := Asm(filepath.Join("fixtures", test.filename), ioSize)
		if err != nil {
			b.Fatalf("asm() err: %v", err)
		}

		b.Run(test.filename, func(b *testing.B) {
			b.StopTimer()

			for n := 0; n < b.N; n++ {
				v := New(ioSize, dataSize, testInputHandler, testOutputHandler)
				v.LoadRoutine(code, data, codeSymbols, dataSymbols)

				b.StartTimer()
				err := v.Run()
				b.StopTimer()

				if err != nil {
					b.Errorf("Run() err: %v", err)
				}
				for dataLoc, wantValue := range test.want {
					if v.data[dataLoc] != wantValue {
						b.Errorf("mem[%d] got: %d, want: %d", dataLoc, v.data[dataLoc], wantValue)
					}
				}
			}
		})
		fmt.Printf("Routine: %s size: %d\n", test.filename, len(code)+len(data))
	}
}
