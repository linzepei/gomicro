package utils

import (
	"crypto/md5"
	"encoding/hex"
)

/* 将url加上 http://IP:PROT/  前缀 */
//http:// + 127.0.0.1 + ：+ 8080 + 请求
func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_img_dns + "/" + url

	return domain_url
}

func Md5String(s string) string {
	//创建1个md5对象
	h := md5.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}
