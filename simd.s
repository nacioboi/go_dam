// amd64 - go 1.23.0

#include "textflag.h"

// TEXT ·asm_std(SB), NOSPLIT, $0
// 	XORQ DX, DX                // Slice offset.
// 	XORQ R9, R9                // Index.
// 	XORQ R8, R8                // Did find.
// 	XORQ R11, R11 		// Mask

// 	MOVQ key+24(FP), AX        // Load `key` parameter into X1.
// 	MOVQ keys+0(FP), BX        // Load address of the slice into BX (keys.data).
// 	MOVQ keys_len+16(FP), CX   // Load the length of the slice into CX (keys.len).

// loop:
// 	CMPQ R9, CX
// 	JE done

// 	MOVQ 0(BX)(DX*1), R11
// 	ADDQ $8, DX

// 	CMPQ AX, R11
// 	JE found

// 	INCQ R9
// 	JMP loop

// found:
// 	MOVQ $1, R8       // Return 1 if found.

// done:
// 	// Set the index where the match was found.
// 	MOVB R9, ret+32(FP)       // Return the index of the found key or 0 if not found.
// 	MOVQ R8, ret+33(FP)       // Return 1 if found.
// 	RET
































































// TEXT ·simd_find_idx_128(SB), NOSPLIT, $0-16

// 	PREFETCHT0 params+0(FP)

// 	// Initialize variables
// 	XORQ DI, DI                 // DI = 0 (element index).

// 	// Load function parameters
// 	MOVQ $128, CX                // CX = len(keys).
// 	MOVQ key_ptr+0(FP), SI       // SI = &keys[0].
// 	MOVQ key+8(FP), X0           // Load key into X0.

// 	VPBROADCASTQ X0, Y0          // Broadcast key across YMM0.

// 	XORQ R15, R15
// 	XORQ R8, R8

// loop_simd:

// 	PREFETCHT0 0(SI)(DI*8)       // Prefetch the next 64 bytes.

// 	VMOVDQU64 (SI)(DI*8), Y1
// 	VPCMPEQQ  Y0, Y1, Y1

// 	// Check if any match was found in the first 4 elements.
// 	VPMOVMSKB Y1, AX
// 	CMPQ AX, R15
// 	CMOVQNE AX, R8
// 	CMOVQNE DI, R9
// 	ADDQ $4, DI

// 	CMPQ DI, CX
// 	JL loop_simd                 // Continue loop if not done.

// done:

// 	TESTQ R8, R8
// 	JZ fail

// success:

// 	BSFQ R8, R8                // Find first set bit in mask (position of match)
// 	SHRQ $3, R8                 // Divide by 8 to get the final index.
// 	ADDQ R9, R8                // Add base index to mask index.
// 	INCQ R8                     // Adjust for Go's 1-based index return.

// 	MOVB R8, ret+16(FP)         // Return index.
// 	RET

// fail:
// 	MOVB $0, ret+16(FP)          // Return 0 if not found.
// 	RET




// // Below is the better performer. The above can perform better when we shave off the overhead between calling each time.
// TEXT ·simd_find_idx_128(SB), NOSPLIT, $0-16

// 	PREFETCHT0 params+0(FP)     

// 	// Initialize variables
// 	MOVQ $0, DI                 // DI = 0 (element index).

// 	// Load function parameters
// 	MOVQ $128, CX     // CX = len(keys).
// 	MOVQ key_ptr+0(FP), SI      // SI = &keys[0].
// 	MOVQ key+8(FP), X0         // Load key into X0.

// 	VPBROADCASTQ X0, Y0         // Broadcast key across YMM0.

// 	PREFETCHT0 0(SI)(DI*8)       // Prefetch the next 64 bytes.

// loop_simd:

// 	// Process 4x 64-bit elements at a time...
// 	VMOVDQU64 (SI)(DI*8), Y1
// 	VPCMPEQQ  Y0, Y1, Y1        // Y1 = (Y1 == Y0)
// 	VPMOVMSKB Y1, AX            // Move mask to AX

// 	// Check if any match was found
// 	TESTQ AX, AX
// 	JNZ success             // Jump if match found

// 	PREFETCHT0 0(SI)(DI*8)       // Prefetch the next 64 bytes.

// 	ADDQ    $4, DI              // Move to next group of four elements
// 	CMPQ DI, CX
// 	JL loop_simd                // Continue loop if not done
// 	JMP fail

// success:

// 	BSFQ AX, R9                // Find first set bit in mask (position of match)
// 	SHRQ $3, R9                 // Divide by 8 to get the final index.
// 	ADDQ DI, R9                // Add base index to mask index.
// 	INCQ R9                     // Adjust for Go's 1-based index return.

// 	MOVB R9, ret+16(FP)         // Return index.
// 	RET

// fail:
// 	MOVB $0, ret+16(FP)         // Return 0 if not found.
// 	RET





TEXT ·avx512_find_idx_64(SB), NOSPLIT|NOFRAME, $0-520

	// Load function parameters
	VPBROADCASTQ 8(SP), Z0    // Load key to search for.

simd:

	PREFETCHT0 16(SP)
	PREFETCHT0 80(SP)
	PREFETCHT0 144(SP)

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 16(SP), Z1
	VPCMPEQQ Z1, Z0, K1        // K1 = (difference between Z1 and Z0). The most significant is how much offset.
	// Check if any match was found...
	KTESTB K1, K1
	JNZ success_1
	// 8 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 80(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_2
	// 16 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 144(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_3
	// 24 elements processed.

	// Prefetch some more but further down...
	// If we were to fetch above then the CPU might evict the cache line before we can use it.
	PREFETCHT0 208(SP)
	PREFETCHT0 272(SP)
	PREFETCHT0 336(SP)
	PREFETCHT0 400(SP)
	PREFETCHT0 464(SP)

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 208(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_4
	// 32 elements processed.s

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 272(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_5
	// 40 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 336(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_6
	// 48 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 400(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_7
	// 56 elements processed.

	// Process 8x 64-bit elements at a time...
	VMOVDQU64 464(SP), Z1
	VPCMPEQQ Z1, Z0, K1
	KTESTB K1, K1
	JNZ success_8
	// 64 elements processed.

	JMP fail

success_8:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $56, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_7:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $48, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_6:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $40, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_5:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $32, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_4:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $24, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_3:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $16, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_2:
	KMOVQ K1, R9
	BSFQ R9, R9
	ADDQ $8, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET
success_1:
	KMOVQ K1, R9
	BSFQ R9, R9
	MOVB R9, 528(SP)
	MOVB $1, 529(SP)
	RET

fail:

	MOVB $0, 529(SP)         // Return failure.
	RET
	






TEXT ·simd_find_idx(SB), NOSPLIT, $0-24

	PREFETCHT0 params+0(FP)     

	// Load function parameters
	MOVQ key_ptr+0(FP), SI      // SI = &keys[0].
	MOVQ key+16(FP), X0         // Load key into X0.
	MOVQ keys_len+8(FP), CX     // CX = len(keys).

	VPBROADCASTQ X0, Y0         // Broadcast key across YMM0.

	// Initialize variables
	MOVQ $0, DI                 // DI = 0 (element index).

loop_simd:

	PREFETCHT0 (SI)(DI*8)       // Prefetch the next 64 bytes.

	// Process 4x 64-bit elements at a time...
	VMOVDQU64 (SI)(DI*8), Y1
	VPCMPEQQ  Y0, Y1, Y1        // Y1 = (Y1 == Y0)
	VPMOVMSKB Y1, AX            // Move mask to AX

	// Check if any match was found
	TESTQ AX, AX
	JNZ success             // Jump if match found

	ADDQ    $4, DI              // Move to next group of four elements
	CMPQ DI, CX
	JL loop_simd                // Continue loop if not done
	JMP fail

success:
	MOVQ DI, R11                // Store the index of the match.
	MOVQ AX, R12                // Store the mask.
	BSFQ R12, R9                // Find first set bit in mask (position of match)
	SHRQ $3, R9                 // Divide by 8 to get the final index.
	ADDQ R11, R9                // Add base index to mask index.
	INCQ R9                     // Adjust for Go's 1-based index return.

	MOVB R9, ret+24(FP)         // Return index.
	RET

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