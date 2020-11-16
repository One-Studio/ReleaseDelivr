package release

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/One-Studio/ReleaseDelivr/p7zip"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/One-Studio/ReleaseDelivr/config"
	"github.com/One-Studio/ReleaseDelivr/util"
)

type Asset struct {
	URL                string `json:"url"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Latest struct {
	URL         string  `json:"url"`
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	Message     string  `json:"message"`
	Assets      []Asset `json:"assets"`
	PublishAt   string  `json:"published_at"` //格式 2020-10-20T12:16:01Z T/Z分割 -/:分割
	ReleaseNote string  `json:"body"`
}

//获取并提取GitHub Release的最近一次给Latest类型变量
func ParseReleaseInfo(owner string, repo string) (Latest, error) {
	//GET请求获得JSON
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest"
	jsonData, err := util.GetHttpData(url)
	if err != nil {
		return Latest{}, err
	}

	//初始化实例并解析JSON
	var latestInst Latest
	err = json.Unmarshal([]byte(jsonData), &latestInst) //第二个参数要地址传递
	if err != nil {
		return Latest{}, err
	}

	//链接有问题也会返回Json，且 "Message": "Not Found"
	if latestInst.Message == "Not Found" {
		return Latest{}, errors.New("Got response but no valid info. Check URL: " + url)
	}

	return latestInst, nil
}

//下载附件，返回所有的文件名
func DownloadAssets(assets []Asset, cfg config.Cfg) ([]string, error) {
	//必过滤"content_type": "application/octet-stream"
	var files []string
	for _, ast := range assets {
		if ast.ContentType == "application/octet-stream" {
			continue
		}
		for _, flt := range cfg.Filter {
			if strings.Contains(ast.Name, flt) {
				err := util.DownloadFile(ast.BrowserDownloadURL, "./"+cfg.DistPath)
				if err != nil {
					return nil, err
				}
				_, fileName := path.Split(ast.BrowserDownloadURL)

				files = append(files, fileName)
				break
			}
		}
	}

	return files, nil
}

//检查当前文件夹的大小是否小于 ~MB
func checkDirSize(dir string,MB int64) error {
	var filesize int64
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			filesize += info.Size()
			fmt.Println(info.Name(), info.Size()/1024/1024, "MB")

			return nil
		})
	if err != nil {
		return err
	} else if filesize/1024/1024 >= MB {
		return errors.New("total file size is above " + strconv.FormatInt(MB, 10) + "MB and JsDelivr won't work")
	} else {
		return nil
	}
}

//自动分割超过20MB的文件，可以根据DeleteFilter适当删除文件
func AutoSplit(files []string, cfg config.Cfg) ([]string, error) {
	//先检查当前目录下所有文件大小之和是否小于50MB
	//if err := checkDirSize(50); err != nil {
	//	return nil, err
	//}

	//再检查所有files，如果有delete过滤的就删除对应文件，如果有超过20MB的就分卷压缩
	var split []string
	for _, file := range files {
		f, err := os.Stat("./" + cfg.DistPath + "/" + file)
		if err != nil {
			return nil, err
		} else if f.Size()/1024/1024 >= 20 {
			//分卷操作
			ext := path.Ext(file)
			filename := strings.TrimSuffix(file, ext)
			//检查文件的格式是不是7z支持的
			ok := 0
			for _, suffix := range p7zip.Support {
				if strings.ToLower(ext) == strings.ToLower(suffix) {
					ok++
					break
				}
			}
			//非7z支持的文件格式则对path做特殊处理，不解压
			if ok == 0 {
				if err := os.Rename(cfg.DistPath+"/"+file, "./temp/"+file); err != nil {
					return nil, err
				}
			} else {
				//先解压到临时目录
				err = p7zip.Un7z(cfg.DistPath+"/"+file, "./temp/"+filename)
				if err != nil {
					return nil, err
				}

				//return nil, nil
				//判断是不是要过滤的
				for _, dflt := range cfg.DeleteFilter {
					if strings.Contains(file, dflt.Index) {
						//删掉List指定的文件/文件夹
						for _, flt := range dflt.List {
							//解决解压后的文件在 filename/另一文件夹名/... 的问题
							if exist, err := util.IsFileExisted("./temp/" + filename + "/" + flt); err != nil {
								return nil, err
							} else if exist == false {
								//出现问题，遍历解决
								subDir := ""
								err := filepath.Walk("./temp",
									func(path string, info os.FileInfo, err error) error {
										if err != nil {
											return err
										} else if info.IsDir() == false {
											return nil
										} else if exist, err := util.IsFileExisted("./temp/" + filename + "/" + info.Name() + "/" + flt); err != nil {
											return err
										} else if exist == true {
											subDir = info.Name()
											fmt.Println("SubDir = " + subDir)
										}

										return nil
									})
								if err != nil {
									return nil, err
								} else {
									//找到了subDir，移动文件位置并删除该文件夹
									if err := os.Rename("./temp/"+filename+"/"+subDir, "./temp/"+filename+"6657"); err != nil {
										return nil, err
									} else {
										if err := os.RemoveAll("./temp/" + filename); err != nil {
											return nil, err
										}
										if err := os.Rename("./temp/"+filename+"6657", "./temp/"+filename); err != nil {
											return nil, err
										}
									}
								}
							}

							err = os.RemoveAll("./temp/" + filename + "/" + flt)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}

			//删除原文件
			err = os.Remove("./" + cfg.DistPath + "/" + file)
			if err != nil {
				return nil, err
			}

			//重新打包，分卷19MB
			if ok == 0 {
				//对不支持的格式，需要替换变量名以适应下面的压缩操作
				filename = file                      // xxx.exe
				file = strings.TrimSuffix(file, ext) // xxx
			}
			file = filename + ".7z"
			if err = p7zip.Do7z("./temp/"+filename, cfg.DistPath+"/"+file, cfg.CompRatio, true, "19m"); err != nil {
				return nil, err
			}

			//删除临时文件夹
			err = os.RemoveAll("./temp")
			if err != nil {
				return nil, err
			}

			//检查分卷后的压缩包个数，只有一个则改名，去掉.001
			sum := 0

			for {
				//补齐到三位数字
				t := strconv.Itoa(sum + 1)
				switch len(t) {
				case 0:
					t = "000"
				case 1:
					t = "00" + t
				case 2:
					t = "0" + t
				}

				fmt.Println(file + "." + t)
				ok, err := util.IsFileExisted(cfg.DistPath + "/" + file + "." + t)
				if err != nil {
					return nil, err
				} else if ok == true {
					sum++
				} else {
					break
				}
			}

			if sum == 1 {
				err = os.Rename(cfg.DistPath+"/"+file+".001", cfg.DistPath+"/"+file)
				if err != nil {
					return nil, err
				}

				split = append(split, file)
			} else if sum > 1 {
				//把分卷后的sum个地址添加到split
				for i := 0; i < sum; i++ {
					//补齐到三位数字
					t := strconv.Itoa(i + 1)
					switch len(t) {
					case 0:
						t = "000"
					case 1:
						t = "00" + t
					case 2:
						t = "0" + t
					}

					split = append(split, file+"."+t)
				}
			} else {
				return nil, errors.New("unexpected error when split files: " + file)
			}

		} else {
			split = append(split, file)
		}
	}

	//再检查当前目录下所有文件大小之和是否小于50MB
	//if err := checkDirSize("./" + cfg.DistPath ,37); err != nil {
	//	return nil, err
	//}

	return split, nil
}

//把文件名转换成最终加速下载的链接
func File2Link(files []string, cfg config.Cfg) []string {
	var links []string
	if cfg.ArchiverGH == true {
		prefix := "https://cdn.jsdelivr.net/gh/" + cfg.ArchiverOwner + "/" + cfg.ArchiverRepo + "@master/" + cfg.DistPath + "/"
		for _, file := range files {
			links = append(links, prefix+file)
		}
	} else {
		prefix := cfg.ArchiverAPI + "/" + cfg.DistPath + "/"
		for _, file := range files {
			links = append(links, prefix+file)
		}
	}

	return links
}

func UpdateVersionList(oldList []string, newVersion string) []string {
	return append([]string{newVersion}, oldList...)
}
