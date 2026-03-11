package http

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vibezz/cms/internal/config"
	"github.com/vibezz/cms/internal/content"
	"github.com/vibezz/cms/internal/http/handlers"
	"github.com/vibezz/cms/internal/http/middleware"
	"github.com/vibezz/cms/internal/media"
	"github.com/vibezz/cms/internal/plugin"
	"github.com/vibezz/cms/plugins"
)

func NewRouter(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS(cfg.AdminURL, cfg.WebURL))

	userService := content.NewUserService(pool)
	ctService := content.NewContentTypeService(pool)
	entryService := content.NewEntryService(pool)
	pageService := content.NewPageService(pool, entryService, ctService)
	mediaService := media.NewService(pool, cfg.UploadDir)

	authHandler := handlers.NewAuthHandler(userService, cfg.JWTSecret)
	ctHandler := handlers.NewContentTypeHandler(ctService)
	entryHandler := handlers.NewEntryHandler(entryService)
	pagesHandler := handlers.NewPagesHandler(pageService, ctService)
	blockTypesHandler := handlers.NewBlockTypesHandler()
	mediaHandler := handlers.NewMediaHandler(mediaService)
	apiBaseURL := fmt.Sprintf("http://localhost:%s", cfg.Port)
	publicHandler := handlers.NewPublicHandler(entryService, ctService, mediaService, pageService, apiBaseURL)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// --- Public API ---
	pub := r.Group("/api/public")
	{
		pub.GET("/routes/*path", publicHandler.ResolveRoute)
		pub.GET("/pages/*path", publicHandler.GetPageByPath)
		pub.GET("/navigation", publicHandler.GetNavigation)
		pub.GET("/media/:id", mediaHandler.Serve)
	}

	// --- Admin Auth ---
	adminAuth := r.Group("/api/admin/auth")
	{
		adminAuth.POST("/login", authHandler.Login)
		adminAuth.POST("/refresh", authHandler.Refresh)
	}

	// --- Admin API (authenticated) ---
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		admin.GET("/me", authHandler.Me)
		admin.GET("/block-types", blockTypesHandler.List)

		ct := admin.Group("/content-types")
		{
			ct.GET("", ctHandler.List)
			ct.POST("", middleware.RoleRequired("admin"), ctHandler.Create)
			ct.GET("/:id", ctHandler.Get)
			ct.PATCH("/:id", middleware.RoleRequired("admin"), ctHandler.Update)
			ct.DELETE("/:id", middleware.RoleRequired("admin"), ctHandler.Delete)
		}

		pages := admin.Group("/pages")
		{
			pages.GET("", pagesHandler.ListTree)
			pages.POST("", pagesHandler.Create)
			pages.GET("/:id", pagesHandler.Get)
			pages.PATCH("/:id", pagesHandler.Update)
			pages.POST("/:id/publish", middleware.RoleRequired("admin", "editor"), pagesHandler.Publish)
			pages.POST("/:id/unpublish", middleware.RoleRequired("admin", "editor"), pagesHandler.Unpublish)
			pages.DELETE("/:id", middleware.RoleRequired("admin"), pagesHandler.Delete)
		}

		entries := admin.Group("/entries")
		{
			entries.GET("", entryHandler.List)
			entries.POST("", entryHandler.Create)
			entries.GET("/:id", entryHandler.Get)
			entries.PATCH("/:id", entryHandler.Update)
			entries.POST("/:id/publish", middleware.RoleRequired("admin", "editor"), entryHandler.Publish)
			entries.POST("/:id/unpublish", middleware.RoleRequired("admin", "editor"), entryHandler.Unpublish)
			entries.DELETE("/:id", middleware.RoleRequired("admin"), entryHandler.Delete)
		}

		mediaRoutes := admin.Group("/media")
		{
			mediaRoutes.GET("", mediaHandler.List)
			mediaRoutes.POST("", mediaHandler.Upload)
			mediaRoutes.GET("/:id", mediaHandler.Get)
			mediaRoutes.DELETE("/:id", middleware.RoleRequired("admin"), mediaHandler.Delete)
		}

		plugin.Load(context.Background(), plugins.All(), pool, cfg, entryService, ctService, mediaService, apiBaseURL, admin, pub)
	}

	return r
}
