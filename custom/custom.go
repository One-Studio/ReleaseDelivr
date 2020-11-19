package custom

import (
	"errors"
	"github.com/One-Studio/ReleaseDelivr/util"
	"log"
	"os"
	"regexp"
)

//DIY获取版本号和下载链接 (version, link, err)
func GetVersionAndLink(api string) (string, string, error) {
	//固定：获取api的内容
	count := 0
	content, err := util.GetHttpData(api)
<<<<<<< Updated upstream
	for ; err != nil && count < 2; count++ {
=======
	for ; err != nil && count < 1; count++ {
>>>>>>> Stashed changes
		content, err = util.GetHttpData(api)
	}
	if err != nil {
		//官方api不稳定，很容易出错，出现问题直接退出放弃
		//保存version文件
		if err = util.WriteFast("./version", "fuck-U-API"); err != nil {
			log.Fatal(err)
		}
		//保存old_version文件
		if err = util.WriteFast("./old_version", "fuck-U-API"); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
		//return "", "", err
	}

	//匹配正则表达式获得版本号
	r := regexp.MustCompile("build: ffmpeg-git-(\\d+)-amd64-static.tar.xz")
	t := r.FindStringSubmatch(content)

	if len(t) == 2 {
		return t[1], "https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz", nil
	}

	//修改：处理得到版本号和链接
	return "", "", errors.New("can't get version, data: " + content)
}
<<<<<<< Updated upstream
=======

//直接调用p7zip不知道为什么有bug，用linux指令处理
func DealWithFile(dist string) ([]string, error) {
	//"ffmpeg-release-amd64-static"
	//xz -d ffmpeg-release-amd64-static.tar.xz
	//tar -xvf ffmpeg-release-amd64-static.tar
	out, err := util.Cmd("xz -d " + dist + "/ffmpeg-release-amd64-static.tar.xz")
	if err != nil {
		fmt.Println(out)
		return nil, err
	}
	out, err = util.Cmd("tar -xvf " + dist + "/ffmpeg-release-amd64-static.tar")
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	err = os.Remove(dist + "/ffmpeg-release-amd64-static.tar")
	if err != nil {
		return nil, err
	}

	out, err = util.Cmd("rm -rf ffmpeg-release-amd64-static/ffprobe")
	if err != nil {
		fmt.Println(out)
		return nil, err
	}
	out, err = util.Cmd("rm -rf ffmpeg-release-amd64-static/manpages")
	if err != nil {
		fmt.Println(out)
		return nil, err
	}
	err = p7zip.Do7z("ffmpeg-release-amd64-static", dist + "/ffmpeg-linux-amd64.7z", 3, false, "")
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll("ffmpeg-release-amd64-static")
	if err != nil {
		return nil, err
	}

	return []string{"ffmpeg-linux-amd64.7z"}, nil
}
>>>>>>> Stashed changes
