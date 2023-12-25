package utils

import "sort"

// IntsUnique 数组 int 去重
func IntsUnique(ss []int) []int {
	m := make(map[int]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	r := make([]int, 0, len(m))
	for s := range m {
		r = append(r, s)
	}
	sort.Sort(sort.IntSlice(r))
	return r
}

// StringsUnique 数组 string 去重
func StringsUnique(ss []string) []string {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	r := make([]string, 0, len(m))
	for s := range m {
		r = append(r, s)
	}
	sort.Sort(sort.StringSlice(r))
	return r
}

// Int64Unique 数组 int64 去重
func Int64Unique(ss []int64) []int64 {
	m := make(map[int64]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	r := make([]int64, 0, len(m))
	for s := range m {
		r = append(r, s)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}

// IsInString 判断是否在数组中
func IsInString(v string, strs []string) bool {
	for _, str := range strs {
		if v == str {
			return true
		}
	}
	return false
}
