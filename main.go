package main

import (
	"fmt"
	"log"

	"github.com/One-Studio/ReleaseDelivr/config"
	"github.com/One-Studio/ReleaseDelivr/util"
)

func main() {
	var defCfg = config.Cfg{
		TargetOwner: "",
		TargetRepo: "",
		TargetGH: true,
		TargetAPI: "",
		ArchiverOwner: "",
		ArchiverRepo: "",
		ArchiverGH: true,
		ArchiverAPI: "",
		ArchiverVersion: "v0.0.1",
		Version: "",
		Checktime: "",
		Format: 1,
		CompRatio: 2,
	}
	var defApi = config.Api{
		Version: "",
		ReleaseTime: "",
		CheckTime: "",
		DownloadLink: []string{},
		Split: false,
		Format: 1,
		ReleaseNote: "",
	}
	configPath, apiPath := "./config.json", "./api.json"

	//读取程序设置
	if ok, err := util.IsFileExisted(configPath); err != nil {
		log.Fatal(err)
	} else if ok == false {
		fmt.Println("配置文件不存在，已生成模版")
		if err = config.WriteConfig(configPath, defCfg); err != nil {
			log.Fatal(err)
		}
	}
	delivrCfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if ok, err := util.IsFileExisted(apiPath); err != nil {
		log.Fatal(err)
	} else if ok == false {
		fmt.Println("API文件不存在，已生成模版")
		if err = config.WriteApi(apiPath, defApi); err != nil {
			log.Fatal(err)
		}
	}
	delivrApi, err := config.ReadApi(apiPath)
	if err != nil {
		log.Fatal(err)
	}

	//欢迎
	fmt.Println("---\nHello, this is ReleaseDelivr～\n---")

	target, current := util.Latest{}
	//使用api读取搬运仓库的版本号和附件
	if delivrCfg.TargetGH == true {
		target, err = util.ParseReleaseInfo(delivrCfg.TargetOwner, delivrCfg.TargetRepo)
		if err != nil {
			log.Fatal(err)
		}
	} else {

	}

	//使用api读取当前仓库的版本号
	current, err = util.ParseReleaseInfo(delivrCfg.ArchiverOwner, delivrCfg.ArchiverRepo)
	if err != nil {
		log.Fatal(err)
	}
	//判断是否要更新
	if res := util.CompareVersion(target.TagName, current.TagName); res == 1 {
		//更新设置信息和下载更新文件

	}

	//保存程序设置和API文件
	if err = config.WriteConfig(configPath, delivrCfg); err != nil {
		log.Fatal(err)
	}
	if err = config.WriteApi(apiPath, delivrApi); err != nil {
		log.Fatal(err)
	}
	//保存version文件
	if err = util.WriteFast("./version", delivrApi.Version); err != nil {
		log.Fatal(err)
	}
	//保存download_link文件

}
