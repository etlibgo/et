package et

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

/**
json封装类
已经隔离框架
实现可以是go内部json解析，也可以是其他第三方json.
可以根据需要切换而不影响业务代码。
建议在代码中不要直接使用第三方类型，比如simplejson
 */

/**
隐藏第三方类型，降低框架耦合度
 */
type ExJsonVal  JsonVal


func (my* ExJsonVal)String() string{
	return ExJsonToStr(my)
}
/**
新建一个空json对象
 */
func exJsonNew() *ExJsonVal {
	return ExJsonRawToPtr(JsonNewNone())
}

/**
新建一个空json对象
*/
func ExJsonNewArray() *ExJsonVal {
	return ExJsonRawToPtr(JsonNewArray())
}

/**
新建一个空json对象
*/
func ExJsonNewObject() *ExJsonVal {
	return ExJsonRawToPtr(JsonNewObject())
}

/**
解析字符串为json对象
 */
func ExJsonParse(strText string) *ExJsonVal {
	return ExJsonRawToPtr(JsonFromStr(strText))
}

/**
把第三方原始对象转换为封装后的类型。
 */
func ExJsonRawToPtr(pRawJson *JsonVal) *ExJsonVal {
	return (*ExJsonVal)(pRawJson)
}

/**
把封装类型的指针转换为第三方原始类型的指针
 */
func ExJsonRawFromPtr(pExJson *ExJsonVal) *JsonVal{
	return (*JsonVal)(pExJson)
}

/**
把第三方原始类型的值转换为封装类型的值
 */
func ExJsonRawToVal(valRawJson JsonVal) ExJsonVal {
	return (ExJsonVal)(valRawJson)
}

/**
把封装类型的值转换为原始类型的值
 */
func ExJsonRawFromVal(valExJson ExJsonVal) JsonVal{
	return (JsonVal)(valExJson)
}

/**
把第三方原始类型转换为字符串
 */
func ExJsonRawToStr(pRawJson *JsonVal) string{
	if pRawJson == nil{
		return ""
	}

	return JsonToStr(pRawJson)
}

/**
把json对象转换为字符串
 */
func ExJsonToStr(pExJson *ExJsonVal) string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return ""
	}

	return JsonToStr(pRawJson)
}

func ExJsonClear(pExJson *ExJsonVal){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return
	}

	if JsonType(pRawJson) == JSON_TYPE_ARRAY{
		JsonArrClear(pRawJson)
	}else if JsonType(pRawJson) == JSON_TYPE_OBJECT{
		JsonObjClear(pRawJson)
	}

}

/**
判断此json节点的类型：空，对象，数组，文本，null，数字，
 */
func ExJsonType(pExJson *ExJsonVal) int{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return JSON_TYPE_NONE
	}

	return JsonType(pRawJson)
}

/**
设置json的一个属性，strKey是键名
 */
func ExJsonSet(pExJson *ExJsonVal, strKey string, xValue interface{}){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return
	}

	if xValue != nil{
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
		if strings.Index(typeVoid.Name(),"ExJsonVal" ) >= 0 {
			if isPtr{
				JsonObjSetVoid(pRawJson,strKey,ExJsonRawFromPtr(xValue.(*ExJsonVal)))
			}else{
				JsonObjSetVoid(pRawJson,strKey,ExJsonRawFromVal(xValue.(ExJsonVal)))
			}
			return
		}
	}

	JsonObjSetVoid(pRawJson,strKey,xValue)
}

func ExJsonSetChild(pExJson *ExJsonVal, strKey string, pExJsonChild *ExJsonVal){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return
	}

	JsonObjSetJson(pRawJson,strKey, ExJsonRawFromPtr(pExJsonChild) )
}

func ExJsonSetByIndex(pExJson *ExJsonVal, iIndex int, xValue interface{}){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return
	}

	if xValue != nil{
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
		if strings.Index(typeVoid.Name(),"ExJsonVal" ) >= 0{
			if isPtr{
				JsonArrSetVoid(pRawJson,iIndex, ExJsonRawFromPtr(xValue.(*ExJsonVal)) )
			}else{
				JsonArrSetVoid(pRawJson,iIndex, ExJsonRawFromVal(xValue.(ExJsonVal)) )
			}

			return
		}
	}

	JsonArrSetVoid(pRawJson,iIndex, xValue )
}

func ExJsonSetChildByIndex(pExJson *ExJsonVal, iIndex int, pExJsonChild *ExJsonVal){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return
	}

	JsonArrSetJson(pRawJson,iIndex, ExJsonRawFromPtr(pExJsonChild) )
}

/**
此json是否为空，空对象或者空数组都是空的
 */
func ExJsonIsEmpty(pExJson *ExJsonVal) bool{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return true
	}

	return JsonIsEmpty(pRawJson)
}

/**
把json字符串解析到现有的son对象中，与parse的区别是，不产生新对象。
 */
func ExJsonUnmarshal(pExJson *ExJsonVal, strJson string) error{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return errors.New("ExJsonUnmarshal param nil")
	}

	return JsonRefresh(pRawJson,strJson)
}

/**
把json对象格式化输出为字符串，会有良好的排版。
 */
func ExJsonToPrettyStr(pExJson *ExJsonVal) string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return ""
	}

	return JsonToFormattedStr(pRawJson)
}

/**
把json对象转为字符串map的数组，是一个2维对象
 */
func ExJsonToArrMapStr(pExJson *ExJsonVal) []map[string]string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return []map[string]string{}
	}

	return JsonToArrMapStr(pRawJson)
}

/**
把json对象转换为空类型map的数组，是一个2维对象
 */
func ExJsonToArrMapVoid(pExJson *ExJsonVal) []map[string]interface{}{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return []map[string]interface{}{}
	}

	return JsonToArrMapVoid(pRawJson)
}

/**
把此json的键遍历出来，得到一个数组
 */
func ExJsonMapKeys(pExJson *ExJsonVal) []string{
	pRawJson := ExJsonRawFromPtr(pExJson)

	return JsonObjKeys(pRawJson)
}

/**
判断此json中是否包含某个键
 */
func ExJsonMapContains(pExJson *ExJsonVal, strKey string) bool{
	pRawJson := ExJsonRawFromPtr(pExJson)

	return JsonObjContainsKey(pRawJson,strKey)
}

/**
若此json为数组，得到此数组的大小
 */
func ExJsonArrSize(pExJson *ExJsonVal) int{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return 0
	}
	return JsonLen(pRawJson)
}

/**
把json对象转换为无类型map
 */
func ExJsonToMapVoid(pExJson *ExJsonVal) map[string]interface{}{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return map[string]interface{}{}
	}

	return JsonToMapVoid(pRawJson)
}

/**
把json对象转换为字符串map
 */
func ExJsonToMapStr(pExJson *ExJsonVal) map[string]string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return map[string]string{}
	}
	return JsonToMapStr(pRawJson)
}

/**
把一个字符串map的内容的所有键遍历读到一个json中
 */
func ExJsonReadMapStr(pExJson *ExJsonVal, mapSrc map[string]string){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil || mapSrc == nil{
		return
	}
	for k,v := range mapSrc {
		JsonObjSetText(pRawJson,k,v)
	}
}

/**
把一个空类型map的内容的所有键遍历读到一个json中
*/
func ExJsonReadMapVoid(pExJson *ExJsonVal, mapSrc map[string]interface{}){
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil || mapSrc == nil{
		return
	}
	for k,v := range mapSrc {
		JsonObjSetVoid(pRawJson,k,v)
	}
}

/**
把json对象转换为字符串数组
 */
func ExJsonToArrStr(pExJson *ExJsonVal) []string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return []string{}
	}

	return JsonToArrStr(pRawJson)
}

/**
把json对象转化为空类型数组
 */
func ExJsonToArrVoid(pExJson *ExJsonVal) []interface{}{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return []interface{}{}
	}

	return JsonToArrVoid(pRawJson)
}

/**
把一个结构体转换为json对象
 */
func ExJsonFromStruct(pExJson *ExJsonVal, xStruct interface{}) error{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if xStruct == nil{
		return errors.New("empty struct")
	}
	strJson := JsonMarshal(xStruct)
	return JsonRefresh(pRawJson,strJson)
}

/**
根据一个结构体，创建新的json对象
 */
func ExJsonFromStructNew(xStruct interface{}) *ExJsonVal{
	if xStruct == nil{
		return nil
	}
	strJson := JsonMarshal(xStruct)
	return ExJsonRawToPtr(JsonFromStr(strJson))
}

/**
把json对象转换到一个已存在的结构体中
 */
func ExJsonToStruct(pExJson *ExJsonVal, pStruct interface{}) error{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return errors.New("ExJsonToStruct json is nil")
	}

	if pStruct == nil{
		return errors.New("ExJsonToStruct pStruct is nil")
	}

	strJson := JsonMarshal(pStruct)
	err := json.Unmarshal([]byte(strJson),pStruct)

	if err != nil{
		return err
	}
	return nil
}

/**
根据字段名而不是tag把一个结构体转化为json
 */
func ExJsonFromStructByFieldName(pExJson *ExJsonVal, xStruct interface{}) error{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return errors.New("ExJsonToStruct json is nil")
	}

	if xStruct == nil{
		return errors.New("ExJsonToStruct pStruct is nil")
	}

	valueVoid := reflect.ValueOf(xStruct)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
	}

	for i := 0; i < typeVoid.NumField(); i++{
		strKey := typeVoid.Field(i).Name
		JsonObjSetVoid(pRawJson,strKey,valueVoid.FieldByName(strKey).Interface())
	}

	return nil
}

/**
根据字段名，把json对象转换到一个结构体中
 */
func ExJsonToStructByFieldName(pExJson *ExJsonVal, pStruct interface{}) error{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return errors.New("ExJsonToStruct json is nil")
	}

	if pStruct == nil{
		return errors.New("ExJsonToStruct pStruct is nil")
	}

	valueVoid := reflect.ValueOf(pStruct)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
	}

	for i := 0; i < typeVoid.NumField(); i++{
		strKey := typeVoid.Field(i).Name
		valueVoid.FieldByName(strKey).Set(reflect.ValueOf( JsonObjGetVoid(pRawJson,strKey)  ))
	}

	return nil
}

/**
给定一个结构体，把原先由tag得到的json转换为以字段名为键
 */
func ExJsonStructTagToField(pExJsonSrc *ExJsonVal, pStruct interface{}, pExJsonDest *ExJsonVal) error{
	pRawJsonSrc := ExJsonRawFromPtr(pExJsonSrc)
	pRawJsonDest := ExJsonRawFromPtr(pExJsonDest)
	if pRawJsonSrc == nil{
		return errors.New("ExJsonToStruct jsonSpSrc is nil")
	}

	if pRawJsonDest == nil{
		return errors.New("ExJsonToStruct jsonSpDest is nil")
	}

	if pStruct == nil{
		return errors.New("ExJsonToStruct pStruct is nil")
	}

	strSrc := ExJsonRawToStr(pRawJsonSrc)
	err := json.Unmarshal([]byte(strSrc), pStruct)
	if err != nil{
		return err
	}
	err = ExJsonFromStructByFieldName(pExJsonDest, pStruct)
	if err != nil{
		return err
	}

	return nil
}

/**
给定一个结构体，把原先由字段名得到的json转换为以为tag键
*/
func ExJsonStructFieldToTag(pExJsonSrc *ExJsonVal, pStruct interface{}, pExJsonDest *ExJsonVal) error{
	pRawJsonDest := ExJsonRawFromPtr(pExJsonDest)
	if pExJsonSrc == nil || pExJsonDest == nil{
		return errors.New("ExJsonStructFieldToTag parameter nil")
	}
	err := ExJsonToStructByFieldName(pExJsonSrc, pStruct)
	if err != nil{
		return err
	}

	strSrc := JsonMarshal(pStruct)

	return JsonRefresh(pRawJsonDest,strSrc)
}

/**
根据键名获得字符串
 */
func ExJsonGetStr(pExJson *ExJsonVal, strKey string) string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return ""
	}

	return JsonObjGetStr(pRawJson,strKey)
}

/**
根据索引获得字符串
*/
func ExJsonGetStrByIndex(pExJson *ExJsonVal, iIndex int) string{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil || iIndex < 0{
		return ""
	}

	return JsonArrGetStr(pRawJson,iIndex)
}

/**
此json节点是否包含某个键
 */
func ExJsonContains(pExJson *ExJsonVal, strKey string) bool{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return false
	}
	return JsonObjContainsKey(pRawJson,strKey)
}

/**
根据键名获得整数
*/
func ExJsonGetInt(pExJson *ExJsonVal, strKey string) int{
	pRawJson := ExJsonRawFromPtr(pExJson)
	iResult := JsonObjGetInt(pRawJson,strKey)
	return iResult
}

/**
根据键名获得bool值
*/
func ExJsonGetBool(pExJson *ExJsonVal, strKey string) bool{
	pRawJson := ExJsonRawFromPtr(pExJson)
	bResult := JsonObjGetBool(pRawJson,strKey)
	return bResult
}

/**
根据键名获得长整型
*/
func ExJsonGetI64(pExJson *ExJsonVal, strKey string) int64{
	pRawJson := ExJsonRawFromPtr(pExJson)
	iResult := JsonObjGetI64(pRawJson,strKey)
	return iResult
}

/**
根据键名获得浮点数
 */
func ExJsonGetFloat(pExJson *ExJsonVal, strKey string) float64{
	pRawJson := ExJsonRawFromPtr(pExJson)
	fResult := JsonObjGetFloat(pRawJson,strKey)
	return fResult
}

/**
根据键名获得下级json节点
*/
func ExJsonGetChild(pExJson *ExJsonVal, strKey string) *ExJsonVal{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return nil
	}

	return ExJsonRawToPtr( JsonObjGetJson(pRawJson,strKey) )
}

/**
根据索引获得整数
*/
func ExJsonGetIndexInt(pExJson *ExJsonVal, iIndex int) int{
	pRawJson := ExJsonRawFromPtr(pExJson)
	iResult := JsonArrGetInt(pRawJson,iIndex)
	return iResult
}

/**
根据索引获得布尔值
*/
func ExJsonGetIndexBool(pExJson *ExJsonVal, iIndex int) bool{
	pRawJson := ExJsonRawFromPtr(pExJson)
	bResult := JsonArrGetBool(pRawJson,iIndex)
	return bResult
}

/**
根据索引获得长整型
*/
func ExJsonGetIndexI64(pExJson *ExJsonVal, iIndex int) int64{
	pRawJson := ExJsonRawFromPtr(pExJson)
	iResult := JsonArrGetI64(pRawJson,iIndex)
	return iResult
}

/**
根据索引获得浮点数
*/
func ExJsonGetIndexFloat(pExJson *ExJsonVal, iIndex int) float64{
	pRawJson := ExJsonRawFromPtr(pExJson)
	fResult := JsonArrGetFloat(pRawJson,iIndex)
	return fResult
}

/**
根据索引获得下级json
*/
func ExJsonGetChildByIndex(pExJson *ExJsonVal, iIndex int) *ExJsonVal{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return nil
	}

	return ExJsonRawToPtr( JsonArrGetJson(pRawJson,iIndex) )
}

/**
把一组struct转换为json
 */
func ExJsonFromStructInGroup(groupStruct interface{}) *ExJsonVal {
	return ExJsonRawToPtr(  ExJsonFromStructInGroupRecur(groupStruct).(*JsonVal) )
}

/**
把一组struct转换为json，递归
*/
func ExJsonFromStructInGroupRecur(groupStruct interface{}) interface{} {
	jsonResult := JsonNewNone()

	valueSrc := reflect.ValueOf(groupStruct)
	typeSrc := valueSrc.Type()
	kindSrc := typeSrc.Kind()

	if kindSrc == reflect.Ptr{
		valueSrc = valueSrc.Elem()
		typeSrc = valueSrc.Type()
		kindSrc = typeSrc.Kind()
	}

	if kindSrc == reflect.Struct{
		err := ExJsonFromStruct(ExJsonRawToPtr(jsonResult),groupStruct)
		if err != nil{
			return nil
		}
	}else if kindSrc == reflect.Map{
		jsonCheckCollectionForObject(jsonResult)
		for _, key := range valueSrc.MapKeys() {
			JsonObjSetVoid(jsonResult, CastVoidToStr(key.Interface()),ExJsonFromStructInGroupRecur( valueSrc.MapIndex(key).Interface() ))
		}
	}else if kindSrc == reflect.Array || kindSrc == reflect.Slice{
		jsonCheckCollectionForArray(jsonResult)
		for i := 0; i < valueSrc.Len(); i++ {
			xTemp := ExJsonFromStructInGroupRecur(valueSrc.Index(i).Interface())
			JsonArrSetVoid(jsonResult,i,xTemp)
		}
	}else{
		return CastVoidToStr(valueSrc.Interface())
	}

	return jsonResult
}

/**
把一组结构体序列化为字符串
 */
func ExJsonFromStructInGroupToStr(groupStruct interface{}) string {
	return ExJsonToStr(ExJsonFromStructInGroup(groupStruct))
}

/**
把字符串反序列化为结构体，利用json做转换
 */
func ExJsonStructFromStr(strJson string,xObj interface{}) error{
	if strJson == ""{
		return errors.New("ExJsonStructFromStr empty json")
	}

	err :=json.Unmarshal([]byte(strJson),xObj)
	if err != nil{
		return err
	}

	return nil
}

/**
把结构体序列化为字符串，
 */
func ExJsonStructToStr(xObj interface{}) string{
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

/**
把字符串转化为json，与parse一样
 */
func ExJsonFromStr(strJson string) *ExJsonVal{
	return ExJsonParse(strJson)
}

/**
把字节流转化为json
 */
func ExJsonFromBuf(bufJson []byte) *ExJsonVal{
	strSrc := string(bufJson)
	return  ExJsonRawToPtr(JsonFromStr(strSrc))
}

/**
根据jsonpath得到一个字符串，path是一个数组，依次是每个节点的键名，如果是数组索引请用[]包装数字
 */
func ExJsonPathGetStr(pExJson *ExJsonVal, arrPath []string) string {
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return ""
	}

	if len(arrPath) == 0{
		return ""
	}

	var jsonParent *JsonVal = pRawJson
	lenArr := len(arrPath)

	for i := 0; i < lenArr; i++{
		strKey := arrPath[i]

		if strings.HasPrefix(strKey,"[") && strings.HasSuffix(strKey,"]"){
			strKey = strKey[1: len(strKey) - 1 ]
			strKey = strings.TrimSpace(strKey)
			if !StrIsInt(strKey){
				return ""
			}
			iNum := CastStrToInt(strKey)
			if i == lenArr - 1{
				return ExJsonGetStrByIndex( ExJsonRawToPtr( jsonParent ),iNum)
			}else{
				jsonParent = JsonArrGetJson(jsonParent,iNum)
			}
		}else{
			if i == lenArr - 1{
				return ExJsonGetStr(ExJsonRawToPtr(jsonParent),strKey)
			}else{
				jsonParent = JsonObjGetJson(jsonParent,strKey)
			}
		}

		if jsonParent == nil{
			break
		}
	}

	return ""
}

func ExJsonPathGetChild(pExJson *ExJsonVal, arrPath []string) *ExJsonVal {
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return nil
	}

	if len(arrPath) == 0{
		return nil
	}

	var jsonParent *JsonVal = pRawJson
	lenArr := len(arrPath)

	for i := 0; i < lenArr; i++{
		strKey := arrPath[i]

		if strings.HasPrefix(strKey,"[") && strings.HasSuffix(strKey,"]"){
			strKey = strKey[1: len(strKey) - 1 ]
			strKey = strings.TrimSpace(strKey)
			if !StrIsInt(strKey){
				return nil
			}
			iNum := CastStrToInt(strKey)
			if i == lenArr - 1{
				return ExJsonGetChildByIndex( ExJsonRawToPtr( jsonParent ),iNum)
			}else{
				jsonParent = JsonArrGetJson(jsonParent,iNum)
			}
		}else{
			if i == lenArr - 1{
				return ExJsonGetChild(ExJsonRawToPtr(jsonParent),strKey)
			}else{
				jsonParent = JsonObjGetJson(jsonParent,strKey)
			}
		}

		if jsonParent == nil{
			break
		}
	}

	return nil
}

/**
根据jsonpath设置json，path是一个数组，依次是每个节点的键名，如果是数组索引请用[]包装数字
 */
func ExJsonPathSetValue(pExJson *ExJsonVal, arrPath []string, xValue interface{}) error {
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return errors.New("json is nil")
	}

	if len(arrPath) == 0{
		return errors.New("jsonpath is empty")
	}

	var jsonParent *JsonVal = pRawJson
	lenArr := len(arrPath)

	for i := 0; i < lenArr; i++{
		strKey := arrPath[i]

		if strings.HasPrefix(strKey,"[") && strings.HasSuffix(strKey,"]"){
			strKey = strKey[1: len(strKey) - 1 ]
			strKey = strings.TrimSpace(strKey)
			if !StrIsInt(strKey){
				return errors.New("index error:" + strKey)
			}
			iNum := CastStrToInt(strKey)
			if i == lenArr - 1{
				JsonArrSetVoid(jsonParent,iNum,xValue)
			}else{
				jsonParent = JsonArrGetJson(jsonParent,iNum)
			}
		}else{
			if i == lenArr - 1{
				JsonObjSetVoid(jsonParent,strKey,xValue)
				break
			}else{
				if ExJsonContains(ExJsonRawToPtr(jsonParent),strKey){
					jsonParent = JsonObjGetJson(jsonParent,strKey)
				}else{
					return errors.New("no this key in json map:" + strKey)
				}
			}
		}

		if jsonParent == nil{
			break
		}
	}

	return nil
}

/**
克隆一个JSON对象
 */
func ExJsonClone(pExJson *ExJsonVal) *ExJsonVal{
	pRawJson := ExJsonRawFromPtr(pExJson)
	if pRawJson == nil{
		return nil
	}

	strJson := JsonToStr(pRawJson)
	return ExJsonRawToPtr(JsonFromStr(strJson))
}
