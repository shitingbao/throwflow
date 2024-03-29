package openDouyin

const (
	BaseDomain            = "https://open.douyin.com"
	ApplicationJson       = "application/json"
	ApplicationForm       = "application/x-www-form-urlencoded"
	RefreshTokenGrantType = "refresh_token"
	AccessTokenGrantType  = "authorization_code"
	ClientTokenGrantType  = "client_credential"
	PageSize20            = 20
	DataType              = 7
)

var (
	ResponseDescription = map[uint64]string{
		10001: "系统错误",
		10002: "参数错误",
		10003: "检查 client_key 参数是否正确",
		10004: "权限不足",
		10005: "参数缺失",
		10007: "授权码过期",
		10010: "refresh_token 已过期",
		10013: "client_key 或者 client_secret 报错",
		10014: "client_key 不匹配",
		10020: "超过刷新次数限制",
		41050: "操作的用户需在通讯录权限范围中",
		40001: "参数错误",
	}
)
