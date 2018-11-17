package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"bufio"

)

func procErr(err error){
	if err != nil{
		panic(err)
	}
} 

func main(){
	fmt.Println("This script is for creating the databse")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("User Name:")
	dbUser, err := reader.ReadString('\n')
	procErr(err)
}
