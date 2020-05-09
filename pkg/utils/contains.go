package utils

import "github.com/thoas/go-funk"

// 判断数组arr是否包含item元素
func Contains(arr interface{}, item interface{}) bool {
	switch arr.(type) {
	case []uint:
		// funk没有强类型是uint数组的方式, 自行实现
		if val, ok := item.(uint); ok {
			return ContainsUint(arr.([]uint), val)
		}
		break
	case []int:
		if val, ok := item.(int); ok {
			return funk.ContainsInt(arr.([]int), val)
		}
		break
	case []string:
		if val, ok := item.(string); ok {
			return funk.ContainsString(arr.([]string), val)
		}
		break
	case []int32:
		if val, ok := item.(int32); ok {
			return funk.ContainsInt32(arr.([]int32), val)
		}
		break
	case []int64:
		if val, ok := item.(int64); ok {
			return funk.ContainsInt64(arr.([]int64), val)
		}
		break
	case []float32:
		if val, ok := item.(float32); ok {
			return funk.ContainsFloat32(arr.([]float32), val)
		}
		break
	case []float64:
		if val, ok := item.(float64); ok {
			return funk.ContainsFloat64(arr.([]float64), val)
		}
		break
	}
	// funk默认使用反射, 性能不如强类型
	return funk.Contains(arr, item)
}

// 判断uint数组是否包含item元素
func ContainsUint(arr []uint, item uint) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
