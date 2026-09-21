package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// minVips is the minimum vips version required for `vips dzsave --layout iiif3`.
var minVips = [3]int{8, 17, 1}

var vipsVersionRe = regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)

// parseVipsVersion extracts major.minor.patch from `vips --version` output
// (e.g. "vips-8.18.4"). A missing patch component is treated as 0.
func parseVipsVersion(out string) ([3]int, error) {
	var v [3]int
	m := vipsVersionRe.FindStringSubmatch(out)
	if m == nil {
		return v, fmt.Errorf("cannot parse vips version from %q", strings.TrimSpace(out))
	}
	for i := range 3 {
		if m[i+1] == "" {
			continue
		}
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return v, err
		}
		v[i] = n
	}
	return v, nil
}

func versionLess(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// checkVips verifies that vips is installed and at least minVips.
func checkVips() error {
	out, err := exec.Command("vips", "--version").Output()
	if err != nil {
		return fmt.Errorf("cannot run vips (is it installed?): %w", err)
	}
	v, err := parseVipsVersion(string(out))
	if err != nil {
		return err
	}
	if versionLess(v, minVips) {
		return fmt.Errorf("vips %d.%d.%d found, but %d.%d.%d or newer is required",
			v[0], v[1], v[2], minVips[0], minVips[1], minVips[2])
	}
	return nil
}
