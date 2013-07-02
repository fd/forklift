package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"os/user"
	"path"
)

func Path(ref string) (string, error) {
	sha := sha1.New()
	sha.Write([]byte(ref))

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(u.HomeDir, ".forklift", "deploypacks", hex.EncodeToString(sha.Sum(nil))), nil
}
