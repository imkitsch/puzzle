package ip

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
