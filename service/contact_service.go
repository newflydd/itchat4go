package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	e "itchat4go/enum"
	m "itchat4go/model"
	t "itchat4go/tools"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

/**
 * 获取所有联系人信息，组装到map中，key为用户的UserName
 * 微信API对此URL使用了Cookie验证
 * GET:https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket=dfLHy%252Fcgw%252BFM1qGhuARU6%252BDGs%252BGmWAD3jZJk6%252BfWcPs%253D&r=1504587952374&seq=0&skey=@crypt_3aaab8d5_c87634a7c5f8f579095cfdceeb8d842a
 */
func GetAllContact(loginMap *m.LoginMap) (map[string]m.User, error) {
	contactMap := map[string]m.User{}

	urlMap := e.GetInitParaEnum()
	urlMap[e.PassTicket] = loginMap.PassTicket
	urlMap[e.R] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
	urlMap["seq"] = "0"
	urlMap[e.SKey] = loginMap.BaseRequest.SKey

	/* 使用Cookie功能，Get数据 */
	u, _ := url.Parse("https://wx.qq.com")

	jar := new(m.Jar)
	jar.SetCookies(u, loginMap.Cookies)

	client := &http.Client{
		Jar: jar}

	resp, err := client.Get(e.GET_ALL_CONTACT_URL + t.GetURLParams(urlMap))
	if err != nil {
		return contactMap, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return contactMap, err
	}

	contactList := m.ContactList{}
	err = json.Unmarshal(bodyBytes, &contactList)
	if err != nil {
		return contactMap, err
	}

	for i := 0; i < contactList.MemberCount; i++ {
		contactMap[contactList.MemberList[i].UserName] = contactList.MemberList[i]
	}

	return contactMap, nil
}

func MapGroupInfo(contactMap map[string]m.User) map[string][]m.User {
	groupMap := map[string][]m.User{}

	for _, user := range contactMap {
		if strings.HasPrefix(user.UserName, "@@") {
			/* 如果该联系人是一个群组，依次判断是否需要加入焦点群列表 */
			for _, key := range e.GetFocusGroupKeywordChildren() {
				reg := regexp.MustCompile(key)
				if reg.MatchString(user.UserName) || reg.MatchString(user.NickName) {
					groupMap[strings.ToLower(key)] = append(groupMap[strings.ToLower(key)], user)
				}
			}
		}
	}

	return groupMap
}
