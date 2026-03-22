package ufw

import (
	"bufio"
	"fmt"
	"strings"
)

func ParseAppList(out string) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	foundHeader := false
	var apps []string
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if !foundHeader {
			// Typically: "Available applications:"
			if strings.HasSuffix(strings.ToLower(line), "applications:") {
				foundHeader = true
			}
			continue
		}
		apps = append(apps, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if !foundHeader {
		return nil, fmt.Errorf("unable to parse ufw app list")
	}
	return apps, nil
}

