package oceanengine

const (
	BaseDomain            = "https://ad.oceanengine.com"
	AwemeAuthDomain       = "https://api.oceanengine.com"
	ApplicationJson       = "application/json"
	PageSize200           = 200
	PageSize1000          = 1000
	PageSize100           = 100
	PageSize10            = 8
	RefreshTokenGrantType = "refresh_token"
	AccessTokenGrantType  = "auth_code"
	AwemeAuthAuthType     = "SELF"
	AwemeAuthEndTime      = "2099-12-31 23:59:59"
)

var (
	ResponseDescription = map[uint32]string{
		0:     "成功",
		40001: "参数错误",
		40002: "没有权限进行相关操作",
		40003: "过滤条件的field字段错误",
		40100: "请求过于频繁",
		40101: "不合法的接入用户",
		40102: "access token过期",
		40103: "refresh token过期",
		40104: "access token为空",
		40105: "access token错误",
		40106: "账户登录异常",
		40107: "refresh token错误",
		40108: "授权类型错误",
		40109: "密码AES加密错误",
		40200: "充值金额太少",
		40201: "账户余额不足",
		40300: "广告主状态不可用",
		40301: "广告主在黑名单中",
		40302: "密码过于简单",
		40303: "邮箱已存在",
		40304: "邮箱不合法",
		40305: "名字已存在",
		40900: "文件签名错误",
		50000: "系统错误",
		61002: "当前开发者账号日累计调用接口次数超限",
		61003: "调用广告主账户和其他同主体广告主账户使用受限",
	}
)
