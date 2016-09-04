package main

import (
	"database/sql"
	_ "github.com/lib/pq"
    "net/http"
    "log"
    "fmt"
    "encoding/json"
    "io/ioutil"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=geogit sslmode=disable")
	if err != nil {
			log.Fatal(err)
	}
  _,  err = db.Query("SELECT * FROM repos")
	if err != nil {
		_,  _err := db.Query("CREATE TABLE repos (URL TEXT NOT NULL, NAME TEXT NOT NULL)")
		if _err != nil {
			log.Fatal(_err)
			panic(fmt.Sprintf("Fuck you and your tables"))
		}
	}

	res, err := http.Get("https://api.github.com/repositories?since=364")
	if err != nil {
		log.Fatal(err)
	} 
	body, err := ioutil.ReadAll(res.Body)

	type GitHubAPIresponse struct {
	    URL string `json:"url"`
	    Name string `json:"name"`
	}
	var new_repo []GitHubAPIresponse

	err = json.Unmarshal(body, &new_repo)
    if(err != nil){
        fmt.Println("whoops:", err)
    }
    
    for each := range new_repo{
    	_,  err = db.Query("INSERT INTO repos (name, url) VALUES ('" + new_repo[each].Name + "', '" + new_repo[each].URL + "');")
    	if err != nil {
			log.Fatal(err)
		}
    }


}
