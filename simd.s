// amd64 - go 1.23.0

#include "textflag.h"

TEXT ·asm_std(SB), NOSPLIT, $0
	XORQ DX, DX                // Slice offset.
	XORQ R9, R9                // Index.
	XORQ R8, R8                // Did find.
	XORQ R11, R11 		// Mask

	MOVQ key+24(FP), AX        // Load `key` parameter into X1.
	MOVQ keys+0(FP), BX        // Load address of the slice into BX (keys.data).
	MOVQ keys_len+16(FP), CX   // Load the length of the slice into CX (keys.len).

loop:
	CMPQ R9, CX
	JE done

	MOVQ 0(BX)(DX*1), R11
	ADDQ $8, DX

	CMPQ AX, R11
	JE found

	INCQ R9
	JMP loop

found:
	MOVQ $1, R8       // Return 1 if found.

done:
	// Set the index where the match was found.
	MOVB R9, ret+32(FP)       // Return the index of the found key or 0 if not found.
	MOVQ R8, ret+33(FP)       // Return 1 if found.
	RET











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
	
	
/*
// Register usage:
// R9: Result and temporary $0 register.
// R15: Temporary register for $1.
// SI: Pointer to the slice.
// CX: Length of the slice.
// X0: Key to find.
// Y0: Broadcasted key.
// R11: Temporary register for index when found.
// R12: Temporary register for mask when found.
// Y1: Temporary register for SIMD comparison.
// AX: Temporary register for SIMD comparison.
// DI: Loop index.
TEXT ·simd_find_idx(SB), NOSPLIT, $0-32

	// Load function parameters
	MOVQ key+24(FP), X0 		// Load key into X0.
	MOVQ keys+0(FP), SI 		// SI = &keys[0].
	MOVQ keys_len+16(FP), CX 	// CX = len(keys).

	VPBROADCASTQ X0, Y0 		// Broadcast key across YMM0.

	// Initialize variables
	XORQ R9, R9 			// R9 = 0 (return value).
	XORQ DI, DI              // DI = 0 (element index)

	XORQ R11, R11

loop_simd:

	PREFETCHT0 (SI)(DI*8)	// Prefetch the next 64 bytes.

	// Process 4x 64-bit elements at a time...
	VMOVDQU64 (SI)(DI*8), Y1
	VPCMPEQQ  Y0, Y1, Y1      // Y1 = (Y1 == Y0)
	VPMOVMSKB Y1, AX         // Move mask to AX

	ADDQ    $4, DI              // Move to next group of four elements

	CMPQ    AX, R9
	CMOVQNE DI, R11
	CMOVQNE AX, R12

	CMPQ DI, CX
	JE   done
	JMP  loop_simd

done:
	TESTQ R11, R11
	JZ fail

success:
	SUBQ $4, R11

	BSFQ R12, R9
	SHRQ $3, R9
	ADDQ R11, R9
	INCQ R9

	MOVB R9, ret+32(FP)      // Return index

ret:
	RET

fail:

	MOVB    $0, ret+32(FP)      // Return index
	JMP ret
*/