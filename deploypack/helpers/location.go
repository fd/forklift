package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"path"

	"github.com/fd/forklift/util/user"
)

func Path(ref string) (string, error) {
	sha := sha1.New()
	sha.Write([]byte(ref))

	home, err := user.Home()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".forklift", "deploypacks", hex.EncodeToString(sha.Sum(nil))), nil
}
