package et

type ArrVoidSortFn func(pLeft interface{}, pRight interface{}) int

func ArrVoidSort(arrVoid []interface{}, fnCompare ArrVoidSortFn ){
	iArrSize := len(arrVoid)
	for i := 0; i < iArrSize - 1; i++{
		for j := i + 1; j < iArrSize; j++{
			if fnCompare(arrVoid[i], arrVoid[j]) > 0{
				ArrVoidSortSwap(arrVoid,i,j)
			}
		}
	}
}

func ArrVoidSortSwap(arrContainer []interface{}, iIndexLeft int, iIndexRight int){
	if arrContainer == nil{
		return
	}

	xTemp := arrContainer[iIndexLeft]
	arrContainer[iIndexLeft] = arrContainer[iIndexRight]
	arrContainer[iIndexRight] = xTemp
}
