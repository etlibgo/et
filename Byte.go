package et

import (
	"errors"
)

/**
关于字节流的操作。
主要提供，对比，截取，组合的操作。
 */

/**
两段字节流是否相等
 */
func ByteEqual(bufLeft []byte, iLeftOffset int, bufRight []byte, iRightOffset int, iSize int) bool{
	if bufLeft == nil && bufRight == nil{
		return true
	}

	if bufLeft == nil || bufRight == nil{
		return false
	}

	if len(bufLeft) < iLeftOffset + iSize{
		return false
	}

	if len(bufRight) < iRightOffset + iSize{
		return false
	}

	for i := 0; i < iSize; i++{
		if bufLeft[iLeftOffset + i] != bufRight[iRightOffset + i]{
			return false
		}
	}

	return true
}


/**
截取字节流中的一部分
 */
func ByteSubBytes(bufSrc []byte, iOffset int, iSize int) []byte{
	if bufSrc == nil{
		return nil
	}

	if len(bufSrc) < iOffset + iSize{
		return nil
	}

	bufDest := bufSrc[iOffset:iOffset+iSize]

	return bufDest
}

/**
把一段字节流中的一部分拷贝到另一个字节流
 */
func ByteCopy(bufSrc []byte, iSrcOffset int, bufDest []byte, iDestOffset int, iSize int) int{
	if bufSrc == nil || bufDest == nil{
		return -1
	}

	if len(bufSrc) < iSrcOffset + iSize{
		return -1
	}

	if len(bufDest) < iDestOffset + iSize{
		return -1
	}

	iResult := copy(bufDest[iDestOffset:],bufSrc[iSrcOffset:iSrcOffset + iSize])

	return iResult
}

/**
查找字节流中是否包含某一部分
 */
func ByteFindFrom(bufSrc []byte, iSrcOffset int, bufFind []byte) int{
	if bufSrc == nil{
		return -1
	}

	if bufFind == nil{
		return -1
	}

	iLenSrc := len(bufSrc)
	if iLenSrc <= iSrcOffset{
		return -1
	}

	iLenFind := len(bufFind)
	if iLenFind == 0{
		return iSrcOffset
	}

	if iLenSrc < iSrcOffset + iLenFind{
		return -1
	}

	for i := iSrcOffset; i < iLenSrc - iLenFind + 1; i++{
		bFound := true
		for j := 0; j < iLenFind; j++{

			if bufSrc[i + j] != bufFind[j]{
				bFound = false
				break
			}
		}
		if bFound{
			return i
		}
	}

	return -1
}

/**
把多个字节流组合成一个
 */
func ByteCombineList(arrBuf [][]byte) []byte{
	if arrBuf == nil{
		return nil
	}

	iTotalSize := 0
	for i := 0; i < len(arrBuf); i++{
		iTotalSize += len(arrBuf[i])
	}
	bufDest := make([]byte,iTotalSize)
	iCurrentSize := 0

	for i := 0; i < len(arrBuf); i++{
		if ByteIsEmpty(arrBuf[i]){
			continue
		}
		iTempSize := ByteCopy(arrBuf[i],0,bufDest,iCurrentSize,len(arrBuf[i]))
		if iTempSize >= 0{
			iCurrentSize += iTempSize
		}
		iTotalSize += len(arrBuf[i])
	}

	return bufDest
}

/**
判断一段字节流是否为空
 */
func ByteIsEmpty(bufSrc []byte)bool{
	return len(bufSrc) == 0
}

/*
批量设置字节流中的值，
相当于memset
 */
func ByteSet(bufSrc []byte, iOffset int, iSize int, byteToSet uint8) error{
	if bufSrc == nil{
		return errors.New("buf nil")
	}

	if len(bufSrc) < iOffset + iSize{
		return errors.New("buf size error")
	}

	for i := 0; i < iSize; i++{
		bufSrc[i + iOffset] = byteToSet
	}

	return nil
}