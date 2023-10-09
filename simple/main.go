package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Geo_Thailand struct {
	Province          string `json:"ProvinceThai"`
	Area              string `json:"Area"`
	DistrictThaiShort string `json:"DistrictThaiShort"`
}

// slmname K1 - K7 // 582546 total rows customer_transaction
type Customer_transaction struct {
	Slmname string `json:"slmname"`
	Docname string `json:"docname"`
	Docdate string `json:"docdate"`
}

func queryDb(db *sql.DB, area string, start time.Time) {
	count := 0
	resultSplit, err := db.Query("SELECT ProvinceThai,Area,DistrictThaiShort FROM Geo_Thailand where Area = ?", area)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for resultSplit.Next() {
		count++
		var tag Geo_Thailand
		err = resultSplit.Scan(&tag.Area, &tag.Province, &tag.DistrictThaiShort)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println(tag)

	}
	fmt.Println(area, time.Since(start), count)
}

func main() {
	now := time.Now()
	// resultCh := make(chan Geo_Thailand)
	dsn := "root:admin12345@tcp(localhost:3389)/vansale_db?timeout=90s&collation=utf8mb4_unicode_ci"
	db, err := sql.Open("mysql", dsn)

	// db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	areaList := []string{"ภาคกลาง", "ภาคใต้", "ภาคเหนือ", "ภาคตะวันออกเฉียงเหนือ", "ภาคตะวันออก", "ภาคตะวันตก"}

	for _, area := range areaList {
		// defer db.Close()
		queryDb(db, area, now)
	}
	fmt.Println(time.Since(now))
}
