package main

import (
	"log"
	"net/http"
	"team-access-control/internal/config"
	"team-access-control/internal/database"
	"team-access-control/internal/handlers"
	"team-access-control/internal/middleware"
	"team-access-control/internal/repository"
	"team-access-control/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){


	if err := godotenv.Load() ; err !=nil{
			log.Println("No .env file found")
	}
	cfg := config.Load()

	db , err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
		defer db.Close()


		//services
		authService := services.NewAuthService(cfg.JWTSecret)

			//user
userRepository := repository.NewUserRepository(db)
//membership
membershipRepository := repository.NewMembershipRepository(db)
//organization 
organizationRepository := repository.NewOrganizationRepository(db)
//role
roleRepository := repository.NewRoleRepository(db)

//service role

serviceRepository :=repository.NewSessionRepository(db)

//registeration service
registrationService := services.NewRegistrationService(
    db,
    userRepository,
    organizationRepository,
    roleRepository,
    membershipRepository,

)
	


		//repo
		rbacREpository := repository.NewRBACRepository(db)
	



authHandler := handlers.NewAuthHandler(
    authService,
    userRepository,
	membershipRepository,
	registrationService,
	serviceRepository,
	
)
			// RBAC service
			rbacService := services.NewRBACService(rbacREpository)
	router := gin.Default()
	router.GET("/health" , func(c *gin.Context){
		c.JSON(http.StatusOK , gin.H{
			"status":"ok",
		})
	})

router.POST("/register", authHandler.Register)
router.POST("/login", authHandler.Login)
router.POST("/refresh",authHandler.Refresh)
	//protected test routeings

	router.GET("/protected",middleware.AuthMiddleware(cfg.JWTSecret) , middleware.RequiredPermission(
		rbacService,	"users.read",
	),


	func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "You have users.read permission",
			})
		},

	


)

_ = authService


if err := 	router.Run(":"+ cfg.Port) ; err != nil {
		log.Fatal(err)
	}
}