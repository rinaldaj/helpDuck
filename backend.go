package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"strings"
)


type service struct{
	Name	string
	Tag	string
	Cost	float64
	Discount	int
	Description	string
	PlaceName	string
	Loc	string
}

var dbMaster *sql.DB

func standardFeilds() string{
	return "description,name,tag,cost,discount,loc"
}

func processServices(res *sql.Rows) (ret []service){
	for res.Next(){
		var cur service
		if err := res.Scan(&cur.Description,&cur.Name,&cur.Tag,&cur.Cost,&cur.Discount,&cur.Loc); err != nil{
			continue
		}
		ret = append(ret,cur)
	}
	return

}

func getServices(db *sql.DB,tags []string, loc string) (ret []service){
	query := fmt.Sprintf("SELECT %s FROM thing",standardFeilds());
	where := false
	for _,tag := range tags{
		if !where {
			query = fmt.Sprintf("%s WHERE tag LIKE '%%%s%%'",query,tag)
			where= true
		} else {
			query = fmt.Sprintf("%s OR tag LIKE '%%%s%%'",query,tag)
		}
	}
	if loc == ""{
		query = fmt.Sprintf("%s;",query)
	} else if !where {
		query = fmt.Sprintf("%s WHERE loc LIKE '%%%s%%';",query,loc)
	} else {
		query = fmt.Sprintf("%s;",query)
	}
	results,err := db.Query(query)
	if err != nil{
		return
	}
	return processServices(results)
}

func main(){
	//TODO: read configuration file
	dbConf,err := ioutil.ReadFile("./.dbconfig")
	if err != nil {
		fmt.Println("Could not read .dbconfig, try running dbSetup.go")
		return
	}
	dbSplits := strings.Split(string(dbConf),"\n")
	dbUser := fmt.Sprintf("%s:%s@/%s",dbSplits[0],dbSplits[1],dbSplits[2])
	port := ":8080"
	dbMaster,err = sql.Open("mysql",dbUser)
	res := getServices(dbMaster,make([]string,0),"")
	if err != nil {
		return
	}
	http.Handle("/",http.FileServer(http.Dir("./frontend")))
	fmt.Printf("Listening on Port %s\n",port)
	for _,i := range res {
		fmt.Printf("%s",i.Name)
	}
	fmt.Println(http.ListenAndServe(port,nil))
}
