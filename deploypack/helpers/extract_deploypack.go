package helpers

import (
	"fmt"
)

func ExtractDeploypack(config map[string]interface{}) (string, error) {
	if vi, ok := config["deploypack"]; ok && vi != nil {
		if vs, ok := vi.(string); ok {
			delete(config, "deploypack")
			return vs, nil
		} else {
			return "", fmt.Errorf("invalid confirguration: deploypack must be a string")
		}
	} else {
		return "", nil
	}
}
