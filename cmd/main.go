package main

import (
	"fmt"
	"log"

	controller "github.com/Orken1119/Ozinshe/internal/controllers"
	pkg "github.com/Orken1119/Ozinshe/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	app, err := pkg.App()

	if err != nil {
		log.Fatal(err)
	}
	defer app.CloseDBConnection()

	ginRouter := gin.Default()

	controller.Setup(app, ginRouter)

	ginRouter.Run(fmt.Sprintf(":%s", 1136))
}
