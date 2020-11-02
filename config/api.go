package config

import (
	"encoding/json"
	"os"

	"github.com/One-Studio/ReleaseDelivr/util"
)

type Api struct {
	Version      string
	VersionList  []string
	ReleaseTime  string
	CheckTime    string
	DownloadLink []string
	Format       int8
	ReleaseNote  string
}

var defApi = Api{
	Version:      "",
	VersionList:  []string{},
	ReleaseTime:  "",
	CheckTime:    "",
	DownloadLink: []string{},
	Format:       1,
	ReleaseNote:  "",
}

func writeJsonApi(path string, Api Api) error {
	JsonData, err := json.MarshalIndent(Api, "", "  ") //第二个参数要地址传递
	if err != nil {
		return err
	}

	err = util.WriteFast(path, string(JsonData))
	if err != nil {
		return err
	}

	return nil
}

func ReadApi(path string) (Api, error) {
	//检查文件是否存在
	if exist, err := util.IsFileExisted(path); err != nil {
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
		//不存在则生成默认文件以供修改
		if _, err = os.Create(path); err != nil {
			return Api{}, err
		}
		if err := writeJsonApi(path, defApi); err != nil {
			return Api{}, nil
		}
		return Api{}, nil
	}
}

func WriteApi(path string, Api Api) error {
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
	}

	if err := writeJsonApi(path, Api); err != nil {
		return err
	}

	return nil
}
