#include "textflag.h"


TEXT ·avx512_find_idx_64i(SB), NOSPLIT|NOFRAME, $0-16

	PREFETCHT0 data+8(FP)

	VPBROADCASTQ x+0(FP), Z0    	// Load key to search for.
	MOVQ data+8(FP), SI 		// Load the pointer to the data.

simd:

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 0(SI), Z1
	VPCMPEQQ Z1, Z0, K1        // K1 = (difference between Z1 and Z0). The most significant is how much offset.
	// Check if any match was found...
	KTESTQ K1, K1
	JNZ success_1
	// 8 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 64(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_2
	// 16 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 128(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_3
	// 24 elements processed.
	
	// Process 8x 64-bit elements at a time...
	VMOVDQU64 192(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_4
	// 32 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 256(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_5
	// 40 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 320(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_6
	// 48 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 384(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_7
	// 56 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 448(SI), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTQ K1, K1
	JNZ success_8
	// 64 elements processed.

	JMP fail

success_1:
	KMOVQ K1, R9
	BSFQ R9, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_2:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $8, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_3:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $16, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_4:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $24, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_5:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $32, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_6:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $40, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_7:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $48, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

success_8:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $56, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

fail:

	MOVB $0, ok+17(FP)
	RET
	