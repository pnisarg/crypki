// Copyright 2019, Oath Inc.
// Licensed under the terms of the Apache License 2.0. Please see LICENSE file in project root for terms.
package config

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()
	cfg := &Config{
		ModulePath:        "/opt/utimaco/lib/libcs_pkcs11_R2.so",
		TLSServerName:     "cortana.corp.yahoo.com",
		TLSCACertPath:     "/opt/crypki/ca.crt",
		TLSClientAuthMode: 4,
		TLSServerCertPath: "/opt/crypki/server.crt",
		TLSServerKeyPath:  "/opt/crypki/server.key",
		TLSPort:           "4443",
		SignersPerPool:    2,
		Keys: []KeyConfig{
			{"key1", 1, "/path/1", "foo", 2, 1, true, "/path/foo", "", "", "", "", "", "My CA"},
			{"key2", 2, "/path/2", "bar", 2, 1, false, "", "", "", "", "", "", ""},
			{"key3", 3, "/path/3", "baz", 2, 1, false, "/path/baz", "", "", "", "", "", ""},
		},
		KeyUsages: []KeyUsage{
			{"/sig/x509-cert", []string{"key1", "key3"}, 3600},
			{"/sig/ssh-host-cert", []string{"key1", "key2"}, 36000},
		},
	}
	testcases := map[string]struct {
		filePath    string
		config      *Config
		expectError bool
	}{
		"good-config": {
			filePath:    "testdata/testconf-good.json",
			config:      cfg,
			expectError: false,
		},
		"bad-config-unknown-endpoint": {
			filePath:    "testdata/testconf-bad-unknown-endpoint.json",
			expectError: true,
		},
		"bad-config-unknown-identifier": {
			filePath:    "testdata/testconf-bad-unknown-identifier.json",
			expectError: true,
		},
		"bad-config-bad-json": {
			filePath:    "testdata/testconf-bad-json.json",
			expectError: true,
		},
		"bad-config-missing": {
			filePath:    "testdata/nonexist.json",
			expectError: true,
		},
	}

	for k, tt := range testcases {
		tt := tt // capture range variable - see https://blog.golang.org/subtests
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			cfg, err := Parse(tt.filePath)
			if err != nil {
				if !tt.expectError {
					t.Errorf("unexpected err: %v", err)
				}
				return
			}
			if tt.expectError {
				t.Error("expected error, got none")
				return
			}
			if !reflect.DeepEqual(cfg, tt.config) {
				t.Errorf("config mismatch, got: \n%+v\n, want: \n%+v\n", cfg, tt.config)
			}
		})
	}
}
