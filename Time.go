package et

import (
	"strings"
	"time"
)

/**
常用的时间操作
 */

/**
本工具包内时间的各个字段的字母标识
 */
const TIME_FORMAT_PREFIX string = "%"
const TIME_FORMAT_YEAR string = "y"
const TIME_FORMAT_MONTH string = "m"
const TIME_FORMAT_DAY string = "d"
const TIME_FORMAT_HOUR string = "h"
const TIME_FORMAT_MINUTE string = "i"
const TIME_FORMAT_SECOND string = "s"
const TIME_FORMAT_MILLI string = "l"
const TIME_FORMAT_WEEK string = "w"


/**
把时间格式化为一个map，键名就是各个字母标识，定义在常量中
 */
func TimeFormatObjToMap(objTime time.Time) map[string]int{

	mapResult := map[string]int{}

	mapResult[TIME_FORMAT_YEAR] = objTime.Year()
	mapResult[TIME_FORMAT_MONTH] = int(objTime.Month())
	mapResult[TIME_FORMAT_DAY] = objTime.Day()
	mapResult[TIME_FORMAT_HOUR] = objTime.Hour()
	mapResult[TIME_FORMAT_MINUTE] = objTime.Minute()
	mapResult[TIME_FORMAT_SECOND] = objTime.Second()
	mapResult[TIME_FORMAT_MILLI] = int( TimeGetStamp13(objTime)  % 1000)
	mapResult[TIME_FORMAT_WEEK] = int(objTime.Weekday())

	return mapResult
}

/**
把13位时间戳转换为map
 */
func TimeFormatStamp13ToMap(iStamp13 int64)  map[string]int{
	return TimeFormatObjToMap(TimeFromStamp13(iStamp13))
}

/**
把时间对象格式化为字符串，格式中的%x 会被替换为实际的时间值，比如%y是年,%m是月
 */
func TimeFormatObjToStr(strFormat string, objTime time.Time) string{
	if strFormat == ""{
		return ""
	}
	mapResult := TimeFormatObjToMap(objTime)
	strResult := strFormat

	strMonth := StrPadLeft( CastIntToStr(mapResult[TIME_FORMAT_MONTH]) ,"0",2)
	strDate := StrPadLeft( CastIntToStr(mapResult[TIME_FORMAT_DAY]),"0",2)
	strHour := StrPadLeft( CastIntToStr(mapResult[TIME_FORMAT_HOUR]),"0",2)
	strMinute := StrPadLeft( CastIntToStr(mapResult[TIME_FORMAT_MINUTE]),"0",2)
	strSecond := StrPadLeft( CastIntToStr(mapResult[TIME_FORMAT_SECOND]),"0",2)
	strMilli := StrPadLeft( CastIntToStr(mapResult[TIME_FORMAT_MILLI]),"0",3)

	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_YEAR, CastIntToStr(mapResult[TIME_FORMAT_YEAR]) )
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_MONTH, strMonth )
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_DAY, strDate )
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_HOUR, strHour )
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_MINUTE, strMinute)
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_SECOND, strSecond )
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_MILLI, strMilli )
	strResult = strings.ReplaceAll(strResult, TIME_FORMAT_PREFIX + TIME_FORMAT_WEEK, CastIntToStr(mapResult[TIME_FORMAT_WEEK]) )

	return strResult
}

/**
根据13为时间戳格式化为字符串
 */
func TimeFormatStamp13ToStr(strFormat string, iStamp13 int64) string{
	return TimeFormatObjToStr(strFormat, TimeFromStamp13(iStamp13))
}

/**
把当前的时间格式化为字符串
 */
func TimeFormatNowToStr(strFormat string) string{
	return TimeFormatObjToStr(strFormat, time.Now())
}

/**
一个时间相关的map中是否包含某个键
 */
func TimeFormattedMapIsContained(mapSmaller map[string]int, mapBigger map[string]int) bool{
	if len(mapSmaller) == 0 || len(mapBigger) == 0{
		return false
	}

	for strFormat,iValue := range mapSmaller{
		if iValue != mapBigger[strFormat]{
			return false
		}
	}

	return true
}

/**
由13位时间戳得到一个时间对象
 */
func TimeFromStamp13(iStamp13 int64) time.Time{
	return time.Unix(int64( iStamp13 / 1000), (iStamp13 % 1000) * 1000000 )
}

/**
获取指定时间对象的13位时间戳
 */
func TimeGetStamp13(xTime time.Time) int64{
	iStamp13 := int64(xTime.UnixNano() / 1000000)
	return iStamp13
}

/**
得到当前毫秒数
 */
func TimeGetMilli(xTime time.Time) int{
	iMilliSecond := int(xTime.UnixNano() / 1000000) % 1000
	return iMilliSecond
}

/**
得到当前的13位时间戳
 */
func TimeNowStamp13() int64{
	return TimeGetStamp13(time.Now())
}

/**
得到当前的10位时间戳
*/
func TimeNowStamp10() int64{
	return time.Now().Unix()
}

