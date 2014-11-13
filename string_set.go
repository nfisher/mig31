package main

type StringSet map[string]bool

func NewStringsSet() (set StringSet) {
	set = StringSet{}
	return
}

// Add adds the value to the set.
func (set StringSet) Add(value string) {
	set[value] = true
}

// Contains checks if the specified value is contained in the set.
func (set StringSet) Contains(value string) (contains bool) {
	contains = set[value]
	return
}

func (set StringSet) Len() (length int) {
	length = len(set)
	return
}

// Diff returns a new set with the values the left set has that are absent from the right set.
func (left StringSet) Diff(right StringSet) (diff StringSet) {
	diff = NewStringsSet()
	for l := range left {
		if !right.Contains(l) {
			diff.Add(l)
		}
	}
	return
}

func (left StringSet) Union(right StringSet) (union StringSet) {
	union = NewStringsSet()
	for l := range left {
		union.Add(l)
	}

	for r := range right {
		union.Add(r)
	}
	return
}

func (left StringSet) Intersection(right StringSet) (intersection StringSet) {
	intersection = NewStringsSet()
	for l := range left {
		if right[l] {
			intersection.Add(l)
		}
	}
	return
}

func (left StringSet) Equal(right StringSet) (isEqual bool) {
	isEqual = left.Len() == right.Len()
	for l := range left {
		isEqual = right.Contains(l)
		if !isEqual {
			return
		}
	}
	return
}
