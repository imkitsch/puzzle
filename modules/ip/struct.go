package ip

type domainResult struct {
	domain string
	A      []string
	cdn    bool
}

type ipResult struct {
	ip      string
	City    string
	Country string
}

type CidrResult struct {
	CIDR     string
	operator string
	count    int
}
