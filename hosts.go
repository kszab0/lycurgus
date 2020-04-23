package main

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

var errParseHosts = errors.New("cannot parse hosts")

// removeComment removes a comment from a line
func removeComment(line string, commentChar string) string {
	split := strings.Split(line, commentChar)
	return strings.TrimSpace(split[0])
}

// getHost extracts a host from a line
func getHost(line string) (string, error) {
	split := strings.Split(line, " ")
	if len(split) == 1 {
		// plain domain list format
		return split[0], nil
	} else if len(split) == 2 {
		// hosts file format
		return split[1], nil
	} else {
		return "", errParseHosts
	}
}

func ignoredHost(host string) bool {
	if host == "localhost" {
		return true
	}
	return false
}

// readLines returns all the lines from an io.Reader
// without comments and empty lines
func readLines(r io.Reader) ([]string, error) {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = removeComment(line, "#")
		// ignore empty lines
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func parseHosts(r io.Reader) ([]string, error) {
	hosts := []string{}
	lines, err := readLines(r)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		host, err := getHost(line)
		if err != nil {
			return nil, err
		}
		if ignoredHost(host) {
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func parseHostsFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseHosts(file)
}

// Getter can make HTTP GET requests
type Getter interface {
	Get(string) (*http.Response, error)
}

func parseHostsURL(getter Getter, url string) ([]string, error) {
	resp, err := getter.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseHosts(resp.Body)
}
