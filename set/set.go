package set

type Set map[string]bool

func New() (set Set) {
	set = Set{}
	return
}
