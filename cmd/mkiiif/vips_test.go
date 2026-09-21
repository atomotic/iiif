package main

import "testing"

func TestParseVipsVersion(t *testing.T) {
	tests := []struct {
		in   string
		want [3]int
		ok   bool
	}{
		{"vips-8.18.4\n", [3]int{8, 18, 4}, true},
		{"vips-8.17.1-Mon Jan 1 00:00:00 UTC 2025", [3]int{8, 17, 1}, true},
		{"vips-8.16", [3]int{8, 16, 0}, true},
		{"garbage", [3]int{}, false},
	}
	for _, tt := range tests {
		got, err := parseVipsVersion(tt.in)
		if (err == nil) != tt.ok || (tt.ok && got != tt.want) {
			t.Errorf("parseVipsVersion(%q) = %v, %v", tt.in, got, err)
		}
	}
}

func TestVersionLess(t *testing.T) {
	if !versionLess([3]int{8, 17, 0}, minVips) {
		t.Error("8.17.0 should be < min")
	}
	if versionLess([3]int{8, 17, 1}, minVips) {
		t.Error("8.17.1 should satisfy min")
	}
	if versionLess([3]int{8, 18, 4}, minVips) {
		t.Error("8.18.4 should satisfy min")
	}
	if !versionLess([3]int{7, 99, 99}, minVips) {
		t.Error("7.99.99 should be < min")
	}
}
