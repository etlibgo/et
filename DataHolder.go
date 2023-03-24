package et

import (
	"fmt"
	"strings"
)

/**
多功能数据存储类
类似于一个tuple
主要用于框架内部通信
业务代码一般不直接操作
 */

const DATA_HOLDER_TYPE_NONE string = "empty"
const DATA_HOLDER_TYPE_BOOL string = "bool"
const DATA_HOLDER_TYPE_INT string = "int"
const DATA_HOLDER_TYPE_LONG string = "int64"
const DATA_HOLDER_TYPE_STR string = "str"
const DATA_HOLDER_TYPE_OBJ string = "obj"
const DATA_HOLDER_TYPE_FLOAT string = "float"
const DATA_HOLDER_TYPE_BUF string = "buf"
const DATA_HOLDER_TYPE_MIXED string = "mixed"


type DataHolder struct{
	bBoolData bool
	iIntData int
	iLongData int64
	strStrData string
	xObjData interface{}
	fFloatData  float64
	bufData []byte
	strType string
}

func DataHolderNew() *DataHolder{
	xDataHolder := DataHolder{}
	xDataHolder.bBoolData = false
	xDataHolder.iIntData = 0
	xDataHolder.iLongData = 0
	xDataHolder.strStrData = ""
	xDataHolder.xObjData = nil
	xDataHolder.fFloatData = 0
	xDataHolder.bufData = nil
	xDataHolder.strType = DATA_HOLDER_TYPE_NONE

	return &xDataHolder
}

func DataHolderNewStr(strData string) *DataHolder{
	xDataHolder := DataHolder{}
	xDataHolder.bBoolData = false
	xDataHolder.iIntData = 0
	xDataHolder.iLongData = 0
	xDataHolder.strStrData = strData
	xDataHolder.xObjData = nil
	xDataHolder.fFloatData = 0
	xDataHolder.bufData = nil
	xDataHolder.strType = DATA_HOLDER_TYPE_STR

	return &xDataHolder
}

func (my *DataHolder)GetType() string{
	return my.strType
}

func (my *DataHolder)SetBool(bData bool){
	my.bBoolData = bData
}

func (my *DataHolder)SetInt(iData int){
	my.iIntData = iData
}

func (my *DataHolder)SetLong(iLong int64){
	my.iLongData = iLong
}

func (my *DataHolder)SetFloat(fData float64){
	my.fFloatData = fData
}

func (my *DataHolder)SetStr(strData string){
	my.strStrData = strData
}

func (my *DataHolder)SetObj(xData interface{}){
	my.xObjData = xData
}

func (my *DataHolder)SetBuf(bufData []byte){
	my.bufData = bufData
}

func (my *DataHolder)GetBool() bool{
	return my.bBoolData
}

func (my *DataHolder)GetInt() int{
	return my.iIntData
}

func (my *DataHolder)GetLong() int64{
	return my.iLongData
}

func (my *DataHolder)GetFloat() float64{
	return my.fFloatData
}

func (my *DataHolder)GetStr() string{
	return my.strStrData
}

func (my *DataHolder)GetObj() interface{}{
	return my.xObjData
}

func (my *DataHolder)GetBuf() []byte{
	return my.bufData
}

func (my *DataHolder)ToJsonStr() string{
	strFormat := "{\"type\":\"$(type)\",\"data\":\"$(data)\"}"
	return my.ToStrByFormat(strFormat)
}

func (my *DataHolder)ToUrlStr() string{
	strFormat := "type=$(type)&data=$(data)"
	return my.ToStrByFormat(strFormat)
}

func (my *DataHolder)ToStrByFormat(strFormat string) string{
	strType := my.strType
	strData := ""

	if my.strType == DATA_HOLDER_TYPE_STR{
		strData = my.strStrData
	}else if my.strType == DATA_HOLDER_TYPE_INT{
		strData = CastIntToStr(my.iIntData)
	}else if my.strType == DATA_HOLDER_TYPE_LONG{
		strData = CastI64ToStr(my.iLongData)
	}else if my.strType == DATA_HOLDER_TYPE_FLOAT{
		strData = CastF64ToStr(my.fFloatData,5)
	}else if my.strType == DATA_HOLDER_TYPE_OBJ{
		strData = fmt.Sprintf("%s",my.xObjData)
	}else if my.strType == DATA_HOLDER_TYPE_BOOL{
		strData = CastBoolToStr(my.bBoolData)
	}else if my.strType == DATA_HOLDER_TYPE_BUF{
		strData = CastBufToStr(my.bufData)
	}else if my.strType == DATA_HOLDER_TYPE_NONE{
		strData = ""
	}else{
		strData = fmt.Sprintf("%s",my.xObjData)
	}

	//strFormat := "type=$(type)&data=$(data)"
	strResult := strFormat
	strResult = strings.ReplaceAll(strResult,"$(type)",strType)
	strResult = strings.ReplaceAll(strResult,"$(data)",strData)

	return strResult
}
