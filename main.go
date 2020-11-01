package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/One-Studio/ReleaseDelivr/config"
	"github.com/One-Studio/ReleaseDelivr/release"
	"github.com/One-Studio/ReleaseDelivr/util"
)

func main() {
	configPath, apiPath := "./config.json", "./api.json"

	//err := util.WriteFast("api.json", "hello world")
	//if err != nil {
	//	fmt.Println(err, "这都能出错")
	//	os.Exit(0)
	//}

	//读取程序设置
	dCfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	dApi, err := config.ReadApi(apiPath)
	if err != nil {
		log.Fatal(err)
	}

	//过滤程序设置，保证必填项非空
	if config.NotEmpty(dCfg) == false {
		fmt.Println("config.json的必填项为空，请正确填写后再打开")
		log.Fatal(errors.New("config.json has not been fully set up"))
	}

	//欢迎
	fmt.Println("---\nHello, this is ReleaseDelivr～\n---")

	//fmt.Println(dCfg, dApi)

	target, current := release.Latest{}, release.Latest{}
	//使用api读取搬运仓库的版本号和附件
	if dCfg.TargetGH == true {
		target, err = release.ParseReleaseInfo(dCfg.TargetOwner, dCfg.TargetRepo)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		target.TagName, err = util.GetHttpData(dCfg.TargetAPI)
		if err != nil {
			log.Fatal(err)
		}
	}

	//使用api读取当前仓库的版本号
	if dCfg.ArchiverGH == true {
		current, err = release.ParseReleaseInfo(dCfg.ArchiverOwner, dCfg.ArchiverRepo)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		current.TagName, err = util.GetHttpData(dCfg.ArchiverAPI)
		if err != nil {
			log.Fatal(err)
		}
	}

	//判断是否要更新
	fmt.Println("最新版本" + target.TagName + "\n当前版本" + current.TagName)
	if res, err := util.CompareVersion(target.TagName, current.TagName); err != nil {
		log.Fatal(err)
	} else if res == 1 {
		//删除dist目录下的所有文件
		err = os.RemoveAll("./" + dCfg.DistPath)
		if err != nil {
			log.Fatal(err)
		}
		//更新设置信息和下载更新文件
		if dCfg.TargetGH == true {
			//对于GitHub项目
			dApi.DownloadLink, err = release.DownloadAssets(target.Assets, dCfg)
			if err != nil {
				log.Fatal(err)
			}
			//更新信息
			//字符串类型转time
			//s4 := "1999年10月19日" //字符串
			//t4, err := time.Parse("2006年01月02日", s4)	//"2006-01-02T15:04Z"
			dCfg.Version = target.TagName
			dCfg.Checktime = time.Now().Format("2006-01-02T15:04Z")
			dApi.Version = target.TagName
			dApi.CheckTime = dCfg.Checktime
			dApi.ReleaseTime = target.PublishAt
			dApi.ReleaseNote = target.ReleaseNote
			dApi.Format = dCfg.Format
			dApi.Split = dCfg.Split
		} else {
			//对于非GitHub网站API，直接用DLink下载
			err = util.DownloadFile(dCfg.TargetDLink, "./"+dCfg.DistPath)
			if err != nil {
				log.Fatal(err)
			}
			//更新信息
			dCfg.Version = target.TagName
			dCfg.Checktime = time.Now().Format("2006-01-02T15:04Z")
			dApi.Version = target.TagName
			dApi.CheckTime = dCfg.Checktime
			dApi.Format = dCfg.Format
			dApi.Split = dCfg.Split
			//更新归档后下载链接
			_, fileName := path.Split(dCfg.TargetDLink)
			var link string
			if dCfg.ArchiverGH == true {
				link = "https://cdn.jsdelivr.net/gh/" + dCfg.ArchiverOwner + "/" + dCfg.ArchiverRepo + "/" + dCfg.DistPath + "/" + fileName
			} else {
				link = dCfg.ArchiverAPI + "/" + dCfg.DistPath + "/" + fileName
			}
			dApi.DownloadLink = []string{link}
			//非GitHub网站无法得知更新时间和更新日志
			//dApi.ReleaseTime = target.PublishAt
			//dApi.ReleaseNote = target.ReleaseNote
		}
	} else if res == 0 {
		fmt.Println("当前版本即是最新版本，无需更新")
		os.Exit(0)
	} else {
		log.Fatal("出现错误，当前版本>最新版本")
	}

	//保存程序设置和API文件
	if err = config.WriteConfig(configPath, dCfg); err != nil {
		log.Fatal(err)
	}
	if err = config.WriteApi(apiPath, dApi); err != nil {
		log.Fatal(err)
	}
	//保存version文件
	if err = util.WriteFast("./version", dApi.Version); err != nil {
		log.Fatal(err)
	}
	//保存download_link文件
	links := strings.Join(dApi.DownloadLink, "\n")
	if err = util.WriteFast("./download_link", links); err != nil {
		log.Fatal(err)
	}
}
