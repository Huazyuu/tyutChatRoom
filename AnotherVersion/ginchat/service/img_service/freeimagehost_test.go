package img_service

import (
	"ginchat/conf"
	"testing"
)

func Test_freeImgUpload(t *testing.T) {
	conf.InitCore("../../conf/settings.yaml")
	// 替换为实际的本地文件路径
	uploadFile := "D:\\goProject\\ginchat\\uploads\\屏幕截图_20241229_212523.png"
	url := freeImgUpload(uploadFile)
	if url != "" {
		t.Logf("上传成功，图片显示 URL: %s\n", url)
	} else {
		t.Fatal("上传失败")
		return
	}
}
