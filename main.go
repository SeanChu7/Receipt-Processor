package main

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type item struct {
	shortDescription string `json:"shortDescription"`
	price            string `json:"price"`
}

type receipt struct {
	Retailer     string `json:"retailer"`
	Purchasedate string `json:"purchaseDate"`
	Purchasetime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

// Calculate points for a given receipt, returns number of points
func calcPoints(r receipt) int {
	var pointTotal int

	//First calculate the points from the retailer
	regex := regexp.MustCompile(`[A-Za-z0-9]`)
	pointTotal = len(regex.FindAllString(r.Retailer, -1))
	//Calculate the amount of points from the total
	if tot, err := strconv.ParseFloat(r.Total, 32); err == nil {
		if res := math.Mod(tot, 1.00); res == 0 {
			pointTotal = pointTotal + 50
		}
		if res := math.Mod(tot, 0.25); res == 0 {
			pointTotal = pointTotal + 25
		}
	}
	//Calculate points from day of purchase
	regex = regexp.MustCompile(`\d{2}$`)
	date := regex.FindString(r.Purchasedate)
	if dateInt, err := strconv.Atoi(date); dateInt%2 == 1 && err == nil {
		pointTotal = pointTotal + 6
	}
	//Calculate points from time of purchase
	regex = regexp.MustCompile(`^.\d`)
	timeStr := regex.FindString(r.Purchasetime)
	if timeInt, err := strconv.Atoi(timeStr); timeInt < 16 && timeInt > 13 && err == nil {
		regex = regexp.MustCompile(`\d{2}$`)
		minuteStr := regex.FindString(r.Purchasetime)
		if timeInt == 14 {
			if minuteInt, err := strconv.Atoi(minuteStr); minuteInt != 0 && err == nil {
				pointTotal = pointTotal + 10 //Make sure to catch the case where purchaseTime is 2:00 PM
			}
		} else {
			pointTotal = pointTotal + 10
		}
	}
	//Calculate points for items
	pointTotal = pointTotal + ((len(r.Items) / 2) * 5)
	for i := 0; i < len(r.Items); i++ {
		desc := strings.Trim(r.Items[i].shortDescription, " ")
		if length := utf8.RuneCountInString(desc); length%3 == 0 {
			if itemPrice, err := strconv.ParseFloat(r.Items[i].price, 32); err == nil {
				num, frac := math.Modf(itemPrice * .2)
				var point int
				if frac > 0 {
					point = int(num + 1)
				} else {
					point = int(num)
				}
				pointTotal = pointTotal + point
			}
		}
	}

	return pointTotal
}
func postReceipts(c *gin.Context) {
	var newReceipt receipt

	// Call BindJSON to bind the received JSON to
	// newReceipt.
	if err := c.BindJSON(&newReceipt); err != nil {
		return
	}
	points := calcPoints(newReceipt)
	fmt.Println(points)
	c.IndentedJSON(http.StatusCreated, newReceipt)
}

func main() {
	//receipts:= make(map[int]int) //Receipts is a dictionary with ID keys and Point values
	router := gin.Default()
	router.POST("/receipts", postReceipts)

	r := receipt{
		Retailer:     "Target",
		Purchasedate: "2022-01-01",
		Purchasetime: "13:01",
		Items: []item{
			{
				shortDescription: "Mountain Dew 12PK",
				price:            "6.49",
			}, {
				shortDescription: "Emils Cheese Pizza",
				price:            "12.25",
			}, {
				shortDescription: "Knorr Creamy Chicken",
				price:            "1.26",
			}, {
				shortDescription: "Doritos Nacho Cheese",
				price:            "3.35",
			}, {
				shortDescription: "   Klarbrunn 12-PK 12 FL OZ   ",
				price:            "12.00",
			}},
		Total: "35.35"}
	/*r2 := receipt{
	Retailer:     "M&M Corner Market",
	Purchasedate: "2022-03-20",
	Purchasetime: "14:33",
	Items: []item{
		{
			shortDescription: "Gatorade",
			price:            "2.25",
		}, {
			shortDescription: "Gatorade",
			price:            "2.25",
		}, {
			shortDescription: "Gatorade",
			price:            "2.25",
		}, {
			shortDescription: "Gatorade",
			price:            "2.25",
		}},
	Total: "9.00"}*/
	points := calcPoints(r)
	fmt.Println(points)
	router.Run("localhost:8080")
}
