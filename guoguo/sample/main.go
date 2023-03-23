package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sample/router"
	"sample/source"
	"sample/source/middleware"
)

var logger *logrus.Logger
var route *gin.Engine
var config *viper.Viper

func init() {
	logger = source.GetLogger()
	route = gin.New()
	config = source.GetConfig()
}

func GetRoute() *gin.Engine {
	return route
}

func main() {
	logger.Infoln("server is running...")
	logger.Infoln(fmt.Sprintf("user config: %s", config.Get("config")))
	route.Use(gin.Recovery(), middleware.LoggerMiddleware())
	router.RegisterRouter(route)
	port := config.GetString("service.port")
	err := route.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infoln("server exited...")
	mongoClient := source.GetMongoClient()
	_ = mongoClient.Disconnect(source.NewCtx())
	return
}