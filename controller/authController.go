package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"test1/controller/dto"
	"test1/infrastructure"
	"test1/model"
	"test1/service"
	"time"
)

type authController struct {
	authService service.AuthService
	rsaService  infrastructure.RSAService
	userService service.UserService
}

type AuthController interface {
	LoginGoogle(c *fiber.Ctx)
	LoginGoogleCallback(c *fiber.Ctx)
	Profile(c *fiber.Ctx)
	LoginwithGoogle(c *fiber.Ctx)
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     "590182213606-34of2ltgaeaa34h9e9ejtl9l6rkcjiip.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"profile", "email"},
	Endpoint:     google.Endpoint,
}

var store = session.New(session.Config{
	Expiration:   3600,
	CookieSecure: false,
})

func (a *authController) LoginGoogle(c *fiber.Ctx) error {

	from := c.Query("from", "/")
	url := googleOauthConfig.AuthCodeURL(from)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}
func (a *authController) LoginGoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	var pathUrl string = "/"
	if state != "" {
		pathUrl = state
	}
	//get code from google
	code := c.Query("code")
	if code == "" {
		return c.SendStatus(401)
	}
	//get token from google
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.SendStatus(401)
	}
	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		return c.SendStatus(500)
	}
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := json.Marshal(userInfo)
	if err != nil {
		panic(err)
	}
	session.Set("user", string(jsonBytes))
	session.Save()

	//cek email ada
	if userInfo["email"] != "tuanpham1ss1@gmail.com" {
		return c.SendStatus(401)
	}
	//generate jwt token
	jwttoken := "ashdhashdsadhasdf"
	return c.JSON(fiber.Map{
		"status": "success",
		"token":  jwttoken,
		"data":   userInfo,
	})
	return c.Redirect(fmt.Sprintf("http://localhost:8080", pathUrl), fiber.StatusTemporaryRedirect)
}
func (a *authController) Profile(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	user := session.Get("user")
	if user == nil {
		return c.SendStatus(401)
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(user.(string)), &data)
	if err != nil {
		panic(err)
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func fetchGoogleUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err.Error())
	}
	defer response.Body.Close()
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err.Error())
	}
	return data, nil
}
func (a *authController) Register(c *fiber.Ctx) error {
	registerDto := dto.RegisterRequest{}
	if err := c.BodyParser(&registerDto); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid request payload",
		})
	}
	//validate
	if err := validator.New().Struct(registerDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//check user exist
	exists := a.userService.CheckUserExist(registerDto.Phone, registerDto.Email)
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User already exists",
		})
	}
	//hash password
	password, err := infrastructure.RsaEncrypt(registerDto.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	//create user
	profile, err := a.userService.CreateUser(registerDto.Username, password, registerDto.Phone, registerDto.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	return c.JSON(Response{
		Status:  true,
		Message: "Register success",
		Data:    profile,
		Errors:  nil,
	})
}
func (a *authController) Login(c *fiber.Ctx) error {
	payload := dto.LoginRequest{}
	if err := c.BodyParser(&payload); err != nil {
		log.Println("Error parsing body:", err)
		c.SendStatus(fiber.StatusBadRequest)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid request",
		})
		return err
	}
	//bước này thuong thuc hien ben phia client
	Encryptpassword, err := infrastructure.RsaEncrypt(payload.Password)
	if err != nil {
		log.Println("Error encrypting password:", err)
		c.SendStatus(fiber.StatusBadRequest)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid password",
		})
		return err
	}
	// nhan thong tin password da ma hoa tu client de giai ma phia server
	password, err := infrastructure.RsaDecrypt(Encryptpassword)
	if err != nil {
		log.Println("Error decrypting password:", err)
		c.SendStatus(fiber.StatusBadRequest)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid password",
		})
		return err
	}

	profile, err := a.userService.CheckUsernameAndPassword(payload.Username, string(password))
	if err != nil {
		log.Println("User not found or incorrect password for username:", payload.Username)
		c.Status(fiber.StatusUnauthorized)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid email or password",
		})
		return err
	}
	tokenDetail, err := a.authService.CreateToken(profile.ID)
	if err != nil {
		log.Println("Error creating token:", err)
		c.Status(fiber.StatusUnprocessableEntity)
		c.JSON(Response{
			Status:  false,
			Message: "Err not identify",
		})
		return err
	}

	if saveErr := a.authService.CreateAuth(profile.ID, tokenDetail); saveErr != nil {
		log.Println("Error saving token to Redis:", saveErr)
		c.Status(fiber.StatusUnprocessableEntity)
		c.JSON(Response{
			Status:  false,
			Message: "Err not identify",
		})
	}

	full_domain := c.Get("Origin")
	cookie_access := fiber.Cookie{
		Name:    "AccessToken",
		Domain:  full_domain,
		Path:    "/",
		Value:   tokenDetail.AccessToken,
		Expires: time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)),
	}
	cookie_refresh := fiber.Cookie{
		Name:    "RefreshToken",
		Domain:  full_domain,
		Path:    "/",
		Value:   tokenDetail.RefreshToken,
		Expires: time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)),
	}
	c.Cookie(&cookie_access)
	c.Cookie(&cookie_refresh)

	tokens := model.TokenLoadResponse{
		Profile:      profile,
		AccessToken:  tokenDetail.AccessToken,
		RefreshToken: tokenDetail.RefreshToken,
	}
	return c.JSON(Response{
		Status:  true,
		Message: "Login success",
		Data:    tokens,
	})
}
func (a *authController) LoginwithGoogle(c *fiber.Ctx) {
	payload := dto.LoginGooglePayload{}
	full_domain := c.Get("Origin")
	if err := c.BodyParser(&payload); err != nil {
		log.Println("Error parsing body:", err)
		c.SendStatus(fiber.StatusBadRequest)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid request",
		})
		return
	}
	//check email exist
	profile, err := a.userService.CheckGoogleToken(payload.Token)
	if err != nil {
		c.SendStatus(fiber.StatusUnauthorized)
		c.JSON(Response{
			Status:  false,
			Message: "Invalid email google",
		})
		return
	}
	tokenDetail, err := a.authService.CreateToken(profile.ID)
	if err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		c.JSON(Response{
			Status:  false,
			Message: "Err not identify",
		})
		return
	}
	if saveErr := a.authService.CreateAuth(profile.ID, tokenDetail); saveErr != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		c.JSON(Response{
			Status:  false,
			Message: "Err not identify",
		})
		return
	}
	cookie_access := fiber.Cookie{
		Name:    "AccessToken",
		Domain:  full_domain,
		Path:    "/",
		Value:   tokenDetail.AccessToken,
		Expires: time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)),
	}
	cookie_refresh := fiber.Cookie{
		Name:    "RefreshToken",
		Domain:  full_domain,
		Path:    "/",
		Value:   tokenDetail.RefreshToken,
		Expires: time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)),
	}
	c.Cookie(&cookie_access)
	c.Cookie(&cookie_refresh)
	tokens := model.TokenLoadResponse{
		Profile: profile,
	}
	c.JSON(Response{
		Status:  true,
		Message: "Login success",
		Data:    tokens,
	})
}
func (a *authController) Logout(c *fiber.Ctx) error {
	full_domain := c.Get("Origin")
	res := &Response{
		Status:  true,
		Message: "Logout success",
	}
	cookie_access := fiber.Cookie{
		Name:   "AccessToken",
		Domain: full_domain,
		Path:   "/",
		Value:  "",
		MaxAge: -1,
	}
	cookie_refresh := fiber.Cookie{
		Name:   "RefreshToken",
		Domain: full_domain,
		Path:   "/",
		Value:  "",
		MaxAge: -1,
	}
	c.Cookie(&cookie_access)
	c.Cookie(&cookie_refresh)
	return c.JSON(res)
}

func NewAuthController() *authController {
	authService := service.NewAuthService()
	rasService := infrastructure.NewRSAService()
	userService := service.NewUserService()
	return &authController{
		authService: authService,
		rsaService:  rasService,
		userService: userService,
	}
}
