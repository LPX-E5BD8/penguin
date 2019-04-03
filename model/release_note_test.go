package model

import (
	"fmt"
	"testing"
)

func TestNewReleaseInfo(t *testing.T) {
	for _, ver := range []string{Version55, Version56, Version57, Version80} {
		_, err := NewReleaseInfo(ver)
		if err != nil {
			t.Error(err)
		}

	}
}

func TestHTTPGetWithCache(t *testing.T) {
	res, err := HTTPGetWithCache("http://www.baidu.com", CacheDir)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(byteToString(res))
}
