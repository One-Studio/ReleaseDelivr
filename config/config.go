package config

import (
	"encoding/json"
	"os"

	"github.com/One-Studio/ReleaseDelivr/util"
)

type DFilter struct {
	Index string   //包含该字符串的会被处理
	List  []string //要删除的文件列表
}

type Cfg struct {
	TargetOwner     string    //目标仓库主
	TargetRepo      string    //目标仓库名
	TargetGH        bool      //目标是否为GitHub仓库
	TargetAPI       string    //目标非GH所用版本号API
	TargetDLink     string    //目标非GH所用下载链接
	ArchiverOwner   string    //归档仓库主
	ArchiverRepo    string    //归档仓库名
	ArchiverGH      bool      //归档是否为GitHub仓库
	ArchiverAPI     string    //归档非GH所用API（链接前缀）
	ArchiverVersion string    //归档工具的版本号
	Version         string    //当前搬运的版本
	VersionList     []string  //历史可用版本
	Checktime       string    //最近一次检查的时间
	Format          int8      //压缩格式
	CompRatio       int8      //压缩率
	DistPath        string    //归档文件夹
	Filter          []string  //更新附件过滤字符串
	DeleteFilter    []DFilter //Delete过滤，单个文件大于20MB时检查，删除部分文件 from DFilter.List
}

var defCfg = Cfg{
	TargetOwner:     "",
	TargetRepo:      "",
	TargetGH:        true,
	TargetAPI:       "",
	TargetDLink:     "",
	ArchiverOwner:   "",
	ArchiverRepo:    "",
	ArchiverGH:      true,
	ArchiverAPI:     "",
	ArchiverVersion: "v0.3.2",
	Version:         "",
	VersionList:     []string{},
	Checktime:       "",
	Format:          1,
	CompRatio:       2,
	DistPath:        "dist",
	Filter:          []string{".zip", ".exe"},
	DeleteFilter: []DFilter{
		{"", []string{""}},
	},
}

//过滤一下设置文件的参数
func NotEmpty(c Cfg) bool {
	if c.TargetGH == true {
		if util.IsEmpty(c.TargetOwner) || util.IsEmpty(c.TargetRepo) {
			return false
		}
	} else if util.IsEmpty(c.TargetAPI) {
		return false
	}

	if c.ArchiverGH == true {
		if util.IsEmpty(c.ArchiverOwner) || util.IsEmpty(c.ArchiverRepo) {
			return false
		}
	} else if util.IsEmpty(c.ArchiverAPI) {
		return false
	}

	return true
}

func writeJsonConfig(path string, config Cfg) error {
	JsonData, err := json.MarshalIndent(config, "", "  ") //第二个参数要地址传递
	if err != nil {
		return err
	}

	err = util.WriteFast(path, string(JsonData))
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig(path string) (Cfg, error) {
	//检查文件是否存在
	if exist, err := util.IsFileExisted(path); err != nil {
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
		//不存在则生成默认文件以供修改
		if _, err = os.Create(path); err != nil {
			return Cfg{}, err
		}
		if err := writeJsonConfig(path, defCfg); err != nil {
			return Cfg{}, nil
		}
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
	}

	if err := writeJsonConfig(path, config); err != nil {
		return err
	}

	return nil
}
