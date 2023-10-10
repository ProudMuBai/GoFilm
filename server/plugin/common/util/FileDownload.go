package util

import (
	"bufio"
	"os"
	"path/filepath"
)

/*
数据请求保存,文件读写
*/

// SaveOnlineFile 保存网络文件, 提供下载url和保存路径, 返回保存后的文件访问url相对路径
func SaveOnlineFile(url, dir string) (err error) {
	// 请求获取文件内容
	r := &RequestInfo{Uri: url}
	ApiGet(r)
	// 创建保存文件的目录
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
	//_, _ = file.Write(r.Resp)

	writer := bufio.NewWriter(file)
	_, err = writer.Write(r.Resp)
	err = writer.Flush()
	return

}
