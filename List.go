package et

import (
	"container/list"
)

/**
traverses the list with the function and every element
@valueOfElement:value of every element
@return:whether to delete this element
 */
type ListWalkFn func(valueOfElement interface{}) bool


func ListCreate() *list.List{
	return list.New()
}

func ListFront(pList *list.List) interface{}{
	if pList == nil{
		return nil
	}

	return pList.Front().Value
}

func ListBack(pList *list.List) interface{}{
	if pList == nil{
		return nil
	}

	return pList.Back().Value
}

func ListLen(pList *list.List) int{
	if pList == nil{
		return 0
	}

	return pList.Len()
}

func ListClear(pList *list.List){
	if pList == nil{
		return
	}

	pList.Init()
}

func ListPeekBack(pList *list.List) interface{}{
	if pList == nil {
		return nil
	}

	if pList.Back() == nil{
		return nil
	}

	return pList.Back().Value
}

func ListPeekFront(pList *list.List) interface{}{
	if pList == nil {
		return nil
	}

	if pList.Front() == nil{
		return nil
	}

	return pList.Front().Value
}

func ListContains(pList *list.List, pValue interface{}) bool{
	if pList == nil{
		return false
	}

	for e := pList.Front(); e != nil; e = e.Next() {
		if e.Value == pValue{
			return true
		}
	}

	return false
}

func ListPushFront(pList *list.List, pValue interface{}) bool{
	if pList == nil{
		return false
	}

	pList.PushFront(pValue)

	return true
}

func ListPushBack(pList *list.List, pValue interface{}) bool{
	if pList == nil{
		return false
	}

	pList.PushBack(pValue)

	return true
}

func ListPopFront(pList *list.List) interface{}{
	if pList == nil{
		return nil
	}

	pElement := pList.Front()
	if pElement == nil{
		return nil
	}

	xValue := pElement.Value
	pList.Remove(pElement)
	return xValue
}

func ListPopBack(pList *list.List) interface{}{
	if pList == nil{
		return nil
	}

	pElement := pList.Back()
	if pElement == nil{
		return nil
	}

	xValue := pElement.Value
	pList.Remove(pElement)
	return xValue
}

func ListRemove(pList *list.List, pValue interface{}) bool{
	if pList == nil{
		return false
	}

	var pElement *list.Element = nil
	for e := pList.Front(); e != nil; e = e.Next() {
		if e.Value == pValue{
			pElement = e
			break
		}
	}

	if pElement != nil{
		pList.Remove(pElement)
		return true
	}

	return true
}


func ListWalk(pList *list.List, fnWalk ListWalkFn) int{
	if pList == nil || fnWalk == nil{
		return 0
	}

	var pElement *list.Element = nil
	for e := pList.Front(); e != nil;  {
		pElement = e
		bDelete := fnWalk(pElement.Value)
		e = e.Next()
		if bDelete{
			pList.Remove(pElement)
		}
	}

	return 0
}
