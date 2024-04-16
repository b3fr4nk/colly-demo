package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type listing struct {
	SourceURL string `json:"URL"`
	Price float32 `json:"Price"`
	Mileage uint32 `json:"Mileage"`
	Location string `json:"Location"`

}

func dollarStringToFloat(s string) (float32, error) {
	cleanedString := strings.ReplaceAll(s, "$", "")
	cleanedString = strings.ReplaceAll(cleanedString, ",", "")

	dollarAmount, err := strconv.ParseFloat(cleanedString, 32)

	return float32(dollarAmount), err
}

func mileageStringToUInt(s string) (uint32, error) {
	cleanedString := strings.ReplaceAll(s, " mi.", "")
	cleanedString = strings.ReplaceAll(cleanedString, ",", "")

	mileage, err := strconv.ParseUint(cleanedString, 10, 32)
	
	return uint32(mileage), err
}

func main() {
	cars := []listing{}

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".vehicle-details", func(e *colly.HTMLElement) {
		item := listing{}
		url := e.ChildAttr("a.vehicle-card-link", "href")
		item.SourceURL = "https://www.cars.com" + url
		item.Location = e.ChildText(".miles-from")
		stringPrice := e.ChildText(".primary-price")
		price, err := dollarStringToFloat(stringPrice)
		if err != nil{
			fmt.Println(err)
		}
		item.Price = price
		
		mileageString := e.ChildText(".mileage")
		mileage, err := mileageStringToUInt(mileageString)
		if err != nil {
			fmt.Println(err)
		}
		item.Mileage = mileage
		
		cars = append(cars, item)
	})

	c.OnRequest((func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	}))

	c.Visit("https://www.cars.com/shopping/results/?dealer_id=&keyword=&list_price_max=30000&list_price_min=&makes[]=porsche&maximum_distance=all&mileage_max=100000&models[]=porsche-cayman&monthly_payment=340&page_size=20&sort=best_match_desc&stock_type=used&transmission_slugs[]=manual&year_max=&year_min=&zip=94912")

	carsJSON, _ := json.MarshalIndent(cars, "", " ")
	os.WriteFile("results/caymans.json", carsJSON, 0666)
	fmt.Println(string(carsJSON))
}