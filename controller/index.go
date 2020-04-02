package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Index 首页
func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello world!")
}

// GetToken 获取token
func GetToken(request *http.Request) (string, error) {
	token := request.Header.Get("Authorization")
	if len(token) == 0 {
		request.ParseForm()
		if len(request.Form["token"]) == 0 {
			return "", errors.New("header Authorization 没有保存 token, url参数也不存在 token， 访问失败 ！")
		}
		token = request.Form["token"][0]
	}
	return token, nil
}

// GetParam 获取url链接中指定名称的参数值
func GetParam(parameterName string, request *http.Request) string {
	if len(request.Form[parameterName]) > 0 {
		return request.Form[parameterName][0]
	}
	return ""
}

// GetBody2Struct 获取POST参数，并转化成指定的struct对象
func GetBody2Struct(request *http.Request, pojo interface{}) error {
	s, _ := ioutil.ReadAll(request.Body)
	if len(s) == 0 {
		return nil
	}
	return json.Unmarshal(s, pojo)
}
