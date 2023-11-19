                ; Version 2
                ; PDP-8 AND
                ; Uses check for high 13-bit rather than OR

                ; MOV memBase+opAddr TO memLoc
                memLoc memLoc
                mmemBase memLoc
                mopAddr memLoc

                tmp tmp
                [memLoc] z
                z tmp
                z z

                ; IF tmp >= masklc JUMP to thbit
                tmp maskl restoreMask
                ; Zero maskl
                maskl maskl

                ; Set high 13-bit
                mmaskl tmp

restoreMask:    mmaskl maskl

andit:          ; Store the AND mask
                andMask andMask
                tmp z
                z andMask
                ; z z  not needed because of z z and13 below

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

.data
z:      0
lm1:    -1
lmdone: 0-done
lmandit: 0-andit

cnt:      0
memLoc:   0
mmemBase: 0-1006    ; TODO: be able to put memBase or similar here
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
maskl:     4096
mmask13:   -8191    ; -0o17777
mmaskl:    -4096    ; -0o10000

andMask:        0
and13ret:       0