package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	dat := Members("recursecenter", 0)
	fmt.Println(dat)
	fmt.Printf("Count: %v\n", len(dat))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func APIRequest(baseRequest string, limit int) []map[string]interface{} {
	file, err := os.Open(fmt.Sprintf("%v/.github_api_key", os.Getenv("HOME")))
	check(err)

	contents, err := ioutil.ReadAll(file)
	check(err)

	token := string(contents)

	total := 0
	page := 0
	done := false
	fin := make([]map[string]interface{}, 0)
	var perPage int
	if limit > 0 && limit < 100 {
		perPage = limit
	} else {
		perPage = 100
	}

	client := &http.Client{Timeout: 5 * time.Second}

	for !done {
		page++

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("https://api.github.com/%v?page=%v&per_page=%v",
				baseRequest, page, perPage),
			nil)
		check(err)
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

		resp, err := client.Do(req)
		check(err)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		check(err)

		js := make([]map[string]interface{}, 0)
		err = json.Unmarshal(body, &js)
		if err != nil {
			panic(string(body))
		}

		if len(js) < perPage {
			done = true
		}

		for _, item := range js {
			fin = append(fin, item)
			total++

			if limit > 0 && total >= limit {
				return fin
			}
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

func Repos(user string, limit int) []string {
	js := APIRequest(fmt.Sprintf("users/%v/repos", user), limit)
	vals := ValuesForKey("name", js)
	return StringifyInterfaceSlice(vals)
}

func Members(org string, limit int) []string {
	js := APIRequest(fmt.Sprintf("orgs/%v/members", org), limit)
	vals := ValuesForKey("login", js)
	return StringifyInterfaceSlice(vals)
}
