package util

type s interface {
	~string | ~int | ~int32 | ~int64
}

// DifferenceOnSlice 寻找两个切片的差集，结果为slice1不在slice2中的元素
func DifferenceOnSlice[T s](slice1, slice2 []T) []T {
	m := make(map[T]bool)
	var diff []T

	for _, item := range slice2 {
		m[item] = true
	}

	for _, item := range slice1 {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}
