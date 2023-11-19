        ; Version 2
        ; SUBLEQ emulator
        ; Store PC including memBase

        ; Fetch operands
fetch:  memA memA
        [pc] z
        z memA
        mmemBase memA
        z z
        ; opB and memB
        lm1 pc
        opB opB
        memB memB
        mmemBase memB
        [pc] z
        z opB
        z memB
        ; z z - moved to bge1000

        ; opC
        ; Store opC as a negative number to make mov to PC quicker
        lm1 pc
        opC opC
        [pc] opC

exec:   ; Execute
        [memA] [memB]

        ; If opB == 1000 THEN halt
        opB l1000 bge1000
        z z ifJmp

bge1000: z z
        l1000 z halt

        ; If mem[opB] <= 0 THEN jump to opC
ifJmp:  l1000 l1000
        lm1000 l1000

        z [memB] jmpC

incPC:  lm1 pc
        z z fetch

jmpC:   pc pc
        mmemBase pc
        opC pc
        z z fetch

        ; HLT
halt:   [memB] z
        z hltVal
        z z
        lm1 0

.data
lm1:    -1
lm1000: -1000
l1000:  1000
z:      0
; TODO: work out why need 0+program instead of just program
pc:     0+program
hltVal: 0
mmemBase: 0-program
opB:     0
opC:     0
memA:    0
memB:    0

; loopuntil_v1 from subleq/fixtures/
program:
15
13
3
16
14
6
16
13
3
16 
1000
12
0
0
sum: 0
4999
-1 