package test

import (
	"fmt"
	e "itchat4go/enum"
	m "itchat4go/model"
	s "itchat4go/service"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

var uuid string
var err error
var loginMap m.LoginMap
var contactMap map[string]m.User

func init() {
}

/* 从微信服务器获取UUID */
func TTestGetUUID(test *testing.T) {
	uuid, err = s.GetUUIDFromWX()

	if err != nil {
		panicErr(err)
	}
}

func TestStr(test *testing.T) {
	fmt.Println(regexp.MustCompile("(编程|全栈)").MatchString("全栈"))
}

/* 根据UUID获取二维码 */
func TTestDir(test *testing.T) {
	err = s.DownloadImagIntoDir(e.QRCODE_URL+uuid, "./qrcode")
	panicErr(err)

	cmd := exec.Command(`cmd`, `/c start ./qrcode/qrcode.jpg`)
	err = cmd.Run()
	panicErr(err)
}

/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
func TTestLogin(test *testing.T) {
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
}

func TTestGetAllContact(test *testing.T) {
	fmt.Println("开始获取联系人信息...")
	contactMap, err = s.GetAllContact(&loginMap)
	if err != nil {
		panicErr(err)
	}
	fmt.Printf("成功获取 %d个 联系人信息\n", len(contactMap))
}

func TTestMsgReceive(test *testing.T) {
	fmt.Println("开始监听消息响应...")
	var retcode, selector int64
	reg := regexp.MustCompile(`^.*@.*丁丁.*$`)
	for {
		retcode, selector, err = s.SyncCheck(&loginMap)
		if err != nil {
			panicErr(err)
		}

		if retcode == 0 && selector == 2 {
			fmt.Println("有消息产生，准备拉取...")
			wxRecvMsges, err := s.WebWxSync(&loginMap)
			panicErr(err)

			for i := 0; i < wxRecvMsges.MsgCount; i++ {
				if wxRecvMsges.MsgList[i].MsgType == 1 {
					fmt.Println(
						contactMap[wxRecvMsges.MsgList[i].FromUserName].NickName+":",
						wxRecvMsges.MsgList[i].Content)

					if reg.MatchString(wxRecvMsges.MsgList[i].Content) {
						/* 有人@我，发个消息回答一下 */
						wxSendMsg := m.WxSendMsg{}
						wxSendMsg.Type = 1
						wxSendMsg.Content = "我是丁丁编写的微信机器人，我已帮你通知我的主人，请您稍等片刻，他会跟您联系"
						wxSendMsg.FromUserName = wxRecvMsges.MsgList[i].ToUserName
						wxSendMsg.ToUserName = wxRecvMsges.MsgList[i].FromUserName
						wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
						wxSendMsg.ClientMsgId = wxSendMsg.LocalID

						s.SendMsg(&loginMap, wxSendMsg)
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
