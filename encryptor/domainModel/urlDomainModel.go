package domainModel

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/G4GANCHAUDHARY/encryption-service/global"
	pq "github.com/lib/pq"
)

func GetUrlFilterString(longUrl string) string {
	return "long_url = '" + longUrl + "'"
}

func GetShortUrlFromUniqKey(key int64) string {
	// using SHA1 for creating uniq strings
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%d", key)))
	return hex.EncodeToString(h.Sum(nil))[:9]
}

func IsUniqueConstraintError(err error) bool {
	// unique_violation
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == global.UniqueConstraintCode
	}
	return false
}
