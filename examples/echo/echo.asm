        ; Echo a character from STDIN to STDOUT
        ;
        ; Copyright (C) 2023 Lawrence Woodman
        ; Licensed under a BSD 0-Clause licence. Please see 0BSD_LICENCE.md for details.

loop:   ch ch

        ; Load ch from STDIN
        1 ch
        z ch chle
        z z done
chle:   ch z loop
        z z outCh

outCh:  ch 1
cont:   z z loop

        ; HLT
done:   lm1 0

.data
z:      0
ch:     0
lm1:    -1
