package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// type Geo_Thailand struct {
// 	Province          string `json:"ProvinceThai"`
// 	Area              string `json:"Area"`
// 	DistrictThaiShort string `json:"DistrictThaiShort"`
// }

type Customer_transaction struct {
	Slmname string `json:"slmname"`
	Docname string `json:"docname"`
	Docdate string `json:"docdate"`
}

func queryDb(db *sql.DB, slmname string, start time.Time) chan Customer_transaction {
	out := make(chan Customer_transaction)

	go func() {
		resultSplit, err := db.Query("SELECT docdate,slmname,docname FROM customer_transaction where slmname = ?", slmname)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		for resultSplit.Next() {
			var tag Customer_transaction
			err = resultSplit.Scan(&tag.Docdate, &tag.Docname, &tag.Slmname)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			out <- tag
		}
		close(out)
	}()
	return out
}

func main() {
	now := time.Now()
	// resultCh := make(chan Customer_transaction)
	dsn := "root:admin12345@tcp(localhost:3389)/vansale_db?timeout=90s&collation=utf8mb4_unicode_ci"
	db, err := sql.Open("mysql", dsn)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}

	// ปิดข้างใน loop
	// defer db.Close()

	areaList := []string{"K1", "K2", "K3", "K4", "K5", "K6", "K7"}

	chList := make([]chan Customer_transaction, len(areaList))

	for i, val := range areaList {
		chList[i] = queryDb(db, val, now)
		// ปิดข้างใน loop
	}
	defer db.Close()
	// queryResult := []Customer_transaction{}
	finalCh := merge(chList...)
	// count := 0
	for data := range finalCh {
		// count++
		// fmt.Println( count)
		fmt.Println(data)
		// queryResult=append(queryResult, data)
	}
	// fmt.Println(queryResult)
	fmt.Println(time.Since(now))

}

func merge(ch ...chan Customer_transaction) <-chan Customer_transaction {
	var wg sync.WaitGroup
	out := make(chan Customer_transaction)
	numberOfCh := len(ch)
	wg.Add(numberOfCh)
	fmt.Println(numberOfCh)

	for _, ch := range ch {
		go func(record <-chan Customer_transaction) {
			for r := range record {
				out <- r
			}
			defer wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)

	}()
	return out
}
