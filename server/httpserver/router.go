package httpserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"luangao/biu"
	controllerauth "luangao/controller/auth"
	controllerfun "luangao/controller/fun"
	handlerauth "luangao/handler/auth"
	handlerfun "luangao/handler/fun"
	"luangao/store"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	return SetupRouterWithFinder(nil, nil)
}

func SetupRouterWithFinder(randomJumpFinder handlerfun.RandomJumpFinder, userStore *store.UserStore) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// ── JWT secret ──────────────────────────────────
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "wuck-local-dev-secret-change-in-production"
		log.Println("[warn] JWT_SECRET not set, using default (insecure for production)")
	}
	biu.SetJWTSecret(jwtSecret)

	// ── Database ─────────────────────────────────────
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/wuck.db"
		if execPath, err := os.Executable(); err == nil {
			altPath := filepath.Join(filepath.Dir(execPath), "data", "wuck.db")
			if _, err := os.Stat(filepath.Dir(altPath)); err == nil {
				dbPath = altPath
			}
		}
	}

	// Ensure the parent directory exists for the db
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Printf("[warn] cannot create data directory: %v", err)
	}

	var err error
	if userStore == nil {
		userStore, err = store.NewUserStore(dbPath)
		if err != nil {
			log.Printf("[warn] failed to initialize user store: %v", err)
		}
	}

	// ── Random jump handler ──────────────────────────
	if randomJumpFinder == nil {
		randomJumpFinder = handlerfun.NewRandomJumpHandler(userStore)
	}

	// ── Controllers ─────────────────────────────────
	randomJumpController := controllerfun.NewRandomJumpController(randomJumpFinder)
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// ── Routes ───────────────────────────────────────
	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(homePage))
	})

	r.GET("/api/biu", randomJumpController.GetRandomJump)
	r.GET("/api/biu/three", randomJumpController.GetThreeRandomJumps)

	// Auth routes
	if userStore != nil {
		authHandler := handlerauth.NewAuthHandler(userStore)
		authController := controllerauth.NewAuthController(authHandler)

		r.POST("/api/auth/register", authController.Register)
		r.POST("/api/auth/login", authController.Login)

		auth := r.Group("/api/auth")
		auth.Use(biu.AuthMiddleware())
		{
			auth.GET("/me", authController.Me)
			auth.POST("/click", authController.RecordClick)
		}
	}

	return r
}
