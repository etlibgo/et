package et

/**
Map数据结构的一些操作
 */


/**
把一个map str的所有键设置给另一个map str
 */
func MapStrSet(mapSrc map[string]string, mapDest map[string]string) int{
	if mapSrc == nil || mapDest == nil{
		return 0
	}

	iCount := 0
	for k,v := range mapSrc{
		iCount++
		mapDest[k] = v
	}

	return iCount
}

/**
是否包含相应的键
 */
func MapStrContains(mapSrc map[string]string, strKey string) bool{
	if _, bExisted := mapSrc[strKey]; bExisted {
		return true
	}else{
		return false
	}
}

/**
是否包含相应的键
*/
func MapVoidContains(mapSrc map[string]interface{}, strKey string) bool{
	if _, bExisted := mapSrc[strKey]; bExisted {
		return true
	}else{
		return false
	}
}

/**
把map void转换为一个url字符串
 */
func MapVoidToUrlStr(mapVoid map[string]interface{}) string{
	if mapVoid == nil{
		return ""
	}
	strResult := ""
	for k,v := range mapVoid{
		strResult += "&" + k + "=" + CastVoidToStr(v)
	}

	return strResult
}

/**
把map str转换为一个url字符串
*/
func MapStrToUrlStr(mapVoid map[string]string) string{
	if mapVoid == nil{
		return ""
	}
	strResult := ""
	for k,v := range mapVoid{
		strResult += "&" + k + "=" + v
	}

	return strResult
}

/**
把map void 转换为map str
 */
func MapVoidToMapStr(mapVoid map[string]interface{}) map[string]string{
	if mapVoid == nil{
		return map[string]string{}
	}

	mapResult := map[string]string{}
	for k,v := range mapVoid{
		mapResult[k] = CastVoidToStr(v)
	}

	return mapResult
}


/**
克隆一个map str
*/
func MapStrClone(mapSrc map[string]string) map[string]string{
	if mapSrc == nil{
		return mapSrc
	}

	mapDest := make(map[string]string)
	for k,v := range mapSrc{
		mapDest[k] = v
	}

	return mapDest
}


/**
根据键名获得字符串
 */
func MapStrGetStr(mapSrc map[string]string, strKey string) string{
	if mapSrc == nil{
		return ""
	}

	return mapSrc[strKey]
}

/**
根据键名获得整数
*/
func MapStrGetInt(mapSrc map[string]string, strKey string) int{
	return CastStrToInt(MapStrGetStr(mapSrc,strKey))
}

/**
根据键名获得长整型
*/
func MapStrGetLong(mapSrc map[string]string, strKey string) int64{
	return CastStrToI64(MapStrGetStr(mapSrc,strKey))
}

/**
根据键名获得布尔型
*/
func MapStrGetBool(mapSrc map[string]string, strKey string) bool{
	return CastStrToBool(MapStrGetStr(mapSrc,strKey))
}

/**
根据键名获得浮点型
*/
func MapStrGetFloat(mapSrc map[string]string, strKey string) float64{
	return CastStrToF64(MapStrGetStr(mapSrc,strKey))
}

/**
根据键名获得interface
*/
func MapVoidGetVoid(mapSrc map[string]interface{}, strKey string) interface{}{
	if mapSrc == nil{
		return nil
	}

	return mapSrc[strKey]
}

/**
根据键名获得字符串
*/
func MapVoidGetStr(mapSrc map[string]interface{}, strKey string) string{
	if mapSrc == nil{
		return ""
	}

	return CastVoidToStr(mapSrc[strKey])
}

/**
根据键名获得数值
*/
func MapVoidGetInt(mapSrc map[string]interface{}, strKey string) int{
	return CastStrToInt(MapVoidGetStr(mapSrc,strKey))
}

/**
根据键名获得长整型
*/
func MapVoidGetLong(mapSrc map[string]interface{}, strKey string) int64{
	return CastStrToI64(MapVoidGetStr(mapSrc,strKey))
}

/**
根据键名获得布尔型
*/
func MapVoidGetBool(mapSrc map[string]interface{}, strKey string) bool{
	return CastStrToBool(MapVoidGetStr(mapSrc,strKey))
}

/**
根据键名获得浮点型
*/
func MapVoidGetFloat(mapSrc map[string]interface{}, strKey string) float64{
	return CastStrToF64(MapVoidGetStr(mapSrc,strKey))
}



