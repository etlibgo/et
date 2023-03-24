package et

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/**
关于标准json的一些操作，
需要获得json的详细操作，请参考ExJson文件
 */

const JSON_TYPE_NONE int = 0
const JSON_TYPE_TEXT int = 1
const JSON_TYPE_NUM int = 2
const JSON_TYPE_BOOL int = 3
const JSON_TYPE_NULL int = 4
const JSON_TYPE_ARRAY int = 5
const JSON_TYPE_OBJECT int = 6

const JSON_INDENT_SIZE int = 2

const JSON_STR_TRUE string = "true"
const JSON_STR_FALSE string = "false"
const JSON_STR_NULL string = "null"
const JSON_STR_LINE string = "\r\n"

const JSON_STR_D_QUOTE string = "\""
const JSON_STR_COLON string = ":"
const JSON_STR_COMMA string = ","
const JSON_STR_SPACE string = " "
const JSON_STR_B_SLASH string = "\\"
const JSON_STR_L_B_L string = "{"
const JSON_STR_L_B_R string = "}"
const JSON_STR_M_B_L string = "["
const JSON_STR_M_B_R string = "]"

const JSON_CH_D_QUOTE byte = '"'
const JSON_CH_COLON byte = ':'
const JSON_CH_COMMA byte = ','
const JSON_CH_SPACE byte = ' '
const JSON_CH_B_SLASH byte = '\\'
const JSON_CH_L_B_L byte = '{'
const JSON_CH_L_B_R byte = '}'
const JSON_CH_M_B_L byte = '['
const JSON_CH_M_B_R byte = ']'

const JSON_ET_TYPE_NAME string = "et.JsonVal"

var jsonArrBlankChars []byte = []byte{' ', '\t', '\r' , '\n'}

type JsonVal struct{
	iValueType int
	xVoidValue interface{}
}

func (my JsonVal)EtTypeName() string{
	return JSON_ET_TYPE_NAME
}

func (my JsonVal)JsonValTag() string{
	return JSON_ET_TYPE_NAME
}

func (my JsonVal)String() string{
	return JsonToStr(&my)
}

func JsonNewArray() *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_ARRAY,xVoidValue: make([]interface{},0)}
}

func JsonNewObject() *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_OBJECT,xVoidValue: make(map[string]interface{})}
}

func jsonNewNull() *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_NULL,xVoidValue: nil}
}

func jsonNewBool(bValue bool) *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_BOOL,xVoidValue: bValue}
}

func jsonNewNum(xValue interface{}) *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_NUM,xVoidValue: xValue}
}

func jsonNewText(strText string) *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_TEXT,xVoidValue: strText}
}

func JsonNewNone() *JsonVal{
	return &JsonVal{iValueType: JSON_TYPE_NONE,xVoidValue: nil}
}

func jsonInitValueAsText(pJson *JsonVal, strText string){
	if pJson == nil{
		return
	}
	pJson.iValueType = JSON_TYPE_TEXT
	pJson.xVoidValue = strText
}

func jsonInitValueAsNull(pJson *JsonVal){
	if pJson == nil{
		return
	}
	pJson.iValueType = JSON_TYPE_NULL
	pJson.xVoidValue = nil
}

func jsonInitValueAsBool(pJson *JsonVal, bValue bool){
	if pJson == nil{
		return
	}
	pJson.iValueType = JSON_TYPE_BOOL
	pJson.xVoidValue = bValue
}

func jsonInitValueAsNum(pJson *JsonVal, xNum interface{}){
	if pJson == nil{
		return
	}
	pJson.iValueType = JSON_TYPE_NUM
	pJson.xVoidValue = xNum
}

func jsonInitType(pJson *JsonVal, iType int){
	if pJson == nil{
		return
	}

	pJson.iValueType = iType
}
func jsonCheckCollectionForObject(pJson *JsonVal){
	if pJson == nil{
		return
	}

	if pJson.iValueType == JSON_TYPE_OBJECT{
		return
	}

	pJson.iValueType = JSON_TYPE_OBJECT
	pJson.xVoidValue = make(map[string]interface{})
}

func jsonCheckCollectionForArray(pJson *JsonVal){
	if pJson == nil{
		return
	}

	if pJson.iValueType == JSON_TYPE_ARRAY{
		return
	}

	pJson.iValueType = JSON_TYPE_ARRAY
	pJson.xVoidValue = make([]interface{},0)
}

func JsonIsEmpty(pJson *JsonVal) bool{
	return JsonLen(pJson) == 0
}

func JsonIsArray(pJson *JsonVal) bool{
	if pJson == nil{
		return false
	}

	return pJson.iValueType == JSON_TYPE_ARRAY
}

func JsonIsObject(pJson *JsonVal) bool{
	if pJson == nil{
		return false
	}

	return pJson.iValueType == JSON_TYPE_OBJECT
}

func JsonIsCollection(pJson *JsonVal) bool{
	if pJson == nil{
		return false
	}

	if pJson.iValueType == JSON_TYPE_ARRAY || pJson.iValueType == JSON_TYPE_OBJECT{
		return true
	}

	return false
}

func JsonLen(pJson *JsonVal) int{
	if pJson == nil{
		return 0
	}

	if pJson.iValueType == JSON_TYPE_ARRAY{
		return len(pJson.xVoidValue.([]interface{}))
	}else if pJson.iValueType == JSON_TYPE_OBJECT{
		return len(pJson.xVoidValue.(map[string]interface{}))
	}

	return 0
}

func JsonArrSetVoid(pJson *JsonVal, iIndex int, xValue interface{}){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	iLen := JsonLen(pJson)
	if iLen < iIndex{
		return
	}

	if xValue == nil{
		JsonArrSetNull(pJson,iIndex)
		return
	}

	isPtr := false
	valueVoid := reflect.ValueOf(xValue)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		isPtr = true
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
		kindVoid = typeVoid.Kind()
	}

	if strings.Index(typeVoid.Name(),"JsonVal" ) >= 0{
		if isPtr{
			JsonArrSetJson(pJson,iIndex,xValue.(*JsonVal))
		}else{
			strJson := fmt.Sprintf("%s",xValue)
			JsonArrSetJson(pJson,iIndex, JsonFromStr(strJson) )
		}
		return
	}

	if kindVoid == reflect.Invalid {
		JsonArrSetNull(pJson,iIndex)
	}else if kindVoid == reflect.Bool{
		JsonArrSetBool(pJson,iIndex,valueVoid.Bool())
	}else if kindVoid == reflect.Float64 || kindVoid == reflect.Float32{
		JsonArrSetFloat(pJson,iIndex,valueVoid.Float())
	}else if kindVoid == reflect.Int || kindVoid == reflect.Int32 || kindVoid == reflect.Int16 || kindVoid == reflect.Int8 || kindVoid == reflect.Uint || kindVoid == reflect.Uint32 || kindVoid == reflect.Uint16 || kindVoid == reflect.Uint8{
		JsonArrSetI64(pJson,iIndex,valueVoid.Int())
	}else if kindVoid == reflect.Int64 || kindVoid == reflect.Uint64{
		JsonArrSetI64(pJson,iIndex,valueVoid.Int())
	}else if kindVoid == reflect.String{
		JsonArrSetText(pJson,iIndex,valueVoid.String())
	}else if kindVoid == reflect.Interface{
		JsonArrSetText(pJson,iIndex,fmt.Sprintf("%s",valueVoid))
	}else if kindVoid == reflect.Map || kindVoid == reflect.Array || kindVoid == reflect.Slice || kindVoid == reflect.Struct{
		strJson := JsonMarshal(xValue)
		jsonCopy := JsonFromStr(strJson)
		JsonArrSetJson(pJson,iIndex,jsonCopy)
	}else if kindVoid == reflect.Complex64 || kindVoid == reflect.Complex128{
		JsonArrSetText(pJson,iIndex,CastF64ToStr(real(valueVoid.Complex()),-1)  + "," + CastF64ToStr(imag(valueVoid.Complex()),-1))
	}else if kindVoid == reflect.Func{
		JsonArrSetText(pJson,iIndex,"(func)")
	}else if kindVoid == reflect.Chan{
		JsonArrSetText(pJson,iIndex,"(chan)")
	}else if kindVoid == reflect.Ptr || kindVoid == reflect.Uintptr || kindVoid == reflect.UnsafePointer{
		JsonArrSetText(pJson,iIndex,"(ptr)")
	}else{
		JsonArrSetText(pJson,iIndex,"(unknown)")
	}
}

func JsonArrClear(pJson *JsonVal){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return
	}

	pJson.xVoidValue = make([]interface{},0)
}

func JsonArrRemoveIndex(pJson *JsonVal, iIndex int){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	iLen := JsonLen(pJson)
	if iLen <= iIndex{
		return
	}

	arrLeft := pJson.xVoidValue.([]interface{})[:iIndex]
	arrRight := pJson.xVoidValue.([]interface{})[iIndex + 1:]
	pJson.xVoidValue = append(arrLeft,arrRight...)
}

func JsonArrSetNum(pJson *JsonVal, iIndex int, xValue interface{}){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	JsonArrFillBlank(pJson,iIndex)

	iLen := len(pJson.xVoidValue.([]interface{}))

	if iLen < iIndex{
		return
	}else{
		pNewJson := jsonNewNum(xValue)
		if iLen == iIndex{
			pJson.xVoidValue = append(pJson.xVoidValue.([]interface{}),pNewJson)
		}else{
			pJson.xVoidValue.([]interface{})[iIndex] = pNewJson
		}
	}
}

func JsonArrSetInt(pJson *JsonVal, iIndex int, iValue int){
	JsonArrSetNum(pJson,iIndex,iValue)
}

func JsonArrSetI64(pJson *JsonVal, iIndex int, iValue int64){
	JsonArrSetNum(pJson,iIndex,iValue)
}

func JsonArrSetFloat(pJson *JsonVal, iIndex int, fValue float64){
	JsonArrSetNum(pJson,iIndex,fValue)
}

func JsonArrSetText(pJson *JsonVal, iIndex int, strText string){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	JsonArrFillBlank(pJson,iIndex)

	iLen := len(pJson.xVoidValue.([]interface{}))
	if iLen < iIndex{
		return
	}else{
		pNewJson := jsonNewText(strText)
		if iLen == iIndex{
			pJson.xVoidValue = append(pJson.xVoidValue.([]interface{}),pNewJson)
		}else{
			pJson.xVoidValue.([]interface{})[iIndex] = pNewJson
		}
	}
}

func JsonArrSetBool(pJson *JsonVal, iIndex int, bValue bool){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	JsonArrFillBlank(pJson,iIndex)

	iLen := len(pJson.xVoidValue.([]interface{}))
	if iLen < iIndex{
		return
	}else{
		pNewJson := jsonNewBool(bValue)
		if iLen == iIndex{
			pJson.xVoidValue = append(pJson.xVoidValue.([]interface{}),pNewJson)
		}else{
			pJson.xVoidValue.([]interface{})[iIndex] = pNewJson
		}
	}
}

func JsonArrSetNull(pJson *JsonVal, iIndex int){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	JsonArrFillBlank(pJson,iIndex)

	iLen := len(pJson.xVoidValue.([]interface{}))
	if iLen < iIndex{
		return
	}else{
		pNewJson := jsonNewNull()
		if iLen == iIndex{
			pJson.xVoidValue = append(pJson.xVoidValue.([]interface{}),pNewJson)
		}else{
			pJson.xVoidValue.([]interface{})[iIndex] = pNewJson
		}
	}
}

func JsonArrSetJson(pJson *JsonVal, iIndex int, pChildJson *JsonVal){
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return
	}

	JsonArrFillBlank(pJson,iIndex)

	iLen := len(pJson.xVoidValue.([]interface{}))
	if iLen < iIndex{
		return
	}else{
		pNewJson := pChildJson
		if pNewJson == nil{
			pNewJson = jsonNewNull()
		}
		if iLen == iIndex{
			pJson.xVoidValue = append(pJson.xVoidValue.([]interface{}),pNewJson)
		}else{
			pJson.xVoidValue.([]interface{})[iIndex] = pNewJson
		}
	}
}

func JsonArrAppendInt(pJson *JsonVal, iValue int){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetInt(pJson,iLen,iValue)
}

func JsonArrAppendI64(pJson *JsonVal, iValue int64){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetI64(pJson,iLen,iValue)
}

func JsonArrAppendFloat(pJson *JsonVal, fValue float64){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetFloat(pJson,iLen,fValue)
}

func JsonArrAppendNum(pJson *JsonVal, numValue interface{}){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetNum(pJson,iLen,numValue)
}

func JsonArrAppendText(pJson *JsonVal, strText string){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetText(pJson,iLen,strText)
}

func JsonArrAppendNull(pJson *JsonVal){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetNull(pJson,iLen)
}

func JsonArrAppendJson(pJson *JsonVal, pChildJson *JsonVal){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	JsonArrSetJson(pJson,iLen,pChildJson)
}

func JsonArrGetJson(pJson *JsonVal, iIndex int) *JsonVal{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return nil
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	return pChildJson
}

func JsonArrGetStr(pJson *JsonVal, iIndex int) string{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return ""
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	if JsonIsCollection(pChildJson){
		return JsonToStr(pChildJson)
	}else{
		return jsonInnerValueToStr(pChildJson)
	}
}

func JsonArrGetVoid(pJson *JsonVal, iIndex int) interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return ""
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	if JsonIsCollection(pChildJson){
		return JsonToStr(pChildJson)
	}else{
		return jsonInnerValueToVoid(pChildJson)
	}
}

func JsonArrGetInt(pJson *JsonVal, iIndex int) int{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return 0
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	return jsonInnerValueToInt(pChildJson)
}

func JsonArrGetI64(pJson *JsonVal, iIndex int) int64{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return 0
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	return jsonInnerValueToI64(pChildJson)
}

func JsonArrGetFloat(pJson *JsonVal, iIndex int) float64{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return 0
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	return jsonInnerValueToFloat(pChildJson)
}

func JsonArrGetBool(pJson *JsonVal, iIndex int) bool{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY || iIndex < 0{
		return false
	}

	pChildJson := pJson.xVoidValue.([]interface{})[iIndex].(*JsonVal)
	return jsonInnerValueToBool(pChildJson)
}

func JsonArrFillBlank(pJson *JsonVal, iMaxIndex int){
	if pJson == nil{
		return
	}
	iLen := JsonLen(pJson)
	if iLen >= iMaxIndex{
		return
	}

	for i := iLen; i < iMaxIndex; i++{
		pNewJson := jsonNewNull()
		pJson.xVoidValue = append(pJson.xVoidValue.([]interface{}),pNewJson)
	}
}

func JsonObjKeys(pJson *JsonVal)[]string{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return nil
	}

	arrKeys := make([]string,0)
	for k,_ := range pJson.xVoidValue.(map[string]interface{}){
		arrKeys = append(arrKeys,k)
	}

	return arrKeys
}

func JsonObjSetVoid(pJson *JsonVal, strKey string, xValue interface{}){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	if xValue == nil{
		JsonObjSetNull(pJson,strKey)
		return
	}

	isPtr := false
	valueVoid := reflect.ValueOf(xValue)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		isPtr = true
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
		kindVoid = typeVoid.Kind()
	}
	
	if strings.Index(typeVoid.Name(),"JsonVal" ) >= 0{
		if isPtr{
			JsonObjSetJson(pJson,strKey,xValue.(*JsonVal))
		}else{
			strJson := fmt.Sprintf("%s",xValue)
			JsonObjSetJson(pJson,strKey, JsonFromStr(strJson) )
		}
		return
	}

	if kindVoid == reflect.Invalid {
		JsonObjSetNull(pJson,strKey)
	}else if kindVoid == reflect.Bool{
		JsonObjSetBool(pJson,strKey,valueVoid.Bool())
	}else if kindVoid == reflect.Float64 || kindVoid == reflect.Float32{
		JsonObjSetFloat(pJson,strKey,valueVoid.Float())
	}else if kindVoid == reflect.Int || kindVoid == reflect.Int32 || kindVoid == reflect.Int16 || kindVoid == reflect.Int8 || kindVoid == reflect.Uint || kindVoid == reflect.Uint32 || kindVoid == reflect.Uint16 || kindVoid == reflect.Uint8{
		JsonObjSetI64(pJson,strKey,valueVoid.Int())
	}else if kindVoid == reflect.Int64 || kindVoid == reflect.Uint64{
		JsonObjSetI64(pJson,strKey,valueVoid.Int())
	}else if kindVoid == reflect.String{
		JsonObjSetText(pJson,strKey,valueVoid.String())
	}else if kindVoid == reflect.Interface{
		JsonObjSetText(pJson,strKey,fmt.Sprintf("%s",valueVoid))
	}else if kindVoid == reflect.Map || kindVoid == reflect.Array || kindVoid == reflect.Slice || kindVoid == reflect.Struct{
		strJson := JsonMarshal(xValue)
		jsonCopy := JsonFromStr(strJson)
		JsonObjSetJson(pJson,strKey,jsonCopy)
	}else if kindVoid == reflect.Complex64 || kindVoid == reflect.Complex128{
		JsonObjSetText(pJson,strKey,CastF64ToStr(real(valueVoid.Complex()),-1)  + "," + CastF64ToStr(imag(valueVoid.Complex()),-1))
	}else if kindVoid == reflect.Func{
		JsonObjSetText(pJson,strKey,"(func)")
	}else if kindVoid == reflect.Chan{
		JsonObjSetText(pJson,strKey,"(chan)")
	}else if kindVoid == reflect.Ptr || kindVoid == reflect.Uintptr || kindVoid == reflect.UnsafePointer{
		JsonObjSetText(pJson,strKey,"(ptr)")
	}else{
		JsonObjSetText(pJson,strKey,"(unknown)")
	}
}

func JsonObjContainsKey(pJson *JsonVal, strKey string) bool{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return false
	}

	if _,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		return true
	}else{
		return false
	}
}

func JsonObjClear(pJson *JsonVal){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	pJson.xVoidValue = make(map[string]interface{})
}

func JsonObjRemoveKey(pJson *JsonVal, strKey string){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	delete(pJson.xVoidValue.(map[string]interface{}),strKey)
}

func JsonObjSetNum(pJson *JsonVal, strKey string, xValue interface{}){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	pNewJson := jsonNewNum(xValue)
	pJson.xVoidValue.(map[string]interface{})[strKey] = pNewJson
}

func JsonObjSetInt(pJson *JsonVal, strKey string, iValue int){
	JsonObjSetNum(pJson,strKey,iValue)
}

func JsonObjSetI64(pJson *JsonVal, strKey string, iValue int64){
	JsonObjSetNum(pJson,strKey,iValue)
}

func JsonObjSetFloat(pJson *JsonVal, strKey string, fValue float64){
	JsonObjSetNum(pJson,strKey,fValue)
}

func JsonObjSetText(pJson *JsonVal, strKey string, strText string){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	pNewJson := jsonNewText(strText)
	pJson.xVoidValue.(map[string]interface{})[strKey] = pNewJson
}

func JsonObjSetBool(pJson *JsonVal, strKey string, bBool bool){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	pNewJson := jsonNewBool(bBool)
	pJson.xVoidValue.(map[string]interface{})[strKey] = pNewJson
}

func JsonObjSetNull(pJson *JsonVal, strKey string){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	pNewJson := jsonNewNull()
	pJson.xVoidValue.(map[string]interface{})[strKey] = pNewJson
}

func JsonObjSetJson(pJson *JsonVal, strKey string, pChildJson *JsonVal){
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return
	}

	pNewJson := pChildJson
	if pNewJson == nil{
		pNewJson = jsonNewNull()
	}
	pJson.xVoidValue.(map[string]interface{})[strKey] = pNewJson
}

func JsonObjGetJson(pJson *JsonVal, strKey string) *JsonVal{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return nil
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		return xValue.(*JsonVal)
	}else{
		return nil
	}

	return nil
}

func JsonObjGetInt(pJson *JsonVal, strKey string) int{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return 0
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		return jsonInnerValueToInt(xValue.(*JsonVal))
	}else{
		return 0
	}

	return 0
}

func JsonObjGetI64(pJson *JsonVal, strKey string) int64{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return 0
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		return jsonInnerValueToI64(xValue.(*JsonVal))
	}else{
		return 0
	}

	return 0
}

func JsonObjGetFloat(pJson *JsonVal, strKey string) float64{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return 0
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		return jsonInnerValueToFloat(xValue.(*JsonVal))
	}else{
		return 0
	}

	return 0
}

func JsonObjGetBool(pJson *JsonVal, strKey string) bool{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return false
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		return jsonInnerValueToBool(xValue.(*JsonVal))
	}else{
		return false
	}

	return false
}

func JsonObjGetStr(pJson *JsonVal, strKey string) string{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return ""
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		subJson := xValue.(*JsonVal)
		if JsonIsCollection(subJson){
			return JsonToStr(subJson)
		}else{
			return jsonInnerValueToStr(subJson)
		}
	}else{
		return ""
	}

	return ""
}

func JsonObjGetVoid(pJson *JsonVal, strKey string) interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return ""
	}

	if xValue,existed := pJson.xVoidValue.(map[string]interface{})[strKey];existed{
		subJson := xValue.(*JsonVal)
		if JsonIsCollection(subJson){
			return JsonToStr(subJson)
		}else{
			return jsonInnerValueToVoid(subJson)
		}
	}else{
		return ""
	}

	return ""
}


func JsonToStr(pJson *JsonVal) string{
	if pJson == nil{
		return ""
	}

	return jsonToStrRecur(pJson)
}

func JsonToFormattedStr(pJson *JsonVal) string{
	if pJson == nil{
		return ""
	}

	return JsonToFormattedStrByLevel(pJson,0)
}

func JsonToFormattedStrByLevel(pJson *JsonVal, iLevel int) string{
	if pJson == nil{
		return ""
	}

	if pJson.iValueType == JSON_TYPE_OBJECT{
		strResult := JSON_STR_L_B_L
		isFirst := true

		for k,v := range pJson.xVoidValue.(map[string]interface{}){
			if isFirst{
				isFirst = false
			}else{
				strResult += ","
			}

			strResult += JSON_STR_LINE
			strResult += strings.Repeat(JSON_STR_SPACE, JSON_INDENT_SIZE * (iLevel + 1))

			strResult += JSON_STR_D_QUOTE + StrAddSlashes( k, JSON_CH_D_QUOTE) + JSON_STR_D_QUOTE
			strResult += JSON_STR_COLON
			strResult += JsonToFormattedStrByLevel(v.(*JsonVal),iLevel + 1)
		}

		strResult += JSON_STR_LINE
		strResult += strings.Repeat(JSON_STR_SPACE, JSON_INDENT_SIZE * (iLevel + 0))

		strResult += JSON_STR_L_B_R
		return strResult
	}else if pJson.iValueType == JSON_TYPE_ARRAY{
		strResult := JSON_STR_M_B_L
		isFirst := true

		for _,v := range pJson.xVoidValue.([]interface{}){
			if isFirst{
				isFirst = false
			}else{
				strResult += ","
			}

			strResult += JSON_STR_LINE
			strResult += strings.Repeat(JSON_STR_SPACE, JSON_INDENT_SIZE * (iLevel + 1))

			strResult += JsonToFormattedStrByLevel(v.(*JsonVal),iLevel + 1)
		}

		strResult += JSON_STR_LINE
		strResult += strings.Repeat(JSON_STR_SPACE, JSON_INDENT_SIZE * (iLevel + 0))

		strResult += JSON_STR_M_B_R
		return strResult
	}else {
		return jsonInnerValueOutputStr(pJson)
	}

	return ""
}

func jsonToStrRecur(pJson *JsonVal) string{
	if pJson == nil{
		return ""
	}

	if pJson.iValueType == JSON_TYPE_OBJECT{
		strResult := JSON_STR_L_B_L
		isFirst := true

		for k,v := range pJson.xVoidValue.(map[string]interface{}){
			if isFirst{
				isFirst = false
			}else{
				strResult += ","
			}

			strResult += JSON_STR_D_QUOTE + StrAddSlashes( k, JSON_CH_D_QUOTE) + JSON_STR_D_QUOTE
			strResult += JSON_STR_COLON
			strResult += jsonToStrRecur(v.(*JsonVal))
		}

		strResult += JSON_STR_L_B_R
		return strResult
	}else if pJson.iValueType == JSON_TYPE_ARRAY{
		strResult := JSON_STR_M_B_L
		isFirst := true

		for _,v := range pJson.xVoidValue.([]interface{}){
			if isFirst{
				isFirst = false
			}else{
				strResult += ","
			}

			strResult += jsonToStrRecur(v.(*JsonVal))
		}

		strResult += JSON_STR_M_B_R
		return strResult
	}else {
		return jsonInnerValueOutputStr(pJson)
	}

	return ""
}

func jsonInnerValueToInt(pJson *JsonVal) int{
	if pJson == nil{
		return 0
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return 0
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return CastStrToInt( pJson.xVoidValue.(string) )
	}else if pJson.iValueType == JSON_TYPE_NUM{
		strNum := CastVoidToStr(pJson.xVoidValue)
		if strNum == ""{
			strNum = "0"
		}
		return CastStrToInt(strNum)
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		if pJson.xVoidValue.(bool) {
			return 1
		}else{
			return 0
		}
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return 0
	}

	return 0
}

func jsonInnerValueToI64(pJson *JsonVal) int64{
	if pJson == nil{
		return 0
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return 0
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return CastStrToI64( pJson.xVoidValue.(string) )
	}else if pJson.iValueType == JSON_TYPE_NUM{
		strNum := CastVoidToStr(pJson.xVoidValue)
		if strNum == ""{
			strNum = "0"
		}
		return CastStrToI64(strNum)
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		if pJson.xVoidValue.(bool) {
			return 1
		}else{
			return 0
		}
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return 0
	}

	return 0
}

func jsonInnerValueToFloat(pJson *JsonVal) float64{
	if pJson == nil{
		return 0
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return 0
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return CastStrToF64( pJson.xVoidValue.(string) )
	}else if pJson.iValueType == JSON_TYPE_NUM{
		strNum := CastVoidToStr(pJson.xVoidValue)
		if strNum == ""{
			strNum = "0"
		}
		return CastStrToF64(strNum)
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		if pJson.xVoidValue.(bool) {
			return 1
		}else{
			return 0
		}
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return 0
	}

	return 0
}

func jsonInnerValueToStr(pJson *JsonVal) string{
	if pJson == nil{
		return ""
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return ""
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return pJson.xVoidValue.(string)
	}else if pJson.iValueType == JSON_TYPE_NUM{
		strNum := CastVoidToStr(pJson.xVoidValue)
		if strNum == ""{
			strNum = "0"
		}
		return strNum
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		if pJson.xVoidValue.(bool) {
			return JSON_STR_TRUE
		}else{
			return JSON_STR_FALSE
		}
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return ""
	}else{
		return JsonToStr(pJson)
	}

	return ""
}

func jsonInnerValueToBool(pJson *JsonVal) bool{
	if pJson == nil{
		return false
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return false
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return CastStrToBool( pJson.xVoidValue.(string) )
	}else if pJson.iValueType == JSON_TYPE_NUM{
		strNum := CastVoidToStr(pJson.xVoidValue)
		if strNum == ""{
			strNum = "0"
		}
		return CastStrToBool(strNum)
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		return pJson.xVoidValue.(bool)
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return false
	}

	return false
}

func jsonInnerValueToVoid(pJson *JsonVal) interface{}{
	if pJson == nil{
		return nil
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return nil
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return pJson.xVoidValue
	}else if pJson.iValueType == JSON_TYPE_NUM{
		return pJson.xVoidValue
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		return pJson.xVoidValue.(bool)
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return nil
	}else{
		return JsonToStr(pJson)
	}

	return nil
}

func jsonInnerValueOutputStr(pJson *JsonVal) string{
	if pJson == nil{
		return ""
	}

	if pJson.iValueType == JSON_TYPE_NONE{
		return JSON_STR_NULL
	}else if pJson.iValueType == JSON_TYPE_TEXT{
		return JSON_STR_D_QUOTE +  StrAddSlashes( pJson.xVoidValue.(string), JSON_CH_D_QUOTE) + JSON_STR_D_QUOTE
	}else if pJson.iValueType == JSON_TYPE_NUM{
		strNum := CastVoidToStr(pJson.xVoidValue)
		if strNum == ""{
			strNum = "0"
		}
		return strNum
	}else if pJson.iValueType == JSON_TYPE_BOOL{
		if pJson.xVoidValue.(bool) {
			return JSON_STR_TRUE
		}else{
			return JSON_STR_FALSE
		}
	}else if pJson.iValueType == JSON_TYPE_NULL{
		return JSON_STR_NULL
	}

	return ""
}

func JsonType(pJson *JsonVal) int{
	if pJson == nil{
		return JSON_TYPE_NONE
	}

	return pJson.iValueType
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//Marshal

func JsonMarshal(xObj interface{}) string{
	if xObj == nil{
		return ""
	}
	buf,err := json.Marshal(xObj)
	if err != nil{
		return ""
	}else{
		return string(buf)
	}
}

func JsonUnmarshal(strJson string, xObj interface{}) error{
	if strJson == "" || xObj == nil{
		return errors.New("JsonUnmarshal empty param")
	}
	err := json.Unmarshal([]byte(strJson),xObj)
	if err != nil{
		return err
	}else{
		return nil
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
//Parser

type JsonParser struct{
	strSrc string
	iPos int
	jsonRoot *JsonVal
	pStack *list.List
	bWrong bool
}

func JsonParserNew(strSrc string, jsonRoot *JsonVal) *JsonParser{
	parserNew := JsonParser{}
	parserNew.jsonRoot = jsonRoot
	parserNew.strSrc = strings.TrimSpace(strSrc)
	parserNew.iPos = 0
	parserNew.pStack = ListCreate()
	parserNew.bWrong = false

	ListPushBack(parserNew.pStack,parserNew.jsonRoot)

	return &parserNew
}

func JsonParserResult(pParser *JsonParser) *JsonVal{
	if pParser == nil{
		return nil
	}

	if pParser.bWrong{
		return nil
	}

	if pParser.jsonRoot == nil{
		return nil
	}

	if !JsonIsCollection(pParser.jsonRoot){
		return nil
	}

	return pParser.jsonRoot
}

func JsonParserStackTopJson(pParser *JsonParser) *JsonVal{
	if pParser == nil{
		return nil
	}

	xTopElement := ListPeekBack(pParser.pStack)
	if xTopElement == nil{
		return nil
	}

	jsonElement := xTopElement.(*JsonVal)

	return jsonElement
}

func JsonParserStackPush(pParser *JsonParser, jsonNew *JsonVal) {
	if pParser == nil{
		return
	}

	ListPushBack(pParser.pStack,jsonNew)
}

func JsonParserStackPop(pParser *JsonParser) {
	if pParser == nil{
		return
	}

	ListPopBack(pParser.pStack)
}

func JsonParserParseByPos(pParser *JsonParser) int{
	if pParser == nil{
		return -1
	}

	if pParser.iPos >= len(pParser.strSrc){
		pParser.bWrong = true
		return -2
	}

	iResult := 0
	iNextNb := -1

	iNextNb = StrFindFirstNonBlank(pParser.strSrc,pParser.iPos)
	if iNextNb < 0{
		pParser.bWrong = true
		return -3
	}

	chFirstNb := pParser.strSrc[iNextNb]
	if chFirstNb == JSON_CH_L_B_L {
		jsonCheckCollectionForObject( JsonParserStackTopJson(pParser) )
		pParser.iPos += 1
		bLoopOver := false
		bEndingWithComma := false
		for {
			if pParser.iPos >= len(pParser.strSrc){
				pParser.bWrong = true
				return -4
			}

			iNextNb = StrFindFirstNonBlank(pParser.strSrc, pParser.iPos)
			if JSON_CH_L_B_R == pParser.strSrc[iNextNb]{
				if bEndingWithComma{
					pParser.bWrong = true
					return -5
				}
				pParser.iPos = iNextNb + 1
				bLoopOver = true
				break
			}else if JSON_CH_D_QUOTE == pParser.strSrc[iNextNb]{
				iDqBegin := iNextNb
				iDqEnd := StrFindFirstWithoutSlash(pParser.strSrc,JSON_CH_D_QUOTE,iDqBegin + 1)
				if iDqEnd < 0{
					pParser.bWrong = true
					return -6
				}

				strObjKey := pParser.strSrc[iDqBegin+1:iDqEnd]
				iColonPos := StrFindFirstNonBlank(pParser.strSrc,iDqEnd + 1)
				if pParser.strSrc[iColonPos] != JSON_CH_COLON{
					pParser.bWrong = true
					return -7
				}

				jsonSub := JsonNewNone()
				JsonObjSetJson(JsonParserStackTopJson(pParser),strObjKey,jsonSub)
				JsonParserStackPush(pParser,jsonSub)
				iNextNb = StrFindFirstNonBlank(pParser.strSrc, iColonPos + 1)
				pParser.iPos = iNextNb
				iResult = JsonParserParseByPos(pParser)
				if iResult < 0{
					return iResult
				}
				JsonParserStackPop(pParser)

				iNextNb = StrFindFirstNonBlank(pParser.strSrc,pParser.iPos)
				if iNextNb < 0{
					pParser.bWrong = true
					return -16
				}
				if pParser.strSrc[iNextNb] == JSON_CH_COMMA {
					pParser.iPos = iNextNb + 1
					bEndingWithComma = true
				}else{
					pParser.iPos = iNextNb
					bEndingWithComma = false
				}
			}
		}
		if bLoopOver{
			return 0
		}
	}else if chFirstNb == JSON_CH_M_B_L {
		jsonCheckCollectionForArray(JsonParserStackTopJson(pParser))
		pParser.iPos += 1

		bLoopOver := false
		bEndingWithComma := false
		for {
			if pParser.iPos >= len(pParser.strSrc){
				pParser.bWrong = true
				return -8
			}

			iNextNb = StrFindFirstNonBlank(pParser.strSrc,pParser.iPos)

			if pParser.strSrc[iNextNb] == JSON_CH_M_B_R{
				if bEndingWithComma{
					pParser.bWrong = true
					return -9
				}

				pParser.iPos = iNextNb + 1
				bLoopOver = true
				break
			}else if pParser.strSrc[iNextNb] == JSON_CH_COMMA || pParser.strSrc[iNextNb] == JSON_CH_COLON || pParser.strSrc[iNextNb] == JSON_CH_L_B_R{
				pParser.bWrong = true
				return -10
			}else{
				jsonSub := JsonNewNone()
				JsonArrAppendJson( JsonParserStackTopJson(pParser),jsonSub )
				JsonParserStackPush( pParser,jsonSub )
				pParser.iPos = iNextNb
				iResult = JsonParserParseByPos(pParser)
				if iResult < 0{
					return iResult
				}
				JsonParserStackPop(pParser)
				iNextNb = StrFindFirstNonBlank(pParser.strSrc,pParser.iPos)
				if iNextNb < 0{
					pParser.bWrong = true
					return -17
				}
				if pParser.strSrc[iNextNb] == JSON_CH_COMMA{
					pParser.iPos = iNextNb + 1
					bEndingWithComma = true
				}else{
					pParser.iPos = iNextNb
					bEndingWithComma = false
				}
			}
		}
		if bLoopOver{
			return 0
		}
	}else if chFirstNb == JSON_CH_D_QUOTE {
		if ListLen( pParser.pStack ) <= 1{
			pParser.bWrong = true
			return -11
		}

		iDqBegin := pParser.iPos
		iDqEnd := StrFindFirstWithoutSlash(pParser.strSrc,JSON_CH_D_QUOTE,iDqBegin + 1)
		if iDqEnd < 0{
			pParser.bWrong = true
			return -12
		}

		strVal := pParser.strSrc[iDqBegin+1:iDqEnd]
		strVal = StrStripSlashes(strVal)
		jsonInitValueAsText( JsonParserStackTopJson(pParser),strVal )
		pParser.iPos = iDqEnd + 1
	}else{
		if ListLen( pParser.pStack ) <= 1{
			pParser.bWrong = true
			return -13
		}

		strInvalidChars := JSON_STR_COMMA + JSON_STR_L_B_R + JSON_STR_M_B_R
		iValBegin := pParser.iPos
		iValEnd := StrFindChars(pParser.strSrc,strInvalidChars,pParser.iPos) - 1
		if iValEnd < iValBegin{
			pParser.bWrong = true
			return -14
		}
		strRaw := pParser.strSrc[iValBegin:iValEnd+1]
		strRaw = strings.TrimSpace(strRaw)
		if strRaw == JSON_STR_NULL{
			jsonInitValueAsNull( JsonParserStackTopJson(pParser) )
		}else if strRaw == JSON_STR_TRUE{
			jsonInitValueAsBool( JsonParserStackTopJson(pParser),true )
		}else if strRaw == JSON_STR_FALSE{
			jsonInitValueAsBool( JsonParserStackTopJson(pParser),false )
		}else if StrIsNum(strRaw){
			if StrIsInt(strRaw){
				jsonInitValueAsNum( JsonParserStackTopJson(pParser),CastStrToI64(strRaw) )
			}else{
				jsonInitValueAsNum( JsonParserStackTopJson(pParser),CastStrToF64(strRaw) )
			}
		}else{
			pParser.bWrong = true
			return -15
		}
		pParser.iPos = iValEnd + 1
	}

	return 0
}

/////////////////////////////////////////////////////////////////
func JsonFromStr(strSrc string) *JsonVal{
	jsonRoot := JsonNewNone()
	parserNew := JsonParserNew( strSrc,jsonRoot )
	iPos := JsonParserParseByPos(parserNew)
	if iPos < 0{
		return nil
	}
	jsonResult := JsonParserResult(parserNew)
	return jsonResult
}

func JsonRefresh(pJson *JsonVal, strSrc string) error{
	if pJson == nil{
		return errors.New("JsonRefresh empty")
	}
	parserNew := JsonParserNew( strSrc,pJson )
	iPos := JsonParserParseByPos(parserNew)
	if iPos < 0{
		return errors.New("JsonRefresh parsed with error")
	}
	jsonResult := JsonParserResult(parserNew)
	if jsonResult == nil{
		return errors.New("JsonRefresh parsed with error")
	}else{
		return nil
	}
}

func JsonToArrVoid(pJson *JsonVal) []interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return nil
	}

	arrVoid := make([]interface{},0)
	for _,v := range pJson.xVoidValue.([]interface{}){
		jsonSub := v.(*JsonVal)
		arrVoid = append(arrVoid, jsonInnerValueToVoid(jsonSub))
	}

	return arrVoid
}

func JsonToArrStr(pJson *JsonVal) []string{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return nil
	}

	arrStr := make([]string,0)
	for _,v := range pJson.xVoidValue.([]interface{}){
		jsonSub := v.(*JsonVal)
		arrStr = append(arrStr, jsonInnerValueToStr(jsonSub))
	}

	return arrStr
}

func JsonToArrJson(pJson *JsonVal) []*JsonVal{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return nil
	}

	arrJson := make([]*JsonVal,0)
	for _,v := range pJson.xVoidValue.([]interface{}){
		jsonSub := v.(*JsonVal)
		arrJson = append(arrJson, jsonSub)
	}

	return arrJson
}

func JsonToArrJsonAsVoid(pJson *JsonVal) []interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return nil
	}

	return pJson.xVoidValue.([]interface{})
}

func JsonToMapVoid(pJson *JsonVal)map[string]interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return nil
	}

	mapVoid := make(map[string]interface{})
	for k,v := range pJson.xVoidValue.(map[string]interface{}){
		jsonSub := v.(*JsonVal)
		mapVoid[k] = jsonInnerValueToVoid(jsonSub)
	}

	return mapVoid
}

func JsonToMapStr(pJson *JsonVal)map[string]string{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return nil
	}

	mapVoid := make(map[string]string)
	for k,v := range pJson.xVoidValue.(map[string]interface{}){
		jsonSub := v.(*JsonVal)
		mapVoid[k] = jsonInnerValueToStr(jsonSub)
	}

	return mapVoid
}

func JsonToMapJson(pJson *JsonVal)map[string]*JsonVal{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return nil
	}

	mapJson := make(map[string]*JsonVal)
	for k,v := range pJson.xVoidValue.(map[string]interface{}){
		jsonSub := v.(*JsonVal)
		mapJson[k] = jsonSub
	}

	return mapJson
}

func JsonToMapJsonAsVoid(pJson *JsonVal)map[string]interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_OBJECT{
		return nil
	}

	return pJson.xVoidValue.(map[string]interface{})
}

func JsonToArrMapStr(pJson *JsonVal)[]map[string]string{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return nil
	}

	arrMapStr := make([]map[string]string,0)

	arrJson := JsonToArrJson(pJson)
	for _,v := range arrJson{
		arrMapStr = append(arrMapStr , JsonToMapStr(v) )
	}

	return arrMapStr
}

func JsonToArrMapVoid(pJson *JsonVal)[]map[string]interface{}{
	if pJson == nil || pJson.iValueType != JSON_TYPE_ARRAY{
		return nil
	}

	arrMapStr := make([]map[string]interface{},0)

	arrJson := JsonToArrJson(pJson)
	for _,v := range arrJson{
		arrMapStr = append(arrMapStr , JsonToMapVoid(v) )
	}

	return arrMapStr
}