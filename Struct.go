package et

import (
	"encoding/json"
	"errors"
	"reflect"
)

/**
结构体常用操作
 */

/**
把结构体转换为一个map
 */
func StructToMapVoid(xStruct interface{}, strTag string) map[string]interface{}{
	if xStruct == nil{
		return nil
	}

	mapResult := map[string]interface{}{}
	valueVoid := reflect.ValueOf(xStruct)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
	}

	for i := 0; i < typeVoid.NumField(); i++{
		strKey := typeVoid.Field(i).Name
		var valTemp reflect.Value = valueVoid.Field(i)
		if strTag != ""{
			strKey = typeVoid.Field(i).Tag.Get(strTag)
		}
		mapResult[strKey] = valTemp.Interface()
	}

	return mapResult
}

/**
结构体数组或者MAP深度复制
 */
func StructInGroupDeepClone(groupStruct interface{}) interface{} {
	valueSrc := reflect.ValueOf(groupStruct)
	typeSrc := valueSrc.Type()
	kindSrc := typeSrc.Kind()

	bPtr := false
	if kindSrc == reflect.Ptr{
		bPtr = true
		valueSrc = valueSrc.Elem()
		typeSrc = valueSrc.Type()
		kindSrc = typeSrc.Kind()
	}

	var xResult interface{}

	if kindSrc == reflect.Struct{
		xResult = StructDeepClone(groupStruct)
	}else if kindSrc == reflect.Map{
		mapNew := make(map[string]interface{})
		for k, v := range valueSrc.Interface().(map[string]interface{}) {
			mapNew[k] = StructInGroupDeepClone(v)
		}
		xResult = mapNew
		if bPtr{
			//xResult = &mapNew
			xResult = reflect.ValueOf(xResult).Addr().Interface()
		}
	}else if kindSrc == reflect.Array || kindSrc == reflect.Slice{
		arrNew := make([]interface{},len(valueSrc.Interface().([]interface{})))
		for k, v := range valueSrc.Interface().([]interface{}) {
			arrNew[k] = StructInGroupDeepClone(v)
		}
		xResult = arrNew
		if bPtr{
			//xResult = &arrNew
			xResult = reflect.ValueOf(xResult).Addr().Interface()
		}
	}else{
		xResult = groupStruct
	}

	return xResult
}

/**
单个结构体深度复制
 */
func StructDeepCopy(pSrc interface{}, pDest interface{}) error{
	if pSrc == nil || pDest == nil{
		return errors.New("StructDeepCopy nil param")
	}

	bufJson,err := json.Marshal(pSrc)
	if err!= nil{
		return err
	}
	err = json.Unmarshal(bufJson,pDest)

	return err
}

/**
结构体克隆，会产生新对象
 */
func StructDeepClone(pSrc interface{}) interface{}{
	if pSrc == nil {
		return nil
	}

	valueSrc := reflect.ValueOf(pSrc)
	typeSrc := valueSrc.Type()
	kindSrc := typeSrc.Kind()


	bufJson,err := json.Marshal(pSrc)
	if err != nil{
		return nil
	}

	if kindSrc == reflect.Ptr{
		pDest := reflect.New( valueSrc.Elem().Type())
		err := json.Unmarshal(bufJson,pDest.Interface())
		if err != nil{
			return nil
		}
		return pDest.Interface()
	}else{
		pDest := reflect.New( typeSrc)
		err := json.Unmarshal(bufJson,pDest.Interface())
		if err != nil{
			return nil
		}
		return pDest.Elem().Interface()
	}

}



