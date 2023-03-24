package et

import (
	"crypto/md5"
	"encoding/hex"
)

/**
获得小写32位的md5
 */
func Md5ToLower32(strSrc string) string{
	pMd5 := md5.New()
	pMd5.Write([]byte(strSrc))
	strResult := hex.EncodeToString(pMd5.Sum(nil))
	return strResult
}
