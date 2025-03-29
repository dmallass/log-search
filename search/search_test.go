package search

import (
	"reflect"
	"strings"
	"testing"
)

func TestRunRipgrep(t *testing.T) {
	var tests = []struct {
		name       string
		logFile    string
		searchTerm string
		searchMode string
		want       []string
		wantErr    bool
		errStr     string
	}{
		{
			name:       "search success with results",
			logFile:    "../testdata/app.log",
			searchTerm: "error",
			searchMode: "fulltext",
			want: []string{
				"2025-03-27T22:32:21.745241-07:00 ERROR GET https://server2-us-central.com/v1/shoppers RATE_LIMIT_EXCEEDED | app=gig-identity-verification | source=firewall | file=applog_2024_03_23.log | datacenter=us-west-1 | ip=235.12.92.47 | user=admin22",
				"2025-03-27T22:32:21.745222-07:00 ERROR DELETE https://server2-us-central.com/v1/shoppers 500 4706 bytes | app=gig-candidates | source=proxy | file=security_alerts.log | datacenter=us-east-1, | ip=239.35.184.51 | user=devOps",
			},
			wantErr: false,
		},
		{
			name:       "search success with empty results, no error",
			logFile:    "../testdata/app.log",
			searchTerm: "v1/search",
			searchMode: "fulltext",
			want:       []string{},
			wantErr:    false,
		},
		{
			name:       "search success with regex, no error",
			logFile:    "../testdata/app.log",
			searchTerm: ".*/v1/candidates",
			searchMode: "regex",
			want: []string{
				"2025-03-27T22:32:21.745239-07:00 INFO PUT https://server1-us-central.com/v1/candidates 200 5046 bytes | app=gig-candidates | source=dns_filter | file=applog_2024_03_23.log | datacenter=us-west-1, | ip=23.212.250.53 | user=root",
				"2025-03-27T22:32:21.745234-07:00 DEBUG POST https://server3-us-central.com/v1/candidates 422 545 bytes | app=gig-candidates | source=vpn_gateway | file=security_alerts.log | datacenter=us-central-1, | ip=154.190.223.224 | user=user123",
				"2025-03-27T22:32:21.745232-07:00 WARN GET https://server1-us-central.com/v1/candidates 422 1413 bytes | app=gig-identity-verification | source=vpn_gateway | file=security_alerts.log | datacenter=us-west-1, | ip=251.153.35.8 | user=guest001",
				"2025-03-27T22:32:21.745229-07:00 WARN GET https://server2-us-central.com/v1/candidates 404 5152 bytes | app=gig-candidates | source=firewall | file=applog_2024_03_23.log | datacenter=us-central-1, | ip=200.128.247.182 | user=root",
				"2025-03-27T22:32:21.74522-07:00 DEBUG PUT https://server3-us-central.com/v1/candidates 301 4894 bytes | app=shopper-onboarding | source=endpoint_security | file=security_alerts.log | datacenter=us-west-1, | ip=119.90.198.187 | user=admin22",
			},
			wantErr: false,
		},
		{
			name:       "search failure case with regex, with error",
			logFile:    "../testdata/app.log",
			searchTerm: "*/v1/candidates",
			searchMode: "regex",
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "test error case, file or dir doesnt exist",
			logFile:    "../testdat/app.log",
			searchTerm: ".*/v1/candidates",
			searchMode: "regex",
			want:       nil,
			wantErr:    true,
			errStr:     "ripgrep error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := RunRipgrep(tt.logFile, tt.searchTerm, tt.searchMode)
			if err != nil && !tt.wantErr {
				t.Errorf("got error, want no error; err=%+v", err)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errStr) {
				t.Errorf("want error: %v, got %v", tt.errStr, err.Error())
			}
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	var tests = []struct {
		name     string
		filepath string
		want     string
		wantErr  bool
	}{
		{
			name:     "success case, valid filepath",
			filepath: "/var/log/app.log",
			want:     "/var/log/app.log",
			wantErr:  false,
		},
		{
			name:     "failure case, invalid filepath",
			filepath: "/tmp/app.log",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "failure case, invalid filepath",
			filepath: "/invalid/path/file.txt\x00.log",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "failure case, invalid filepath",
			filepath: "/",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ValidateFilePath(tt.filepath)
			if err != nil && !tt.wantErr {
				t.Errorf("got error, want no error; err=%+v", err)
			}
			if res != tt.want {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}
}
