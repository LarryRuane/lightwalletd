// Copyright (c) 2019-2020 The Zcash developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or https://www.opensource.org/licenses/mit-license.php .
package cmd

import (
	"testing"
)

func TestFileExists(t *testing.T) {
	if fileExists("nonexistent-file") {
		t.Fatal("fileExists unexpected success")
	}
	// If the path exists but is a directory, should return false
	if fileExists(".") {
		t.Fatal("fileExists unexpected success")
	}
	// The following file should exist, it's what's being tested
	if !fileExists("root.go") {
		t.Fatal("fileExists failed")
	}
}

func TestDetectBackend(t *testing.T) {
	tests := []struct {
		name   string
		subver string
		want   string
	}{
		{"zebrad", "/Zebra:2.1.0/", "zebrad"},
		{"zcashd", "/MagicBean:6.3.0/", "zcashd"},
		{"zakura", "/Zakura:0.1.0/", "zakura"},

		// A BIP 14 user agent may carry several tokens. When zakura
		// advertises zcashd compatibility the leading token is the real
		// node, so it must win over the MagicBean token it also carries --
		// otherwise lightwalletd would demand zcashd's experimental
		// features from a node that has none.
		{"zakura advertising zcashd compat", "/Zakura:0.1.0/MagicBean:6.3.0/", "zakura"},
		{"zakura advertising zebra compat", "/Zakura:0.1.0/Zebra:2.1.0/", "zakura"},

		// Unrecognized backends yield "", which is what makes startServer
		// refuse to run without --no-backend-check.
		{"empty", "", ""},
		{"unknown node", "/SomeOtherNode:1.0.0/", ""},
		{"name without the BIP 14 colon", "Zakura", ""},
		{"bare MagicBean, no slash-colon form", "MagicBean", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectBackend(tt.subver); got != tt.want {
				t.Errorf("detectBackend(%q) = %q, want %q", tt.subver, got, tt.want)
			}
		})
	}
}
