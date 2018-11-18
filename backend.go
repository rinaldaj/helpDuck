package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"strings"
	"encoding/json"
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

type place struct{
	Name	string
	Location	string
	Services	[]service
}

var dbMaster *sql.DB

func serviceFeilds() string{
	return "description,name,tag,cost,discount,loc"
}

func locationFeilds() string{
	return "location,name"
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
func processPlaces(res *sql.Rows) (ret []place){
	for res.Next(){
		var cur place
		if err := res.Scan(&cur.Location,&cur.Name,); err != nil{
			continue
		}
		ret = append(ret,cur)
	}
	return

}

func getServices(db *sql.DB,tags []string, loc string) (ret []service){
	query := fmt.Sprintf("SELECT %s FROM thing",serviceFeilds());
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


func getPlaces(db *sql.DB) (ret []place){
	query := fmt.Sprintf("SELECT %s FROM place;",locationFeilds())
	res,err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	mid := processPlaces(res)
	for _,i := range mid{
		i.Services = getServices(db,make([]string,0),i.Name)
		ret = append(ret,i)
	}
	return
}

func serviceToJson(services []service) (ret [][]byte){
	for _,i := range services{
		tmp,err := json.Marshal(i)
		if err != nil {
			continue
		}
		ret = append(ret,tmp)
	}
	return ret;
}
func placeToJson(services []place) (ret [][]byte){
	for _,i := range services{
		tmp,err := json.Marshal(i)
		if err != nil {
			continue
		}
		ret = append(ret,tmp)
	}
	return ret;
}


func main(){
	dbConf,err := ioutil.ReadFile("./.dbconfig")
	if err != nil {
		fmt.Println("Could not read .dbconfig, try running dbSetup.go")
		return
	}
	dbSplits := strings.Split(string(dbConf),"\n")
	dbUser := fmt.Sprintf("%s:%s@/%s",dbSplits[0],dbSplits[1],dbSplits[2])
	port := ":8080"
	dbMaster,err = sql.Open("mysql",dbUser)
	if err != nil {
		return
	}
	defer dbMaster.Close()
	http.Handle("/",http.FileServer(http.Dir("./frontend")))
	fmt.Printf("Listening on Port %s\n",port)
	for _,i := range placeToJson(getPlaces(dbMaster)){
		fmt.Printf("%s\n",i)
	}
	fmt.Println(http.ListenAndServe(port,nil))
}
