package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
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
func ParsePriKeyBytes(buf []byte) (any, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("private key parse  error")
	}
	return x509.ParsePKCS8PrivateKey(p.Bytes)
}

// ParsePubKeyBytes 解析公钥
func ParsePubKeyBytes(buf []byte) (any, error) {
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse publicKey content nil")
	}
	pubKey, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return nil, errors.New("x509.ParsePKIXPublicKey error")
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

// ValidPwd 校验密码
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

// TruncateBySep 截断字符串,保留指定数量的结果
func TruncateBySep(s string, limit int) string {
	// 如果保留数量小于等于0则返回空值
	if len(s) <= 0 || limit <= 0 {
		return ""
	}
	// 先强制对不同的分割符进行统一替换为 ,
	s = regexp.MustCompile(`[$&#%]`).ReplaceAllString(s, ",")
	// 使用 strings.Split 分割字符串
	// Split 会在分隔符连续出现或出现在首尾时产生空字符串，这通常符合预期
	parts := strings.Split(s, ",")
	// 片段数量小于或等于限制，直接返回原字符串
	if len(parts) <= limit {
		return strings.Join(parts, ",")
	}
	// 返回原字符串是为了保留原始的格式（比如末尾是否有分隔符）
	// 即使不截断也重新 Join 一遍（去除多余的空片段等)
	return strings.Join(parts[:limit], ",")
}

// CleanFilmName 清洗影片名称，只保留主体
func CleanFilmName(name string) string {
	if name == "" {
		return ""
	}
	// 1. 去除常见的前缀 (方括号、圆括号内的内容) 匹配 [xxx], 【xxx】, (xxx), （xxx）
	//rePrefix := regexp.MustCompile(`^\s*[\[【\(（][^\]】\)）]*[\]】\)）]\s*`)
	//for rePrefix.MatchString(name) {
	//	name = rePrefix.ReplaceAllString(name, "")
	//}
	// 2.定义需要清洗的特殊标识关键字集合
	var noisePatterns = []string{
		`第 [零一二三四五六七八九十\d]+ 季`, `第 [零一二三四五六七八九十\d]+ 话`, `第 [零一二三四五六七八九十\d]+ 集`,
		`Season\s*\d+`, `S\d+`, `Ep\d+`, `\d{1,3}\s*(话 | 集)`,
		`\s+(II|III|IV|V|VI|VII|VIII|IX|X)\s*$`,
		`剧场版`, `电影版`, `OVA`, `OAD`, `SP`, `特别篇`, `总集篇`, `外传`, `序`, `破`, `急`, `终章`,
		`\d{3,4}[Pp]`, `HD`, `FHD`, `UHD`, `4K`, `BD`, `BluRay`, `BDRip`, `HEVC`, `H264`, `H265`,
		`GB`, `MB`, `MP4`, `MKV`, `AVI`, `RMVB`,
		`字幕组`, `动漫`, `动画`, `新版`, `重制版`, `连载`, `更新`, `全集`, `合集`,
		`Uncensored`, `NoCen`, `Dubbed`, `Subbed`, `Raw`, `生肉`, `熟肉`,
	}
	// 3. 处理拼接完整的正则表达式
	fullPattern := `(?i)(?:\s+|\.+|_+|-+) (` + strings.Join(noisePatterns, "|") + `).*$`
	cutRegex := regexp.MustCompile(fullPattern)
	// 去除满足匹配集的子串
	name = cutRegex.ReplaceAllString(name, "")
	// 特殊处理 "之" 字结构 (仅当 "之" 后紧跟噪音词时切除)	之\s*(噪音词)
	reRegex := regexp.MustCompile(`(?i) 之\s* (` + strings.Join(noisePatterns, "|") + `).*$`)
	name = reRegex.ReplaceAllString(name, "")

	// 修剪 - 去除末尾残留的符号和空白
	name = strings.TrimRight(name, " \t\n\r._-])：:")

	return name
}

// FormatSpecialChar 格式化特殊字符, 统一替换为逗号
func FormatSpecialChar(src string) string {
	// 执行替换
	return strings.Map(func(r rune) rune {
		switch r {
		case '#', '/', '$', '&', '%', '^', '*', '-':
			return ','
		default:
			return r
		}
	}, src)
}
