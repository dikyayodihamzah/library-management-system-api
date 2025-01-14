package lib

import (
	"strings"

	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
)

func AssignMinioPrefix(image string, customPrefix ...string) string {
	if IsEmptyString(image) {
		return image
	}

	prefix := utils.GetString("MINIO_PREFIX")
	bucket := utils.GetString("MINIO_BUCKET")

	if len(customPrefix) == 2 {
		prefix = customPrefix[0]
		bucket = customPrefix[1]
	} else if len(customPrefix) != 0 && len(customPrefix) != 2 {
		return image
	}

	if !strings.Contains(image, "http") {
		return AppendStr("/", prefix, bucket, image)
	}
	return image
}

func TrimMinioPrefix(image string, customPrefix ...string) string {
	if IsEmptyString(image) {
		return image
	}

	prefix := utils.GetString("MINIO_PREFIX")
	bucket := utils.GetString("MINIO_BUCKET")

	if len(customPrefix) == 2 {
		prefix = customPrefix[0]
		bucket = customPrefix[1]
	}

	// check if image has minio prefix
	if strings.HasPrefix(image, AppendStr("/", prefix, bucket)) {
		return strings.TrimPrefix(image, AppendStr("/", prefix, bucket)+"/")
	}

	return image
}
