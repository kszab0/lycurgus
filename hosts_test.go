package main

import (
	"strings"
	"testing"
)

func TestRemoveComment(t *testing.T) {

	tt := []struct {
		in, expected string
	}{
		{"", ""},
		{"a", "a"},
		{"# this is a comment", ""},
		{"a # comment", "a"},
		{"a# comment", "a"},
		{"a #comment", "a"},
		{"a		#comment", "a"},
		{"127.0.0.1 crash.163.com # comment comment comment", "127.0.0.1 crash.163.com"},
	}

	for _, tc := range tt {
		if out := removeComment(tc.in, "#"); out != tc.expected {
			t.Errorf("Output should be: %v; got: %v", tc.expected, out)
		}
	}
}

func TestReadLines(t *testing.T) {
	tt := []struct {
		in            string
		expectedErr   error
		expectedLines []string
	}{
		{
			in:            ``,
			expectedErr:   nil,
			expectedLines: []string{},
		},
		{
			in: `asdf
asdf`,
			expectedErr: nil,
			expectedLines: []string{
				"asdf",
				"asdf",
			},
		},
		{
			in: `asdf
			
asdf`,
			expectedErr: nil,
			expectedLines: []string{
				"asdf",
				"asdf",
			},
		},
		{
			in: `asdf
#comment
#comment comment
asdf
# comment
# comment comment
asdf #comment
asdf #comment comment
asdf # comment
asdf # comment comment`,
			expectedErr: nil,
			expectedLines: []string{
				"asdf",
				"asdf",
				"asdf",
				"asdf",
				"asdf",
				"asdf",
			},
		},
	}

	for _, tc := range tt {
		lines, err := readLines(strings.NewReader(tc.in))
		if err != tc.expectedErr {
			t.Errorf("Error should be %s; got: %s", tc.expectedErr, err)
		}
		if len(lines) != len(tc.expectedLines) {
			t.Errorf("Length of hosts should be %d; got: %d", len(tc.expectedLines), len(lines))
		}
	}
}

func TestParseHosts(t *testing.T) {

	tt := []struct {
		text          string
		expectedErr   error
		expectedHosts []string
	}{
		{
			text: `0x1f4b0.com
1q2w3.website
2giga.download`,
			expectedErr: nil,
			expectedHosts: []string{
				"0x1f4b0.com",
				"1q2w3.website",
				"2giga.download",
			},
		},
		{
			text: `127.0.0.1 analytics.163.com
127.0.0.1 mt.analytics.163.com
127.0.0.1 crash.163.com
127.0.0.1 crashlytics.163.com`,
			expectedErr: nil,
			expectedHosts: []string{
				"2giga.download",
				"mt.analytics.163.com",
				"crash.163.com",
				"crashlytics.163.com",
			},
		},
		{
			text: `127.0.0.1 localhost
::1 localhost
127.0.0.1 analytics.163.com
127.0.0.1 mt.analytics.163.com
127.0.0.1 crash.163.com
127.0.0.1 crashlytics.163.com`,
			expectedErr: nil,
			expectedHosts: []string{
				"2giga.download",
				"mt.analytics.163.com",
				"crash.163.com",
				"crashlytics.163.com",
			},
		},
		{
			text:          ``,
			expectedErr:   nil,
			expectedHosts: []string{},
		},
		{
			text: `127.0.0.1 analytics.163.com asdfasdf.com
127.0.0.1 crash.163.com`,
			expectedErr:   errParseHosts,
			expectedHosts: nil,
		},
	}

	for _, tc := range tt {
		hosts, err := parseHosts(strings.NewReader(tc.text))
		if err != tc.expectedErr {
			t.Errorf("Error should be %s; got: %s", tc.expectedErr, err)
		}
		if len(hosts) != len(tc.expectedHosts) {
			t.Errorf("Length of hosts should be %d; got: %d", len(tc.expectedHosts), len(hosts))
		}
	}
}
