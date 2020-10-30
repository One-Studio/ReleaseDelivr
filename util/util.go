package util

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func ReadAll(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	//对内容的操作
	//ReadFile返回的是[]byte字节切片，要用string()方法转变成字符串
	//去除内容结尾的换行符
	str := strings.TrimRight(string(content), "\n")
	return str, nil
}

//文件写入 先清空再写入 利用ioutil
func WriteFast(filePath string, content string) error {
	dir, _ := path.Split(filePath)
	exist, err := IsFileExisted(dir)
	if err != nil {
		return err
	} else if exist == false {
		os.Mkdir(dir, os.ModePerm)
	}
	err = ioutil.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//判断文件/文件夹是否存在
func IsFileExisted(path string) (bool, error) {
	//返回 true, nil = 存在
	//返回 false, nil = 不存在
	//返回 _, !nil = 位置错误，无法判断
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//利用HTTP Get请求获得数据
func GetHttpData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()

	return string(data), nil
}

//分割版本号 v1.2.3-stable -> ['v', '1', '2', '3', '-stable']
func splitVersion(version string) []string {
	//TODO: 实现版本号分割
	return nil
}

//对比版本号 返回int8 1->前者更大 -1->后者更大 0->相等 注意：stable等后缀只按串对比大小
func CompareVersion(v1 string, v2 string) int8 {
	//先分割版本号，比较两者共同可比较的部分，然后选择分割切片长度更长的
	s1, s2, n := splitVersion(v1), splitVersion(v2), 0
	if len(s1) < len(s2) {
		n = len(s1)
	} else {
		n = len(s2)
	}
	for  i := 0 ; i < n ; i++ {
		if  t :=strings.Compare(s1[i], s2[i]); t == 1 {
			return 1
		} else if t == -1 {
			return -1
		}
	}
	if len(s1) > len(s2) {
		return 1
	} else if len(s1) < len(s2) {
		return -1
	} else {
		return 0
	}
}
