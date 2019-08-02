package utils

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"up_and_down/entity"
	"up_and_down/utils/stringUtils"
)

func getUrlPreFix() string {
	requestProtocol := beego.AppConfig.DefaultString("RequestProtocol", "http")
	return requestProtocol + "://"
}

func GetRequestWithHeader(headerMap map[string]string, url string, params map[string]string, container ...interface{}) (entity.ApiResult, error) {

	if strings.Index(url, "http") < 0 {
		url = getUrlPreFix() + url
	}

	//beego.Info("GetRequestWithHeader url", url)

	var result entity.ApiResult
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		beego.Error(err.Error())
		return result, err
	}

	for k, v := range headerMap {
		request.Header.Set(k, v)
	}

	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	request.URL.RawQuery = query.Encode()
	results, _, err := getResponse(request, false, container...)
	return results, err
}

// 发起get请求
func GetRequest(url string, params map[string]string, container ...interface{}) (entity.ApiResult, error) {
	return GetRequestWithHeader(make(map[string]string), url, params, container...)
}

func PostRequestWithHeader(headMap map[string]string, url string, params map[string]string, container ...interface{}) (entity.ApiResult, error) {
	if strings.Index(url, "http") < 0 {
		url = getUrlPreFix() + url
	}

	var result entity.ApiResult
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range params {
		e := writer.WriteField(k, v)
		if e != nil {
			return result, e
		}
	}
	err := writer.Close()
	if err != nil {
		return result, err
	}
	beego.Info("PostRequest url", url)
	request, err := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	for k, v := range headMap {
		request.Header.Set(k, v)
	}
	apiResult, _, err := getResponse(request, false, container...)
	return apiResult, err
}

//发起post请求
func PostRequest(url string, params map[string]string, container ...interface{}) (entity.ApiResult, error) {
	return PostRequestWithHeader(make(map[string]string), url, params, container...)
}

func GetRequestFile(rawUrl string, params map[string]string, filePath string) (apiResult entity.ApiResult, err error) {
	if strings.Index(rawUrl, "http") < 0 {
		rawUrl = getUrlPreFix() + rawUrl
	}

	u, _ := url.Parse(rawUrl)
	q := u.Query()

	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, e := http.Get(u.String())
	if e != nil {
		err = e
		return
	}

	defer resp.Body.Close()
	beego.Debug(stringUtils.JsonString(resp.Header))
	body, _ := ioutil.ReadAll(resp.Body)

	path := filepath.Dir(filePath)

	mkdirE := os.MkdirAll(path, 0775)
	if mkdirE != nil {
		err = mkdirE
		return
	}

	writeE := ioutil.WriteFile(filePath, body, 0644)
	if writeE != nil {
		err = writeE
		return
	}

	apiResult.Code = 200
	return apiResult, err
}

// 读取响应信息
func getResponse(request *http.Request, pure bool, container ...interface{}) (entity.ApiResult, []byte, error) {
	var apiResult entity.ApiResult
	var result []byte
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		beego.Warn(err.Error())
		if response != nil {
			response.Body.Close()
		}
		return apiResult, result, err
	}

	readBytes := make([]byte, 1024)
	n := 1
	for n > 0 {
		n, err = response.Body.Read(readBytes)
		result = append(result, readBytes[:n]...)
	}

	if response != nil {
		response.Body.Close()
	}

	if len(result) == 0 {
		return apiResult, result, err
	}

	if pure {
		return apiResult, result, err
	}

	//beego.Debug("RequestUtils", string(result))

	if len(container) > 0 {
		unmarshalE := json.Unmarshal(result, container[0])
		if unmarshalE != nil {
			apiResult.Code = 202
			apiResult.Message = "format failed"
			apiResult.Data = string(result)
			return apiResult, result, unmarshalE
		}
		apiResult.Code = 200
		apiResult.Message = ""
		apiResult.Data = container
		return apiResult, result, unmarshalE
	} else {
		unmarshalE := json.Unmarshal(result, &apiResult)
		if unmarshalE != nil {
			apiResult.Code = 202
			apiResult.Message = "format failed"
			apiResult.Data = string(result)
			beego.Error("invalid format resuilt ", string(result))
			return apiResult, result, unmarshalE
		}
	}
	return apiResult, result, nil
}

// 发起json请求
func JsonRequestWithHeader(method string, headerMap map[string]string, url string, params map[string]interface{}, container ...interface{}) (result entity.ApiResult, err error) {
	bs, e := json.Marshal(params)
	if e != nil {
		err = e
		return
	}
	return JsonBytesRequestWithHeader(method, headerMap, url, bs, container...)
}

// 发起json请求
func JsonBytesRequestWithHeader(method string, headerMap map[string]string, url string, bs []byte, container ...interface{}) (result entity.ApiResult, err error) {

	if strings.Index(url, "http") < 0 {
		url = getUrlPreFix() + url
	}

	beego.Info("JsonRequestWithHeader url", url)

	var request *http.Request
	if len(bs) > 0 {
		reader := bytes.NewReader(bs)
		request, err = http.NewRequest(method, url, reader)
	} else {
		request, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return
	}

	if len(bs) > 0 {
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	for k, v := range headerMap {
		request.Header.Set(k, v)
	}

	results, _, err := getResponse(request, false, container...)
	return results, err
}
