package main

// 342000 rows all K  1.347313958s seconds

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// type Geo_Thailand struct {
// 	Province          string `json:"ProvinceThai"`
// 	Area              string `json:"Area"`
// 	DistrictThaiShort string `json:"DistrictThaiShort"`
// }

// slmname K1 - K7 // 582546 total rows customer_transaction
type Customer_transaction struct {
	Slmname string `json:"slmname"`
	Docname string `json:"docname"`
	Docdate string `json:"docdate"`
}

// func queryDb(db *sql.DB, slmname string,  i int, box *[7]int) {
func queryDb(db *sql.DB, slmname string) {
	// count := 0
	// resultSplit, err := db.Query("SELECT ProvinceThai,Area,DistrictThaiShort FROM Geo_Thailand where Area = ?", slmname)
	resultSplit, err := db.Query("SELECT docdate,slmname,docname FROM customer_transaction where slmname = ?", slmname)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for resultSplit.Next() {
		// count++
		var tag Customer_transaction
		err = resultSplit.Scan(&tag.Docdate, &tag.Docname, &tag.Slmname)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println(tag)

	}
	// (box)[i] = count
	// fmt.Println(slmname, count)
	fmt.Println(slmname)
}

func main() {
	now := time.Now()
	// resultCh := make(chan Geo_Thailand)
	dsn := "root:admin12345@tcp(localhost:3389)/vansale_db?timeout=90s&collation=utf8mb4_unicode_ci"
	db, err := sql.Open("mysql", dsn)
	var box [7]int

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	areaList := []string{"K1", "K2", "K3", "K4", "K5", "K6", "K7"}
	// areaList := []string{"ภาคกลาง", "ภาคใต้", "ภาคเหนือ", "ภาคตะวันออกเฉียงเหนือ", "ภาคตะวันออก", "ภาคตะวันตก"}

	// for i, area := range areaList {
	for _, area := range areaList {
		// defer db.Close()
		// queryDb(db, area, i, &box)
		queryDb(db, area)
	}
	fmt.Println(time.Since(now),"finished")
	fmt.Println(box)
}
