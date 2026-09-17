package spider

// Engine 测绘引擎需实现的证书检索能力。
// 引擎内部自行处理查询打包/翻页/key轮换/查询限速,
// 对上层只暴露"输入目标列表,输出统一searchResult"的语义
type Engine interface {
	Name() string
	// QueryCert 按根域列表做证书检索
	QueryCert(domains []string) []*searchResult
}

// IPEngine 在证书检索之外还支持按IP检索的引擎(如fofa)。
// quake因api可用量较少且证书检索比fofa全面,仅保留证书检索,不实现该接口
type IPEngine interface {
	Engine
	// QueryIP 按IP列表检索关联资产
	QueryIP(ips []string) []*searchResult
}
