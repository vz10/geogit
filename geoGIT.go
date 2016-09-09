package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/tomnomnom/linkheader"
	"time"
)

func main() {

	type GitHubAPIresponse struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	var newRepo []GitHubAPIresponse
	count := 0
	// type GitHubAPIerror struct {
	//     message string `json:"message"`
	// }
	// var github_message GitHubAPIerror

	// db, err := sql.Open("postgres", "user=postgres dbname=geogit sslmode=disable")

    // TODO use ENV variables for user, pass, db_name
	db, err := sql.Open("postgres", "user=docker password=docker dbname=geogit host=db sslmode=disable")
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	_, err = db.Query("SELECT * FROM repos")
	if err != nil {
		_, _err := db.Query("CREATE TABLE repos (URL TEXT NOT NULL, NAME TEXT NOT NULL)")
		if _err != nil {
			log.Fatal(_err)
			panic(fmt.Sprintf("Fuck you and your tables"))
		}
	}
	rel := "next"
	link := "https://api.github.com/repositories"
	for rel == "next" {
		res, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		if len(res.Header["Link"]) > 0 {
			links := linkheader.Parse(res.Header["Link"][0])
			if len(links) > 0 {
				link = links[0].URL
				rel = links[0].Rel
			} else {
				break
			}
		} else {
			break
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(body[1])
		err = json.Unmarshal(body, &newRepo)
		if err != nil {
			fmt.Println("whoops:", err)
		}

		for each := range newRepo {
			fmt.Println(newRepo[each].Name, newRepo[each].URL)
			row, err := db.Query("INSERT INTO repos (name, url) VALUES ('" + newRepo[each].Name + "', '" + newRepo[each].URL + "');")
			if err != nil {
				log.Fatal(err)
			}
			row.Close()
			count++

		}
		time.Sleep(60000 * time.Millisecond)
	}
	fmt.Println(count, " repos added")
}
