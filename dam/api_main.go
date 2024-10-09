/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type I_Large_Positive_Integer interface {
	uint16 | uint32 | uint64
}

type I_Positive_Integer interface {
	uint8 | uint16 | uint32 | uint64
}

// Warning: This function must be called with keys satisfying the following conditions:
//
// - len(keys) % 4 == 0
//
// - no duplicate keys.
//
// - len(*p) >= 8
//
//go:noescape
//go:nosplit
func simd_find_idx_128i(key uint64, p *uint64) (uint8, bool)

// Warning: This function must be called with keys satisfying the following conditions:
//
// - len(*p) % 8 == 0
//
// - no duplicate keys.
//
// - len(*p) >= 8
//go:noescape
//go:nosplit
func avx512_find_idx_64i(query uint64, p *uint64) (uint8, bool)

// Warning: This function must be called with keys satisfying the following conditions:
//
// - len(*p) % 8 == 0
//
// - no duplicate keys.
//
// - len(*p) >= 8
//
//go:noescape
//go:nosplit
func avx512_find_idx_8i(query uint64, p *uint64) (uint8, bool)
