package main

import "testing"

func TestSanitize(t *testing.T) {
	sane := sanitize("\"[]\\/$I?=&")
	if sane != "I" {
		t.Errorf("sanitize(\"\"[]\\/$I?=&\") = '%s'; want 'I'", sane)
	}
	sane = sanitize("../../../../../etc/passwd")
	if sane != ".etcpasswd" {
		t.Errorf("sanitize(\"../../../../../etc/passwd\") = '%s'; want '.etcpasswd'", sane)
	}

}
