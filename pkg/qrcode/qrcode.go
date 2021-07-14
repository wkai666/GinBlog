package qrcode

import (
	"ginApp/pkg/file"
	"ginApp/pkg/setting"
	"ginApp/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
)

type QrCode struct{
	URL string
	Width int
	Height int
	Ext string
	Level qr.ErrorCorrectionLevel
	Mode qr.Encoding
}

const EXT_JPG = ".jpg"

func NewQrCode(url string, wight, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL: url,
		Width: wight,
		Height: height,
		Level: level,
		Mode: mode,
		Ext: EXT_JPG,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetQrCodePath()
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

func (q *QrCode ) GetCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckCode(path string) bool {
	src := path + GetQrCodeFullPath() + q.GetCodeExt()
	if true == file.CheckNotExist(src) {
		return false
	}

	return true
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetCodeExt()
	src := path + name
	if file.CheckNotExist(src) {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}

		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}