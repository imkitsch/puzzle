package util

import (
	"path/filepath"
	"testing"
)

func TestCdnRangeGenerate(t *testing.T) {
	t.Skip("需要访问外部CDN厂商接口,手动执行")
	path := filepath.Join(t.TempDir(), "sources_data.txt")
	CdnRangeGenerate(path)
}
