package chatComm

// File 代表一个文件对象
type File struct {
	Path string `json:"path"`
	Name string `json:"name"` // 文件名
	Type string `json:"type"` // 文件类型，例如 "image/jpeg", "application/pdf" 等
	Size int64  `json:"size"` // 文件大小（字节数）
}
