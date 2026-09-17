package qqwry

type ResultQQwry struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}

type fileData struct {
	Data     []byte
	FilePath string
	IPNum    int64
}

// QQwry 纯真ip库
type QQwry struct {
	Data   *fileData
	Offset int64
}
