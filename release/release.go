package release

import (
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/One-Studio/ReleaseDelivr/config"
	"github.com/One-Studio/ReleaseDelivr/util"
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
	URL         string  `json:"url"`
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	Message     string  `json:"message"`
	Assets      []Asset `json:"assets"`
	PublishAt   string  `json:"published_at"` //格式 2020-10-20T12:16:01Z T/Z分割 -/:分割
	ReleaseNote string  `json:"note"`
}

//获取并提取GitHub Release的最近一次给Latest类型变量
func ParseReleaseInfo(owner string, repo string) (Latest, error) {
	//GET请求获得JSON
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest"
	jsonData, err := util.GetHttpData(url)
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
		return Latest{}, errors.New("Got response but no valid info. Check URL: " + url)
	}

	return latestInst, nil
}

func DownloadAssets(assets []Asset, cfg config.Cfg) ([]string, error) {
	//必过滤"content_type": "application/octet-stream"
	var links []string
	for _, ast := range assets {
		if ast.ContentType == "application/octet-stream" {
			continue
		}
		for _, flt := range cfg.Filter {
			if strings.Contains(ast.Name, flt) {
				err := util.DownloadFile(ast.BrowserDownloadURL, "./"+cfg.DistPath)
				if err != nil {
					return nil, err
				}
				_, fileName := path.Split(ast.BrowserDownloadURL)

				if cfg.ArchiverGH == true {
					links = append(links, "https://cdn.jsdelivr.net/gh/"+cfg.ArchiverOwner+"/"+cfg.ArchiverRepo+"/"+cfg.DistPath+"/"+fileName)
				} else {
					links = append(links, cfg.ArchiverAPI+"/"+cfg.DistPath+"/"+fileName)
				}
				break
			}
		}
	}

	return links, nil
}