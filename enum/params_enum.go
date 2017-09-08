package enum

const (
	APPID     = "appid"
	FUN       = "fun"
	Lang      = "lang"
	LangValue = "zh_CN"
	TimeStamp = "_"
	UUID      = "uuid"
	R         = "r"

	/* 以下信息会存储在loginMap中 */
	Ret          = "ret"
	Message      = "message"
	SKey         = "skey"
	WXSid        = "wxsid"
	WXUin        = "wxuin"
	PassTicket   = "pass_ticket"
	IsGrayscale  = "isgrayscale"
	DeviceID     = "DeviceID"
	SelfUserName = "UserName"
	SelfNickName = "NickName"
	SyncKeyStr   = "synckeystr"

	Sid         = "sid"
	Uin         = "uin"
	DeviceId    = "deviceid"
	SyncKey     = "synckey"
	BaseRequest = "BaseRequest"
)

var (
	uuidParaEnum = map[string]string{
		APPID:     "wx782c26e4c19acffb",
		FUN:       "new",
		Lang:      LangValue,
		TimeStamp: ""}

	loginParaEnum = map[string]string{
		"loginicon": "true",
		"tip":       "0",
		UUID:        "",
		R:           "",
		TimeStamp:   ""}

	initParaEnum = map[string]string{
		R:          "",
		Lang:       LangValue,
		PassTicket: ""}
)

func GetUUIDParaEnum() map[string]string {
	return uuidParaEnum
}

func GetLoginParaEnum() map[string]string {
	return loginParaEnum
}

func GetInitParaEnum() map[string]string {
	return initParaEnum
}
