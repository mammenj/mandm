package main

import (
	//"embed"
	"html/template"
	"regexp"

	//"io/fs"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/casbin/casbin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mammenj/mandm/handlers"
	"github.com/mammenj/mandm/models"
	"github.com/mammenj/mandm/security"
	"github.com/mammenj/mandm/storage"
	"github.com/mammenj/mandm/validators"
)

/*go:embed templates*/
//var templateFS embed.FS

/*go:embed static*/
//var staticFiles embed.FS

// var templateFS fs.FS
//var staticFiles fs.FS

var rootTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/index.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var contactTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/contact.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var groomTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/grooms.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var brideTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/brides.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var placeAdsTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/placead.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var loginTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/loginregister.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var tncTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/tnc.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var messageTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/chat.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var aboutUsTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/aboutus.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var myAdsTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/myads.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var activateTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/activate.html", "templates/menu.html", "templates/header.html", "templates/footer.html"))

var chatTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/chatbox.html"))

var validEmailTemplate *template.Template = template.Must(template.ParseFiles(
	"templates/validemail.html"))

type pageData struct {
	User  models.User
	AdMap map[string][]models.Ad
}

type pageDataMessages struct {
	User          models.User
	AdMap         map[string][]models.Ad
	AdMessagesMap map[string][]models.AdMessages
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cookie_secret := []byte(os.Getenv("COOKIE_SECRET"))
	r := gin.Default()

	store := cookie.NewStore([]byte(cookie_secret))
	r.Use(sessions.Sessions("jwt-session", store))

	r.Use(func(c *gin.Context) {
		c.Header("User-Agent", "Unreal-Minna_Minny")
	})

	r.StaticFS("/static", http.Dir("static"))

	e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv", true)

	auth := security.NewJwtAuth(e)

	r.Use(security.NewJwtAuthorizer(e))

	r.GET("/", func(c *gin.Context) {
		var page pageData

		user := auth.GetLoggedInUser(c)
		if user != nil {

			page = pageData{*user, nil}
		} else {
			page = pageData{models.User{Name: ""}, nil}
		}
		rootTemplate.Execute(c.Writer, page)
	})

	r.GET("/contact.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		contactTemplate.Execute(c.Writer, page)
	})

	r.GET("/messages.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		messageTemplate.Execute(c.Writer, page)
	})

	r.GET("/activate.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		activateTemplate.Execute(c.Writer, page)
	})

	r.GET("/myads.html", func(c *gin.Context) {
		var page pageDataMessages
		user := auth.GetLoggedInUser(c)
		if user != nil {
			log.Println("In my ads, user is not nil.....")
			adMsgStore := storage.NewSqliteAdMessageStore()
			adStore := storage.NewSqliteAdsStore()
			ads, err := adStore.GetMyAds(user.ID)

			if err != nil {
				log.Println("Error in getting Ads, ", err.Error())
			}
			admap := map[string][]models.Ad{"Ads": ads}
			msgs, err := adMsgStore.GetMessagesToID(user.ID)
			if err != nil {
				log.Println("Error in getting messages, ", err.Error())
			}
			admsgmap := map[string][]models.AdMessages{"AdMessages": msgs}
			log.Println("Ad message map:: ", admsgmap)
			page = pageDataMessages{*user, admap, admsgmap}
		}
		myAdsTemplate.Execute(c.Writer, page)
	})

	r.GET("/aboutus.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		aboutUsTemplate.Execute(c.Writer, page)
	})

	r.GET("/grooms.html", func(c *gin.Context) {
		adStore := storage.NewSqliteAdsStore()
		offset := c.Query("offset")
		offsetInt, _ := strconv.Atoi(offset)
		log.Println(" Groom offset ", offset)
		ads, err := adStore.GetSection("Groom Wanted", offsetInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		admap := map[string][]models.Ad{"Ads": ads}
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			log.Println("In groom user !=null ")
			page = pageData{*user, admap}
		} else {
			page = pageData{models.User{}, admap}
		}

		groomTemplate.Execute(c.Writer, page)
	})

	r.GET("/brides.html", func(c *gin.Context) {
		var page pageData
		adStore := storage.NewSqliteAdsStore()
		offset := c.Query("offset")
		offsetInt, _ := strconv.Atoi(offset)
		log.Println(" Bride offset ", offset)
		ads, err := adStore.GetSection("Bride Wanted", offsetInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		admap := map[string][]models.Ad{"Ads": ads}
		//brideTemplate.Execute(c.Writer, admap)
		user := auth.GetLoggedInUser(c)
		if user != nil {
			log.Println("In bride user !=null ")
			page = pageData{*user, admap}
		} else {
			page = pageData{models.User{}, admap}
		}
		//log.Println("Page in brides ", page)
		brideTemplate.Execute(c.Writer, page)
	})

	r.GET("/ads.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		placeAdsTemplate.Execute(c.Writer, page)
	})

	r.GET("/login.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}

		log.Println("Logged in Page is :::::", page)
		loginTemplate.Execute(c.Writer, page)
	})

	r.GET("/tnc.html", func(c *gin.Context) {
		var page pageData
		user := auth.GetLoggedInUser(c)
		if user != nil {
			page = pageData{*user, nil}
		}
		tncTemplate.Execute(c.Writer, page)
	})

	/// TEST CODE FOR EMBED END

	userHandler := handlers.CreateNewUserHandler()
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.PATCH("/users", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/login", security.Login)
	r.PATCH("/users/activate", userHandler.ActivateUser)
	r.POST("/users/sendmessage", func(c *gin.Context) {
		message := c.Request.FormValue("Message")
		log.Println("Body from send messge ", message)
		adid := c.Request.FormValue("ad-id")
		// hack to remove [] brackets
		adIdStr := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(adid, "")

		log.Println("ad-id from send messge ", adIdStr)
		adDescription := c.Request.FormValue("ad-description")
		log.Println("ad-description from send messge ", adDescription)
		user := auth.GetLoggedInUser(c)
		admsgStore := storage.NewSqliteAdMessageStore()
		loggedID := user.ID
		adStore := storage.NewSqliteAdsStore()
		toID, err := adStore.GetUserIDbyAdId(adIdStr)
		if err != nil {
			c.String(http.StatusOK, "div class=\"mx-1 bg-info\"> Invalid AD!</div>", nil)
			return
		}
		admessages := &models.AdMessages{FromUser: loggedID, ToUser: toID, AdID: adIdStr, Message: message}
		_, err = admsgStore.Create(admessages)
		if err != nil {
			c.String(http.StatusOK, err.Error(), nil)
			return
		}
		c.String(http.StatusOK, "<div class=\"mx-1 bg-info\"> Message sent!</div>")
	})
	r.GET("/users/message", func(c *gin.Context) {
		params := c.Request.URL.Query()
		log.Println("Params :: ", params)
		chatTemplate.Execute(c.Writer, params)
	})
	r.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("jwt")
		session.Clear()
		session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0
		session.Save()
		c.Header("HX-Location", "/")
		c.String(http.StatusOK, "Redirecting to Home...")
	})

	adHandler := handlers.CreateNewAdHandler()
	r.POST("/ads", adHandler.CreateAd)
	r.GET("/ads", adHandler.GetAds)
	r.PATCH("/ads", adHandler.UpdateAd)
	r.DELETE("/ads/:id", adHandler.DeleteAd)
	r.GET("/ads/:id", adHandler.GetAd)
	r.POST("/messages", func(c *gin.Context) {
		msgStore := storage.NewSqliteMessageStore()
		var input models.Messages
		if err := c.Bind(&input); err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		ID, err := msgStore.Create(&input)
		if err != nil {
			c.String(http.StatusOK, err.Error(), nil)
			return
		}

		log.Println("Messages CREATED ID: ", ID)
		c.String(http.StatusOK, "Messages sent, your feedback ID :: "+ID, "Keep the ID for reference")
	})

	r.POST("/validatemail", func(c *gin.Context) {
		email := c.Request.FormValue("Email")
		log.Println("Email is ", email)
		isValid, err := validators.ValidateEmail(email)
		if err != nil {
			log.Println("Email err is ", err)
			validEmailTemplate.Execute(c.Writer, err.Error())

		}
		if isValid {
			log.Println("Email valid ")
			validEmailTemplate.Execute(c.Writer, email)
		}
	})
	r.Run()
}
