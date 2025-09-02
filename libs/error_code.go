package libs

type ErrorInfo struct {
	Code int
	Data any
	Msg  string
}

var ErrorCode = map[string]*ErrorInfo{
	"LoginFailed":                      {Code: 0, Data: "", Msg: "账号异常"},
	"LoginSuccessful":                  {Code: 200, Data: "", Msg: "登录成功"},
	"NetworkTimeout":                   {Code: 1, Data: "", Msg: "网络超时"},
	"AccountHasBeenTakenOffline":       {Code: 2, Data: "", Msg: "账号已下线，请重新登录"},
	"OtherError":                       {Code: 3, Data: "", Msg: "其他错误"},
	"TransferSuccessful":               {Code: 200, Data: "", Msg: "转账成功"},
	"TransferFailed":                   {Code: 4, Data: "", Msg: "转账失败"},
	"AccountRequiresEmailVerification": {Code: 0, Data: "", Msg: "账号需要邮箱验证"},
	"AccountRequiresSMSVerification":   {Code: 0, Data: "", Msg: "账号需要短信验证"},
	"AccountRequiresRobotVerification": {Code: 0, Data: "", Msg: "账号需要机器人验证"},
	"BrowserOpenFail":                  {Code: 0, Data: "", Msg: "浏览器开启失败"},
	"AccountFailed":                    {Code: 0, Data: "", Msg: "账号不存在"},
}
