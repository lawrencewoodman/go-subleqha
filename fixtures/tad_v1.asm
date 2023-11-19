                ; Version 1
                ; PDP-8 TAD
                ; Assumes lac will only every contain a 13-bit value and
                ; memory will only ever contain a 12-bit value
                ; and checks for overflow rather than using AND

                ; MOV memBase+opAddr TO memLoc
                memLoc memLoc
                mmemBase memLoc
                mopAddr memLoc

                ; ADD I memLoc lac
                [memLoc] z
                z lac
                z z

                ; IF lac >= 8192 JUMP to overflow
                lac l8192 overflow
                ; ELSE JUMP to done
                z z done

overflow:       ; Remove overflow amount
                l8192c lac

done:           ; Restore l8192
                l8192 l8192
                lm8192 l8192

                ; HLT
                lm1 0

.data
z:      0
l8192:  8192
l8192c: 8192    ; Used because l8192 gets temporarily corrupted
lm1:    -1
lm8192: -8192

mmemBase: 0-1005   ; TODO: be able to put memBase or similar here
mopAddr:  0-4
memLoc:  0
lac:     9
val:     23