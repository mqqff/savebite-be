package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	Analysishandler "github.com/mqqff/savebite-be/internal/app/analysis/interface/rest"
	AnalysisRepo "github.com/mqqff/savebite-be/internal/app/analysis/repository"
	AnalysisUsecase "github.com/mqqff/savebite-be/internal/app/analysis/usecase"
	AuthHandler "github.com/mqqff/savebite-be/internal/app/auth/interface/rest"
	AuthUsecase "github.com/mqqff/savebite-be/internal/app/auth/usecase"
	UserHandler "github.com/mqqff/savebite-be/internal/app/user/interface/rest"
	UserRepo "github.com/mqqff/savebite-be/internal/app/user/repository"
	UserUsecase "github.com/mqqff/savebite-be/internal/app/user/usecase"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"github.com/mqqff/savebite-be/internal/infra/gemini"
	"github.com/mqqff/savebite-be/internal/middlewares"
	"github.com/mqqff/savebite-be/pkg/jwt"
	"github.com/mqqff/savebite-be/pkg/log"
	"github.com/mqqff/savebite-be/pkg/markdown"
	"github.com/mqqff/savebite-be/pkg/oauth"
	"github.com/mqqff/savebite-be/pkg/supabase"
	"gorm.io/gorm"
	"time"
)

type HTTPServerItf interface {
	Start(socket string)
	MountMiddlewares()
	MountRoutes(db *gorm.DB)
	GetApp() *fiber.App
}

type HTTPServer struct {
	app *fiber.App
}

func NewHTTPServer() HTTPServerItf {
	config := fiber.Config{
		AppName:       env.AppEnv.AppName,
		CaseSensitive: true,
		ServerHeader:  env.AppEnv.AppEnv,
		IdleTimeout:   10 * time.Second,
		JSONDecoder:   sonic.Unmarshal,
		JSONEncoder:   sonic.Marshal,
	}

	app := fiber.New(config)

	return &HTTPServer{app}
}

func (s *HTTPServer) GetApp() *fiber.App {
	return s.app
}

func (s *HTTPServer) Start(socket string) {
	err := s.app.Listen(socket)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[HTTPServer][Start] failed to start server")
	}
}

func validateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	hashedAPIKey := sha256.Sum256([]byte(env.AppEnv.APIKey))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}

	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func (s *HTTPServer) MountMiddlewares() {
	s.app.Use(middlewares.Logger())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization,X-API-KEY,x-forwarded-for",
	}))

	s.app.Use(helmet.New())
	s.app.Use(compress.New())
	s.app.Use(middlewares.RequireAPIKey())
}

func (s *HTTPServer) MountRoutes(db *gorm.DB) {
	app := s.app

	api := app.Group("/api")
	v1 := api.Group("/v1")

	jwt := jwt.JWT
	oauth := oauth.OAuth
	validator := validator.New()
	supabase := supabase.Supabase
	gemini := gemini.Gemini
	md := markdown.Markdown

	middleware := middlewares.NewMiddleware(jwt)

	userRepo := UserRepo.NewUserRepo(db)
	analysisRepo := AnalysisRepo.NewAnalysisRepo(db)

	authUsecase := AuthUsecase.NewAuthUsecase(userRepo, oauth, jwt)
	userUsecase := UserUsecase.NewUserUsecase(userRepo)
	analysisUsecase := AnalysisUsecase.NewAnalysisUsecase(analysisRepo, supabase, gemini, md)

	AuthHandler.NewAuthHandler(v1, authUsecase, validator)
	UserHandler.NewUserHandler(v1, middleware, userUsecase)
	Analysishandler.NewAnalysisHandler(v1, middleware, analysisUsecase)
}
