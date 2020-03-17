package model

// taobao 查询数据结构
type TaobaoInfo struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Isp       string `json:"isp"`
}
type TaobaoIP struct {
	Code int        `json:"code"`
	Data TaobaoInfo `json:"data"`
}

// tencent 查询数据结构， 1w/日
type TencentInfo struct {
	IP       string `json:"ip"`
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	AdInfo struct {
		Nation   string `json:"nation"`
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Adcode   int    `json:"adcode"`
	} `json:"ad_info"`
}
type TencentIP struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  TencentInfo `json:"result"`
}

// 百度IP
type BaiduInfo struct {
	Location    string `json:"location"`
	OrigIP      string `json:"origip"`
	OrigipQuery string `json:"origipquery"`
}
type BaiduIP struct {
	Status    string      `json:"status"`
	T         string      `json:"t"`
	CacheTime string      `json:"set_cache_time"`
	Data      []BaiduInfo `json:"data"`
}

// my ipinfo
/*
需要判断两个，包括 area以及city
在阿里数据库中，
	香港地区 area='香港',city='XX'
	大陆地区 area='省级名称'，city='具体城市'

*/
type IPInfo struct {
	IP       string  `json:"ip"`       //ip地址
	Country  string  `json:"country"`  //国家
	Area     string  `json:"area"`     //区域,province,
	City     string  `json:"city"`     // 城市名
	District string  `json:"district"` //街道名
	Lng      float64 `json:"lng"`      //纬度
	Lat      float64 `json:"lat"`      //经度
	Isp      string  `json:"isp"`      //提供商
}
