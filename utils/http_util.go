package utils

import (
	"io/ioutil"
	"net/http"
)

func GetUrl(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		Logger.Printf("GetUrl(%s) error!\n", url)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// panic(err)
		Logger.Printf("ioutil.ReadAll error: %s\n", err.Error())
		return ""
	}
	// fmt.Println(string(body))
	return string(body)
}
