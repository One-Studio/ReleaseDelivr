package p7zip

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/One-Studio/ReleaseDelivr/util"
)

const (
	path string = "./7z"
)

var Support = [4]string{".7z", ".zip", ".gzip", ".tar"} //7z、ZIP、GZIP、BZIP2 和 TAR

func check7z() bool {
	if ok, err := util.IsFileExisted(path); err != nil {
		return false
	} else if ok == true {
		return true
	} else {
		return false
	}
}

//7z压缩
func Do7z(from string, to7z string, ratio int8, split bool, size string) error {
	//检查7z是否存在
	if ok := check7z(); ok == false {
		return errors.New("7z executable file is not existed")
	}
	//初始化
	//cmd := exec.Command(path)
	var a []string
	var r string
	//设置7z位置
	a = append(a, path)
	//设置"压缩"模式
	a = append(a, "a")
	//设置压缩算法=LZMA2
	a = append(a, "-mm=LZMA2")
	//设置压缩率
	switch ratio {
	case 1:
		r = "1" //快速
	case 2:
		r = "5" //标准
	case 3:
		r = "9" //极限
	}
	a = append(a, "-mx"+r)
	//设置分卷
	if split == true {
		//正则表达式匹配
		if match, err := regexp.Match("^\\d+[bkmgBKMG]$", []byte(size)); err != nil {
			return err
		} else if match == false {
			return errors.New("size of volumes is not correct, eg: 19m or 1g, current: " + size)
		}
		a = append(a, "-v"+size)
	}
	//设置压缩包位置
	if strings.HasSuffix(to7z, ".7z") {
		a = append(a, strconv.Quote(to7z))
	} else {
		return errors.New("7z path/name is incorrect: " + to7z)
	}
	//设置文件位置
	a = append(a, strconv.Quote(from))
	//出现问题先打印混合输出再返回error
	str := strings.Join(a, " ")
	out, err := util.Cmd(str)
	if err != nil {
		fmt.Println(out)
	}
	return err
}

//7z解压
func Un7z(from7z string, to string) error {
	//检查7z是否存在
	if ok := check7z(); ok == false {
		return errors.New("7z executable file is not existed")
	}
	//检查from7z是否以.7z结尾
	if !strings.HasSuffix(from7z, ".7z") && !strings.HasSuffix(from7z, ".001") {
		return errors.New("7z path/name is incorrect: " + from7z)
	}
	if ok, err := util.IsFileExisted(from7z + ".001"); err != nil {
		return err
	} else if ok == true {
		from7z += ".001"
	}

	//初始化
	var a []string
	//设置7z位置
	a = append(a, path)
	//设置"解压"模式 x 在 xxx/文件名/ | e 在 xxx/ !!!e慎用，所有子文件和文件夹都到一个目录了
	a = append(a, "x")
	//设置压缩包位置
	a = append(a, strconv.Quote(from7z))
	//设置文件位置
	a = append(a, "-o"+strconv.Quote(to))
	//出现问题先打印混合输出再返回error
	str := strings.Join(a, " ")
	out, err := util.Cmd(str)
	if err != nil {
		fmt.Println(out)
	}
	return err
}
