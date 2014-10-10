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

// Contains checks if the specified value is contained in the set.
func (set Set) Contains(value string) (contains bool) {
	contains = set[value]
	return
}

func (set Set) Len() (length int) {
	length = len(set)
	return
}

// Diff returns a new set with the values the left set has that are absent from the right set.
func (left Set) Diff(right Set) (diff Set) {
	diff = New()
	for l := range left {
		if !right.Contains(l) {
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

func (left Set) Equal(right Set) (isEqual bool) {
	isEqual = left.Len() == right.Len()
	for l := range left {
		isEqual = right.Contains(l)
		if !isEqual {
			return
		}
	}
	return
}
