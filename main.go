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

	//util.SplitVersion("v1.2.3.4 beta")
	//util.SplitVersion("v1.2.3.4 beta")
	//util.SplitVersion("v1.2.3 stable")
	//util.SplitVersion("v1.2alpha")
	//util.SplitVersion("v1.3c")
	//util.SplitVersion("v1.2 space")
	//util.SplitVersion("v1.2.3.4alpha")
	//util.SplitVersion("v1 GM")
	//
	//fmt.Println(util.CompareVersion("v1.1.2", "v1.1.3"))
	//fmt.Println(util.CompareVersion("1.2.3", "1.3.2"))
	//fmt.Println(util.CompareVersion("v1.100.1", "1.99.2"))
	//fmt.Println(util.CompareVersion("v1.2.6", "v1.3.2"))
	//fmt.Println(util.CompareVersion("v1.2.3 alpha", "v1.2.3 beta"))
	//return
	//向Actions环境变量输出两个相同的版本号，以防出错时仍然发布release
	//err := util.UpdateVerInActions("v0", "v0")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//读取程序设置

	configPath, apiPath := "./config.json", "./api.json"

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
		//!! Archiver仓库要先发布一个v0.0.0的release
		current, err = release.ParseReleaseInfo(dCfg.ArchiverOwner, dCfg.ArchiverRepo)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//!! archiverAPI要先备好一个v0.0.0版本信息
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

		fmt.Println("dist目录已删除")

		//更新设置信息和下载更新文件
		var files []string
		if dCfg.TargetGH == true {
			//对于GitHub项目
			files, err = release.DownloadAssets(target.Assets, dCfg)
			if err != nil {
				log.Fatal(err)
			}
			//更新信息
			dApi.ReleaseTime = target.PublishAt
			dApi.ReleaseNote = target.ReleaseNote
		} else {
			//对于非GitHub网站API，直接用DLink下载
			err = util.DownloadFile(dCfg.TargetDLink, "./"+dCfg.DistPath)
			if err != nil {
				log.Fatal(err)
			}
			//更新归档后的所有文件名
			_, fileName := path.Split(dCfg.TargetDLink)
			files = append(files, fileName)
		}
		//处理总文件大小，单个文件大小，进行自动分卷，并且引入精简filter
		files, err = release.AutoSplit(files, dCfg)
		if err != nil {
			log.Fatal(err)
		}
		files, err = release.Rename(files, dCfg)
		if err != nil {
			log.Fatal(err)
		}
		dApi.DownloadLink = release.File2Link(files, dCfg)
	} else if res == 0 {
		fmt.Println("当前版本即是最新版本，无需更新")
		//保存version文件
		if err = util.WriteFast("./version", target.TagName); err != nil {
			log.Fatal(err)
		}
		//保存old_version文件
		if err = util.WriteFast("./old_version", current.TagName); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	} else {
		log.Fatal("出现错误，当前版本>最新版本")
	}

	//更新信息
	//向Actions环境变量输出版本号
	//err = util.UpdateVerInActions(target.TagName, dCfg.Version)
	//if err != nil {
	//	log.Fatal(err)
	//}
	dCfg.Version = target.TagName
	dCfg.VersionList = release.UpdateVersionList(dCfg.VersionList, target.TagName)
	dCfg.CheckTime = time.Now().Format("2006-01-02T15:04Z")
	dApi.Version = target.TagName
	dApi.VersionList = release.UpdateVersionList(dApi.VersionList, target.TagName)
	dApi.CheckTime = dCfg.CheckTime
	dApi.Format = dCfg.Format

	//字符串类型转time
	//s4 := "1999年10月19日" //字符串
	//t4, err := time.Parse("2006年01月02日", s4)	//"2006-01-02T15:04Z"

	//保存程序设置和API文件
	if err = config.WriteConfig(configPath, dCfg); err != nil {
		log.Fatal(err)
	}
	if err = config.WriteApi(apiPath, dApi); err != nil {
		log.Fatal(err)
	}
	//保存version文件
	if err = util.WriteFast("./version", target.TagName); err != nil {
		log.Fatal(err)
	}
	//保存old_version文件
	if err = util.WriteFast("./old_version", current.TagName); err != nil {
		log.Fatal(err)
	}
	//保存download_link文件
	links := strings.Join(dApi.DownloadLink, "\n")
	if err = util.WriteFast("./download_link", links); err != nil {
		log.Fatal(err)
	}
}
