// Copyright (C) 2017 Librato, Inc. All rights reserved.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidServiceKey(t *testing.T) {
	valid1 := "ae38315f6116585d64d82ec2455aa3ec61e02fee25d286f74ace9e4fea189217:Go"

	invalid1 := ""
	invalid2 := "abc:Go"
	invalid3 := `
ae38315f6116585d64d82ec2455aa3ec61e02fee25d286f74ace9e4fea189217:
Go0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
`
	invalid4 := "1234567890abcdef"
	invalid5 := "1234567890abcdef:"
	invalid6 := ":Go"
	invalid7 := "abc:123:Go"

	keyPairs := map[string]bool{
		valid1:   true,
		invalid1: false,
		invalid2: false,
		invalid3: false,
		invalid4: false,
		invalid5: false,
		invalid6: false,
		invalid7: false,
	}

	for key, valid := range keyPairs {
		assert.Equal(t, valid, IsValidServiceKey(key))
	}
}

func TestMaskServiceKey(t *testing.T) {
	keyPairs := map[string]string{
		"1234567890abcdef:Go": "1234********cdef:Go",
		"abc:Go":              "abc:Go",
		"abcd1234:Go":         "abcd1234:Go",
	}

	for key, masked := range keyPairs {
		assert.Equal(t, masked, maskServiceKey(key))
	}
}

func TestIsValidTracingMode(t *testing.T) {
	assert.Equal(t, true, IsValidTracingMode("always"))
	assert.Equal(t, true, IsValidTracingMode("never"))
	assert.Equal(t, false, IsValidTracingMode("abc"))
	assert.Equal(t, false, IsValidTracingMode(""))
	assert.Equal(t, true, IsValidTracingMode("ALWAYS"))
}

func TestIsValidReporterType(t *testing.T) {
	assert.Equal(t, true, IsValidReporterType("udp"))
	assert.Equal(t, true, IsValidReporterType("ssl"))
	assert.Equal(t, true, IsValidReporterType("Udp"))
	assert.Equal(t, false, IsValidReporterType("xxx"))
	assert.Equal(t, false, IsValidReporterType(""))
	assert.Equal(t, false, IsValidReporterType("udpabc"))
}

func TestConverters(t *testing.T) {
	assert.Equal(t, int64(1), ToInt64("1"))
	assert.Equal(t, "ssl", ToReporterType("ssl").(string))
	assert.Equal(t, "never", ToTracingMode("never").(string))
}
