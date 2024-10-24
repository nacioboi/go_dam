#include "textflag.h"



TEXT ·sse_find_idx_uint32_4aat(SB), NOSPLIT, $0

	// Clear registers that will be used
	XORQ CX, CX
	XORQ R9, R9

	// Move query to XMM0.
	XORQ BX, BX
	MOVL x+0(FP), BX
	VPBROADCASTD BX, Y0
	
	// Load the length and pointer to of the array...
	MOVL x+4(FP), DX
	MOVQ x+8(FP), SI

simd:

	PREFETCHT0 0(SI)(CX*4)

	// Process 8x 32-bit elements at a time using SSE...
	VMOVDQU32 0(SI)(CX*4), Y1        // Load 4x 32-bit elements into XMM1.
	VPCMPEQD Y1, Y0, Y1       // Compare elements in XMM1 and XMM0.
	VPMOVMSKB Y1, R9        // Move comparison results to R9.
	TESTQ R9, R9
	JNZ success
	// 4 elements processed.

	ADDQ $8, CX
	CMPL CX, DX
	JL simd

fail:

	MOVB $0, ok+17(FP)
	RET

success:

	BSFQ R9, R9              // Find the first set bit in R9.
	SHRQ $2, R9
	ADDQ CX, R9

	MOVB R9, v+16(FP)
	MOVB $1, ok+17(FP)
	RET
