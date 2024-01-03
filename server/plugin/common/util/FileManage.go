package util

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"server/config"
)

/*
数据请求保存,文件读写
*/

// SaveOnlineFile 保存网络文件, 提供下载url和保存路径, 返回保存后的文件访问url相对路径
func SaveOnlineFile(url, dir string) (path string, err error) {
	// 请求获取文件内容
	r := &RequestInfo{Uri: url}
	ApiGet(r)
	// 如果请求结果为空则直接跳过当前图片的同步, 等待后续触发时重试
	if len(r.Resp) <= 0 {
		err = errors.New("SyncPicture Failed: response is empty")
		return
	}
	// 成功拿到图片数据 则创建保存文件的目录
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}
	// 通过保存路径和url得到保存的具体的文件全路径
	fileName := filepath.Join(dir, filepath.Base(url))
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	// 将文件内容写入到file
	writer := bufio.NewWriter(file)
	_, err = writer.Write(r.Resp)
	err = writer.Flush()
	return filepath.Base(fileName), err

}

func CreateBaseDir() error {
	// 如果不存在指定目录则创建该目录
	if _, err := os.Stat(config.FilmPictureUploadDir); os.IsNotExist(err) {
		return os.MkdirAll(config.FilmPictureUploadDir, os.ModePerm)
	}
	return nil
}

func RemoveFile(path string) error {
	err := os.Remove(path)
	return err
}
