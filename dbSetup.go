package main

import (
	 "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"strings"

)

func procErr(err error){
	if err != nil{
		panic(err)
	}
} 

func main(){
	fmt.Println("This script is for creating the databse")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("root user Name:")
	dbRoot, err := reader.ReadString('\n')
	procErr(err)
	fmt.Println("Password:")
	rootPass, err := reader.ReadString('\n')
	procErr(err)
	dbRoot = strings.TrimSpace(dbRoot)
	rootPass = strings.TrimSpace(rootPass)
	rootUser := fmt.Sprintf("%s:%s@/",dbRoot,rootPass)


	fmt.Println("user userName:")
	dbUser, err := reader.ReadString('\n')
	procErr(err)
	fmt.Println("Password:")
	dbPass, err := reader.ReadString('\n')
	procErr(err)
	dbUser = strings.TrimSpace(dbUser)
	dbPass = strings.TrimSpace(dbPass)
	userUser := fmt.Sprintf("%s:%s@/duckhelp",dbUser,dbPass)

	db,err := sql.Open("mysql",rootUser)
	procErr(err)
	defer db.Close()
	_,err = db.Query("CREATE DATABASE duckhelp")
	procErr(err)
	_,err = db.Query(fmt.Sprintf("CREATE USER %q@'localhost' IDENTIFIED BY %q;",dbUser,dbPass))
	procErr(err)
	_,err = db.Query(fmt.Sprintf("GRANT ALL PRIVILEGES ON duckhelp.* TO %q@'localhost';",dbUser))
	procErr(err)
	_,err = db.Query("FLUSH PRIVILEGES")
	procErr(err)
	db.Close()

	db,err = sql.Open("mysql",userUser)
	procErr(err)
	_,err = db.Query("CREATE TABLE place(location TEXT NOT NULL,name NVARCHAR(200) NOT NULL, PRIMARY KEY(name));")
	procErr(err)
	_,err = db.Query("CREATE TABLE thing(description TEXT NOT NULL,name TEXT NOT NULL, tag TEXT NOT NULL, cost FLOAT(10,2) DEFAULT 0.00, discount INT DEFAULT 0, loc NVARCHAR(200) NOT NULL,FOREIGN KEY(loc) REFERENCES place(name) );")
	procErr(err)

	confs := fmt.Sprintf("%s\n%s\n",dbUser,dbPass)
	fmt.Println(confs)
	err = ioutil.WriteFile("./.dbconfig",[]byte(confs), 0644)
	procErr(err)
}
