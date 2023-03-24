package et

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

/**
字符串的基本操作
 */

const STR_HEX_CHAR_LOWER string = "0123456789abcdef"
const STR_HEX_CHAR_UPPER string = "0123456789ABCDEF"
const STR_BLANK_CH string = " \r\n\t"

/**
只要strSrc中出现strChars中的任意一个字符，则返回出现的位置
 */
func StrFindChars(strSrc string, strChars string, iFromPos int) int{
	if strSrc == "" || strChars == ""{
		return -1
	}
	if iFromPos < 0{
		iFromPos = 0
	}
	iResult := -1
	iLenSrc := len(strSrc)
	if iFromPos >= iLenSrc{
		return COMMON_NOT_FOUND
	}
	iLenChars := len(strChars)
	for iSrcPos := iFromPos; iSrcPos < iLenSrc; iSrcPos++{
		chSrc := strSrc[iSrcPos]
		for iCharsPos := 0; iCharsPos < iLenChars; iCharsPos++{
			if chSrc == strChars[iCharsPos]{
				iResult = iSrcPos
				break
			}
		}
		if iResult != -1{
			break
		}
	}

	return iResult
}

/**
删除两端的字符串。比如两端的引号。
 */
func StrRemoveStrEnds(strSrc string, strOnEnds string) string{
	if strSrc == "" || strOnEnds == ""{
		return strSrc
	}

	if strings.HasPrefix(strSrc,strOnEnds) && strings.HasSuffix(strSrc,strOnEnds){
		iLenSrc := len(strSrc)
		iLenEnds := len(strOnEnds)

		strResult := strSrc[iLenEnds:iLenSrc-iLenEnds]
		return strResult
	}

	return strSrc
}

/**
encodeuri,转义url字符串
 */
func StrEncodeUriComponent(strSrc string) string{
	strResult := url.QueryEscape(strSrc)
	strResult = strings.ReplaceAll(strResult,"+","%20")
	return strResult
}

/**
decodeuri,反转义url字符串
*/
func StrDecodeUriComponent(strSrc string) string{
	strResult,err := url.QueryUnescape(strSrc)
	if err != nil{
		strResult = strSrc
	}
	return strResult
}

/**
在字符串的左边补齐字符
 */
func StrPadLeft(strSrc string, strPad string, iTotalCount int) string{
	lenSrc := len(strSrc)

	if lenSrc >= iTotalCount{
		return strSrc
	}

	strResult := strings.Repeat(strPad,iTotalCount - lenSrc) + strSrc
	return strResult
}

/**
在字符串右边补齐字符
 */
func StrPadRight(strSrc string, strPad string, iTotalCount int) string{
	lenSrc := len(strSrc)

	if lenSrc >= iTotalCount{
		return strSrc
	}

	strResult := strSrc + strings.Repeat(strPad,iTotalCount - lenSrc)
	return strResult
}

/**
判断字符串的内容是否都限于某些字符
 */
func StrIsWithIn(strSrc string, strCollection string) bool{
	if strSrc == ""{
		return true
	}

	if strCollection == ""{
		return false
	}

	srcLen := len(strSrc)
	for i := 0; i < srcLen; i++{
		if !strings.Contains(strCollection,strSrc[i:i+1] ){
			return false
		}
	}

	return true
}

/**
是否为整型
 */
func StrIsInt(strSrc string) bool{
	const MaxLenOfInt64 int = 20
	if strSrc == "" || len(strSrc) > MaxLenOfInt64{
		return false
	}
	match, _ := regexp.MatchString(`^[\+-]?\d+$`, strSrc)
	return match
}

/**
是否为数字
 */
func StrIsNum(strSrc string) bool {
	// 去除首尾空格
	for i := 0; i < len(strSrc); i++ {
		// 存在 e 或 E, 判断是否为科学计数法
		if strSrc[i] == 'e' || strSrc[i] == 'E' {
			return StrIsSciNum(strSrc[:i], strSrc[i+1:])
		}
	}
	// 否则判断是否为整数或小数
	return StrIsInt(strSrc) || StrIsDec(strSrc)
}

//is it a scientific number?
func StrIsSciNum(strNumLeft, strNumRight string) bool {
	// e 前后字符串长度为0 是错误的
	if len(strNumLeft) == 0 || len(strNumRight) == 0 {
		return false
	}
	// e 后面必须是整数，前面可以是整数或小数  4  +
	return (StrIsInt(strNumLeft) || StrIsDec(strNumLeft)) && StrIsInt(strNumRight)
}

//is it a decimal number
func StrIsDec(strSrc string) bool {
	// eg: 11.15, -0.15, +10.15, 3., .15,
	// err: +. 0..
	match1, _ := regexp.MatchString(`^[\+-]?\d*\.\d+$`, strSrc)
	match2, _ := regexp.MatchString(`^[\+-]?\d+\.\d*$`, strSrc)
	return match1 || match2
}

/**
给某字符加入反斜杠
 */
func StrAddSlashes(strSrc string, chSpecial uint8) string{
	if strSrc == ""{
		return ""
	}
	if chSpecial == 0{
		return strSrc
	}
	strResult := ""
	srcLen := len(strSrc)
	for i := 0; i < srcLen; i++{
		chTemp := strSrc[i]
		if chTemp == '\\'{
			strResult = strResult + "\\\\"
		}else if chTemp == chSpecial{
			strResult = strResult + "\\" + string(chTemp)
		}else{
			strResult += string(chTemp)
		}
	}

	return strResult
}

/**
给某字符加入反斜杠后，反复原
*/
func StrStripSlashes(strSrc string) string{
	if strSrc == ""{
		return ""
	}

	if !strings.Contains(strSrc,"\\"){
		return strSrc
	}

	strResult := ""

	srcLen := len(strSrc)

	for i := 0; i < srcLen; i++{
		chTemp := strSrc[i]
		if chTemp == '\\'{
			if i == srcLen - 1{
				break
			}
			strResult += string(strSrc[i + 1])
			i++
		}else{
			strResult += string(chTemp)
		}
	}

	return strResult
}

/**
把16进制转换为字节流
 */
func StrHexToBuf(strHex string) []byte{
	if strHex == ""{
		return []byte{}
	}

	bufLen := len(strHex) / 2
	bufResult := make([]byte,bufLen)

	for i := 0; i < bufLen; i++{
		idx := i * 2
		tempInt,err := strconv.ParseInt(strHex[idx:idx+2],16,32)
		if err != nil{
			bufResult[i] = '0'
		}else{
			bufResult[i] = uint8(tempInt)
		}
	}

	return bufResult
}

/**
把字节流转换为16进制
 */
func StrHexFromBuf(bufSrc []byte) string{
	if len(bufSrc) == 0{
		return ""
	}

	iBufLen := len(bufSrc)
	strResult := ""

	for i := 0; i < iBufLen; i++{
		strResult += StrByteToHex(bufSrc[i])
	}

	return strResult
}

/**
把单个字节转换为16进制字符串
 */
func StrByteToHex(bySrc uint8) string{
	iLeft := bySrc >> 4
	iRight := bySrc & 0x0f
	chLeft := STR_HEX_CHAR_LOWER[iLeft]
	chRight := STR_HEX_CHAR_LOWER[iRight]
	strResult := string(chLeft) + string(chRight)
	return strResult
}

/**
把16进制字符串转换为一个字节
 */
func StrByteFromHex(strHex string) uint8{
	if strHex == ""{
		return 0
	}

	iMaxLen := 2
	if len(strHex) < 2{
		iMaxLen = len(strHex)
	}

	tempInt,err := strconv.ParseInt(strHex[0:iMaxLen],16,32)
	if err != nil{
		tempInt = 0
	}

	byResult := uint8(tempInt)
	return byResult
}

/**
找到第一个非空字符串
 */
func StrFindFirstNonBlank(strSrc string, iFromPos int) int{
	if strSrc == ""{
		return COMMON_NOT_FOUND
	}

	if iFromPos < 0{
		iFromPos = 0
	}

	iSrcLen := len(strSrc)

	if iFromPos >= iSrcLen{
		return COMMON_NOT_FOUND
	}

	iNonBlankPos := COMMON_NOT_FOUND
	for i := iFromPos; i < iSrcLen; i++{
		strTemp := string(strSrc[i])

		if !strings.Contains(STR_BLANK_CH,strTemp){
			iNonBlankPos = i
			break
		}
	}

	return iNonBlankPos
}

/**
找到最后一个非空字符
 */
func StrFindLastNonBlank(strSrc string) int{
	if strSrc == ""{
		return COMMON_NOT_FOUND
	}

	iSrcLen := len(strSrc)
	iNonBlankPos := COMMON_NOT_FOUND
	for i := iSrcLen - 1; i >= 0; i--{
		strTemp := string(strSrc[i])

		if !strings.Contains(STR_BLANK_CH,strTemp){
			iNonBlankPos = i
			break
		}
	}

	return iNonBlankPos
}

/**
找到第一个某字符，前面不是反斜杠的，也就是非专业
 */
func StrFindFirstWithoutSlash(strSrc string, chToFind uint8, iFromPos int) int{
	if strSrc == ""{
		return COMMON_NOT_FOUND
	}

	if iFromPos < 0{
		iFromPos = 0
	}

	iSrcLen := len(strSrc)
	if iFromPos >= iSrcLen{
		return COMMON_NOT_FOUND
	}

	iNoSlashPos := COMMON_NOT_FOUND
	for i := iFromPos; i < iSrcLen; i++{
		chTemp := strSrc[i]
		if chTemp == chToFind{
			if i == iFromPos || strSrc[i-1] != '\\'{
				iNoSlashPos = i
				break
			}
		}
	}

	return iNoSlashPos
}

func StrFindFrom(strSrc string, strFind string, iFromPos int) int{
	iLenFind := len(strFind)
	if iLenFind == 0{
		return 0
	}

	iLenSrc := len(strSrc)
	if iLenSrc == 0{
		return COMMON_NOT_FOUND
	}

	if iFromPos < 0{
		iFromPos = 0
	}

	if iFromPos >= iLenSrc{
		return COMMON_NOT_FOUND
	}

	strNew := strSrc[iFromPos:]
	iIndexNew := strings.Index(strNew,strFind)
	if iIndexNew < 0{
		return COMMON_NOT_FOUND
	}else{
		return iFromPos + iIndexNew
	}
}

/**
把首字母转为大写
 */
func StrCapitalize(strSrc string) string{
	if strSrc == ""{
		return ""
	}

	chFirst := strSrc[0]
	if chFirst >= 97 && chFirst <= 122{
		chFirst = chFirst - 32
		strResult := string(chFirst) + strSrc[1:]
		return strResult
	}else{
		return strSrc
	}
}

/**
把一个字符串按两种切割符号切两轮，形成一个字符map，比如&与=
Xsv的含义是Csv，只是不一定是用逗号
 */
func StrXsvToMap(strXsv string, strSplit1 string, strSplit2 string) map[string]string{
	if strXsv == ""{
		return map[string]string{}
	}
	mapResult := map[string]string{}
	arrSplit1 := strings.Split(strXsv,strSplit1)
	for _,strItem1 := range arrSplit1{
		arrSplit2 := strings.Split(strItem1,strSplit2)
		if len(arrSplit2) < 2{
			continue
		}
		mapResult[arrSplit2[0]] = arrSplit2[1]
	}
	return mapResult
}

/**
把一个字符MAP序列化为xsv,分隔符是split1,和split2
*/
func StrXsvFromMap(mapXsv map[string]string, strSplit1 string, strSplit2 string) string{
	if mapXsv == nil{
		return ""
	}
	strResult := ""
	bFirst := true
	for k1,v1 := range mapXsv{
		if bFirst{
			bFirst = false
		}else{
			strResult += strSplit1
		}

		strResult += k1 + strSplit2 + v1
	}
	return strResult
}



