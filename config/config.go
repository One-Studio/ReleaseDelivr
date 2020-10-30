package config

import (
	"encoding/json"
	"os"

	"github.com/One-Studio/ReleaseDelivr/util"
)

type Cfg struct {
	TargetOwner	string
	TargetRepo string
	TargetGH bool
	TargetAPI string
	ArchiverOwner string
	ArchiverRepo string
	ArchiverGH bool
	ArchiverAPI string
	ArchiverVersion string
	Version string
	Checktime string
	Format int8
	CompRatio int8
}

//TODO: 过滤一下设置文件的参数
func configFilter(config Cfg) (bool, error) {

	return true, nil
}

func ReadConfig(path string) (Cfg, error) {
	//检查文件是否存在
	if exist, err := util.IsFileExisted(path) ;err != nil {
		return Cfg{}, err
	} else if exist == true {
		//存在则读取文件
		content, err := util.ReadAll(path)
		if err != nil {
			return Cfg{}, err
		}

		//初始化实例并解析JSON
		var CfgInst Cfg
		err = json.Unmarshal([]byte(content), &CfgInst) //第二个参数要地址传递
		if err != nil {
			return Cfg{}, err
		}
		//CfgInst.Files = nil //清空API，防止累加

		return CfgInst, nil
	} else {

		return Cfg{}, nil
	}
}

func WriteConfig(path string, config Cfg) error {
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
		//TODO: 保存使用具体的Cfg变量
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
