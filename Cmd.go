package et

import (
	"errors"
	"os/exec"
	"time"
)

/**
简单的执行一个系统命令
 */
func CmdExec(strCmd string, arg... string) error{
	pCmd := exec.Command(strCmd,arg...)
	if pCmd == nil{
		return errors.New("et:illegal cmd string:" + strCmd)
	}

	err := pCmd.Run()
	if err != nil{
		return err
	}

	return nil

}

/**
执行一个系统命令，并获得返回值，
返回值中包含普通输出流和错误输出流中的内容
 */
func CmdOutput(strCmd string, arg... string) (string,error) {
	pCmd := exec.Command(strCmd,arg...)
	if pCmd == nil{
		return "",errors.New("cmd building error")
	}

	bufResult,err := pCmd.CombinedOutput()
	if err != nil{
		return "",err
	}

	strResult := string(bufResult)

	return strResult,nil
}

/**
执行一个系统命令，并获得返回值，带有超时时间，以毫秒计算
返回值中包含普通输出流和错误输出流中的内容
*/
func CmdOutputTimeout(iMIllisec int, strCmd string, arg... string) (bool,string,error) {
	if iMIllisec <= 0{
		iMIllisec = 1000*5
	}

	isTimeout := false
	var chState chan int = make(chan int, 1)
	var pErrResult *string = new(string)
	var pStrResult *string = new(string)

	go CmdOutputState(chState, pErrResult,pStrResult, strCmd,arg...)

	xTicker := time.After( time.Millisecond * time.Duration(iMIllisec) )
	select {
	case <-chState :
		isTimeout = false
	case <-xTicker :
		isTimeout = true
			// case
	}

	if isTimeout{
		return isTimeout,"",nil
	}else{
		return isTimeout,*pStrResult,errors.New(*pErrResult)
	}
}

/**
执行一个系统命令，并获得返回值，带一个通道，用于同步状态
返回值中包含普通输出流和错误输出流中的内容
*/
func CmdOutputState(chState chan int, pStrError *string, pStrResult *string, strCmd string, arg... string) {
	pCmd := exec.Command(strCmd,arg...)
	if pCmd == nil{
		*pStrError = "cmd building error"
		*pStrResult = ""
		return
	}

	bufResult,err := pCmd.CombinedOutput()
	if err != nil{
		*pStrError = err.Error()
		*pStrResult = ""
		return
	}

	*pStrError = ""
	*pStrResult = string(bufResult)

	chState <- 1
}
