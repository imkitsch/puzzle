package hydra

import "puzzle/util"

var ProtocolList = []string{
	"ssh", "rdp", "ftp", "smb", "telnet",
	"mysql", "mssql", "oracle", "postgresql", "mongodb", "redis",
	//110:   "pop3",
	//995:   "pop3",
	//25:    "smtp",
	//994:   "smtp",
	//143:   "imap",
	//993:   "imap",
	//389:   "ldap",
	//23:   "telnet",
	//50000: "db2",
}

func Ok(protocol string) bool {
	if util.IsDuplicate(ProtocolList, protocol) {
		return true
	}
	return false
}
