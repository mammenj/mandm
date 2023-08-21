package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/mammenj/mandm/models"
	"github.com/mammenj/mandm/security"
	"github.com/mammenj/mandm/storage"
)

type AdHandler struct {
	store *storage.AdSqlliteStore
}

func CreateNewAdHandler() *AdHandler {
	return &AdHandler{
		store: storage.NewSqliteAdsStore(),
	}
}

func (ah *AdHandler) GetAds(c *gin.Context) {
	log.Println("IN GET handler")

	ads, err := ah.store.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("ADS List", ads)
}

func (ah *AdHandler) UpdateAd(c *gin.Context) {
	log.Println("IN PATCH AD handler")
	var input models.Ad
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("IN PATCH AD handler ", &input)
	ID, err := ah.store.Update(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER Updated ID: ", ID)

}

func (ah *AdHandler) DeleteAd(c *gin.Context) {
	log.Println("IN Delete handler")

	id := c.Param("id")

	ID, err := ah.store.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER deleted ID: ", ID)
}
func (ah *AdHandler) GetAd(c *gin.Context) {
	log.Println("IN GET one AD handler")
	id := c.Param("id")

	user, err := ah.store.GetOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("...... Get Ad: ", user)
}

func (ah *AdHandler) CreateAd(c *gin.Context) {
	log.Println("IN Create AD handler")
	c.Request.ParseForm()
	for key, value := range c.Request.PostForm {
		fmt.Println(key, value)
	}
	e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv", true)

	auth := security.NewJwtAuth(e)
	user := auth.GetLoggedInUser(c)
	log.Println("Before In POST AD hanlder user IS :", user)
	if user == nil {
		c.JSON(http.StatusOK, "You must be logged in to post an Ad")
		return
	}
	log.Println("In POST AD hanlder user IS :", user)
	var input models.Ad
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.UserID = user.ID
	log.Println("In POST AD hanlder INPUT IS :", input)

	ID, err := ah.store.Create(&input)
	if err != nil {
		c.JSON(http.StatusOK,  err.Error())
		return
	}
	c.JSON(http.StatusOK, "Ad created with ID:: "+ ID)
	fmt.Println("AD CREATED ID: ", ID)
}
