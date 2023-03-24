package et

import "math"

/**
浮点数操作，主要是比较
 */

const FLOAT_COMPARE_DIFF = 0.00001

/**
两个浮点数是否相等，判断误差是否在可接受范围内
 */
func FloatEqual(numLeft float64 ,numRight float64) bool{
	if math.Abs(numLeft - numRight) > FLOAT_COMPARE_DIFF{
		return false
	}else{
		return true
	}
}