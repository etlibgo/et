package et

import (
	"reflect"
)

/**
反射的一些基本操作
 */

/**
根据名称执行对象的方法，参数是一个字符串，返回值也是一个字符串
 */
func ReflectCallForStrByStr(pObj interface{}, strMethodName string, strParam string) string{
	if pObj == ""{
		return ""
	}

	strResult := ""

	var valObj reflect.Value = reflect.ValueOf(pObj)

	if strMethodName == "" || valObj.IsNil() || !valObj.IsValid(){
		return ""
	}

	xMethod := valObj.MethodByName(strMethodName)
	if xMethod.IsNil() || !xMethod.IsValid(){
		return ""
	}

	arrInput := make([]reflect.Value,1)
	arrInput[0] = reflect.ValueOf(strParam)
	arrOutput := xMethod.Call(arrInput)
	if len(arrOutput) >= 1{
		strResult = arrOutput[0].String()
	}

	return strResult
}
