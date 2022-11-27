package utils

import (
	//"github.com/pkg/errors"
)

func MessageMap(status bool, message string) map[string]interface{} {
	str := "error"
	if status {
		str = "success"
	}
	return map[string]interface{}{"status": str, "message": message}
}