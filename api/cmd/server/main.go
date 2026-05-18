package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/optea-tech/api/internal/config"
	"github.com/optea-tech/api/internal/database"
	"github.com/optea-tech/api/internal/handlers"
	"github.com/optea-tech/api/internal/middleware"
	"github.com/optea-tech/api/internal/repository"
	"github.com/optea-tech/api/internal/services"
)

func main() {
	cfg := config.Load()
	app := fiber.New(fiber.Config{
		AppName: "optea.tech API",
	})

	app.Use(recover.New())
	app.Use(logger.New(middleware.Logger()))
	app.Use(middleware.SecurityHeaders())
	app.Use(cors.New(middleware.CORS(cfg.AllowedOrigins)))
	app.Use(limiter.New(middleware.GlobalRateLimit()))

	var dbPool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)
		if err != nil {
			log.Printf("postgres disabled, fallback to in-memory store: %v", err)
		} else {
			dbPool = pool
		}
	}
	if dbPool != nil {
		defer dbPool.Close()
	}

	store := handlers.NewStore()
	projectsHandler := handlers.NewProjectsHandler(store, dbPool)
	servicesHandler := handlers.NewServicesHandler(store, dbPool)
	testimonialsHandler := handlers.NewTestimonialsHandler(store, dbPool)
	messagesHandler := handlers.NewMessagesHandler(store, dbPool)
	contactHandler := handlers.NewContactHandler(messagesHandler)
	authHandler := handlers.NewAuthHandler(store, dbPool, cfg.JWTSecret, cfg.JWTExpiresSeconds)
	requestsRepo := repository.NewRequestsRepo(dbPool)
	portalRepo := repository.NewPortalRepo(dbPool)
	accessLogRepo := repository.NewAccessLogRepo(dbPool)
	emailService := services.NewEmailService(cfg.ResendAPIKey, cfg.EmailFrom, cfg.AdminEmail)
	requestsHandler := handlers.NewRequestsHandler(requestsRepo, emailService, cfg.PortalBaseURL)
	clientPortalHandler := handlers.NewClientPortalHandler(requestsRepo, portalRepo, accessLogRepo, emailService, cfg.PortalBaseURL)
	adminRequestsHandler := handlers.NewAdminRequestsHandler(requestsRepo, portalRepo, emailService, cfg.PortalBaseURL)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "version": "1.0.0"})
	})

	api := app.Group("/api")
	api.Get("/projects", projectsHandler.ListPublic)
	api.Get("/projects/:slug", projectsHandler.GetBySlug)
	api.Get("/services", servicesHandler.ListPublic)
	api.Get("/testimonials", testimonialsHandler.ListPublic)

	api.Post("/contact", limiter.New(middleware.ContactRateLimit()), contactHandler.Submit)
	api.Post("/requests", limiter.New(middleware.ClientRequestRateLimit()), requestsHandler.CreateRequest)
	api.Post("/client/request-new-link", limiter.New(middleware.ClientRequestNewLinkRateLimit()), clientPortalHandler.RequestNewToken)

	// Public aliases used by Next.js /api/go proxy in the no-Nest architecture.
	public := api.Group("/public")
	public.Post("/requests", limiter.New(middleware.ClientRequestRateLimit()), requestsHandler.CreateRequest)
	public.Post("/client/request-new-link", limiter.New(middleware.ClientRequestNewLinkRateLimit()), clientPortalHandler.RequestNewToken)
	public.Get("/client/dashboard", middleware.ClientAuth(requestsRepo, accessLogRepo), clientPortalHandler.GetDashboard)
	public.Post("/client/messages", middleware.ClientAuth(requestsRepo, accessLogRepo), limiter.New(middleware.ClientMessageRateLimit()), clientPortalHandler.SendMessage)
	public.Get("/client/deliverables/:id/download", middleware.ClientAuth(requestsRepo, accessLogRepo), clientPortalHandler.DownloadDeliverable)
	public.Post("/client/quote/accept", middleware.ClientAuth(requestsRepo, accessLogRepo), clientPortalHandler.AcceptQuote)

	client := api.Group("/client", middleware.ClientAuth(requestsRepo, accessLogRepo))
	client.Get("/dashboard", clientPortalHandler.GetDashboard)
	client.Post("/messages", limiter.New(middleware.ClientMessageRateLimit()), clientPortalHandler.SendMessage)
	client.Get("/deliverables/:id/download", clientPortalHandler.DownloadDeliverable)
	client.Post("/quote/accept", clientPortalHandler.AcceptQuote)

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	admin := api.Group("/admin", middleware.AdminJWT(cfg.JWTSecret))
	admin.Get("/projects", projectsHandler.ListAdmin)
	admin.Post("/projects", projectsHandler.Create)
	admin.Put("/projects/:id", projectsHandler.Update)
	admin.Delete("/projects/:id", projectsHandler.Delete)
	admin.Get("/services", servicesHandler.ListAdmin)
	admin.Post("/services", servicesHandler.Create)
	admin.Put("/services/:id", servicesHandler.Update)
	admin.Delete("/services/:id", servicesHandler.Delete)
	admin.Get("/messages", messagesHandler.List)
	admin.Put("/messages/:id/read", messagesHandler.MarkRead)
	admin.Delete("/messages/:id", messagesHandler.Delete)
	admin.Get("/requests", adminRequestsHandler.List)
	admin.Get("/requests/:id", adminRequestsHandler.Get)
	admin.Put("/requests/:id/status", adminRequestsHandler.UpdateStatus)
	admin.Put("/requests/:id/progress", adminRequestsHandler.UpdateProgress)
	admin.Post("/requests/:id/milestones", adminRequestsHandler.CreateMilestone)
	admin.Put("/requests/:id/milestones/:mid", adminRequestsHandler.UpdateMilestone)
	admin.Post("/requests/:id/messages", adminRequestsHandler.SendMessage)
	admin.Post("/requests/:id/deliverables", adminRequestsHandler.AddDeliverable)
	admin.Post("/requests/:id/quote", adminRequestsHandler.SetQuote)
	admin.Post("/requests/:id/revoke-token", adminRequestsHandler.RevokeToken)
	admin.Post("/requests/:id/regenerate-token", adminRequestsHandler.RegenerateToken)
	admin.Get("/testimonials", testimonialsHandler.ListAdmin)
	admin.Post("/testimonials", testimonialsHandler.Create)
	admin.Put("/testimonials/:id", testimonialsHandler.Update)
	admin.Delete("/testimonials/:id", testimonialsHandler.Delete)
	admin.Get("/dashboard", messagesHandler.Dashboard)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
