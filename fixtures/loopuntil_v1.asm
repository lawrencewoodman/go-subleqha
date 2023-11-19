        ; Version 1
        ; LOOP UNTIL

        ; MOV l4999 to cnt, cnt := 0--4999 = 0+4999 (Loop 5000 times)
        l4999 cnt
loop:
        ; ADD l1 to sum  (Take -1 from sum)
        lm1 sum

        ; ADD l1 to cnt and JUMP to loop if <= 0
        lm1 cnt loop

        ; HLT
done:   lm1 0

.data
z:      0
cnt:    0
sum:    0
l4999:  4999
lm1:    -1
