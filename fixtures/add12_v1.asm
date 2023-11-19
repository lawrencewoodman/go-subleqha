                ; Version 1
                ; Assumes a and b will only ever contain a 12-bit value
                ; Uses and overflow check rather than AND
                ; ADD12

                ; ADD a to b
                a z
                z b
                z z

                ; IF b >= 4096 JUMP to overflow
                b l4096 overflow
                ; ELSE JUMP to done
                z z done

overflow:       ; Remove overflow amount
                l4096c b

done:           ; Restore l4096
                l4096 l4096
                lm4096 l4096

                ; HLT
                lm1 0

.data
z:      0
a:      4094
b:      6
l4096:  4096
l4096c: 4096    ; Used because l4096 temporarily gets corrupted
lm1:    -1
lm4096: -4096
