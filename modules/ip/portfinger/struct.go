package portfinger

// Result  定义返回结果
type Result struct {
	Addr          string
	ServiceName   string
	ProbeName     string
	VendorProduct string
	Version       string
}

type Options struct {
	Proxy  string
	Debug  bool
	Host   string
	Ports  string
	Thread int
}
