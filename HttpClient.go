package et

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

/**
http客户端的一些实用工具函数
 */

/**
用get得到字符串
 */
func HttpClientGetStr(strUrl string, iTimeoutMilli int, mapHeader map[string]string) string{
	buf,err := HttpClientGetBuf(strUrl,iTimeoutMilli,mapHeader)
	if err != nil{
		return ""
	}else{
		return CastBufToStr(buf)
	}
}

/**
用get得到字节流，一般是下载
 */
func HttpClientGetBuf(strUrl string, iTimeoutMilli int, mapHeader map[string]string)([]byte, error){
	if strUrl == ""{
		return nil,errors.New("url empty")
	}

	if iTimeoutMilli < 0{
		iTimeoutMilli = 0
	}

	strMethod := "GET"
	pClient := &http.Client{}
	pClient.Timeout = time.Millisecond * time.Duration(iTimeoutMilli)

	pNewRequest,err := http.NewRequest(strMethod,strUrl,nil)
	if err != nil {
		//fmt.Println("new request err")
		return nil,errors.New("http request new err")
	}

	//fmt.Println("HttpClientGetBuf 222")
	for strKey,strVal := range mapHeader {
		pNewRequest.Header.Add(strKey,strVal)
	}

	pNewResponse,err := pClient.Do(pNewRequest)
	if err != nil {
		return nil,errors.New("response err")
	}
	bufData,err := ioutil.ReadAll(pNewResponse.Body)
	defer pNewResponse.Body.Close()

	if err != nil {
		//fmt.Println("read response err")
		return nil,errors.New("read response err")
	}else{
		//fmt.Println("HttpClientGetBuf 555")
		return bufData,nil
	}
}

/**
用post得到字符串
 */
func HttpClientPostStr(strUrl string, strPostData string, iTimeoutMilli int, mapHeader map[string]string)string {
	bufResult,err := HttpClientPostBuf(strUrl,[]byte(strPostData),iTimeoutMilli,mapHeader)
	if err != nil{
		return ""
	}else{
		return string(bufResult)
	}
}

/**
用post得到字节流
 */
func HttpClientPostBuf(strUrl string, bufRequest []byte, iTimeoutMilli int, mapHeader map[string]string)([]byte, error){
	if strUrl == ""{
		return nil,errors.New("url empty")
	}

	if iTimeoutMilli < 0{
		iTimeoutMilli = 0
	}

	strMethod := "POST"
	pClient := &http.Client{}
	pClient.Timeout = time.Millisecond * time.Duration(iTimeoutMilli)


	pNewRequest,err := http.NewRequest(strMethod,strUrl,bytes.NewReader(bufRequest))
	if err != nil {
		//fmt.Println("new request err")
		return nil,errors.New("http request new err")
	}

	for strKey,strVal := range mapHeader {
		pNewRequest.Header.Add(strKey,strVal)
	}

	pNewResponse,err := pClient.Do(pNewRequest)
	if err != nil {
		return nil,err
	}


	bufResult,err := ioutil.ReadAll(pNewResponse.Body)
	defer pNewResponse.Body.Close()

	if err != nil {
		return nil,errors.New("read response err")
	}else{
		return bufResult,nil
	}
}

/**
用get得到一个response响应对象
 */
func HttpClientGetResp(strUrl string, iTimeoutMilli int, mapHeader map[string]string)(*http.Response, error){

	if strUrl == ""{
		return nil,errors.New("url empty")
	}

	if iTimeoutMilli < 0{
		iTimeoutMilli = 0
	}
	strMethod := "GET"
	pClient := &http.Client{}
	pClient.Timeout = time.Millisecond * time.Duration(iTimeoutMilli)


	pNewRequest,err := http.NewRequest(strMethod,strUrl,nil)
	if err != nil {
		return nil,errors.New("http request new err")
	}

	for strKey,strVal := range mapHeader {
		pNewRequest.Header.Add(strKey,strVal)
	}

	pNewResponse,err := pClient.Do(pNewRequest)
	if err != nil {
		return nil,errors.New("response err")
	}

	return pNewResponse,nil
}

/**
用get得到一个响应对象，需要设置是否跟随302，病情是否需要转发原始头
 */
func HttpClientGetRespWithRedirect(strUrl string, iTimeoutMilli int, mapHeader map[string]string, toRedirect bool, redirectHeader bool)(*http.Response, int, error){

	if strUrl == ""{
		return nil,0,errors.New("url empty")
	}

	if iTimeoutMilli < 0{
		iTimeoutMilli = 0
	}
	strMethod := "GET"
	pClient := &http.Client{}
	pClient.Timeout = time.Millisecond * time.Duration(iTimeoutMilli)

	iRedirectCount := 0

	pClient.CheckRedirect = func(pReq *http.Request, via []*http.Request)error{
		iRedirectCount += 1
		if toRedirect{
			if redirectHeader{
				for strKey,strVal := range mapHeader {
					pReq.Header.Add(strKey,strVal)
				}
			}
			return nil
		}else{
			return http.ErrUseLastResponse
		}
	}

	pNewRequest,err := http.NewRequest(strMethod,strUrl,nil)
	if err != nil {
		return nil,iRedirectCount,errors.New("http request new err")
	}

	for strKey,strVal := range mapHeader {
		pNewRequest.Header.Add(strKey,strVal)
	}

	pNewResponse,err := pClient.Do(pNewRequest)
	if err != nil {
		return nil,iRedirectCount,errors.New("response err")
	}

	return pNewResponse,iRedirectCount,nil
}


