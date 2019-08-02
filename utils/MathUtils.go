package utils

import "errors"

func FTIntersection(from1, to1, from2, to2 int64) (resStart, resEnd int64, err error) {
	if from1 > to2 || from2 > to1 {
		err = errors.New("no intersection")
		return
	}

	if from1 >= from2 {
		resStart = from1
	} else {
		resStart = from2
	}

	if to1 <= to2 {
		resEnd = to1
	} else {
		resEnd = to2
	}
	return
}

type SESubRes struct {
	ResStart int64
	ResEnd   int64
}

func FTSub(oriFrom, oriTo, subFrom, subTo int64) (result []SESubRes, err error) {

	resStart, resEnd, e := FTIntersection(oriFrom, oriTo, subFrom, subTo)
	if e != nil {
		res := SESubRes{
			ResStart: oriFrom,
			ResEnd:   oriTo,
		}
		result = append(result, res)
		return
	}

	if resStart > oriFrom {
		result = append(result, SESubRes{
			ResStart: oriFrom,
			ResEnd:   resStart,
		})
	}

	if oriTo > resEnd {
		result = append(result, SESubRes{
			ResStart: resEnd,
			ResEnd:   oriTo,
		})
	}

	return
}
