package utils

import "testing"

func TestGetInternalIP(t *testing.T) {
	ip, err := GetInternalIP()
	if err != nil {
		t.Fatalf("get internal ip:%s", err.Error())
	}
	t.Logf("internal ip: %s", ip)
}
