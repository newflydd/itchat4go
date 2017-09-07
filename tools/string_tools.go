package tools

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

/**
 * 有序(或者无序)地从一个map中按照index的顺序构造URL中的params
 * 加上有序的目的是为了防止有些环境下需要params根据key的ASC大小排序后进行签名加密
 */
func GetURLParams(values ...interface{}) string {
	var result = "?"
	if len(values) == 1 {
		maap := values[0].(map[string]string)
		for key, value := range maap {
			if key != "" && value != "" {
				result += fmt.Sprintf("%s=%s&", key, url.QueryEscape(value))
			}
		}
	} else if len(values) == 2 {
		index := values[1].([]string)
		maap := values[0].(map[string]string)
		for _, key := range index {
			if key != "" && maap[key] != "" {
				result += fmt.Sprintf("%s=%s&", key, url.QueryEscape(maap[key]))
			}
		}
	}

	return result[:len(result)-1]
}

/**
 *  生成随机字符串
 *  index：取随机序列的前index个
 *  0-9:10
 *  0-9a-z:10+24
 *  0-9a-zA-Z:10+24+24
 *  length：需要生成随机字符串的长度
 */
func GetRandomString(index int, length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(index)])
	}
	return string(result)
}
