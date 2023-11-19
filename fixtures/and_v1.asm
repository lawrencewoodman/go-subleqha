                ; Version 1
                ; PDP-8 AND

                ; MOV memBase+opAddr TO memLoc
                memLoc memLoc
                mmemBase memLoc
                mopAddr memLoc

                tmp tmp
                [memLoc] z
                z tmp
                z z

                ; Store the OR orA and orBLoc values
                orA orA
                orBLoc orBLoc
                mmaskl orA
                mtmp orBLoc

                ; Store the return location
                or13ret or13ret
                lmandit or13ret

                ; Jump to or13
                z z or13

                ; Store the AND mask
andit:          andMask andMask
                tmp z
                z andMask
                z z

                ; Store the return location
                and13ret and13ret
                lmdone and13ret

                ; Jump to and13
                z z and13

                ; HLT
done:           lm1 0


                ;--------------------------
                ; AND lac
                ; andMask contains value
                ; to mask with
                ;--------------------------

and13:          ; Zero res
                res res

                ; Loop 14 times
                cnt cnt
                l13 cnt
aloop:
                ; m := 0
                m m

                ; res << 2
                res z
                z res
                z z

                ; IF andMask >= lhbitval JUMP to mhbit
                andMask lhbitvalc mhbit
                ; Zero lhbitvalc and JUMP to acheck
                lhbitvalc lhbitvalc acheck

mhbit:          ; INC m
                lm1 m
                ; andMask -= lhbitval
                lhbitval andMask
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc

acheck:         lmhbitval lhbitvalc
                ; IF lac >= lhbitval JUMP to ahbit
                lac lhbitvalc ahbit
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc
                lmhbitval lhbitvalc
                z z acont

ahbit:          ; lac -= lhbitval
                lhbitval lac
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc
                lmhbitval lhbitvalc
                ; IF m <= 0 JUMP to acont
                z m acont

                ; High bits match
                ; res++
                lm1 res

acont:          ; andMask << 2
                andMask z
                z andMask
                z z

                ; lac << 2
                lac z
                z lac
                z z

                ; ADD l1 to cnt and JUMP to loop if <= 0
                lm1 cnt aloop

                ; Copy result to lac
                lac lac
                res z
                z lac

                ; Return and zero z
                z z [and13ret]

                ;------------
                ; End of AND
                ;------------

                ;------------------------------
                ; OR A BLoc
                ; orA contains value
                ; to OR the value at location
                ; orBLoc
                ;------------------------------


or13:           ; Zero res
                res res

                ; Loop 14 times
                cnt cnt
                l13 cnt
oloop:
                ; m := 0
                m m

                ; res << 2
                res z
                z res
                z z

                ; IF orA >= lhbitval JUMP to orahbit
                orA lhbitvalc orahbit
                ; Zero lhbitvalc and JUMP to bcheck
                lhbitvalc lhbitvalc bcheck

orahbit:        ; INC m
                lm1 m
                ; orA -= lhbitval
                lhbitval orA
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc

bcheck:         lmhbitval lhbitvalc
                ; IF [orBLoc] >= lhbitval JUMP to bhbit
                [orBLoc] lhbitvalc bhbit
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc
                lmhbitval lhbitvalc
                z z mcheck

bhbit:          ; INC m
                lm1 m
                ; [orBLoc] -= lhbitval
                lhbitval [orBLoc]
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc
                lmhbitval lhbitvalc


mcheck:         ; IF m <= 0 JUMP to cont
                z m ocont

                ; High bit present in one of the values
                ; res++
                lm1 res

ocont:          ; orA << 2
                orA z
                z orA
                z z

                ; [orBLoc] << 2
                [orBLoc] z
                z [orBLoc]
                z z

                ; ADD l1 to cnt and JUMP to loop if <= 0
                lm1 cnt oloop

                ; Copy result to [orBloc]
                [orBLoc] [orBLoc]
                res z
                z [orBLoc]

                ; Return and zero z
                z z [or13ret]

                ;------------
                ; End of OR
                ;------------

.data
z:      0
lm1:    -1
lmdone: 0-done
lmandit: 0-andit

cnt:      0
memLoc:   0
mtmp:     0-tmp
mmemBase: 0-1007    ; TODO: be able to put memBase or similar here
mopAddr:  0-4
lac:      4503
tmp:      0
val:      3003

res:    0
m:      0
lhbitval:   8192
lhbitvalc:  8192
l13:        13
lmhbitval: -8192
mmask13:   -8191    ; -0o17777
mmaskl:    -4096    ; -0o10000

andMask:        0
and13ret:       0

orA:            0
orBLoc:         0
or13ret:        0
