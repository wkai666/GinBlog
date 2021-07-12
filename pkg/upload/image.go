package upload

import (
	"fmt"
	"ginApp/pkg/file"
	"ginApp/pkg/logging"
	"ginApp/pkg/setting"
	"ginApp/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd error: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" +src)
	if err != nil {
		return fmt.Errorf("file.IsNotExist error: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission error: %s", src)
	}

	return nil
}