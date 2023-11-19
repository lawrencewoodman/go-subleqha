                ; Version 1
                ; PDP-8 ISZ
                ; Assumes memory will only ever contain a 12-bit value
                ; memory will only ever contain a 12-bit value

                ; MOV memBase+opAddr TO memLoc
                ; [memLoc] is will referenced as v in labels below
                memLoc memLoc
                memBase z
                opAddr z
                z memLoc
                z z

                ; Increment [memLoc]
                lm1 [memLoc]

                ; IF [memLoc] >= 4096 JUMP to voverflow
                [memLoc] l4096 voverflow
                ; ELSE JUMP to setit
                z z vresl4096

voverflow:      ; Remove overflow amount
                l4096c t

                ; Restore l4096
vresl4096:      l4096 l4096
                lm4096 l4096

                ; IF [memLoc] <= 0 JUMP to vle
                z [memLoc] vle
                ; ELSE JUMP to done
                z z done

vle:            ; IF z <= [memLoc] JUMP to teq (equivalent to IF [memLoc] == 0)
                [memLoc] z veq
                ; ELSE
                z z done

                ; [memLoc] == 0, so increment PC
veq:            lm1 pc

                ; IF pc >= 4096 JUMP to pcoverflow
                pc l4096 pcoverflow
                ; ELSE JUMP to done
                z z pcresl4096

pcoverflow:     ; Remove overflow amount
                l4096c pc

pcresl4096:     ; Restore l4096
                l4096 l4096
                lm4096 l4096

done:           ; HLT
                lm1 0

.data
z:      0
t:      0
l4096:  4096
l4096c: 4096    ; Used because l4096 temporarily gets corrupted
lm1:    -1
lm4096: -4096
memBase: 1006
opAddr:  4
memLoc:  0
pc:      9
val:     23