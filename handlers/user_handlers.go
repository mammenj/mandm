package handlers

import (
	"fmt"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mammenj/mandm/messages"
	"github.com/mammenj/mandm/models"
	"github.com/mammenj/mandm/security"
	"github.com/mammenj/mandm/storage"
	"github.com/mammenj/mandm/validators"
)

type UserHandler struct {
	store *storage.UserSqlliteStore
}

func CreateNewUserHandler() *UserHandler {
	return &UserHandler{
		store: storage.NewSqliteUserStore(),
	}
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	log.Println("IN GET handler")

	users, err := uh.store.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//fmt.Println("USER List", users)
	c.JSON(http.StatusOK, &users)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	log.Println("IN PATCH  handler")
	var user models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Password = security.HashAndSalt([]byte(user.Password))
	log.Println("IN PATCH  handler ", &user)
	ID, err := uh.store.Update(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("USER Updated ID: ", ID)
	c.JSON(http.StatusOK, &user)
}

func (uh *UserHandler) ActivateUser(c *gin.Context) {
	log.Println("IN ActivateUser  handler")
	var user models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("IN ActivateUser  handler ", &user)
	// get the user and compare the UUID

	if err := uh.store.DB.Where("email = ? AND uuid = ?", user.Email, user.UUID).First(&user).Error; err != nil {
		log.Println("Failed to GetUser in db")
		//c.AbortWithStatus(http.StatusNotFound)
		//return nil, err
		c.String(http.StatusOK, "invalid email or activation code")
		return
	}

	// if matched update user to "active"
	var userActivated models.User
	userActivated.Status = "active"
	userActivated.Email = user.Email
	userActivated.ID = user.ID
	ID, err := uh.store.Update(&userActivated)
	if err != nil {
		c.String(http.StatusOK, "unable to activate, contact support")
		return
	}
	fmt.Println("USER Activated ID: ", ID)
	c.Header("HX-Location", "/login.html")
	c.String(http.StatusOK, "User Activated", nil)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	log.Println("IN Delete handler")

	id := c.Param("id")

	ID, err := uh.store.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("USER deleted ID: ", ID)
	c.JSON(http.StatusOK, "Deleted User ID: "+ID)
}
func (uh *UserHandler) GetUser(c *gin.Context) {
	log.Println("IN GET one handler")
	id := c.Param("id")
	value := c.GetHeader("Authorization")
	log.Println("GetHeader in UserHandler controller :: ", value)
	user, err := uh.store.GetOne(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("...... Get user: ", user)
	c.JSON(http.StatusOK, &user)
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	log.Println("IN Create handler User")

	log.Println("CreateUser in db")

	var user models.User
	if err := c.Bind(&user); err != nil {
		//c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		c.String(http.StatusOK, err.Error(), nil)
		return
	}
	email := user.Email
	_, err := validators.ValidateEmail(email)
	if err != nil {
		c.String(http.StatusOK, "Invalid Email: "+err.Error())
		return
	}
	log.Println("CreateUser bound user", user)
	user.Password = security.HashAndSalt([]byte(user.Password))
	ID, err := uh.store.Create(&user)
	log.Println("CreateUser user hashed password", user.Password)
	if err != nil {
		c.String(http.StatusOK, err.Error(), nil)
		return
	}
	//body := "Hello " + user.Name + ", you have registered for a Minnaminny account. If you havent registered, please ignore. Your activation code is " + user.UUID + ". Please visit http://localhost:8080/activate.html to activate your account."

	multi_body := ` Hello %s, you have registered for a MinnaMinny account.
	If you have not registered, please ignore. If you did register, this is your activation code: %s.
	Please visit our site at http://localhost:8080/activate.html to activate and use our website.
	Thanks and all the best in finding your spouse`

	filled_body := fmt.Sprintf(multi_body, user.Name, user.UUID)

	log.Println(">> multi_body is :: ", multi_body)

	log.Println("USER CREATED ID: ", ID)
	go messages.SendMail("mammenj@live.com", user.Email, "Activation test", filled_body)
	// if err != nil {
	// 	log.Fatal(">> Error Sending mail:: ", err.Error())
	// 	c.String(http.StatusOK, "User created ID:: "+ID, " but error sending email: "+err.Error())
	// 	return
	// }
	c.String(http.StatusOK, "User created ID:: "+ID)
}
