package custom

import (
	"encoding/json"
	"errors"
	"github.com/One-Studio/ReleaseDelivr/util"
)

//利用 https://www.sojson.com/json/json2go.html
//把对应json格式转换成golang的如下格式
type Custom struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Version  string `json:"version"`
	Size     int    `json:"size"`
	Download struct {
		SevenZ struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
			Sig  string `json:"sig"`
		} `json:"7z"`
		Zip struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
			Sig  string `json:"sig"`
		} `json:"zip"`
	} `json:"download"`
}

//利用获得上述json的链接，DIY获取版本号和下载链接 (version, link, err)
func GetVersionAndLink(api string) (string, string, error) {
	//固定：获取api的内容
	content, err := util.GetHttpData(api)
	if err != nil {
		return "", "", err
	}

	//固定：初始化实例并解析JSON
	var CustomInst Custom
	err = json.Unmarshal([]byte(content), &CustomInst) //第二个参数要地址传递
	if err != nil {
		return "", "", err
	}

	//固定：过滤失败的情况
	if util.IsEmpty(CustomInst.Version) || util.IsEmpty(CustomInst.Download.SevenZ.URL) {
		return "", "", errors.New("can't get version or download link")
	} else {
		//修改：处理得到版本号和链接
		return CustomInst.Version, CustomInst.Download.SevenZ.URL, nil
	}
}
