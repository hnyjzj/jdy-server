package common

import (
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
		return nil, errors.New("file is nil")
	}

	if up.File != nil {
		up.Files = append(up.Files, up.File)
	}

	if up.Files == nil {
		return nil, errors.New("file is nil")
	}

	switch up.conf.Type {
	case config.StorageTypeLocal:
		err := up.saveLocal()
		if err != nil {
			return nil, err
		}
		return up, nil
	default:
		return nil, errors.New("storage type error")
	}
}

func (up *Upload) saveLocal() error {
	for _, file := range up.Files {
		if file == nil {
			return errors.New("file is nil")
		}
		contentType := up.File.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, up.Type.String()) { // 判断文件类型
			log.Println(contentType, up.Type.String(), up.File)
			return errors.New("file type error")
		}
		// 获取文件扩展名
		var name string
		if up.Prefix != "" {
			name = "u" + up.Prefix + "_" + utils.RandomAlphanumeric(6) + filepath.Ext(up.File.Filename)
		} else {
			name = utils.GetCurrentMilliseconds() + "_" + utils.RandomAlphanumeric(6) + filepath.Ext(up.File.Filename)
		}
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
		if err := up.Ctx.SaveUploadedFile(up.File, pwd); err != nil {
			return err
		}
		// 文件访问路径
		up.Uris = append(up.Uris, filepath.Join(up.conf.Prefix, path))
	}

	return nil
}
