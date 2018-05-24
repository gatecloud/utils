package utils

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractToken(t *testing.T) {
	header := make(http.Header)
	want := "Hello_gatecloud"
	header.Add("Authorization", "Bearer "+want)
	got, err := ExtractToken(header)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, want, got)
}

func TestGetURLPath(t *testing.T) {
	want := "Hello-World"
	s1 := "Hello World"
	s2 := "Hello &$@*#$#@World"
	s3 := "Hello      - World   "
	got1, err := GetURLPath(s1)
	got2, err := GetURLPath(s2)
	got3, err := GetURLPath(s3)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, want, got1)
	assert.Equal(t, want, got2)
	assert.Equal(t, want, got3)
}

func TestGetRandomString(t *testing.T) {
	s1 := GetRandomString(10)
	time.Sleep(1 * time.Nanosecond)
	s2 := GetRandomString(10)
	assert.NotEqual(t, s1, s2)

	s1 = GetRandomString(9)
	time.Sleep(1 * time.Nanosecond)
	s2 = GetRandomString(10)
	assert.NotEqual(t, s1, s2)
}

func TestGetDistance(t *testing.T) {
	want := 0.00
	got1 := GetDistance(-34.994709, 173.464301, -34.994709, 173.464301)
	assert.Equal(t, want, got1)
	got2 := GetDistance(-34.994710, 173.464301, -34.994709, 173.464301)
	assert.NotEqual(t, want, got2)
}
