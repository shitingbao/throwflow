package kuaidi

import (
	"common/internal/pkg/kuaidi"
	"encoding/json"
)

type Data struct {
	Context    string `json:"context"`    // 内容
	Time       string `json:"time"`       // 时间，原始格式
	Ftime      string `json:"ftime"`      // 格式化后时间
	Status     string `json:"status"`     // 本数据元对应的物流状态名称或者高级状态名称，实时查询接口中提交resultv2=1或者resultv2=4标记后才会出现
	StatusCode string `json:"statusCode"` // 本数据元对应的高级物流状态值，实时查询接口中提交resultv2=4标记后才会出现
	AreaCode   string `json:"areaCode"`   // 本数据元对应的行政区域的编码，实时查询接口中提交resultv2=1或者resultv2=4标记后才会出现
	AreaName   string `json:"areaName"`   // 本数据元对应的行政区域的名称，实时查询接口中提交resultv2=1或者resultv2=4标记后才会出现
	AreaCenter string `json:"areaCenter"` // 本数据元对应的行政区域经纬度，实时查询接口中提交resultv2=4标记后才会出现
	Location   string `json:"location"`   // 本数据元对应的快件当前地点，实时查询接口中提交resultv2=4标记后才会出现
	AreaPinYin string `json:"areaPinYin"` // 本数据元对应的行政区域拼音，实时查询接口中提交resultv2=4标记后才会出现
}

type GetKuaidiResponse struct {
	Message    string `json:"message"`   // 消息体，请忽略
	State      string `json:"state"`     // 快递单当前状态，默认为0在途，1揽收，2疑难，3签收，4退签，5派件，8清关，14拒签等10个基础物流状态，如需要返回高级物流状态，请参考 resultv2 传值
	Status     string `json:"status"`    // 通讯状态，请忽略
	Condition  string `json:"condition"` // 快递单明细状态标记，暂未实现，请忽略
	Ischeck    string `json:"ischeck"`   // 是否签收标记，0未签收，1已签收，请忽略，明细状态请参考state字段
	Com        string `json:"com"`       // 快递公司编码,一律用小写字母
	Nu         string `json:"nu"`        // 单号
	Data       []Data `json:"data"`      // 最新查询结果，数组，包含多项，全量，倒序（即时间最新的在最前），每项都是对象，对象包含字段请展开
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
}

func (gkr *GetKuaidiResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gkr); err != nil {
		return kuaidi.NewKuaidiError(kuaidi.BaseDomain+"/poll/query.do", "解析json失败："+err.Error(), response)
	} else {
		if len(gkr.ReturnCode) > 0 {
			return kuaidi.NewKuaidiError(kuaidi.BaseDomain+"/poll/query.do", gkr.Message, response)
		}
	}

	return nil
}
