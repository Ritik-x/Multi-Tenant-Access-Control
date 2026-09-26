package main

import (
	"log"
	"net/http"
	"time"

	"team-access-control/internal/config"
	"team-access-control/internal/database"
	"team-access-control/internal/handlers"
	"team-access-control/internal/middleware"
	"team-access-control/internal/repository"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	// =========================
	// SERVICES
	// =========================

	authService := services.NewAuthService(cfg.JWTSecret)

	// =========================
	// REPOSITORIES
	// =========================

	// User
	userRepository := repository.NewUserRepository(db)

	// Membership
	membershipRepository := repository.NewMembershipRepository(db)

	// Organization
	organizationRepository := repository.NewOrganizationRepository(db)

	// Role
	roleRepository := repository.NewRoleRepository(db)

	// Session
	sessionRepository := repository.NewSessionRepository(db)

	// Invitation
	invitationRepository := repository.NewInviationRepository(db)

	auditLogRepository := repository.NewAuditLogRepository(db)

auditLogService := services.NewAuditLogService(
	auditLogRepository,
)


auditLogHandler := handlers.NewAuditLogHandler(
		auditLogService,
)

	// RBAC
	rbacRepository := repository.NewRBACRepository(db)

	// =========================
	// SERVICES
	// =========================

	// Invitation Service
	invitationService := services.NewInvitationService(
	db,
	invitationRepository,
	roleRepository,
	membershipRepository,
	authService,
	auditLogService,
)

	// Session Service
	sessionService := services.NewSessionService(
		db,
		sessionRepository,
		authService,
	)

	// Registration Service
	registrationService := services.NewRegistrationService(
		db,
		userRepository,
		organizationRepository,
		roleRepository,
		membershipRepository,
	)

	// RBAC Service
	rbacService := services.NewRBACService(
		rbacRepository,
	)

	// ---------------------------------------


// rate limiters 


rateLimiter := middleware.NewRateLimiter(5 , time.Minute)
postRateLimiter := middleware.NewRateLimiter(5 , time.Minute)
loginRateLimiter := middleware.NewRateLimiter(
	5,
	time.Minute,
)
refreshRateLimiter := middleware.NewRateLimiter(
	10,
	time.Minute,
)
go loginRateLimiter.Cleanup()
go postRateLimiter.Cleanup()
go refreshRateLimiter.Cleanup()



	// =========================
	// HANDLERS
	// =========================

	authHandler := handlers.NewAuthHandler(
		authService,
		userRepository,
		membershipRepository,
		registrationService,
		sessionRepository,
		sessionService,

	)

	invitationHandler := handlers.NewServicesHandler(
		invitationService,
	)

	// =========================
	// ROUTER
	// =========================

	router := gin.Default()

	// Health
	router.GET("/health", rateLimiter.Middleware(),func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// =========================
	// AUTH ROUTES
	// =========================

	router.POST("/register", postRateLimiter.Middleware(),authHandler.Register)
	router.POST("/login",loginRateLimiter.Middleware(), authHandler.Login)
	router.POST("/refresh",refreshRateLimiter.Middleware(), authHandler.Refresh)
	router.POST("/logout", authHandler.Logout)




	// =========================
	// SESSION ROUTES
	// =========================

	router.GET(
		"/sessions",
		middleware.AuthMiddleware(cfg.JWTSecret),
		authHandler.GetSessions,
	)

	router.DELETE(
		"/sessions/:id",
		middleware.AuthMiddleware(cfg.JWTSecret),
		authHandler.DeleteRevoke,
	)

	router.POST(
		"/sessions-revoke-all",
		middleware.AuthMiddleware(cfg.JWTSecret),
		authHandler.RevokeAllsessions,
	)

	// =========================
	// PROTECTED TEST ROUTE
	// =========================

	router.GET(
		"/protected",
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RequiredPermission(
			rbacService,
			"users.read",
		),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "You have users.read permission",
			})
		},
	)

	// =========================
	// INVITATION ROUTES
	// =========================

	router.POST(
		"/organizations/:organizationID/invitations",
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RequiredPermission(
			rbacService,
			"team.invite",
		),
		invitationHandler.CreateInvitation,
	)

	router.POST("/invitations/accept" , middleware.AuthMiddleware(cfg.JWTSecret) , invitationHandler.AcceptInvitation)



	router.GET("/organizations/:organizationID/audit-logs",middleware.AuthMiddleware(cfg.JWTSecret), middleware.RequiredPermission(rbacService ,"audit .read" ), auditLogHandler.GetAuditLogs)
	// =========================
	// SERVER
	// =========================

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}