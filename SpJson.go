package et

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"reflect"
	"strings"
)

func SpJsonToStr(jsonSp *simplejson.Json) string{
	if jsonSp == nil{
		return ""
	}

	buf,err := jsonSp.MarshalJSON()
	if err != nil{
		return ""
	}

	strResult := string(buf)
	return strResult
}

func SpJsonType(jsonSp *simplejson.Json) int{
	if jsonSp == nil{
		return JSON_TYPE_NONE
	}

	bytes,err := jsonSp.MarshalJSON()
	if err != nil{
		return JSON_TYPE_NONE
	}
	strJson := string(bytes)

	if strJson == "" || strJson == "null"{
		return JSON_TYPE_NONE
	}

	if strings.HasPrefix(strJson,"{"){
		return JSON_TYPE_OBJECT
	}

	if strings.HasPrefix(strJson,"["){
		return JSON_TYPE_ARRAY
	}

	if strings.HasPrefix(strJson,"\""){
		return JSON_TYPE_TEXT
	}

	if strings.EqualFold(strJson,"null"){
		return JSON_TYPE_NULL
	}

	if strings.EqualFold(strJson,"true") || strings.EqualFold(strJson,"false"){
		return JSON_TYPE_BOOL
	}

	if StrIsNum(strJson){
		return JSON_TYPE_NUM
	}

	return JSON_TYPE_NONE
}

func SpJsonIsEmpty(jsonSp *simplejson.Json) bool{
	if jsonSp == nil{
		return true
	}

	strJson := SpJsonToStr(jsonSp)
	if strJson == ""{
		return true
	}

	if strings.HasPrefix(strJson,"{"){
		if len(strJson) > 2{
			return false
		}
	} else if strings.HasPrefix(strJson,"["){
		if len(strJson) > 2{
			return false
		}
	}

	return true
}

func SpJsonToPrettyStr(jsonSp  *simplejson.Json) string{
	if jsonSp == nil{
		return ""
	}

	bufResult,err := jsonSp.EncodePretty()
	if err != nil{
		return ""
	}
	strResult := string(bufResult)
	return strResult

}

func SpJsonToArrMapStr(jsonArr *simplejson.Json) []map[string]string{
	if jsonArr == nil{
		return []map[string]string{}
	}
	var jsonArrCopy *simplejson.Json = nil
	arrVoid,err := jsonArr.Array()
	if err != nil{
		jsonArrCopy = SpJsonFromStr(SpJsonToStr(jsonArr))
		arrVoid = jsonArrCopy.MustArray()
	}else{
		jsonArrCopy = jsonArr
	}
	lenArr := len(arrVoid)
	arrMapStr := make([]map[string]string,lenArr)
	for i := 0; i < lenArr; i++{
		jsonTemp := jsonArrCopy.GetIndex(i)
		mapStr := SpJsonToMapStr(jsonTemp)
		arrMapStr[i] = mapStr
	}
	return arrMapStr
}

func SpJsonToArrMapVoid(jsonArr *simplejson.Json) []map[string]interface{}{
	if jsonArr == nil{
		return []map[string]interface{}{}
	}
	var jsonArrCopy *simplejson.Json = nil
	arrVoid,err := jsonArr.Array()
	if err != nil{
		jsonArrCopy = SpJsonFromStr(SpJsonToStr(jsonArr))
		arrVoid = jsonArrCopy.MustArray()
	}else{
		jsonArrCopy = jsonArr
	}
	lenArr := len(arrVoid)
	arrMapVoid := make([]map[string]interface{},lenArr)
	for i := 0; i < lenArr; i++{
		jsonTemp := jsonArrCopy.GetIndex(i)
		mapVoid := SpJsonToMapVoid(jsonTemp)
		arrMapVoid[i] = mapVoid
	}
	return arrMapVoid
}

func SpJsonMapKeys(jsonSp *simplejson.Json) []string{
	arrKeys := make([]string,0)
	if jsonSp == nil{
		return arrKeys
	}
	var jsonSpCopy *simplejson.Json = nil
	mapVoid,err := jsonSp.Map()
	if err != nil{
		jsonSpCopy = SpJsonFromStr(SpJsonToStr(jsonSp))
		mapVoid = jsonSpCopy.MustMap()
	}else{
		jsonSpCopy = jsonSp
	}
	for k,_ := range mapVoid {
		arrKeys = append(arrKeys,k)
	}

	return arrKeys
}

func SpJsonMapContains(jsonSp *simplejson.Json, strKey string) bool{
	if jsonSp == nil{
		return false
	}
	var jsonSpCopy *simplejson.Json = nil
	mapVoid,err := jsonSp.Map()
	if err != nil{
		jsonSpCopy = SpJsonFromStr(SpJsonToStr(jsonSp))
		mapVoid = jsonSpCopy.MustMap()
	}else{
		jsonSpCopy = jsonSp
	}

	if _, bExisted := mapVoid[strKey]; bExisted {
		return true
	}else{
		return false
	}
}

func SpJsonArrSize(jsonSp *simplejson.Json) int{
	if jsonSp == nil{
		return 0
	}
	var jsonSpCopy *simplejson.Json = nil
	arrVoid,err := jsonSp.Array()
	if err != nil{
		jsonSpCopy = SpJsonFromStr( SpJsonToStr(jsonSp) )
		arrVoid = jsonSpCopy.MustArray()
	}else{
		jsonSpCopy = jsonSp
	}
	return len(arrVoid)
}

func SpJsonToMapVoid(jsonSp *simplejson.Json) map[string]interface{}{
	if jsonSp == nil{
		return map[string]interface{}{}
	}
	mapVoid,err := jsonSp.Map()
	if err != nil{
		jsonSpCopy := SpJsonFromStr( SpJsonToStr(jsonSp) )
		mapVoid = jsonSpCopy.MustMap()
	}

	return mapVoid
}

func SpJsonToMapStr(jsonSp *simplejson.Json) map[string]string{
	if jsonSp == nil{
		return map[string]string{}
	}
	mapVoid,err := jsonSp.Map()
	var jsonSpCopy *simplejson.Json = nil
	if err != nil{
		jsonSpCopy = SpJsonFromStr( SpJsonToStr(jsonSp) )
		mapVoid = jsonSpCopy.MustMap()
	}else{
		jsonSpCopy = jsonSp
	}

	mapStrResult := make(map[string]string)
	for k,_ := range mapVoid {
		mapStrResult[k] = SpJsonGetStr(jsonSpCopy,k)
	}

	return mapStrResult
}

func SpJsonReadMapStr(jsonSp *simplejson.Json, mapSrc map[string]string){
	if jsonSp == nil || mapSrc == nil{
		return
	}
	for k,v := range mapSrc {
		jsonSp.Set(k,v)
	}
}

func SpJsonReadMapVoid(jsonSp *simplejson.Json, mapSrc map[string]interface{}){
	if jsonSp == nil || mapSrc == nil{
		return
	}
	for k,v := range mapSrc {
		jsonSp.Set(k,v)
	}
}

func SpJsonToArrStr(jsonSp *simplejson.Json) []string{
	if jsonSp == nil{
		return []string{}
	}

	var jsonSpCopy *simplejson.Json = nil
	arrVoid,err := jsonSp.Array()
	if err != nil{
		jsonSpCopy = SpJsonFromStr( SpJsonToStr(jsonSp) )
		arrVoid = jsonSpCopy.MustArray()
	}else{
		jsonSpCopy = jsonSp
	}

	arrStr := make([]string,len(arrVoid))
	for i := 0; i < len(arrVoid); i++{
		arrStr[i] = CastVoidToStr(arrVoid[i])
	}

	return arrStr
}

func SpJsonToArrVoid(jsonSp *simplejson.Json) []interface{}{
	if jsonSp == nil{
		return []interface{}{}
	}

	var jsonSpCopy *simplejson.Json = nil
	arrVoid,err := jsonSp.Array()
	if err != nil{
		jsonSpCopy = SpJsonFromStr( SpJsonToStr(jsonSp) )
		arrVoid = jsonSpCopy.MustArray()
	}else{
		jsonSpCopy = jsonSp
	}

	return arrVoid
}

func SpJsonFromStruct(jsonSp *simplejson.Json, stObj interface{}) error{
	if stObj == nil{
		return errors.New("empty struct")
	}
	strJson := JsonMarshal(stObj)
	err := jsonSp.UnmarshalJSON([]byte(strJson))
	if err != nil{
		return err
	}
	return nil
}

func SpJsonFromStructNew(stObj interface{}) *simplejson.Json{
	if stObj == nil{
		return nil
	}
	jsonNew := simplejson.New()
	strJson := JsonMarshal(stObj)
	err := jsonNew.UnmarshalJSON([]byte(strJson))
	if err != nil{
		return nil
	}
	return jsonNew
}

func SpJsonToStruct(jsonSp *simplejson.Json, stObj interface{}) error{
	if jsonSp == nil{
		return errors.New("SpJsonToStruct json is nil")
	}

	if stObj == nil{
		return errors.New("SpJsonToStruct stObj is nil")
	}

	strJson := JsonMarshal(stObj)
	err := json.Unmarshal([]byte(strJson),stObj)

	if err != nil{
		return err
	}
	return nil
}

func SpJsonFromStructByFieldName(jsonSp *simplejson.Json, xStruct interface{}) error{
	if jsonSp == nil{
		return errors.New("SpJsonToStruct json is nil")
	}

	if xStruct == nil{
		return errors.New("SpJsonToStruct stObj is nil")
	}

	valueVoid := reflect.ValueOf(xStruct)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
		kindVoid = typeVoid.Kind()
	}

	for i := 0; i < typeVoid.NumField(); i++{
		strKey := typeVoid.Field(i).Name
		jsonSp.Set(strKey,valueVoid.FieldByName(strKey).Interface())
	}

	return nil
}

func SpJsonToStructByFieldName(jsonSp *simplejson.Json, xStruct interface{}) error{
	if jsonSp == nil{
		return errors.New("SpJsonToStruct json is nil")
	}

	if xStruct == nil{
		return errors.New("SpJsonToStruct stObj is nil")
	}

	valueVoid := reflect.ValueOf(xStruct)
	typeVoid := valueVoid.Type()
	kindVoid := typeVoid.Kind()

	if kindVoid == reflect.Ptr{
		valueVoid = valueVoid.Elem()
		typeVoid = valueVoid.Type()
		kindVoid = typeVoid.Kind()
	}

	for i := 0; i < typeVoid.NumField(); i++{
		strKey := typeVoid.Field(i).Name
		valueVoid.FieldByName(strKey).Set(reflect.ValueOf(jsonSp.Get(strKey).Interface()))
	}

	return nil
}

func SpJsonStructTagToField(jsonSpSrc *simplejson.Json, xStruct interface{}, jsonSpDest *simplejson.Json) error{
	if jsonSpSrc == nil{
		return errors.New("SpJsonToStruct jsonSpSrc is nil")
	}

	if jsonSpDest == nil{
		return errors.New("SpJsonToStruct jsonSpDest is nil")
	}

	if xStruct == nil{
		return errors.New("SpJsonToStruct stObj is nil")
	}

	strSrc := SpJsonToStr(jsonSpSrc)
	err := json.Unmarshal([]byte(strSrc),xStruct)
	if err != nil{
		return err
	}
	err = SpJsonFromStructByFieldName(jsonSpDest,xStruct)
	if err != nil{
		return err
	}

	return nil
}

func SpJsonStructFieldToTag(jsonSpSrc *simplejson.Json, xStruct interface{}, jsonSpDest *simplejson.Json) error{
	err := SpJsonToStructByFieldName(jsonSpDest,xStruct)
	if err != nil{
		return err
	}
	strSrc := JsonMarshal(xStruct)
	err = jsonSpDest.UnmarshalJSON([]byte(strSrc))
	if err != nil{
		return err
	}

	return nil
}

func SpJsonGetStr(jsonSp *simplejson.Json, strKey string) string{
	if jsonSp == nil{
		return ""
	}

	jsonSub := jsonSp.Get(strKey)
	if jsonSub == nil{
		return ""
	}

	bytes,err := jsonSub.MarshalJSON()
	if err != nil{
		return ""
	}

	strJson := string(bytes)
	if strJson == "null"{
		return ""
	}
	strJson = StrRemoveStrEnds(strJson,"\"")

	return strJson
}

func SpJsonGetIndexStr(jsonSp *simplejson.Json, iIndex int) string{
	if jsonSp == nil || iIndex < 0{
		return ""
	}

	jsonSub := jsonSp.GetIndex(iIndex)
	if jsonSub == nil{
		return ""
	}

	bytes,err := jsonSub.MarshalJSON()
	if err != nil{
		return ""
	}

	strJson := string(bytes)
	if strJson == "null"{
		return ""
	}

	strJson = StrRemoveStrEnds(strJson,"\"")

	return strJson
}

func SpJsonContains(jsonSp *simplejson.Json, strKey string) bool{
	return SpJsonMapContains(jsonSp,strKey)
}

func SpJsonGetInt(jsonSp *simplejson.Json, strKey string) int{
	iResult := CastStrToInt(SpJsonGetStr(jsonSp,strKey))
	return iResult
}

func SpJsonGetBool(jsonSp *simplejson.Json, strKey string) bool{
	bResult := CastStrToBool(SpJsonGetStr(jsonSp,strKey))
	return bResult
}

func SpJsonGetI64(jsonSp *simplejson.Json, strKey string) int64{
	iResult := CastStrToI64(SpJsonGetStr(jsonSp,strKey))
	return iResult
}

func SpJsonGetFloat(jsonSp *simplejson.Json, strKey string) float64{
	fResult := CastStrToF64(SpJsonGetStr(jsonSp,strKey))
	return fResult
}

func SpJsonGetChild(jsonSp *simplejson.Json, strKey string) *simplejson.Json{
	if jsonSp == nil{
		return nil
	}

	jsonChild := jsonSp.Get(strKey)

	return jsonChild
}

func SpJsonGetIndexInt(jsonSp *simplejson.Json, iIndex int) int{
	iResult := CastStrToInt(SpJsonGetIndexStr(jsonSp,iIndex))
	return iResult
}

func SpJsonGetIndexBool(jsonSp *simplejson.Json, iIndex int) bool{
	bResult := CastStrToBool(SpJsonGetIndexStr(jsonSp,iIndex))
	return bResult
}

func SpJsonGetIndexI64(jsonSp *simplejson.Json, iIndex int) int64{
	iResult := CastStrToI64(SpJsonGetIndexStr(jsonSp,iIndex))
	return iResult
}

func SpJsonGetIndexFloat(jsonSp *simplejson.Json, iIndex int) float64{
	fResult := CastStrToF64(SpJsonGetIndexStr(jsonSp,iIndex))
	return fResult
}

func SpJsonGetIndexChild(jsonSp *simplejson.Json, iIndex int) *simplejson.Json{
	if jsonSp == nil{
		return nil
	}

	jsonChild := jsonSp.GetIndex(iIndex)

	return jsonChild
}

func SpJsonFromStructInGroup(groupStruct interface{}) *simplejson.Json {
	return SpJsonFromStructInGroupRecur(groupStruct).(*simplejson.Json)
}

func SpJsonFromStructInGroupRecur(groupStruct interface{}) interface{} {
	jsonResult := simplejson.New()

	valueSrc := reflect.ValueOf(groupStruct)
	typeSrc := valueSrc.Type()
	kindSrc := typeSrc.Kind()

	if kindSrc == reflect.Ptr{
		valueSrc = valueSrc.Elem()
		typeSrc = valueSrc.Type()
		kindSrc = typeSrc.Kind()
	}

	if kindSrc == reflect.Struct{
		SpJsonFromStruct(jsonResult,groupStruct)
	}else if kindSrc == reflect.Map{
		for _, key := range valueSrc.MapKeys() {
			jsonResult.Set( CastVoidToStr(key.Interface()) ,SpJsonFromStructInGroupRecur( valueSrc.MapIndex(key).Interface() ))
		}
	}else if kindSrc == reflect.Array || kindSrc == reflect.Slice{
		arrObj := make([]interface{},0)
		for i := 0; i < valueSrc.Len(); i++ {
			xTemp := SpJsonFromStructInGroupRecur(valueSrc.Index(i).Interface())
			arrObj = append(arrObj,xTemp)
		}
		jsonTemp := simplejson.New()
		jsonTemp.Set("temp_arr",arrObj)
		return jsonTemp.Get("temp_arr")
	}else{
		return CastVoidToStr(valueSrc.Interface())
	}

	return jsonResult
}

func SpJsonFromStructInGroupToStr(groupStruct interface{}) string {
	return SpJsonToStr(SpJsonFromStructInGroup(groupStruct))
}

func SpJsonStructFromStr(strJson string,xObj interface{}) error{
	if strJson == ""{
		return errors.New("SpJsonStructFromStr empty json")
	}

	err :=json.Unmarshal([]byte(strJson),xObj)
	if err != nil{
		return err
	}

	return nil
}

func SpJsonStructToStr(xObj interface{}) string{
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

func SpJsonFromStr(strJson string) *simplejson.Json{
	if strJson == ""{
		return nil
	}

	if !strings.HasPrefix(strJson,"[") && !strings.HasPrefix(strJson,"{"){
		return nil
	}

	if !strings.HasSuffix(strJson,"]") && !strings.HasSuffix(strJson,"}"){
		return nil
	}

	jsonNew , err := simplejson.NewJson([]byte(strJson))
	if err != nil{
		return nil
	}

	return jsonNew
}

func SpJsonFromBuf(bufJson []byte) *simplejson.Json{
	jsonNew , err := simplejson.NewJson(bufJson)
	if err != nil{
		return nil
	}
	return jsonNew
}

func SpJsonPathGetStr(jsonSp *simplejson.Json, strPath string) string {
	if jsonSp == nil || strPath == ""{
		return ""
	}

	var jsonParent *simplejson.Json = jsonSp
	strResult := ""
	strPathNew := strPath

	for {
		strKey := ""
		if strPathNew == ""{
			break
		}

		if strPathNew[0] == '"'{
			strPathNew = strPathNew[1:]
			iQuoteSlashPos := strings.Index(strPathNew,"\"/")
			if iQuoteSlashPos < 0{
				strKey = strPathNew[0:len(strPathNew) -1]
				strPathNew = ""
			}else{
				strKey = strPathNew[0:iQuoteSlashPos]
				strPathNew = strPathNew[iQuoteSlashPos + 2:]
			}
		}else{
			iSlashPos := strings.Index(strPathNew,"/")
			if iSlashPos < 0{
				strKey = strPathNew[0:]
				strPathNew = ""
			}else{
				strKey = strPathNew[0:iSlashPos]
				strPathNew = strPathNew[iSlashPos + 1:]
			}
		}

		if strKey == ""{
			break
		}

		if strings.HasPrefix(strKey,"[") && strings.HasSuffix(strKey,"]"){
			strNum := strKey[1: len(strKey) - 1 ]
			iNum := CastStrToInt(strNum)
			jsonParent = jsonParent.GetIndex(iNum)
		}else{
			if SpJsonContains(jsonParent,strKey){
				jsonParent = jsonParent.Get(strKey)
			}else{
				return ""
			}

		}

		if strPathNew == ""{
			strResult = SpJsonToStr(jsonParent)
			break
		}

		if jsonParent == nil{
			return ""
		}
	}

	return strResult
}

func SpJsonPathSetValue(jsonSp *simplejson.Json, strPath string, xValue interface{}) error {
	if jsonSp == nil{
		return errors.New("json is nil")
	}

	if strPath == ""{
		return errors.New("jsonpath is empty")
	}

	var jsonParent *simplejson.Json = jsonSp
	var jsonParentLast *simplejson.Json = nil
	strPathNew := strPath
	strLastKey := ""
	isLastIndex := false

	iLevel := 0
	for {
		iLevel++
		strKey := ""
		isIndex := false
		if strPathNew == ""{
			break
		}

		if strPathNew[0] == '"'{
			strPathNew = strPathNew[1:]
			iQuoteSlashPos := strings.Index(strPathNew,"\"/")
			if iQuoteSlashPos < 0{
				strKey = strPathNew[0:len(strPathNew) -1]
				strPathNew = ""
			}else{
				strKey = strPathNew[0:iQuoteSlashPos]
				strPathNew = strPathNew[iQuoteSlashPos + 2:]
			}
		}else{
			iSlashPos := strings.Index(strPathNew,"/")
			if iSlashPos < 0{
				strKey = strPathNew[0:]
				strPathNew = ""
			}else{
				strKey = strPathNew[0:iSlashPos]
				strPathNew = strPathNew[iSlashPos + 1:]
			}
		}

		if strKey == ""{
			return errors.New("json path error")
		}

		if strings.HasPrefix(strKey,"[") && strings.HasSuffix(strKey,"]"){
			isIndex = true
			strKey = strKey[1: len(strKey) - 1 ]
			iNum := CastStrToInt(strKey)
			if strPathNew == ""{
				if iLevel == 1{
					return errors.New("can not change top level value by index")
				}else{
					if isLastIndex{
						arrValues := jsonParentLast.GetIndex(CastStrToInt(strLastKey)).MustArray()
						if iNum > len(arrValues){
							return errors.New("array index error")
						}
						arrValues[iNum] = xValue
					}else{
						arrValues := jsonParentLast.Get(strLastKey).MustArray()
						if iNum > len(arrValues){
							return errors.New("array index error")
						}
						arrValues[iNum] = xValue
					}
				}
			}else{
				jsonParentLast = jsonParent
				jsonParent = jsonParent.GetIndex(iNum)
			}
		}else{
			if strPathNew == ""{
				jsonParent.Set(strKey,xValue)
				break
			}else{
				jsonParentLast = jsonParent
				if SpJsonContains(jsonParent,strKey){
					jsonParent = jsonParent.Get(strKey)
				}else{
					return errors.New("no this key in json map:" + strKey)
				}

			}
		}

		strLastKey = strKey
		isLastIndex = isIndex

		if strPathNew == ""{
			break
		}

		if jsonParent == nil{
			break
		}
	}

	return nil
}