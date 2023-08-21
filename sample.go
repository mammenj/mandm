package main

import (
	"fmt"

	"github.com/mammenj/mandm/models"
)

type MMData interface {
	models.Ad | models.User
}

type pageData1[T MMData] struct {
	User models.User
	Data map[string][]T
}

func main2() {

	ad := models.Ad{}
	ads := make([]models.Ad, 3)
	ads = append(ads, ad)
	admap := map[string][]models.Ad{"Ads": ads}
	//page := pageData1{models.User{}, "Ads", ads}

	data := pageData1[models.Ad]{
		User: models.User{}, // Initialize User with appropriate data
		//Key:  "Ads",
		//Data: ads, // Initialize Ads with MMData instances
		Data: admap,
	}

	fmt.Println(data)

	//log.Println("ads ", page)

}
