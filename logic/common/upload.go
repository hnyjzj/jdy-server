package common

import (
	"context"
	"errors"
	"jdy/config"
	"jdy/types"
	"jdy/utils"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Upload struct {
	conf *config.Storage

	Ctx   *gin.Context
	File  *multipart.FileHeader
	Files []*multipart.FileHeader

	Model  types.UploadModel
	Type   types.UploadType
	Prefix string

	Uris []string
}

func (up *Upload) Save() (*Upload, error) {

	up.conf = &config.Config.Storage

	if up.File == nil && up.Files == nil {
		return nil, errors.New("文件不存在")
	}

	if up.File != nil {
		up.Files = append(up.Files, up.File)
	}

	if up.Files == nil {
		return nil, errors.New("文件不存在")
	}

	switch up.conf.Type {
	case config.StorageTypeLocal:
		// 本地存储
		if err := up.saveLocal(); err != nil {
			return nil, err
		}
		return up, nil
	case config.StorageTypeTencentCloud:
		// 腾讯云存储
		if err := up.saveTencentCloud(); err != nil {
			return nil, err
		}
		return up, nil
	default:
		return nil, errors.New("存储类型错误")
	}
}

func (up *Upload) saveLocal() error {
	for _, file := range up.Files {
		if file == nil {
			return errors.New("文件不存在")
		}
		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, up.Type.String()) { // 判断文件类型
			log.Println(contentType, up.Type.String(), file)
			return errors.New("文件类型错误")
		}

		// 生成文件名
		name := up.getName(file)
		// 生成文件路径
		path := filepath.Join(up.Model.String(), name)
		// 生成文件保存路径
		pwd := filepath.Join(up.conf.Local.Root, path)
		// 检查路径是否合法
		cleanPath := filepath.Clean(pwd)
		if !strings.HasPrefix(cleanPath, filepath.Clean(up.conf.Local.Root)) {
			return errors.New("invalid path")
		}
		// 保存文件
		if err := up.Ctx.SaveUploadedFile(file, pwd); err != nil {
			return err
		}
		// 文件访问路径
		up.Uris = append(up.Uris, filepath.Join(up.conf.Prefix, path))
	}

	return nil
}

func (up *Upload) saveTencentCloud() error {
	for _, file := range up.Files {
		if file == nil {
			return errors.New("file is nil")
		}
		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, up.Type.String()) { // 判断文件类型
			log.Println(contentType, up.Type.String(), file)
			return errors.New("file type error")
		}

		// 生成文件名
		name := up.getName(file)
		// 生成文件路径
		path := filepath.Join(up.Model.String(), name)
		// 生成文件保存路径
		pwd := filepath.Join(up.conf.TencentCloud.Root, path)

		// 打开上传的文件
		fd, err := file.Open()
		if err != nil {
			return err
		}
		defer fd.Close()

		// 初始化客户端
		client := config.NewTencentCloudStorage()

		_, err = client.Object.Put(context.Background(), pwd, fd, nil)
		if err != nil {
			return err
		}
		// 文件访问路径
		up.Uris = append(up.Uris, pwd)
	}
	return nil
}

// 生成文件名
func (up *Upload) getName(file *multipart.FileHeader) string {
	// 生成文件名
	var name string
	// 获取前缀
	if up.Prefix != "" {
		name = "u" + up.Prefix + "_" + utils.RandomAlphanumeric(6) + filepath.Ext(file.Filename)
	} else {
		name = utils.GetCurrentMilliseconds() + "_" + utils.RandomAlphanumeric(6) + filepath.Ext(file.Filename)
	}

	return name
}
