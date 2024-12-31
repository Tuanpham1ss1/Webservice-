package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"test1/controller"
	"test1/infrastructure"
	"time"
)

func Router() *fiber.App {
	// Create a new Fiber instance
	r := fiber.New()
	// Middleware
	r.Use(logger.New())    // Logger middleware
	r.Use(requestid.New()) // Request ID middleware
	r.Use(recover.New())   // Recover middleware
	r.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	})) // Compress middleware
	r.Use((cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           int(5 * time.Minute.Seconds()),
	})))
	//api swagger for develope mode
	r.Get("/api/swagger/*", swagger.New(swagger.Config{
		URL:          infrastructure.GetHTTPSwagger(),
		DocExpansion: "none",
	}))

	authcontroller := controller.NewAuthController()

	// Routes
	r.Get("/auth/google/callback", authcontroller.LoginGoogleCallback)

	r.Route("/api/user", func(router fiber.Router) {
		router.Route("/auth", func(access fiber.Router) {
			access.Post("/login", authcontroller.Login)
			access.Post("/register", authcontroller.Register)
			access.Post("/logout", authcontroller.Logout)
			access.Post("/login/google", authcontroller.LoginGoogle)
			//access.Post("login/facebook", authcontroller.LoginFacebook)
		})
		router.Get("/profile", authcontroller.Profile)
		//router.Put("/profile", authcontroller.UpdateProfile)
		//router.Put("/change-password", authcontroller.ChangePassword)
	})

	return r
}
