package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"server/config"
	"server/logic"
	"server/model/system"
	"server/plugin/common/util"
	"strconv"
)

// SingleUpload 单文件上传, 暂定为图片上传
func SingleUpload(c *gin.Context) {
	// 获取执行操作的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("上传失败, 当前用户信息异常", c)
		return
	}
	// 结合搜文件内容
	file, err := c.FormFile("file")
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}

	// 生成文件名, 保存文件到服务器
	fileName := fmt.Sprintf("%s/%s%s", config.FilmPictureUploadDir, util.RandomString(8), filepath.Ext(file.Filename))
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}

	uc := v.(*system.UserClaims)
	// 记录图片信息到系统表中, 并获取返回的图片访问路径
	link := logic.FileL.SingleFileUpload(fileName, int(uc.UserID))
	// 返回图片访问地址以及成功的响应
	system.Success(link, "上传成功", c)

}

// MultipleUpload 批量文件上传
func MultipleUpload(c *gin.Context) {
	// 获取执行操作的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("上传失败, 当前用户信息异常", c)
		return
	}
	// 解析表单数据
	form, err := c.MultipartForm()
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 获取文件列表
	files := form.File["files"]

	// 解析当前登录的用户信息
	uc := v.(*system.UserClaims)

	// 遍历文件列表
	var fileNames []string
	for _, file := range files {
		// 生成文件名, 保存文件到服务器
		fileName := fmt.Sprintf("%s/%s%s", config.FilmPictureUploadDir, util.RandomString(8), filepath.Ext(file.Filename))
		err = c.SaveUploadedFile(file, fileName)
		if err != nil {
			system.Failed(err.Error(), c)
			return
		}
		// 记录图片信息到系统表中, 并获取返回的图片访问路径
		fileNames = append(fileNames, logic.FileL.SingleFileUpload(fileName, int(uc.UserID)))
	}

	// 返回图片访问地址以及成功的响应
	system.Success(fileNames, "上传成功", c)
}

// DelFile 删除文件
func DelFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.DefaultQuery("id", ""), 10, 64)
	if err != nil {
		system.Failed("操作失败, 未获取到需删除的文件标识信息", c)
		return
	}
	if e := logic.FileL.RemoveFileById(uint(id)); e != nil {
		system.Failed(fmt.Sprint("删除失败", e.Error()), c)
		return
	}
	system.SuccessOnlyMsg("文件已删除", c)
}

// PhotoWall 照片墙数据
func PhotoWall(c *gin.Context) {
	current, err := strconv.Atoi(c.DefaultQuery("current", "1"))
	if err != nil {
		system.Failed("图片分页数据获取失败, 分页参数异常", c)
		return
	}
	// 获取系统保存的文件的图片分页数据
	page := system.Page{PageSize: 39, Current: current}
	// 获取分页数据
	pl := logic.FileL.GetPhotoPage(&page)
	system.Success(gin.H{"list": pl, "page": page}, "图片分页数据获取成功", c)
}
