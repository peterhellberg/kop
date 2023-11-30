package list

import "cmp"

type Store[E cmp.Ordered] interface {
	Add(...E)
	Remove(...E)
	Contains(E) bool
	Members() []E
	Clear()
}
