/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type T_Option_Type uint8

const (
	OPTION_TYPE__WITH_HASH_FUNC T_Option_Type = iota
	OPTION_TYPE__WITH_PERFORMANCE_PROFILE
	OPTION_TYPE__WITH_EXPERIMENTAL_BATCHED_GETS
)

type T_Option[KT I_Positive_Integer, VT any] struct {
	t     T_Option_Type
	f     func(*DAM[KT, VT])
	other interface{}
}

func With_Hash_Func[KT I_Positive_Integer, VT any](f func(KT) uint64) T_Option[KT, VT] {
	return T_Option[KT, VT]{
		t: OPTION_TYPE__WITH_HASH_FUNC,
		f: func(m *DAM[KT, VT]) {
			m.users_chosen_hash_func = f
			m.using_users_hash_func = true
		},
	}
}

type T_Performance_Profile uint8

const (
	PERFORMANCE_PROFILE__FAST T_Performance_Profile = iota
	PERFORMANCE_PROFILE__NORMAL
	PERFORMANCE_PROFILE__SAVE_MEMORY
)

func With_Performance_Profile[KT I_Positive_Integer, VT any](p T_Performance_Profile) T_Option[KT, VT] {
	return T_Option[KT, VT]{
		t:     OPTION_TYPE__WITH_PERFORMANCE_PROFILE,
		other: p,
	}
}
