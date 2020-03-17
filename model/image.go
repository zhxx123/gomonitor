package model

// Image 图片
type Image struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	Title        string `json:"title"`
	OrignalTitle string `json:"orignalTitle"`
	URL          string `json:"url"`
	Width        uint   `json:"width"`
	Height       uint   `json:"height"`
	Mime         string `json:"mime"`
}

// ImageUploadedInfo 图片上传后的相关信息(目录、文件路径、文件名、UUIDName、请求URL)
type ImageUploadedInfo struct {
	UploadDir      string
	UploadFilePath string
	Filename       string
	UUIDName       string
	ImgURL         string
}
