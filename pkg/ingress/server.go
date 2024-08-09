package ingress

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rokn/notifications-manager/pkg/config"
	_ "github.com/rokn/notifications-manager/pkg/ingress/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Notifications Manager Ingress API
// @version 1.0
// @description This is the API for the ingress service of the notifications manager

// @BasePath /api/v1

// @securityDefinitions.basic  BasicAuth

type Server interface {
	Start() error
}

type server struct {
	router   *gin.Engine
	port     int
	validate *validator.Validate
	log      *zap.Logger
	svc      Service
}

func NewServer(port int, profile config.ProfileType, svc Service, logger *zap.Logger) Server {
	if profile == config.ProfileProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	s := &server{
		router:   router,
		port:     port,
		validate: validator.New(),
		svc:      svc,
		log:      logger.With(zap.String("server", "ingress")),
	}
	public := router.Group("/api/v1")
	public.GET("/channels", s.getChannels)
	public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	private := router.Group("/api/v1")
	private.Use(gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))
	private.POST("/notification", s.postNotification)
	return s
}

func (s *server) Start() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

// @Summary Create a notification
// @Description Submits a new notification for the given channels which will be sent out by the system.
// @Param notification body NotificationDTO true "Notification object that needs to be sent"
// @Accept json
// @Produce json
// @Success 200 {object} NotificationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401
// @Security BasicAuth
// @Router /notification [post]
func (s *server) postNotification(c *gin.Context) {
	c.MustGet(gin.AuthUserKey)
	notification := NotificationDTO{}
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	s.log.Debug("received notification", zap.Any("notification", notification))
	err := s.svc.TransmitNotification(c.Request.Context(), notification)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "notification received"})
}

// @Summary Get the list of channels
// @Description Returns the list of channels that are available in the system.
// @Produce json
// @Success 200 {object} ChannelsResponse
// @Failure 500 {object} ErrorResponse
// @Router /channels [get]
func (s *server) getChannels(c *gin.Context) {
	channels, err := s.svc.GetChannels(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"channels": channels})
}
