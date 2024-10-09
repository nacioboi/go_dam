#include "textflag.h"

TEXT ·avx512_find_idx_8i(SB), NOSPLIT|NOFRAME, $0-16

	PREFETCHT0 data+8(FP)

	VPBROADCASTQ x+0(FP), Z0    // Load key to search for.
	MOVQ data+8(FP), SI

simd:

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 (SI), Z1
	VPCMPEQQ Z1, Z0, K1        // K1 = (difference between Z1 and Z0). The most significant is how much offset.
	// Check if any match was found...
	KTESTQ K1, K1
	JNZ success_1
	// 8 elements processed.

	JMP fail

success_1:
	KMOVQ K1, R9
	BSFQ R9, R9
	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET

fail:

	MOVB $0, ok+17(FP)
	RET
	