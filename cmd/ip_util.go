package cmd

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func ObtainExternalIp() string {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(body))
}
