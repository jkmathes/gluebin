    .section .rodata

    .global thing1
    .type   thing1, @object
    .balign 4
thing1:
    .incbin "thing1.bin"
thing1_end:

    .global thing1_size
    .type   thing1_size, @object
    .balign 4
thing1_size:
    .int    thing1_end - thing1
