package utils

import (
	"fmt"
	"os"
)

// @Title PathExists
// @Description 判断文件夹是否存在
// 输入参数 path:路径
// 返回参数：true , error
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// @Title PathExists
// @Description 判断文件夹是否存在,如果存在返回true文件，如果不存在创建文件夹
// 输入参数 path:路径
// 返回参数：是否保存文件成功
func SavePathFolder(path string) (isSuccess bool) {

	exist := PathExists(path)
	if !exist {
		// 创建文件夹
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			isSuccess = false
			return
		} else {
			//fmt.Printf("mkdir success!\n")

		}
	}
	isSuccess = true
	return
}
