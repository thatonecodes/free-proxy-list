package internal

import (
	"bufio"
	"bytes"
	"log"
	"strings"
	"fmt"
	"time"
)

func Load(proto string, content []byte) error {

	s := bufio.NewScanner(bytes.NewReader(content))

	var line, src string
	var transformer Transformer
	var parser Parser
	for s.Scan() {
		line = strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "https://") || strings.HasPrefix(line, "http://") {
			src, transformer, parser = parseLine(line)

			if src == "" {
				continue
			}

			log.Printf("> %v %s", Fetch(proto, src, transformer, parser), src)
		}

	}

	return nil
}

func roundToNearestJoyDeployHour(currentHour int) int {
	// joy-deploy uses: 00, 06, 12, 18 (4-hour increments)
	joyDeployHours := []int{0, 6, 12, 18}
	
	closest := 0
	minDiff := 24
	
	for _, joyHour := range joyDeployHours {
		diff := currentHour - joyHour
		if diff < 0 {
			diff = -diff
		}

		if diff > 12 {
			diff = 24 - diff
		}
		
		if diff < minDiff {
			minDiff = diff
			closest = joyHour
		}
	}
	
	return closest
}

func replaceDateTimeTokens(url string) string {
	now := time.Now()

	hour := now.Hour()
	minute := now.Minute()
	if strings.Contains(url, "joy-deploy/free-proxy-list") {
		hour = roundToNearestJoyDeployHour(hour)
		minute = 0
	}

	replacements := map[string]string{
		"{YYYY}": fmt.Sprintf("%04d", now.Year()),
		"{MM}":   fmt.Sprintf("%02d", now.Month()),
		"{DD}":   fmt.Sprintf("%02d", now.Day()),
		"{HH}":   fmt.Sprintf("%02d", hour),
		"{mm}":   fmt.Sprintf("%02d", minute),
		"{M}":    fmt.Sprintf("%d", now.Month()),
	}

	for token, val := range replacements {
		url = strings.ReplaceAll(url, token, val)
	}

	return url
}

func parseLine(line string) (string, Transformer, Parser) {

	if strings.HasPrefix(line, "https://") || strings.HasPrefix(line, "http://") {
		items := strings.Split(line, ",")

		var src string
		transformer := FromRaw
		parser := ParseProxyURL

		if len(items) > 0 {
			src = replaceDateTimeTokens(strings.TrimSpace(items[0]))
		}

		if len(items) > 1 {
			transformer = GetTransformer(strings.TrimSpace(items[1]))
		}

		if len(items) > 2 {
			parser = GetParser(strings.TrimSpace(items[2]))
		}

		return src, transformer, parser
	}

	return "", nil, nil
}
