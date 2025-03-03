package img_service

import "ginchat/conf"

type ImgUploadInterface interface {
	Upload(filename string) string
}

var serveMap = map[string]ImgUploadInterface{
	"freeimg": &FreeImageService{},
	// "sm": &sm_app.SmAppService{},
}

func ImgCreate() ImgUploadInterface {
	return serveMap[conf.GlobalConf.PictureBed.Type]
}
