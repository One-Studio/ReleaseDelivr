package custom

import (
	"errors"
	"github.com/One-Studio/ReleaseDelivr/util"
	"regexp"
	"strings"
)

//利用获得上述json的链接，DIY获取版本号和下载链接 (version, link, err)
func GetVersionAndLink(api string) (string, string, error) {
	//固定：获取api的内容
	count := 0
	version, link := "", ""
	content, err := util.GetHttpData(api)
	for ; err != nil && count < 2; count++ {
		content, err = util.GetHttpData(api)
	}
	if err != nil {
		return "", "", err
	}

	//匹配正则表达式获得版本号
	r := regexp.MustCompile(">(x265-Yuuki-\\S+.7z)<")
	t := r.FindAllStringSubmatch(content, -1)

	var files []string
	for _, v := range t {
		if len(v) > 0 {
			files = append(files, v[1])
		}
	}

	if len(files) == 0 {
		return "", "", errors.New("can't get target file names")
	}

	var latest string
	latest = files[0]
	for _, file := range files {
		if strings.Compare(file, latest) == 1 {
			latest = file
		}
	}

	link = "https://down.7086.in/x265-Yuuki-Asuna/" + latest

	r = regexp.MustCompile("x265-Yuuki-([\\S]+).7z")
	res := r.FindStringSubmatch(latest)
	if len(t) > 0 {
		version = res[1]
	}

	return version, link, nil
	//修改：处理得到版本号和链接
	//return "", "", errors.New("can't get version, data: " + content)
}
