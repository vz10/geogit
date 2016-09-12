package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"
	"github.com/tomnomnom/linkheader"
	"time"
	"strconv"
)

type Configuration struct {
    Client_id    string
    Client_secret   string
}

type GitHubAPIresponse struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Message string `json:"message"`
}

func checkCount(rows *sql.Rows) (count int) {
 	for rows.Next() {
    	err := rows.Scan(&count)
    	if err != nil {
   			panic(err)
    	}
    }   
    return count
}

func main() {

	var newRepo []GitHubAPIresponse
	var configuration Configuration
	count := 0

	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
	  fmt.Println("error:", err)
	}
	
	// db, err := sql.Open("postgres", "user=postgres dbname=geogit sslmode=disable")

    // TODO use ENV variables for user, pass, db_name
	db, err := sql.Open("postgres", "user=docker password=docker dbname=geogit host=db sslmode=disable")
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM repos")
	if err != nil {
		_, _err := db.Query("CREATE TABLE repos (URL TEXT NOT NULL, NAME TEXT NOT NULL,  IS_GEO BOOLEAN DEFAULT FALSE, CONSTRAINT u_constraint UNIQUE (URL))")
		if _err != nil {
			log.Fatal(_err)
			panic(fmt.Sprintf("Fuck you and your tables"))
		}
	} else {
		rows.Close()
	}
	rel := "next"
	rows, err = db.Query("SELECT COUNT(*) as count FROM repos")
 	repos_downloaded := checkCount(rows)
	link := "https://api.github.com/repositories?client_id="+configuration.Client_id+"&client_secret="+configuration.Client_secret+"&since="+strconv.Itoa(repos_downloaded)
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
			fmt.Println(res.Header)
			fmt.Println("*************")
			fmt.Println(res.Body)
			break
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &newRepo)
		if err != nil {
			fmt.Println("whoops:", err)
		}

		for each := range newRepo {
			fmt.Println(newRepo[each].Name, newRepo[each].URL)
			row, err := db.Query("INSERT INTO repos (name, url) VALUES ('" + newRepo[each].Name + "', '" + newRepo[each].URL + "') ON CONFLICT (URL) DO NOTHING;")
			if err != nil {
				log.Fatal(err)
			}
			row.Close()
			count++
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println(count, " repos added")
}
