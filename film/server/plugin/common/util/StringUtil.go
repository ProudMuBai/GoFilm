package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
)

// GenerateUUID 生成UUID
func GenerateUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

// RandomString 生成指定长度两倍的随机字符串
func RandomString(length int) (uuid string) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid = fmt.Sprintf("%x", b)
	return
}

// GenerateSalt 生成 length为16的随机字符串
func GenerateSalt() (uuid string) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid = fmt.Sprintf("%X", b)
	return
}

// PasswordEncrypt 密码加密 , (password+salt) md5 * 3
func PasswordEncrypt(password, salt string) string {
	b := []byte(fmt.Sprint(password, salt)) // 将字符串转换为字节切片
	var r [16]byte
	for i := 0; i < 3; i++ {
		r = md5.Sum(b) // 调用md5.Sum()函数进行加密
		b = []byte(hex.EncodeToString(r[:]))
	}
	return hex.EncodeToString(r[:])
}

// ParsePriKeyBytes 解析私钥
func ParsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("private key parse  error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

// ParsePubKeyBytes 解析公钥
func ParsePubKeyBytes(buf []byte) (*rsa.PublicKey, error) {
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse publicKey content nil")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(p.Bytes)
	if err != nil {
		return nil, errors.New("x509.ParsePKCS1PublicKey error")
	}
	return pubKey, nil
}

// ValidDomain 域名校验(http://example.xxx)
func ValidDomain(s string) bool {
	return regexp.MustCompile(`^(http|https)://[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*\.[a-z]{2,6}(:[0-9]{1,5})?$`).MatchString(s)
}

// ValidIPHost 校验是否符合http|https//ip 格式
func ValidIPHost(s string) bool {
	return regexp.MustCompile(`^(http|https)://(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(:[0-9]{1,5})?$`).MatchString(s)
}

// ValidURL 校验http链接是否是符合规范的URL
func ValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return true
}

func ValidPwd(s string) error {
	if len(s) < 8 || len(s) > 12 {
		return fmt.Errorf("密码长度不符合规范, 必须为8-10位")
	}
	// 分别校验数字 大小写字母和特殊字符
	num := `[0-9]{1}`
	l := `[a-z]{1}`
	u := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, s); !b || err != nil {
		return errors.New("密码必须包含数字 ")
	}
	if b, err := regexp.MatchString(l, s); !b || err != nil {
		return errors.New("密码必须包含小写字母")
	}
	if b, err := regexp.MatchString(u, s); !b || err != nil {
		return errors.New("密码必须包含大写字母")
	}
	if b, err := regexp.MatchString(symbol, s); !b || err != nil {
		return errors.New("密码必须包含特殊字")
	}
	return nil
}
