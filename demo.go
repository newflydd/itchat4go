package main

import (
	"fmt"
	e "itchat4go/enum"
	m "itchat4go/model"
	s "itchat4go/service"
	"os/exec"
	"regexp"
	"time"
)

var (
	uuid       string
	err        error
	loginMap   m.LoginMap
	contactMap map[string]m.User
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
	fmt.Printf("成功获取 %d个 联系人信息\n", len(contactMap))

	fmt.Println("开始监听消息响应...")
	var retcode, selector int64
	regAt := regexp.MustCompile(`^.*@.*丁丁.*$`)
	regGroup := regexp.MustCompile(`^@@.+`)
	regAd := regexp.MustCompile(`(朋友圈|点赞)+`)
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
