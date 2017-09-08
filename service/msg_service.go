package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	e "itchat4go/enum" /* 取个别名 */
	m "itchat4go/model"
	t "itchat4go/tools"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
 * 微信心跳检查 -> 根据retcode和selector 逻辑判断是否需要拉取数据
 * Get:https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/synccheck?
 * 	r=1504582260660&
 * 	skey=%40crypt_3aaab8d5_4f38b712dd654bff65bef8406d020771&
 * 	sid=k4FG05I6vi26YyeW&uin=154158775&
 * 	deviceid=e637087162023800&
 * 	synckey=1_696275959%7C2_696276056%7C3_696275995%7C11_696275997%7C13_696160001%7C201_1504582232%7C203_1504580411%7C1000_1504572121%7C1001_1504572152&
 * 	_=1504582198153
 * 这里依然要用到Cookies
 *
 * 根据 https://github.com/yaphone/itchat4j/blob/master/src/main/java/cn/zhouyafeng/itchat4j/utils/enums/RetCodeEnum.java
 * retcode：
 * 	NORMAL("0", "普通"),
 *	LOGIN_OUT("1102", "退出"),
 *	LOGIN_OTHERWHERE("1101", "其它地方登陆"),
 *	MOBILE_LOGIN_OUT("1102", "移动端退出"),
 *	UNKOWN("9999", "未知")
 *
 * 自己分析下来
 * selector：
 * 	4:有某位联系人更新了个人信息
 * 	2:有人发消息了 拉取后，MsgType = 1为文字消息，提取content得到聊天文本
 * 	每次拉取，需要更新SyncKey，见WebWxSync函数
 **/
func SyncCheck(loginMap *m.LoginMap) (int64, int64, error) {
	urlMap := map[string]string{}
	urlMap[e.R] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
	urlMap[e.SKey] = loginMap.BaseRequest.SKey
	urlMap[e.Sid] = loginMap.BaseRequest.Sid
	urlMap[e.Uin] = loginMap.BaseRequest.Uin
	urlMap[e.DeviceId] = loginMap.BaseRequest.DeviceID
	urlMap[e.SyncKey] = loginMap.SyncKeyStr
	urlMap[e.TimeStamp] = urlMap[e.R]

	u, _ := url.Parse("https://wx.qq.com")
	timeout := time.Duration(30 * time.Second)

	jar := new(m.Jar)
	jar.SetCookies(u, loginMap.Cookies)

	client := &http.Client{
		Jar:     jar,
		Timeout: timeout}

	resp, err := client.Get(e.SYNC_CHECK_URL + t.GetURLParams(urlMap))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	/* 根据正则得到selector => window.synccheck={retcode:"0",selector:"0"}*/
	reg := regexp.MustCompile(`^window.synccheck={retcode:"(\d+)",selector:"(\d+)"}$`)
	matches := reg.FindStringSubmatch(string(respBytes))

	retcode, err := strconv.ParseInt(matches[1], 10, 64) /* 取第二个数据为retcode值 */
	if err != nil {
		return 0, 0, errors.New("解析微信心跳数据失败:\n" + err.Error() + "\n" + string(respBytes))
	}

	selector, err := strconv.ParseInt(matches[2], 10, 64) /* 取第三个数据为selector值 */
	if err != nil {
		return 0, 0, errors.New("解析微信心跳数据失败:\n" + err.Error() + "\n" + string(respBytes))
	}

	if retcode != 0 {
		return retcode, selector, errors.New("retcode异常，程序将退出")
	}

	return retcode, selector, nil
}

/**
 * 微信同步拉取消息，最主要的消息响应逻辑
 * POST:https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync?
 * 		sid=u6j1MoqYJgbI66kO&
 * 		skey=@crypt_3aaab8d5_1f2b5690081714b78fc9b70ee3e4538c&
 * 		pass_ticket=M3%252BbNp65pZ53GJwvd7aq%252FQuS5Rd7lctHU0qilHs7Fjw%253D
 * BODY:{
 * 	"BaseRequest":{"Uin":154158775,"Sid":"u6j1MoqYJgbI66kO","Skey":"@crypt_3aaab8d5_1f2b5690081714b78fc9b70ee3e4538c","DeviceID":"e215003169713295"},
 * 	"SyncKey":{"Count":8,"List":[{"Key":1,"Val":696275959},{"Key":2,"Val":696276275},{"Key":3,"Val":696276260},{"Key":11,"Val":696275997},{"Key":13,"Val":696160001},{"Key":201,"Val":1504604303},{"Key":1000,"Val":1504603672},{"Key":1001,"Val":1504572152}]},
 * 	"rr":-1365833882
 * }
 * @param {[type]} loginMap m.LoginMap [description]
 */
func WebWxSync(loginMap *m.LoginMap) (m.WxRecvMsges, error) {
	wxMsges := m.WxRecvMsges{}

	urlMap := map[string]string{}
	urlMap[e.Sid] = loginMap.BaseRequest.Sid
	urlMap[e.SKey] = loginMap.BaseRequest.SKey
	urlMap[e.PassTicket] = loginMap.PassTicket

	webWxSyncJsonData := map[string]interface{}{}
	webWxSyncJsonData["BaseRequest"] = loginMap.BaseRequest
	webWxSyncJsonData["SyncKey"] = loginMap.SyncKeys
	webWxSyncJsonData["rr"] = -time.Now().Unix()

	jsonBytes, err := json.Marshal(webWxSyncJsonData)
	if err != nil {
		return wxMsges, err
	}

	resp, err := http.Post(e.WEB_WX_SYNC_URL+t.GetURLParams(urlMap), e.JSON_HEADER, strings.NewReader(string(jsonBytes)))
	if err != nil {
		return wxMsges, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wxMsges, err
	}

	/* 解析组装消息对象 */
	err = json.Unmarshal(bodyBytes, &wxMsges)
	if err != nil {
		return wxMsges, err
	}

	/* 更新SyncKey */
	loginMap.SyncKeys = wxMsges.SyncKeys
	loginMap.SyncKeyStr = wxMsges.SyncKeys.ToString()

	return wxMsges, nil
}

func SendMsg(loginMap *m.LoginMap, wxSendMsg m.WxSendMsg) error {
	urlMap := map[string]string{}
	urlMap[e.Lang] = e.LangValue
	urlMap[e.PassTicket] = loginMap.PassTicket

	wxSendMsgMap := map[string]interface{}{}
	wxSendMsgMap[e.BaseRequest] = loginMap.BaseRequest
	wxSendMsgMap["Msg"] = wxSendMsg
	wxSendMsgMap["Scene"] = 0

	jsonBytes, err := json.Marshal(wxSendMsgMap)
	if err != nil {
		return err
	}

	// TODO: 发送微信消息时暂不处理返回值
	_, err = http.Post(e.WEB_WX_SENDMSG_URL+t.GetURLParams(urlMap), e.JSON_HEADER, strings.NewReader(string(jsonBytes)))
	if err != nil {
		return err
	}

	return nil
}

/* 邀请联系人加入群 */
func InviteMember(loginMap *m.LoginMap, memberUserName string, chatRoomUserName string) error {
	urlMap := map[string]string{}
	urlMap["fun"] = "invitemember"

	wxUpdateChatRoomMap := map[string]interface{}{}
	wxUpdateChatRoomMap[e.BaseRequest] = loginMap.BaseRequest
	wxUpdateChatRoomMap["InviteMemberList"] = memberUserName
	wxUpdateChatRoomMap["ChatRoomName"] = chatRoomUserName
	jsonBytes, err := json.Marshal(wxUpdateChatRoomMap)
	if err != nil {
		return err
	}

	//TODO:发送群聊邀请暂不做反馈解析
	u, _ := url.Parse("https://wx.qq.com")
	timeout := time.Duration(30 * time.Second)

	jar := new(m.Jar)
	jar.SetCookies(u, loginMap.Cookies)

	client := &http.Client{
		Jar:     jar,
		Timeout: timeout}

	_, err = client.Post(e.WEB_WX_UPDATECHATROOM_URL+t.GetURLParams(urlMap), e.JSON_HEADER, strings.NewReader(string(jsonBytes)))

	return nil
}
