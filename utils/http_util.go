package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetUrl(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	return string(body)

}
