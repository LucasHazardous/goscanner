package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./HashDB.db")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM 'HashDB'")
	checkErr(err)

	var rowCount string
	rowCountQuery, err := db.Query("SELECT COUNT(*) FROM 'HashDB'")
	checkErr(err)
	rowCountQuery.Next()
	err = rowCountQuery.Scan(&rowCount)
	checkErr(err)
	fmt.Println("Loaded", rowCount, "hashes")

	fmt.Print("Enter name of a file to scan: ")
	var filename string
	fmt.Scanf("%s", &filename)
	filehash := calculateHash(filename)

	var hash string
	var name string
	var i int32
	fmt.Printf("0/%s", rowCount)

	for rows.Next() {
		err = rows.Scan(&hash, &name)
		checkErr(err)
		if hash == filehash {
			fmt.Println("\n" + name)
			return
		}
		i++
		fmt.Printf("\033[2K\r%d/%s", i, rowCount)
	}
	fmt.Println("\nNo threats found")
}

func calculateHash(filename string) string {
	file, err := os.Open(filename)

	checkErr(err)

	defer file.Close()

	filehash := md5.New()
	_, err = io.Copy(filehash, file)

	checkErr(err)
	return fmt.Sprintf("%x", filehash.Sum(nil))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
