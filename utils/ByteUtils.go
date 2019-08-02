package utils

import (
	"encoding/hex"
	"strconv"
	"up_and_down/utils/stringUtils"
)

func ByteToInt(b byte) (i int, err error) {
	b16Str := hex.EncodeToString([]byte{b})
	i64, e := strconv.ParseInt(b16Str, 16, 64)
	if e != nil {
		err = e
		return
	}

	i = stringUtils.Int64ToInt(i64)
	return
}

func StringToByte(s string) (b byte, err error) {
	bytes, e := hex.DecodeString(s)
	err = e
	if len(bytes) == 0 {
		return
	}
	return bytes[0], err
}
