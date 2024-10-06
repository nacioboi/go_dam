#include "textflag.h"

TEXT ·simd_find_idx(SB), NOSPLIT, $0-24

	// Load function parameters
	MOVQ key_ptr+0(FP), SI      // SI = &keys[0].
	MOVQ key+16(FP), X0         // Load key into X0.
	MOVQ keys_len+8(FP), CX     // CX = len(keys).

	VPBROADCASTQ X0, Y0         // Broadcast key across YMM0.

	// Initialize variables
	XORQ DI, DI                 // DI = 0 (element index).
	XORQ R12, R12               // Clear mask register.

loop_simd:

	PREFETCHT0 (SI)(DI*8)       // Prefetch the next 64 bytes.

	// Process 4x 64-bit elements at a time...
	VMOVDQU64 (SI)(DI*8), Y1
	VPCMPEQQ  Y0, Y1, Y1        // Y1 = (Y1 == Y0)
	VPMOVMSKB Y1, AX            // Move mask to AX

	// Check if any match was found
	TESTQ AX, AX
	JNZ match_found             // Jump if match found

	ADDQ    $4, DI              // Move to next group of four elements
	CMPQ DI, CX
	JL loop_simd                // Continue loop if not done

done:
	TESTQ R12, R12
	JZ fail

success:
	BSFQ R12, R9                // Find first set bit in mask (position of match)
	SHRQ $3, R9                 // Divide by 8 to get the final index.
	ADDQ R11, R9                // Add base index to mask index.
	INCQ R9                     // Adjust for Go's 1-based index return.

	MOVB R9, ret+24(FP)         // Return index.
	RET

match_found:
	MOVQ DI, R11                // Store the index of the match.
	MOVQ AX, R12                // Store the mask.
	JMP done

fail:
	MOVB $0, ret+24(FP)         // Return 0 if not found.
	RET
	