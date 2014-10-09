package set

type Set map[string]bool

func New() (set Set) {
	set = Set{}
	return
}

// Add adds the value to the set.
func (set Set) Add(value string) {
	set[value] = true
}

// Diff returns the values that the left set has that the right set does not.
func (left Set) Diff(right Set) (diff Set) {
	diff = New()
	for l := range left {
		if !right[l] {
			diff.Add(l)
		}
	}
	return
}

func (left Set) Union(right Set) (union Set) {
	union = New()
	for l := range left {
		union.Add(l)
	}

	for r := range right {
		union.Add(r)
	}
	return
}

func (left Set) Intersection(right Set) (intersection Set) {
	intersection = New()
	for l := range left {
		if right[l] {
			intersection.Add(l)
		}
	}
	return
}
