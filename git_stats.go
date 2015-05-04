package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	repos := ReposForUser("maryrosecook")
	fmt.Println(repos)
	fmt.Println(len(repos))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ReposForUser(user string) []string {
	page := 0
	perPage := 100
	done := false
	j := make([]string, 0)

	for !done {
		page++

		// Get the response
		resp, err := http.Get(
			fmt.Sprintf("https://api.github.com/users/%s/repos?page=%v&per_page=%v",
				user, page, perPage))
		check(err)
		defer resp.Body.Close()

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		// Convert to JSON
		pg := make([]map[string]interface{}, 0)
		err = json.Unmarshal(body, &pg)
		check(err)

		// If the last page fetched less than 30 (the default repos per page)
		if len(pg) < perPage {
			done = true
		}

		// Find the names of all the repo maps and append them
		for _, repo := range pg {
			// JSON hash map[string]interface{} so we ensure that we are getting strings
			if str, ok := repo["name"].(string); ok {
				j = append(j, str)
			} else {
				panic(fmt.Sprintf("Non-string value returned from API: %v", repo["name"]))
			}
		}
	}

	return j
}
