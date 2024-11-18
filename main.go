package main

import (
	"api_journal/core/server"
	"api_journal/core/server/shared"
	"api_journal/core/service"
	_ "api_journal/docs"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Golang Documentação API 2.5
// @version 1.0
// @description Esta é uma API simples usando Gin e Swagger.
// @host localhost:4000
func main() {

	go func() {
		r := gin.Default()

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		if err := r.Run(":2345"); err != nil {
			log.Fatalf("Erro ao iniciar o servidor Gin: %v", err)
		}
	}()

	go func() {
		meuBanco, _ := service.InitConnection()
		defer meuBanco.Close()
		shared.SetDB(meuBanco)
		// go func() {
		// 	for {
		// 		scrapperService.RunScrapper()
		// 		time.Sleep(8 * time.Hour)
		// 	}
		// }()
		server.MainServer()
	}()

	select {}
}
