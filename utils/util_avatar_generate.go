package utils

import (
	"github.com/disintegration/letteravatar"
	"github.com/golang/freetype"
	"image/png"
	"os"
	"path"
	"unicode/utf8"
)

func DrawImage(name string) (string, error) {
	dir := "uploads/avatar"
	// 创建目录
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	fontFile, err := os.ReadFile("./uploads/font/STHUPO.TTF")
	if err != nil {
		return "", err
	}

	font, err := freetype.ParseFont(fontFile)
	if err != nil {
		return "", err
	}

	options := &letteravatar.Options{
		Font: font,
	}

	// 绘制文字
	firstLetter, _ := utf8.DecodeRuneInString(name)
	img, err := letteravatar.Draw(140, firstLetter, options)
	if err != nil {
		return "", err
	}

	// 保存
	filePath := path.Join(dir, name+".png")
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
