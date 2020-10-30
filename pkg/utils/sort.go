package utils

// 用于uint数组类型排序
type UintSort []uint

func (s UintSort) Len() int           { return len(s) }
func (s UintSort) Less(i, j int) bool { return s[i] < s[j] }
func (s UintSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
