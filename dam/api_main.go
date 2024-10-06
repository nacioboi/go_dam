/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type I_Positive_Integer interface {
	uint8 | uint16 | uint32 | uint64
}

type t_bucket_entry[KT I_Positive_Integer, VT any] struct {
	key   KT
	value VT
}

type t_bucket_std[KT I_Positive_Integer, VT any] struct {
	entries []t_bucket_entry[KT, VT]
}
