package mini

const (
	BaseDomain                 = "https://api.weixin.qq.com"
	Code2SessionTokenGrantType = "authorization_code"
	GetAccessTokenGrantType    = "client_credential"
)

var (
	ResponseDescription = map[int32]string{
		-1:    "系统繁忙，此时请开发者稍候再试",
		40001: "获取 access_token 时 AppSecret 错误，或者 access_token 无效。请开发者认真比对 AppSecret 的正确性，或查看是否正在为恰当的公众号调用接口",
		40013: "不合法的 AppID ，请开发者检查 AppID 的正确性，避免异常字符，注意大小写",
		40029: "js_code无效",
		40129: "最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）",
		40226: "高风险等级用户，小程序登录拦截 。风险等级详见用户安全解方案",
		41004: "缺少 secret 参数",
		41002: "缺少 appid 参数",
		41030: "page路径不正确：根路径前不要填加 /，不能携带参数（参数请放在scene字段里），需要保证在现网版本小程序中存在，与app.json保持一致。\n设置check_path=false可不检查page参数。",
		45011: "API 调用太频繁，请稍候再试",
		48006: "api 禁止清零调用次数，因为清零次数达到上限",
		85096: "scancode_time为系统保留参数，不允许配置",
	}
)
