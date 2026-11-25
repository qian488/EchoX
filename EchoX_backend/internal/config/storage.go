package config

import (
	"crypto/rand"
	"encoding/hex"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// UploadDir 为本地媒体存储的根目录，相对于服务运行时的工作目录。
const UploadDir = "./uploads"

// MaxUploadSize 为单个上传文件的大小上限：100MB。
const MaxUploadSize int64 = 100 << 20

var imageMimeTypes = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".webp": "image/webp",
}

var videoMimeTypes = map[string]string{
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mov":  "video/quicktime",
	".avi":  "video/x-msvideo",
}

// GetUploadDir 返回媒体存储目录（零配置，无需 .env 变量）。
func GetUploadDir() string {
	return UploadDir
}

// EnsureUploadDir 确保存储目录存在，运行时按需创建。
func EnsureUploadDir() error {
	return os.MkdirAll(GetUploadDir(), 0755)
}

// normalizeExt 将扩展名统一为小写并补齐前导点。
func normalizeExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

// IsAllowedExtension 判断扩展名是否在图片/视频白名单内。
func IsAllowedExtension(ext string) bool {
	ext = normalizeExt(ext)
	if _, ok := imageMimeTypes[ext]; ok {
		return true
	}
	_, ok := videoMimeTypes[ext]
	return ok
}

// GetMimeType 根据扩展名返回 MIME 类型，优先使用内置映射，再回退到标准库探测。
// 不使用客户端上传时声明的 Content-Type，避免被伪造。
func GetMimeType(ext string) string {
	ext = normalizeExt(ext)
	if mt, ok := imageMimeTypes[ext]; ok {
		return mt
	}
	if mt, ok := videoMimeTypes[ext]; ok {
		return mt
	}
	if mt := mime.TypeByExtension(ext); mt != "" {
		return mt
	}
	return "application/octet-stream"
}

// GenerateFilename 生成形如 20261125_abc123.png 的存储文件名。
// 原始文件名仅用于提取白名单内的扩展名，路径部分与其余字符一律丢弃，
// 因此不会产生路径穿越，也不会覆盖已有文件。
func GenerateFilename(orig string) string {
	ext := normalizeExt(filepath.Ext(orig))
	if !IsAllowedExtension(ext) {
		ext = ""
	}
	return time.Now().Format("20060102") + "_" + randomToken(3) + ext
}

// randomToken 返回 n 字节随机数据的十六进制字符串；极端情况下回退到时间戳。
func randomToken(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(buf)
}
