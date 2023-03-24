package et

import "encoding/base64"

/**
base64的编解码
使用的时候，请注意是否要用于url。
若要用于url，请使用带ForUrl后缀的版本
 */

func Base64EncBufToStr(bufSrc []byte) string{
	if len(bufSrc) == 0{
		return ""
	}

	strResult := base64.StdEncoding.EncodeToString(bufSrc)

	return strResult
}

func Base64DecStrToBuf(strSrc string) []byte{
	if strSrc == ""{
		return nil
	}

	bufResult,err := base64.StdEncoding.DecodeString(strSrc)
	if err != nil{
		bufResult = nil
	}

	return bufResult
}


func Base64EncStrToStr(strSrc string) string{
	return Base64EncBufToStr([]byte(strSrc))
}

func Base64DecStrToStr(strSrc string) string{
	return string(Base64DecStrToBuf(strSrc))
}


func Base64EncBufToStrForUrl(bufSrc []byte) string{
	if len(bufSrc) == 0{
		return ""
	}

	strResult := base64.URLEncoding.EncodeToString(bufSrc)

	return strResult
}

func Base64DecStrToBufForUrl(strSrc string) []byte{
	if strSrc == ""{
		return nil
	}

	bufResult,err := base64.URLEncoding.DecodeString(strSrc)
	if err != nil{
		bufResult = nil
	}

	return bufResult
}

func Base64EncStrToStrForUrl(strSrc string) string{
	return Base64EncBufToStrForUrl([]byte(strSrc))
}

func Base64DecStrToStrForUrl(strSrc string) string{
	return string(Base64DecStrToBufForUrl(strSrc))
}