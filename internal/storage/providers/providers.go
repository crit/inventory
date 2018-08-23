package providers

import (
	"path/filepath"

	"github.com/tuhlz/growtv/src/internal/storage"
	"github.com/tuhlz/growtv/src/internal/storage/providers/bolt"
	"github.com/tuhlz/growtv/src/internal/storage/providers/dynamodb"
	"github.com/tuhlz/growtv/src/internal/storage/providers/mock"
	"github.com/tuhlz/growtv/src/internal/storage/providers/s3"
)

func Mock(reporter func(db *mock.DB)) storage.Storage {
	return mock.New(reporter)
}

func DynamoDB(region, table string) (storage.Storage, error) {
	return dynamodb.New(region, table)
}

func Bolt(path string) (storage.Storage, error) {
	return bolt.New(path)
}

type S3Config *s3.Config

func S3(config S3Config) storage.Storage {
	return s3.New(config)
}

func FileDecoder(key string) (contentType string) {
	found := filepath.Ext(key)

	switch found {
	case ".aac":
		return "audio/aac"
	case ".abw":
		return "application/x-abiword"
	case ".arc":
		return "application/octet-stream"
	case ".avi":
		return "video/x-msvideo"
	case ".azw":
		return "application/vnd.amazon.ebook"
	case ".bin":
		return "application/octet-stream"
	case ".bmp":
		return "image/bmp"
	case ".bz":
		return "application/x-bzip"
	case ".bz2":
		return "application/x-bzip2"
	case ".csh":
		return "application/x-csh"
	case ".css":
		return "text/css"
	case ".csv":
		return "text/csv"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".epub":
		return "application/epub+zip"
	case ".es":
		return "application/ecmascript"
	case ".gif":
		return "image/gif"
	case ".htm":
		return "text/html"
	case ".html":
		return "text/html"
	case ".ico":
		return "image/x-icon"
	case ".ics":
		return "text/calendar"
	case ".jar":
		return "application/java-archive"
	case ".jpeg":
		return "image/jpeg"
	case ".jpg":
		return "image/jpeg"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".mid":
		return "audio/midi"
	case ".midi":
		return "audio/midi"
	case ".mpeg":
		return "video/mpeg"
	case ".mpkg":
		return "application/vnd.apple.installer+xml"
	case ".odp":
		return "application/vnd.oasis.opendocument.presentation"
	case ".ods":
		return "application/vnd.oasis.opendocument.spreadsheet"
	case ".odt":
		return "application/vnd.oasis.opendocument.text"
	case ".oga":
		return "audio/ogg"
	case ".ogv":
		return "video/ogg"
	case ".ogx":
		return "application/ogg"
	case ".otf":
		return "font/otf"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".rar":
		return "application/x-rar-compressed"
	case ".rtf":
		return "application/rtf"
	case ".sh":
		return "application/x-sh"
	case ".svg":
		return "image/svg+xml"
	case ".swf":
		return "application/x-shockwave-flash"
	case ".tar":
		return "application/x-tar"
	case ".tif":
		return "image/tiff"
	case ".tiff":
		return "image/tiff"
	case ".ts":
		return "application/typescript"
	case ".ttf":
		return "font/ttf"
	case ".txt":
		return "text/plain"
	case ".vsd":
		return "application/vnd.visio"
	case ".wav":
		return "audio/wav"
	case ".weba":
		return "audio/webm"
	case ".webm":
		return "video/webm"
	case ".webp":
		return "image/webp"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".xhtml":
		return "application/xhtml+xml"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xml":
		return "application/xml"
	case ".xul":
		return "application/vnd.mozilla.xul+xml"
	case ".zip":
		return "application/zip"
	case ".7z":
		return "application/x-7z-compressed"
	}

	return "application/octet-stream"
}
