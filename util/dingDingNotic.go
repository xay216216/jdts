package util

import (
	"net/http"
	"net/url"
	"strings"
)

// 发送提醒
func DingDingNotice(msgType, title, msg string) {
	client := &http.Client{}
	postData := url.Values{}
	postData.Add("msgType", msgType)
	postData.Add("title", title)
	postData.Add("msg", msg)
	postUrl := "https://97.64.40.172:8089/api/v1/sendmsg/dingding"
	req, _ := http.NewRequest("POST", postUrl, strings.NewReader(postData.Encode()))
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)
}
