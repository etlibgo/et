package et

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/**
定时器工具类
主要使用 ClockManager 这个类的 AddInterval或者AddAbsolute方法，
用于添加绝对时间触发或者按间隔时间触发。
触发的操作可以是访问URL，执行系统命令，调用回调函数。
使用的时候记得Start，Stop
 */

type ClockCallback interface{
	OnTime(strEventName string, xParam interface{}, iCurrentStamp13 int64)
}

const CLOCK_ABSOLUTE_FORMAT_DATE string = "0000-00-00 00:00:00"
const CLOCK_ABSOLUTE_FORMAT_WEEK string = "0 00:00:00"

const CLOCK_TIME_TYPE_ABSOLUTE string = "absolute"
const CLOCK_TIME_TYPE_INTERVAL string = "interval"

const CLOCK_ACTION_TYPE_URL string = "url"
const CLOCK_ACTION_TYPE_CMD string = "cmd"
const CLOCK_ACTION_TYPE_CALLBACK string = "callback"

const CLOCK_CHECK_RESULT_SLEEP int = 0
const CLOCK_CHECK_RESULT_RUN int = 1
const CLOCK_CHECK_RESULT_REMOVE int = 2

const CLOCK_ENTITY_HTTP_TIMEOUT int = 30*1000
const CLOCK_MANAGER_INTERVAL int64 = 100

const CLOCK_FIELD_TIME_TYPE string = "TimeType"
const CLOCK_FIELD_ACTION_TYPE string = "ActionType"
const CLOCK_FIELD_EVENT_NAME string = "EventName"
const CLOCK_FIELD_ACTION_CONTENT string = "ActionContent"
const CLOCK_FIELD_MAX_COUNT string = "MaxCount"
const CLOCK_FIELD_FIRST_DELAY string = "FirstDelay"
const CLOCK_FIELD_INTERVAL string = "Interval"
const CLOCK_FIELD_ABSOLUTE_PERIOD string = "AbsolutePeriod"
const CLOCK_FIELD_ABSOLUTE_MIN_SPAN string = "AbsoluteMinSpan"
const CLOCK_FIELD_ABSOLUTE_FORMAT string = "AbsoluteFormat"
const CLOCK_FIELD_ABSOLUTE_MAP string = "AbsoluteMap"
const CLOCK_FIELD_CURRENT_COUNT string = "CurrentCount"
const CLOCK_FIELD_LAST_TIME string = "LastTime"
const CLOCK_FIELD_SETUP_TIME string = "SetupTime"
const CLOCK_FIELD_TO_RUN_IN_THE_SAME_THREAD string = "ToRunInTheSameThread"
const CLOCK_FIELD_PARAM_OBJECT string = "ParamObject"
const CLOCK_FIELD_CALLBACK_IMPL string = "CallbackImpl"

/**
定时器事件实体类
一个定时管理器含有多个事件实体
一般不需要操作此类
 */
type ClockEntity struct{
	TimeType string		`Required:"1" Enum:"interval|absolute" Remark:"定时类型 按间隔，按绝对时间"`
	ActionType string	`Required:"1" Enum:"url|cmd|callback" Remark:"操作类型 访问URL，执行系统命令，调用回调函数"`
	EventName string	`Required:"1" Remark:"给定时器定一个名称，确保相对唯一即可"`
	ActionContent string	`Required:"1" Remark:"定时器操作的内容，根据操作类型，可能是URL，也可能是系统命令，也可能是备注信息"`
	MaxCount int	`Required:"0" Remark:"最大执行次数"`
	FirstDelay int64	`Required:"0" Remark:"要推迟多久才执行第一次"`
	Interval int64	`Required:"0" Remark:"间隔时间"`
	AbsolutePeriod    string	`Required:"0" Remark:"绝对时间的周期，是年，月，日，时，分，秒钟的一种"`
	AbsoluteMinSpan int64	`Required:"0" Remark:"绝对时间的最小触发间隔"`
	AbsoluteFormat string	`Required:"0" Remark:"绝对时间的格式，有星期(0 00:00:00)和非星期两种模式(0000-00-00 00:00:00),用星号*代表忽略"`
	AbsoluteMap map[string]int `Required:"0" Remark:"解析完绝对时间格式后得到的一个表"`
	CurrentCount int	`Required:"0" Remark:"当前已经执行的次数"`
	LastTime int64	`Required:"0" Remark:"上次执行的时间戳"`
	SetupTime int64	`Required:"0" Remark:"启动时间"`
	ToRunInTheSameThread bool	`Required:"0" Remark:"是否要在同一个线程内执行，将决定是否序列化执行，用于前后定时器可能会交错的情况"`
	ParamObject  interface{}	`Required:"0" Remark:"用于函数回调的情况，将被设置为函数的参数"`
	CallbackImpl ClockCallback	`Required:"0" Remark:"被设置的回调函数"`
}

func ClockEntityNew(strEventName string, strActionType string, strTimeType string) *ClockEntity{
	xClockEntity := ClockEntity{}
	xClockEntity.EventName = strEventName
	xClockEntity.ActionType = strActionType
	xClockEntity.TimeType = strTimeType
	xClockEntity.ActionContent = ""
	xClockEntity.MaxCount = 0
	xClockEntity.FirstDelay = 0
	xClockEntity.Interval = 0
	xClockEntity.AbsolutePeriod = ""
	xClockEntity.AbsoluteMap = map[string]int{}
	xClockEntity.AbsoluteMinSpan = 0
	xClockEntity.AbsoluteFormat = "0000-00-00 00:00:00"
	xClockEntity.CurrentCount = 0
	xClockEntity.LastTime = 0
	xClockEntity.SetupTime = 0
	xClockEntity.ToRunInTheSameThread = false
	xClockEntity.ParamObject = nil
	xClockEntity.CallbackImpl = nil
	return &xClockEntity
}

func (pThis *ClockEntity)AnalyzeFormat() error{
	if pThis.AbsoluteFormat == "" || len(pThis.AbsoluteFormat) < 10{
		return errors.New("empty absolute time format")
	}

	iFormatYear := 0
	iFormatMonth := 0
	iFormatDay := 0
	iFormatHour := 0
	iFormatMinute := 0
	iFormatSecond := 0
	iFormatWeek := 0

	if len(pThis.AbsoluteFormat) == 19{
		//date 0000-00-00 00:00:00
		if pThis.AbsoluteFormat[14:16] == "**"{
			pThis.AbsolutePeriod = TIME_FORMAT_SECOND
			iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[17:])
			pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
		}else if pThis.AbsoluteFormat[11:13] == "**"{
			pThis.AbsolutePeriod = TIME_FORMAT_MINUTE
			iFormatMinute = CastStrToInt(pThis.AbsoluteFormat[14:16])
			iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[17:])
			pThis.AbsoluteMap[TIME_FORMAT_MINUTE] = iFormatMinute
			pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
		}else if pThis.AbsoluteFormat[8:10] == "**"{
			pThis.AbsolutePeriod = TIME_FORMAT_HOUR
			iFormatHour = CastStrToInt(pThis.AbsoluteFormat[11:13])
			iFormatMinute = CastStrToInt(pThis.AbsoluteFormat[14:16])
			iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[17:])
			pThis.AbsoluteMap[TIME_FORMAT_HOUR] = iFormatHour
			pThis.AbsoluteMap[TIME_FORMAT_MINUTE] = iFormatMinute
			pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
		}else if pThis.AbsoluteFormat[5:7] == "**"{
			pThis.AbsolutePeriod = TIME_FORMAT_DAY
			iFormatDay = CastStrToInt(pThis.AbsoluteFormat[8:10])
			iFormatHour = CastStrToInt(pThis.AbsoluteFormat[11:13])
			iFormatMinute = CastStrToInt(pThis.AbsoluteFormat[14:16])
			iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[17:])
			pThis.AbsoluteMap[TIME_FORMAT_DAY] = iFormatDay
			pThis.AbsoluteMap[TIME_FORMAT_HOUR] = iFormatHour
			pThis.AbsoluteMap[TIME_FORMAT_MINUTE] = iFormatMinute
			pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
		}else if pThis.AbsoluteFormat[0:4] == "****"{
			pThis.AbsolutePeriod = TIME_FORMAT_MONTH
			iFormatMonth = CastStrToInt(pThis.AbsoluteFormat[5:7])
			iFormatDay = CastStrToInt(pThis.AbsoluteFormat[8:10])
			iFormatHour = CastStrToInt(pThis.AbsoluteFormat[11:13])
			iFormatMinute = CastStrToInt(pThis.AbsoluteFormat[14:16])
			iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[17:])
			pThis.AbsoluteMap[TIME_FORMAT_MONTH] = iFormatMonth
			pThis.AbsoluteMap[TIME_FORMAT_DAY] = iFormatDay
			pThis.AbsoluteMap[TIME_FORMAT_HOUR] = iFormatHour
			pThis.AbsoluteMap[TIME_FORMAT_MINUTE] = iFormatMinute
			pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
		}else{
			pThis.AbsolutePeriod = TIME_FORMAT_YEAR
			iFormatYear = CastStrToInt(pThis.AbsoluteFormat[0:4])
			iFormatMonth = CastStrToInt(pThis.AbsoluteFormat[5:7])
			iFormatDay = CastStrToInt(pThis.AbsoluteFormat[8:10])
			iFormatHour = CastStrToInt(pThis.AbsoluteFormat[11:13])
			iFormatMinute = CastStrToInt(pThis.AbsoluteFormat[14:16])
			iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[17:])
			pThis.AbsoluteMap[TIME_FORMAT_YEAR] = iFormatYear
			pThis.AbsoluteMap[TIME_FORMAT_MONTH] = iFormatMonth
			pThis.AbsoluteMap[TIME_FORMAT_DAY] = iFormatDay
			pThis.AbsoluteMap[TIME_FORMAT_HOUR] = iFormatHour
			pThis.AbsoluteMap[TIME_FORMAT_MINUTE] = iFormatMinute
			pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
		}
	}else if len(pThis.AbsoluteFormat) == 10{
		//week 0 00:00:00
		pThis.AbsolutePeriod = TIME_FORMAT_WEEK
		iFormatWeek = CastStrToInt(pThis.AbsoluteFormat[0:1])
		iFormatHour = CastStrToInt(pThis.AbsoluteFormat[2:4])
		iFormatMinute = CastStrToInt(pThis.AbsoluteFormat[5:7])
		iFormatSecond = CastStrToInt(pThis.AbsoluteFormat[8:])
		pThis.AbsoluteMap[TIME_FORMAT_WEEK] = iFormatWeek
		pThis.AbsoluteMap[TIME_FORMAT_HOUR] = iFormatHour
		pThis.AbsoluteMap[TIME_FORMAT_MINUTE] = iFormatMinute
		pThis.AbsoluteMap[TIME_FORMAT_SECOND] = iFormatSecond
	}else{
		return errors.New("empty absolute time format")
	}

	return nil
}

func (pThis *ClockEntity)CheckBeforeRun(iStamp13 int64) int{
	iResult := CLOCK_CHECK_RESULT_SLEEP

	if pThis.MaxCount > 0{
		if pThis.CurrentCount >= pThis.MaxCount{
			return CLOCK_CHECK_RESULT_REMOVE
		}
	}

	if pThis.TimeType == CLOCK_TIME_TYPE_INTERVAL {
		iResult = pThis.CheckBeforeRunByInterval(iStamp13)
	}else if pThis.TimeType == CLOCK_TIME_TYPE_ABSOLUTE {
		iResult = pThis.CheckBeforeRunByAbsolute(iStamp13)
	}

	if iResult == CLOCK_CHECK_RESULT_RUN {
		pThis.LastTime = iStamp13
		pThis.CurrentCount += 1
	}

	return iResult
}

func (pThis *ClockEntity)CheckBeforeRunByAbsolute(iStamp13 int64) int{
	var iResult int = CLOCK_CHECK_RESULT_RUN

	tNow := TimeFromStamp13(iStamp13)

	mapCurrent := TimeFormatObjToMap(tNow)
	isMapContained := TimeFormattedMapIsContained(pThis.AbsoluteMap,mapCurrent)

	if !isMapContained {
		return CLOCK_CHECK_RESULT_SLEEP
	}else if pThis.CurrentCount > 0{
		if pThis.AbsoluteMinSpan > 0{
			if iStamp13 - pThis.LastTime < pThis.AbsoluteMinSpan{
				iResult = CLOCK_CHECK_RESULT_SLEEP
			}
		}else{
			return CLOCK_CHECK_RESULT_REMOVE
		}
	}

	if iResult == CLOCK_CHECK_RESULT_RUN {
		if pThis.AbsolutePeriod == TIME_FORMAT_WEEK{
			tNext := tNow.AddDate(0,0,7)
			pThis.AbsoluteMinSpan = TimeGetStamp13(tNext) - TimeGetStamp13(tNow)
		}else if pThis.AbsolutePeriod == TIME_FORMAT_YEAR{
			pThis.AbsoluteMinSpan = 0
		}else if pThis.AbsolutePeriod == TIME_FORMAT_MONTH{
			tNext := tNow.AddDate(1,0,0)
			pThis.AbsoluteMinSpan = TimeGetStamp13(tNext) - TimeGetStamp13(tNow)
		}else if pThis.AbsolutePeriod == TIME_FORMAT_DAY{
			tNext := tNow.AddDate(0,1,0)
			pThis.AbsoluteMinSpan = TimeGetStamp13(tNext) - TimeGetStamp13(tNow)
		}else if pThis.AbsolutePeriod == TIME_FORMAT_HOUR{
			tNext := tNow.AddDate(0,0,1)
			pThis.AbsoluteMinSpan = TimeGetStamp13(tNext) - TimeGetStamp13(tNow)
		}else if pThis.AbsolutePeriod == TIME_FORMAT_MINUTE{
			tNext := tNow.Add(time.Hour * 1)
			pThis.AbsoluteMinSpan = TimeGetStamp13(tNext) - TimeGetStamp13(tNow)
		}else if pThis.AbsolutePeriod == TIME_FORMAT_SECOND{
			tNext := tNow.Add(time.Minute * 1)
			pThis.AbsoluteMinSpan = TimeGetStamp13(tNext) - TimeGetStamp13(tNow)
		}else{
			return CLOCK_CHECK_RESULT_REMOVE
		}
	}

	return iResult
}

func (pThis *ClockEntity)CheckBeforeRunByInterval(iStamp13 int64) int{
	var iResult int = CLOCK_CHECK_RESULT_RUN

	if pThis.FirstDelay > 0 && pThis.CurrentCount == 0{
		if iStamp13 - pThis.SetupTime < pThis.FirstDelay{
			iResult = CLOCK_CHECK_RESULT_SLEEP
		}
	}else{
		if pThis.LastTime != 0{
			if iStamp13 - pThis.LastTime < pThis.Interval{
				iResult = CLOCK_CHECK_RESULT_SLEEP
			}
		}else{
			if iStamp13 - pThis.SetupTime < pThis.Interval{
				iResult = CLOCK_CHECK_RESULT_SLEEP
			}
		}
	}

	return iResult
}

func (pThis *ClockEntity)Run(iCurrentStamp13 int64){

	if pThis.ToRunInTheSameThread{
		pThis.RunSync(iCurrentStamp13)
	}else{
		go pThis.RunSync(iCurrentStamp13)
	}
}

func (pThis *ClockEntity)RunSync(iCurrentStamp13 int64){

	if pThis.ActionType == CLOCK_ACTION_TYPE_URL {
		_,err := HttpClientGetResp(pThis.ActionContent, CLOCK_ENTITY_HTTP_TIMEOUT,nil)
		if err != nil{
			fmt.Printf("ClockEntity.RunSync %s\r\n",err.Error())
		}
	}else if pThis.ActionType == CLOCK_ACTION_TYPE_CMD {
		err := CmdExec(pThis.ActionContent)
		if err != nil{
			fmt.Printf("ClockEntity.RunSync %s\r\n",err.Error())
		}
	}else if pThis.ActionType == CLOCK_ACTION_TYPE_CALLBACK {
		if pThis.CallbackImpl != nil{
			pThis.CallbackImpl.OnTime(pThis.EventName,pThis.ParamObject,iCurrentStamp13)
		}
	}
}

type ClockManager struct{
	Running bool
	EntityMap map[string]*ClockEntity
	Interval int64
	Mutex sync.Mutex
}

var ClockManagerStatic = struct{
	DefaultInstance *ClockManager
}{
	DefaultInstance: nil,
}

/**
新建一个定时器管理器
不建议一个应用内含有多个定时器，
请使用ClockManagerGetDefaultInstance()方法
 */
func ClockManagerNew() *ClockManager{
	xClockManager := ClockManager{}
	xClockManager.Interval = CLOCK_MANAGER_INTERVAL
	xClockManager.Running = false
	xClockManager.EntityMap = map[string]*ClockEntity{}
	xClockManager.Mutex = sync.Mutex{}
	return &xClockManager
}

/**
获得默认的系统定时管理器，建议使用此方法，而不是new
 */
func ClockManagerGetDefaultInstance() *ClockManager{
	if ClockManagerStatic.DefaultInstance == nil{
		ClockManagerStatic.DefaultInstance = ClockManagerNew()
	}

	return ClockManagerStatic.DefaultInstance
}

func (pThis *ClockManager)Loop(){
	for{
		if !pThis.IsRunning(){
			break
		}
		time.Sleep(time.Millisecond * time.Duration( pThis.Interval ) )

		arrToRemove := make([]string,0)
		iNowStamp13 := TimeNowStamp13()
		for strEventName,pEntity := range pThis.EntityMap{
			iCheckResult := pEntity.CheckBeforeRun(iNowStamp13)
			if iCheckResult == CLOCK_CHECK_RESULT_SLEEP {
				continue
			}else if iCheckResult == CLOCK_CHECK_RESULT_REMOVE {
				arrToRemove = append(arrToRemove,strEventName)
			}else if iCheckResult == CLOCK_CHECK_RESULT_RUN {
				pEntity.Run(iNowStamp13)
			}else{
				continue
			}
		}
		for _,strEventName := range arrToRemove{
			delete(pThis.EntityMap,strEventName)
		}
	}
}

/**
启动定时器
 */
func (pThis *ClockManager)Start() error{
	if pThis.IsRunning(){
		return nil
	}

	pThis.SetRunning(true)

	go pThis.Loop()

	return nil
}

/**
关闭定时器
 */
func (pThis *ClockManager)Stop(){
	if !pThis.IsRunning(){
		return
	}

	pThis.SetRunning(false)
	time.Sleep(time.Millisecond * time.Duration(CLOCK_MANAGER_INTERVAL))
}

func (pThis *ClockManager)IsRunning() bool{
	bRunning := false

	pThis.Mutex.Lock()
	bRunning = pThis.Running
	pThis.Mutex.Unlock()

	return bRunning
}

func (pThis *ClockManager)AddEntity(pEntity *ClockEntity) error{
	if pEntity == nil{
		return errors.New("nil entity")
	}

	if pEntity.TimeType == "" || pEntity.ActionType == ""{
		return errors.New("entity format error")
	}
	pEntity.SetupTime = TimeNowStamp13()

	pThis.Mutex.Lock()
	pThis.EntityMap[pEntity.EventName] = pEntity
	pThis.Mutex.Unlock()
	return nil
}

func (pThis *ClockManager)RemoveEntity(strEventName string) error{
	pThis.Mutex.Lock()
	delete(pThis.EntityMap,strEventName)
	pThis.Mutex.Unlock()
	return nil
}

/**
清空所有定时器
 */
func (pThis *ClockManager)Clear(){

	pThis.Mutex.Lock()
	pThis.EntityMap = map[string]*ClockEntity{}
	pThis.Mutex.Unlock()
}

func (pThis *ClockManager)SetRunning(bRunning bool){
	pThis.Mutex.Lock()
	pThis.Running = bRunning
	pThis.Mutex.Unlock()
}

/**
按照绝对时间的模式添加一个定时事件
*/
func (pThis *ClockManager)AddAbsolute(strEventName string, strActionType string, strActionContent string, implCallback ClockCallback, strAbsoluteFormat string, objParam interface{}, toRunInTheSameThread bool) error{
	if strActionType == ""{
		return errors.New("invalid action type")
	}
	pEntity := ClockEntityNew(strEventName,strActionType, CLOCK_TIME_TYPE_ABSOLUTE)
	pEntity.ActionContent = strActionContent
	pEntity.CallbackImpl = implCallback
	pEntity.AbsoluteFormat = strAbsoluteFormat
	pEntity.ParamObject = objParam
	pEntity.ToRunInTheSameThread = toRunInTheSameThread

	err := pEntity.AnalyzeFormat()
	if err != nil{
		return err
	}

	err = pThis.AddEntity(pEntity)
	if err != nil{
		return err
	}

	return nil
}

/**
按照间隔的模式添加一个定时事件
 */
func (pThis *ClockManager)AddInterval(strEventName string, strActionType string, strActionContent string, implCallback ClockCallback, iInterval int64, iFirstDelay int64, iMaxCount int, objParam interface{}, toRunInTheSameThread bool) error{

	if strActionType == ""{
		return errors.New("invalid action type")
	}
	pEntity := ClockEntityNew(strEventName,strActionType, CLOCK_TIME_TYPE_INTERVAL)
	pEntity.ActionContent = strActionContent
	pEntity.CallbackImpl = implCallback
	pEntity.Interval = iInterval
	pEntity.FirstDelay = iFirstDelay
	pEntity.MaxCount = iMaxCount
	pEntity.ParamObject = objParam
	pEntity.ToRunInTheSameThread = toRunInTheSameThread

	err := pThis.AddEntity(pEntity)
	if err != nil{
		return err
	}

	return nil
}






