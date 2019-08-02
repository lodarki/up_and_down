package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

func ExecuteShellFile(fileType int, fullFilePath string, args ...string) error {

	binName := ""
	switch fileType {
	case 1:
		binName = "/bin/bash"
		break
	case 2:
		binName = "/usr/bin/python3"
		break
	default:
		binName = "/bin/bash"
		break
	}

	return execute(binName, fullFilePath, args...)
}

func execute(binName string, fullFilePath string, args ...string) error {
	exists := utils.FileExists(fullFilePath)
	if !exists {
		return errors.New(fmt.Sprintf("file %v is not exist", fullFilePath))
	}

	if len(args) > 0 {
		param := strings.Join(args, " ")
		fullFilePath = fullFilePath + " " + param
	}

	cmd := exec.Command(binName, "-c", fullFilePath)

	output, err := cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("Execute Shell:%s failed with error:%s", fullFilePath, err.Error()))
	}
	beego.Info(fmt.Sprintf("Execute Shell:%s finished with output:\n%s", fullFilePath, string(output)))
	return nil
}
