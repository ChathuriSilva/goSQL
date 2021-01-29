package main

import (
	"fmt"
	"net/url"

	"database/sql"

	"github.com/gocolly/colly"

	_ "github.com/go-sql-driver/mysql"
)

func saveData(item string, district string, title string, price string) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/ikman")
	fmt.Println(item, district, title, price)

	if err != nil {
		panic(err.Error())
	}
	insert, err := db.Query("INSERT INTO products VALUES( '" + item + "', '" + district + "', '" + title + "')")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	defer db.Close()
	fmt.Println("inserted")
}

func getData(url *url.URL, title string, pprice string, item string, dis string) {

	col := colly.NewCollector()

	col.OnHTML(".amount--3NTpl", func(e *colly.HTMLElement) {
		price := e.Attr("class")
		fmt.Println(e.Text)
		pprice = e.Text
		col.Visit(e.Request.AbsoluteURL(price))
	})
	col.OnHTML(".word-break--2nyVq", func(e *colly.HTMLElement) {
		information := e.Attr("class")
		fmt.Println(e.Text)
		col.Visit(e.Request.AbsoluteURL(information))
	})
	saveData(item, dis, title, pprice)
	col.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})

	col.Visit(url.String())

}

func main() {
	fmt.Println("Enter Item Type:")
	var item string
	fmt.Scanln(&item)

	fmt.Println("Enter District:")
	var dis string
	fmt.Scanln(&dis)

	var Ptitle string
	var Pprice string

	///////////////////////////////////////////////

	c := colly.NewCollector()

	c.OnHTML(".gtm-normal-ad a", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("class"))
		title := e.Attr("href")
		Ptitle = e.Text
		fmt.Println(e.Text)
		//fmt.Printf("Title Found: %q -> %s\n", e.Text, title)
		c.Visit(e.Request.AbsoluteURL(title))
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
		fmt.Println("===================================")
		url := (r.URL)
		getData(url, Ptitle, Pprice, item, dis)
		fmt.Println("===================================")
	})

	c.Visit("https://ikman.lk/en/ads/" + dis + "/" + item)

}
