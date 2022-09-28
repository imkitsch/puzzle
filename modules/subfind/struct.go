package subfind

type domainResult struct {
	domain string
	a      []string
	cname  string
	cdn    bool
}

type Item struct {
	Domain      string
	Name        string
	Description string
}
