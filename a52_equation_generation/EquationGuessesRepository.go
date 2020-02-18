package a52_equation_generation

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func InsertEquations(multipleEquations string){

	db := OpenConnection()

	defer db.Close()

	_, err := db.Exec("INSERT INTO Equation VALUES " + multipleEquations)

	if err != nil {
		fmt.Println("Error inserting row: ", err)
		panic(err.Error())
	}
}

func OpenConnection() *sql.DB{
	db, err := sql.Open("mysql", "root:a52a52a52@tcp(127.0.0.1:3306)/equationguesses")

	if err != nil{
		fmt.Println("Error establishing connection to the database: ", err)
		panic(err.Error())
	}

	return db
}

