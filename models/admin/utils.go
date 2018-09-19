package admin

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/rand"
	"regexp"
	"time"
)

// IsPhone 11位数字
func IsPhone(str string) bool {
	matchPhone := regexp.MustCompile(`^\d{11}$`)
	return matchPhone.MatchString(str)
}

// HashPassword HashPassword
func HashPassword(pass, salt string) string {
	h := sha1.New()
	h.Write([]byte(pass + salt))
	checksum := h.Sum(nil)
	passHashed := base64.StdEncoding.EncodeToString([]byte(string(checksum) + salt))
	return passHashed
}

// Md5 Md5Sign
func Md5(src string) string {
	h := md5.New()
	io.WriteString(h, src)
	return hex.EncodeToString(h.Sum(nil))
}

// IsEmail 只判断是否有@符号
func IsEmail(src string) bool {
	matchPhone := regexp.MustCompile(`^.+@.+$`)
	return matchPhone.MatchString(src)
}

// IsValidNickname 用户名是否有效
func IsValidNickname(name string) bool {
	b := []byte(name)

	if len(b) < 4 || len(b) > 25 {
		return false
	}

	m := regexp.MustCompile(`^[-\p{Han}\w]+$`)
	return m.MatchString(name)
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// Krand 随机字符串
func Krand(size int, rands ...int) []byte {
	kind := KC_RAND_KIND_ALL
	if len(rands) > 0 {
		kind = rands[0]
	}
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

// RandString 随机字符串
func RandString(size int) string {
	return string(Krand(size, KC_RAND_KIND_ALL))
}
