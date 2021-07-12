package export

import (
	"ginApp/pkg/file"
	"ginApp/pkg/setting"
)

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	excelPath := setting.AppSetting.RuntimeRootPath + GetExcelPath()
	if  file.CheckNotExist(excelPath) {
		file.IsNotExistMkDir(excelPath)
	}
	return excelPath
}
