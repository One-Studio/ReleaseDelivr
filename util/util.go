package util

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func IsEmpty(s string) bool {
	if len(s) > 0 {
		return false
	} else {
		return true
	}
}

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

	if len(dir) != 0 {
		if exist, err := IsFileExisted(dir); err != nil {
			return err
		} else if exist == false {
			if err = os.Mkdir(dir, os.ModePerm); err != nil {
				return err
			}
		}
	}
	if err := ioutil.WriteFile(filePath, []byte(content), 0666); err != nil {
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
	//client := &http.Client{}
	//    //提交请求
	//    reqest, err := http.NewRequest("GET", url, nil)
	//
	//    //增加header选项
	//    reqest.Header.Add("Cookie", "xxxxxx")
	//    reqest.Header.Add("User-Agent", "xxx")
	//    reqest.Header.Add("Time-Zone", "Asia/Shanghai")	//*** 设置时区
	//
	//    if err != nil {
	//        panic(err)
	//    }
	//    //处理返回结果
	//    response, _ := client.Do(reqest)
	//    defer response.Body.Close()
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

//下载文件 (下载地址，存放位置)
func DownloadFile(url string, location string) error {
	//利用HTTP下载文件并读取内容给data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errorInfo := "http failed, check if file exists, HTTP Status Code:" + strconv.Itoa(resp.StatusCode)
		return errors.New(errorInfo)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	//确保下载位置存在
	_, fileName := path.Split(url)
	ok, err := IsFileExisted(location)
	if err != nil {
		return err
	} else if ok == false {
		err := os.Mkdir(location, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//删除已有同名文件
	ok, err = IsFileExisted(location + "/" + fileName)
	if err != nil {
		return err
	} else if ok == true {
		err = os.Remove(location + "/" + fileName)
		if err != nil {
			return err
		}
	}

	//文件写入 先清空再写入 利用ioutil
	err = ioutil.WriteFile(location+"/"+fileName, data, 0666)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//分割版本号 v1.2.3-stable -> 'v', [1, 2, 3], '-stable'   prefix, version slice, suffix
func SplitVersion(version string) (string, []int, string) {
	//分离出prefix
	r := regexp.MustCompile("^([a-zA-Z]?)(\\d+[\\s\\S]*)")
	t := r.FindStringSubmatch(version)

	//错误处理
	if len(t) != 3 {
		log.Println("版本分割出错，step1")
		return "", nil, ""
	}
	prefix, version := t[1], t[2]

	//循环分离出每个版本号，转换成int，赋给切片
	var v []int
	r = regexp.MustCompile("^(\\d+)\\.?([\\s\\S]*)")

	for t = r.FindStringSubmatch(version); len(t) == 3; t = r.FindStringSubmatch(version){
		//debug
		//fmt.Println(t)
		i, err := strconv.Atoi(t[1])
		if err != nil {
			log.Println("版本分割出错，step2")
			return "", nil, ""
		}
		v = append(v, i)
		version = t[2]
	}
	//fmt.Println(r.MatchString(version), verSlice)

	//处理suffix的左侧空格
	suffix := strings.TrimSpace(version)

	//debug
	//fmt.Println("函数返回值", prefix, v, suffix)

	return prefix, v, suffix
}

//对比版本号 返回int8 1->前者更大 -1->后者更大 0->相等 注意：stable等后缀只按串对比大小
func CompareVersion(ver1 string, ver2 string) (int8, error) {
	pre1, v1, suf1 := SplitVersion(ver1)
	pre2, v2, suf2 := SplitVersion(ver2)
	n := 0

	//检测版本号是否出错
	if v1 == nil || v2 == nil {
		return 0, errors.New("version string is null or not matched")
	}
	//检测前缀
	if strings.Compare(pre1, pre2) != 0 {
		return 0, errors.New("prefix of 2 versions are different: " + pre1 + "!=" + pre2)
	}
	//先分割版本号，比较两者共同可比较的部分，然后选择分割切片长度更长的
	if len(v1) < len(v2) {
		n = len(v1)
	} else {
		n = len(v2)
	}
	for i := 0; i < n; i++ {
		if v1[i] > v2[i] {
			return 1, nil
		} else if v1[i] > v2[i] {
			return -1, nil
		}
	}
	if len(v1) > len(v2) {
		return 1, nil
	} else if len(v1) < len(v2) {
		return -1, nil
	} else {
		//版本号长度一致，比较后缀
		return int8(strings.Compare(suf1, suf2)), nil
	}
}

//执行指令
//func Exec(command string, arg ...string) (string, error) {
//	c := exec.Command(command, arg...)
//	//cmd.Args = a
//	out, err := c.CombinedOutput()
//	return string(out), err
//}

//更新环境变量里的版本号
//func UpdateVerInActions(tVer string, cVer string) error {
//	str := strconv.Quote("tVer="+tVer) + " >> $GITHUB_ENV"
//	_, err := Exec("echo", str)
//	if err != nil {
//		return err
//	}
//	//fmt.Println(out)
//
//	str = strconv.Quote("cVer="+cVer) + " >> $GITHUB_ENV"
//	_, err = Exec("echo", str)
//	if err != nil {
//		return err
//	}
//	//fmt.Println(out)
//
//	return nil
//}

//执行一次command指令
func Cmd(command string) (string, error) {
	c := exec.Command("/bin/bash", "-c", command)
	//cmd.Args = a
	out, err := c.CombinedOutput()
	return string(out), err
}

/* 获取指定路径下的所有文件，只搜索当前路径，不进入下一级目录，可匹配后缀过滤（suffix为空则不过滤）*/
//func ListDir(dir, suffix string) (files []string, err error) {
//	files = []string{}
//
//	_dir, err := ioutil.ReadDir(dir)
//	if err != nil {
//		return nil, err
//	}
//
//	suffix = strings.ToLower(suffix) //匹配后缀
//
//	for _, _file := range _dir {
//		if _file.IsDir() {
//			continue //忽略目录
//		}
//		if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(_file.Name()), suffix) {
//			//文件后缀匹配
//			files = append(files, path.Join(dir, _file.Name()))
//		}
//	}
//
//	return files, nil
//}

//实时获取cmd输出
//func CmdAndChangeDirToShow(dir string, commandName string, params []string) error {
//  cmd := exec.Command(commandName, params...)
//  fmt.Println("CmdAndChangeDirToFile", dir, cmd.Args)
//  //StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
//  stdout, err := cmd.StdoutPipe()
//  if err != nil {
//    fmt.Println("cmd.StdoutPipe: ", err)
//    return err
//  }
//  cmd.Stderr = os.Stderr
//  cmd.Dir = dir
//  err = cmd.Start()
//  if err != nil {
//    return err
//  }
//  //创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
//  reader := bufio.NewReader(stdout)
//  //实时循环读取输出流中的一行内容
//  for {
//    line, err2 := reader.ReadString('\n')
//    if err2 != nil || io.EOF == err2 {
//      break
//    }
//    fmt.Println(line)
//  }
//  err = cmd.Wait()
//  return err
//}
