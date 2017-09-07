package model

import (
	"net/http"
	"net/url"
)

/* 获取联系人列表时需要带入Cookie信息，实现CookieJar接口 */
type Jar struct {
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}

type LoginMap struct {
	PassTicket  string
	BaseRequest BaseRequest /* 将涉及登陆有关的验证数据封装成对象 */

	SelfNickName string
	SelfUserName string

	SyncKeys   SyncKeysJsonData /* 同步消息时需要验证的Keys */
	SyncKeyStr string           /* Keys组装成的字符串 */

	Cookies []*http.Cookie /* 微信相关API需要用到的Cookies */
}
