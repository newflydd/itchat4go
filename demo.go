package main

import (
	"fmt"
	e "itchat4go/enum"
	m "itchat4go/model"
	s "itchat4go/service"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	uuid       string
	err        error
	loginMap   m.LoginMap
	contactMap map[string]m.User
	groupMap   map[string][]m.User /* 关键字为key的，群组数组 */
)

func main() {
	/* 从微信服务器获取UUID */
	uuid, err = s.GetUUIDFromWX()
	if err != nil {
		panicErr(err)
	}

	/* 根据UUID获取二维码 */
	err = s.DownloadImagIntoDir(e.QRCODE_URL+uuid, "./qrcode")
	panicErr(err)
	cmd := exec.Command(`cmd`, `/c start ./qrcode/qrcode.jpg`)
	err = cmd.Run()
	panicErr(err)

	/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
	for {
		fmt.Println("正在验证登陆... ...")
		status, msg := s.CheckLogin(uuid)

		if status == 200 {
			fmt.Println("登陆成功,处理登陆信息...")
			loginMap, err = s.ProcessLoginInfo(msg)
			if err != nil {
				panicErr(err)
			}

			fmt.Println("登陆信息处理完毕,正在初始化微信...")
			err = s.InitWX(&loginMap)
			if err != nil {
				panicErr(err)
			}

			fmt.Println("初始化完毕,通知微信服务器登陆状态变更...")
			err = s.NotifyStatus(&loginMap)
			if err != nil {
				panicErr(err)
			}

			fmt.Println("通知完毕,本次登陆信息：")
			fmt.Println(e.SKey + "\t\t" + loginMap.BaseRequest.SKey)
			fmt.Println(e.PassTicket + "\t\t" + loginMap.PassTicket)
			break
		} else if status == 201 {
			fmt.Println("请在手机上确认")
		} else if status == 408 {
			fmt.Println("请扫描二维码")
		} else {
			fmt.Println(msg)
		}
	}

	fmt.Println("开始获取联系人信息...")
	contactMap, err = s.GetAllContact(&loginMap)
	if err != nil {
		panicErr(err)
	}
	fmt.Printf("成功获取 %d个 联系人信息,开始整理群组信息...\n", len(contactMap))

	groupMap = s.MapGroupInfo(contactMap)
	groupSize := 0
	for _, v := range groupMap {
		groupSize += len(v)
	}
	fmt.Printf("整理完毕，共有 %d个 群组是焦点群组，它们是：\n", groupSize)
	for _, v := range groupMap {
		for _, user := range v {
			fmt.Println(user.NickName)
		}
	}

	fmt.Println("开始监听消息响应...")
	var retcode, selector int64
	regAt := regexp.MustCompile(`^@.*@.*丁丁.*$`) /* 群聊时其他人说话时会在前面加上@XXX */
	regGroup := regexp.MustCompile(`^@@.+`)
	regAd := regexp.MustCompile(`(朋友圈|点赞)+`)
	secretCode := "cmd"
	regFather := regexp.MustCompile(`^:[^:]{2,4}$`)
	regChildren := regexp.MustCompile(`^::.{2,}$`)
	for {
		retcode, selector, err = s.SyncCheck(&loginMap)
		if err != nil {
			fmt.Println(retcode, selector)
			printErr(err)
			continue
		}

		if retcode == 0 && selector != 0 {
			fmt.Printf("selector=%d,有新消息产生,准备拉取...\n", selector)
			wxRecvMsges, err := s.WebWxSync(&loginMap)
			panicErr(err)

			for i := 0; i < wxRecvMsges.MsgCount; i++ {
				if wxRecvMsges.MsgList[i].MsgType == 1 {
					/* 普通文本消息 */
					fmt.Println(
						contactMap[wxRecvMsges.MsgList[i].FromUserName].NickName+":",
						wxRecvMsges.MsgList[i].Content)

					if regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName) && regAt.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人在群里@我，发个消息回答一下 */
						wxSendMsg := m.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = "我是丁丁编写的微信机器人，我已帮你通知我的主人，请您稍等片刻，他会跟您联系"
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						go s.SendMsg(&loginMap, wxSendMsg)
					} else if (!regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName)) && regAd.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人私聊我，并且内容含有「朋友圈」、「点赞」等敏感词，则回复 */
						wxSendMsg := m.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = "我是丁丁编写的微信机器人，我已经感应到一些敏感字符，我的主人对微商、朋友圈集赞、砍价等活动不感兴趣，请不要再发这些请求给他了，Sorry！"
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						go s.SendMsg(&loginMap, wxSendMsg)
					} else if (!regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName)) && strings.EqualFold(wxRecvMsges.MsgList[i].Content, secretCode) {
						/* 有人私聊我，并且内容是密语，输出加群菜单 */
						wxSendMsg := m.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = fmt.Sprintf("我是丁丁编写的微信群聊助手，我为您提供了以下分组关键词：\n\n%s\n\n您可以输入半角的':' + 关键词获取更详细的群聊目录，比如您可以输入[:编程]，我会为您细化系统内所有与IT技能相关的领域，您可以进一步选择加入该领域的微信群聊。", e.GetFatherKeywordsStr())
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						go s.SendMsg(&loginMap, wxSendMsg)
					} else if (!regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName)) && regFather.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人私聊我，并且内容是:+一级目录的形式，输出子目录 */
						description, exampleStr, keywords := e.GetChildKeywordsInfo(wxRecvMsges.MsgList[i].Content[1:])
						var content string
						if keywords == "" {
							content = fmt.Sprintf("没有查询到您输入关键字对应的分组信息，目前提供了以下分组关键词：\n\n%s\n\n您可以输入半角的':' + 关键词获取更详细的群聊目录，比如您可以输入[:编程]，我会为您细化系统内所有与IT技能相关的领域，您可以进一步选择加入该领域的微信群聊。", e.GetFatherKeywordsStr())
						} else {
							content = fmt.Sprintf(`%s\n\n%s\n\n您可以输入两个半角的:: + 关键词，我会为您寻找系统内的微信群组并邀请您加入，比如您可以输入%s`, description, keywords, exampleStr)
						}

						wxSendMsg := m.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = content
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						go s.SendMsg(&loginMap, wxSendMsg)
					} else if (!regGroup.MatchString(wxRecvMsges.MsgList[i].FromUserName)) && regChildren.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人私聊我，并且内容是::+二级目录的形式，邀请对方加入群聊 */
						if len(groupMap[strings.ToLower(wxRecvMsges.MsgList[i].Content[1:])]) == 0 {
							wxSendMsg := m.WxSendMsg{}
							wxSendMsg.Type = 1
							wxSendMsg.Content = fmt.Sprintf("Sorry，没有找到您查询的群组。我为您提供了以下分组关键词：\n\n%s\n\n您可以输入半角的':' + 关键词获取更详细的群聊目录，比如您可以输入[:编程]，我会为您细化系统内所有与IT技能相关的领域，您可以进一步选择加入该领域的微信群聊。", e.GetFatherKeywordsStr())
							wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
							wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
							wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
							wxSendMsg.ClientMsgId = wxSendMsg.LocalID

							go s.SendMsg(&loginMap, wxSendMsg)
						} else {
							go s.InviteMember(&loginMap, wxRecvMsges.MsgList[i].FromUserName, groupMap[strings.ToLower(wxRecvMsges.MsgList[i].Content[1:])][0].UserName)
						}
					}
				}
			}
		}
	}
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
