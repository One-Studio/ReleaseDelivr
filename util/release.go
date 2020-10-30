package util

import (
	"encoding/json"
	"errors"
)

type Asset struct {
	URL                string `json:"url"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Latest struct {
	URL     string  `json:"url"`
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Assets  []Asset `json:"assets"`
}

//获取并提取GitHub Release的最近一次给Latest类型变量
func ParseReleaseInfo(owner string, repo string) (Latest, error) {
	//GET请求获得JSON
	jsonData, err := GetHttpData("https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest")
	if err != nil {
		return Latest{}, err
	}

	//初始化实例并解析JSON
	var latestInst Latest
	err = json.Unmarshal([]byte(jsonData), &latestInst) //第二个参数要地址传递
	if err != nil {
		return Latest{}, err
	}

	//链接有问题也会返回Json，且 "Message": "Not Found"
	if latestInst.Message == "Not Found" {
		return Latest{}, errors.New("got Json but no valid. Check URL")
	}

	return latestInst, nil
}