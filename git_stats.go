package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	dat := Members("recursecenter")
	fmt.Println(dat)
	fmt.Printf("Count: %v\n", len(dat))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func APIRequest(baseRequest string) []map[string]interface{} {
	page := 0
	perPage := 100
	done := false
	fin := make([]map[string]interface{}, 0)

	for !done {
		resp, err := http.Get(
			fmt.Sprintf("https://api.github.com/%v?page=%v&per_page=%v", baseRequest, page, perPage))
		check(err)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		check(err)

		js := make([]map[string]interface{}, 0)
		err = json.Unmarshal(body, &js)
		check(err)

		if len(js) < perPage {
			done = true
		}

		for _, item := range js {
			fin = append(fin, item)
		}
	}

	return fin
}

func ValuesForKey(key string, js []map[string]interface{}) []interface{} {
	fin := make([]interface{}, 0)
	for _, item := range js {
		fin = append(fin, item[key])
	}

	return fin
}

func StringifyInterfaceSlice(slc []interface{}) []string {
	fin := make([]string, 0)

	for _, v := range slc {
		if str, ok := v.(string); ok {
			fin = append(fin, str)
		} else {
			panic(fmt.Sprintf("Non-string value in JSON: %v\n", v))
		}
	}

	return fin
}

func Repos(user string) []string {
	js := APIRequest(fmt.Sprintf("users/%v/repos", user))
	vals := ValuesForKey("name", js)
	return StringifyInterfaceSlice(vals)
}

func Members(org string) []string {
	js := APIRequest(fmt.Sprintf("orgs/%v/members", org))
	vals := ValuesForKey("login", js)
	return StringifyInterfaceSlice(vals)
}
