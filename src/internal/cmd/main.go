package main

import (
	"github.com/D-D-EINA-Calendar/CalendarServer/docs"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/services/issue"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/services/monitoring"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/services/space"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/services/users"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/handlers"
	issuerepositorymemory "github.com/D-D-EINA-Calendar/CalendarServer/src/internal/repositories/Memory/IssueRepository"
	spacerepositorymemory "github.com/D-D-EINA-Calendar/CalendarServer/src/internal/repositories/Memory/spaceRepository"
	usersrepositorymemory "github.com/D-D-EINA-Calendar/CalendarServer/src/internal/repositories/Memory/usersRepository"
	monitoringrepositoryrabbitamq "github.com/D-D-EINA-Calendar/CalendarServer/src/internal/repositories/RabbitAMQ/monitoringRepository"

	connection "github.com/D-D-EINA-Calendar/CalendarServer/src/pkg/connect"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/pkg/constants"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var rabbitConn connection.Connection

func config() (handlers.HTTPHandler, error) {
	//Conexión con rabbit
	var err error
	rabbitConn, err = connection.New(constants.AMQPURL)
	if err != nil {
		//TODO
	}
	chMonitoring, err := rabbitConn.NewChannel()
	if err != nil {
		//TODO
	}
	monitoringRepo := monitoringrepositoryrabbitamq.New(chMonitoring)
	//TODO canal
	_, err = rabbitConn.NewChannel()
	if err != nil {
		//TODO
	}
	//spaceRepoAMQ, _ := spacerepositoryrabbitamq.New(chSpaces)
	spaceRepo := spacerepositorymemory.New()
	usersRepo := usersrepositorymemory.New()
	issuesRepo := issuerepositorymemory.New()

	return handlers.HTTPHandler{
		Monitoring: monitoring.New(monitoringRepo),
		Users:      users.New(usersRepo),
		Spaces:     space.New(spaceRepo),
		Issues:     issue.New(issuesRepo),
	}, nil

}

//SetupRouter is a func which bind each uri with a handler function
func SetupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler, err := config()
	if err != nil {
		//TODO
	}
	r.GET(constants.PING_URL, handler.Ping)
	r.GET(constants.LOGIN, handler.Login)
	r.GET(constants.FILTER_SPACES, handler.FilterBy)
	r.GET(constants.REQUEST_INFO_SLOTS, handler.RequestInfoSlots)
	r.POST(constants.RESERVE_SPACE, handler.Reserve)
	r.GET(constants.RESERVE_BATCH, handler.ReserveBatch)

	r.GET(constants.CANCEL_RESERVE, handler.CancelReserve)
	r.GET(constants.DELETE_ISSUE, handler.DeleteIssue)
	r.POST(constants.CREATE_ISSUE, handler.CreateIssue)
	r.GET(constants.MODIFY_ISSUE, handler.ChangeStateIssue)
	r.GET(constants.GET_ALL_ISSUES, handler.GetAllIssues)
	r.GET(constants.GET_RESERVES_USER, handler.GetReservesOwner)

	return r
}

func main() {
	// · Swagger ·
	docs.SwaggerInfo.Title = "API UNIZAR calendar and schedule"
	docs.SwaggerInfo.Description = "This is API for managing and visulizating the calendar and schedule of Unizar."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := SetupRouter()
	defer rabbitConn.Disconnect()

	r.Run(":8080")
}
