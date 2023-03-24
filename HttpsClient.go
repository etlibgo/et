package et

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

/**
https操作类，一般不需要用，因为go http内置工具会在底层判断https。
 */
func HttpsClientGetStr(strUrl string, iTimeoutMilli int, mapHeader map[string]string) string{
	buf,err := HttpsClientGetBuf(strUrl,iTimeoutMilli,mapHeader)
	if err != nil{
		return ""
	}else{
		return CastBufToStr(buf)
	}
}

func HttpsClientGetBuf(strUrl string, iTimeoutMilli int, mapHeader map[string]string)([]byte, error){
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

	//fmt.Println("HttpClientGetBuf 222")
	for strKey,strVal := range mapHeader {
		pNewRequest.Header.Add(strKey,strVal)
	}

	pTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Millisecond  * time.Duration(iTimeoutMilli/2),
		}).Dial,
		MaxIdleConnsPerHost:   0,
		ResponseHeaderTimeout: time.Millisecond  * time.Duration(iTimeoutMilli/2),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	pClient.Transport = pTransport

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