package et

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/**
对各种类型的转换
void表示interface{}
long表示int64
 */

/**
// 整型格式标记：
// 'b' (-ddddp±ddd，二进制指数)
// 'e' (-d.dddde±dd，十进制指数)
// 'E' (-d.ddddE±dd，十进制指数)
// 'f' (-ddd.dddd，没有指数)
// 'g' ('e':大指数，'f':其它情况)
// 'G' ('E':大指数，'f':其它情况)
*/


func CastI64ToStr(iNum int64) string{
	return strconv.FormatInt(iNum,10)
}

func CastIntToStr(iNum int) string{
	return strconv.Itoa(iNum)
}

func CastIntToBool(iNum int) bool{
	if iNum == 0{
		return false
	}else{
		return true
	}
}

func CastStrToInt(strNum string) int{
	iNum,err := strconv.Atoi(strNum)
	if err != nil{
		iNum = 0
	}

	return iNum
}

func CastStrToI64(strNum string) int64{
	iNum,err := strconv.ParseInt(strNum,10,64)
	if err != nil{
		iNum = 0
	}

	return iNum
}

func CastStrToF64(strNum string) float64{
	fNum,err := strconv.ParseFloat(strNum,64)
	if err != nil{
		fNum = 0
	}

	return fNum
}

func CastF64ToStr(fNum float64, iPrecision int) string{
	return strconv.FormatFloat(fNum,'f',iPrecision,64)
}

/**
把字符串转换为bool，
会识别，true,false,0,1,
其他的情况下，非空即为真
 */
func CastStrToBool(strBool string) bool{
	if strBool == ""{
		return false
	}

	if strings.ToUpper(strBool) == "FALSE"{
		return false
	}

	if strings.ToUpper(strBool) == "TRUE"{
		return true
	}

	if strBool == "0"{
		return false
	}

	if strBool == "1"{
		return true
	}

	iNum := CastStrToInt(strBool)
	if iNum == 0{
		return false
	}else{
		return true
	}


}

func CastBoolToStr(bValue bool) string{
	if bValue{
		return "true"
	}else{
		return "false"
	}
}

func CastBufToStr(bufData []byte) string{
	if bufData == nil{
		return ""
	}else{
		return string(bufData)
	}
}

func CastBoolToNumStr(bValue bool) string{
	if bValue{
		return "1"
	}else{
		return "0"
	}
}

func CastErrToStr(objErr error) string{
	if objErr == nil{
		return ""
	} else{
		return objErr.Error()
	}
}

/**
把未知类型转化为相应的字符串
如果是结构体，会调用其String方法。
若无，则返回空字符串
 */
func CastVoidToStr(xValue interface{}) string{
	if xValue == nil{
		return ""
	}

	strResult := ""
	valueVoid := reflect.ValueOf(xValue)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
		kindVoid = typeVoid.Kind()
	}

	if kindVoid == reflect.Invalid {
		strResult = ""
	}else if kindVoid == reflect.Bool{
		strResult = CastBoolToStr(valueVoid.Bool())
	}else if kindVoid == reflect.Float64 || kindVoid == reflect.Float32{
		strResult = CastF64ToStr(valueVoid.Float(),-1)
	}else if kindVoid == reflect.Int || kindVoid == reflect.Int32 || kindVoid == reflect.Int16 || kindVoid == reflect.Int8 || kindVoid == reflect.Uint || kindVoid == reflect.Uint32 || kindVoid == reflect.Uint16 || kindVoid == reflect.Uint8{
		strResult = CastI64ToStr(valueVoid.Int())
	}else if kindVoid == reflect.Int64 || kindVoid == reflect.Uint64{
		strResult = CastI64ToStr(valueVoid.Int())
	}else if kindVoid == reflect.String{
		strResult = valueVoid.String()
	}else if kindVoid == reflect.Interface{
		strResult = fmt.Sprintf("%s",valueVoid)
	}else if kindVoid == reflect.Map || kindVoid == reflect.Array || kindVoid == reflect.Slice || kindVoid == reflect.Struct{
		strResult = fmt.Sprintf("%s",valueVoid)
	}else if kindVoid == reflect.Complex64 || kindVoid == reflect.Complex128{
		strResult = CastF64ToStr(real(valueVoid.Complex()),-1)  + "," + CastF64ToStr(imag(valueVoid.Complex()),-1)
	}else if kindVoid == reflect.Func{
		strResult = "(func)"
	}else if kindVoid == reflect.Chan{
		strResult = "(chan)"
	}else if kindVoid == reflect.Ptr || kindVoid == reflect.Uintptr || kindVoid == reflect.UnsafePointer{
		strResult = "(ptr)"
	}else{
		strResult = ""
	}

	return strResult
}

