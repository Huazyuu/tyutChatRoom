package fileComm

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	UserID   string `json:"user_id"`
	TargetID string `json:"target_id"`
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	FileSize int    `json:"file_size"`
	FileType string `json:"file_type"`
}
