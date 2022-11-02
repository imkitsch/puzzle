package subfind

type domainResult struct {
	Domain  string
	Cdn     bool
	Address string
	Cname   string
}

type Item struct {
	Domain      string
	Name        string
	Description string
}
