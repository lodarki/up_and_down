package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"up_and_down/constants"
	"up_and_down/entity"
	"up_and_down/utils"
)

var globalSession *session.Manager

type BaseController struct {
	beego.Controller
	DefaultData interface{}
	o           orm.Ormer
}

// 获取session
func (c *BaseController) GetSession() session.Store {
	globalSession := beego.GlobalSessions
	sess, err := globalSession.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		beego.Error("err", err.Error())
	}
	return sess
}

func (c *BaseController) GetDateTime(key string, def ...time.Time) (t time.Time, err error) {
	tStr := c.GetString(key)
	if len(tStr) == 0 {
		return
	}

	t, err = utils.ParseDateTime(tStr)
	if err != nil && len(def) > 0 {
		t = def[0]
	}
	return
}

func (c *BaseController) GetDate(key string, def ...time.Time) (t time.Time, err error) {
	tStr := c.GetString(key)
	if len(tStr) == 0 {
		return
	}

	t, err = utils.ParseDate(tStr)
	if err != nil && len(def) > 0 {
		t = def[0]
	}
	return
}

func (c *BaseController) GetAuthorization() string {
	tokenStr := c.Ctx.Input.Header("Authorization")
	// TODO 未来补充 通过参数来获取 token
	if tokenStr == "" {
		tokenStr = c.GetString("token")
	}

	return tokenStr
}

// 获取IP
func (c *BaseController) GetIp() string {
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	return ip[0]
}

// 提交
func (c *BaseController) Commit() {
	err := c.o.Commit()
	beego.Error("commit", err)
}

// 回滚
func (c *BaseController) Rollback() {
	err := c.o.Rollback()
	beego.Error("rollback", err)
}

// 回调函数
func (c *BaseController) FailFn() {
}

// 成功回到函数
func (c *BaseController) SuccessFn() {
}

// OkResult 返回正确结果
func (c *BaseController) Success(data interface{}, encoding ...bool) {
	if data == nil {
		data = "{}"
	}
	Message := "success"
	result := entity.ApiResult{Code: constants.Success, Data: data, Message: Message}
	c.Data["json"] = &result

	c.SuccessFn()
	cLog := beego.AppConfig.DefaultInt("ControllerLog", 1)
	if cLog == 1 &&
		!strings.Contains(c.Ctx.Request.URL.String(), "/command-record") &&
		!strings.Contains(c.Ctx.Request.URL.String(), "/live") {
	}

	encode := true
	if len(encoding) > 0 {
		encode = encoding[0]
	}
	c.ServeJSON(encode)
}

// ErrResult 错误返回结果
func (c *BaseController) Fail(code int, data interface{}, message string) {
	result := entity.ApiResult{Code: code, Data: data, Message: message}
	c.Data["json"] = &result
	c.FailFn()
	c.ServeJSON(true)
}

// 请求错误
func (c *BaseController) BadRequest(message string) {
	c.Fail(constants.BadRequest, c.DefaultData, message)
}

// 服务器内部错误
func (c *BaseController) ServerError(message string) {
	c.Fail(constants.ServerError, c.DefaultData, message)
}

// 网关错误
func (c *BaseController) GateError(message string) {
	c.Fail(constants.GateError, c.DefaultData, message)
}

// 没有权限
func (c *BaseController) NotPermission(message string) {
	c.Fail(constants.NotPermission, c.DefaultData, message)
}

// 没有找到 404
func (c *BaseController) NotFound(message string) {
	c.Fail(constants.NotFound, c.DefaultData, message)
}

// 数据库执行错误
func (c *BaseController) SQLError(message string) {
	c.Fail(constants.SQLError, c.DefaultData, message)
}

// 获取当前URL
func (c *BaseController) CurrentUrl() string {
	uri := c.Ctx.Input.URI()
	return c.Ctx.Input.Site() + uri
}

// 获取当前URL
func (c *BaseController) Host() string {
	return c.Ctx.Input.Site()
}

// 解析json
func (c *BaseController) JsonDecodeRequestBody(data interface{}) error {
	return json.Unmarshal(c.Ctx.Input.RequestBody, data)
}

type JsonParam map[string]interface{}

func (c *BaseController) JsonParamMap() JsonParam {
	var paramMap = make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &paramMap)
	return paramMap
}

func (c *BaseController) DownloadFile(filePath string) error {

	file, e := os.Open(filePath)
	if e != nil {
		return e
	}

	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", path.Base(filePath))) // 文件名
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	// 最主要的一句
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
	file.Close()
	return nil
}

func (c *BaseController) GetJsonParam(key string, d interface{}, jsonParams ...JsonParam) error {
	var jsonParam JsonParam
	if len(jsonParams) > 0 {
		jsonParam = jsonParams[0]
	} else {
		jsonParam = c.JsonParamMap()
	}

	bytes, marshalE := json.Marshal(jsonParam[key])
	if marshalE != nil {
		return marshalE
	}
	return json.Unmarshal(bytes, d)
}

func (c *BaseController) GetBindStrings(key string) (result []string, bindE error) {
	bindE = c.Ctx.Input.Bind(&result, key)
	if len(result) == 0 {
		keyJson := c.GetString(key)
		bindE = json.Unmarshal([]byte(keyJson), &result)
	}
	return
}

func (c *BaseController) GetBindInt64s(key string) (result []int64, bindE error) {
	bindE = c.Ctx.Input.Bind(&result, key)
	if len(result) == 0 {
		keyJson := c.GetString(key)
		bindE = json.Unmarshal([]byte(keyJson), &result)
	}
	return
}

func (c *BaseController) GetBindInts(key string) (result []int, bindE error) {
	bindE = c.Ctx.Input.Bind(&result, key)
	if len(result) == 0 {
		keyJson := c.GetString(key)
		bindE = json.Unmarshal([]byte(keyJson), &result)
	}
	return
}
