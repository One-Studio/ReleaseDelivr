package config

import (
	"encoding/json"
	util "github.com/One-Studio/ReleaseDelivr/util"
	"os"
)

type Api struct {
	Version string
	ReleaseTime string
	CheckTime string
	DownloadLink []string
	Split bool
	Format int8
	ReleaseNote string
}

func ReadApi(path string) (Api, error) {
	//检查文件是否存在
	if exist, err := util.IsFileExisted(path) ;err != nil {
		return Api{}, err
	} else if exist == true {
		//存在则读取文件
		content, err := util.ReadAll(path)
		if err != nil {
			return Api{}, err
		}

		//初始化实例并解析JSON
		var ApiInst Api
		err = json.Unmarshal([]byte(content), &ApiInst) //第二个参数要地址传递
		if err != nil {
			return Api{}, err
		}
		//ApiInst.Files = nil //清空API，防止累加

		return ApiInst, nil
	} else {

		return Api{}, nil
	}
}

func WriteApi(path string, config Api) error {
	//检查文件是否存在
	exist, err := util.IsFileExisted(path)
	if err != nil {
		return err
	} else if exist == true {
		//存在则删除文件
		ok, err := util.IsFileExisted(path)
		if err != nil {
			return err
		} else if ok == true {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		//TODO: 保存使用具体的Api变量
		JsonData, err := json.Marshal(config) //第二个参数要地址传递
		if err != nil {
			return err
		}

		err = util.WriteFast(path, string(JsonData))
		if err != nil {
			return err
		}
	}

	return nil
}