package portscan

import "testing"

func TestParseSmoke(t *testing.T) {
	cases := []struct {
		in      string
		wantLen int
	}{
		{"192.168.1.1", 1},
		{"192.168.1.0/24", 256},
		{"192.168.1.10-20", 11},
		{"192.168.1-255", 0},             // 畸形输入,不panic即通过
		{"10.1.255.0-10.2.0.255", 512},   // 跨八位组区间
		{"999.1.1.1", 0},                 // 非法IP
		{"192.168.1.20-10", 0},           // 起点大于终点
		{"255.255.255.254-255.255.255.255", 2}, // 区间终点为最大值,验证uint32回绕
	}
	for _, c := range cases {
		got, err := ParseIP(c.in)
		if err != nil {
			t.Fatalf("%s: unexpected err %v", c.in, err)
		}
		if len(got) != c.wantLen {
			t.Errorf("%s: got %d want %d", c.in, len(got), c.wantLen)
		}
	}
}

func TestPortsSmoke(t *testing.T) {
	p, err := ParsePorts("")
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1000 {
		t.Errorf("empty ports: got %d want 1000", len(p))
	}
}
