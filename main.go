package main

import (
	"flag"
	"fmt"
	tinyUrl "test_amartha_muhammad_huzair/internal/tiny_url"
	"test_amartha_muhammad_huzair/pkg/response"
	"test_amartha_muhammad_huzair/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	//initiate flagging on the command line argument
	//using this example "./tiny_url -port=<PORT>" or "./tiny_url -port <PORT>"
	//or using the go run method example "go run main.go -port=<PORT>" or "go run main.go -port <PORT>"
	port := flag.String("port", "8080", "Port")
	//and then parse the flagging option
	flag.Parse()
	//init the router engine
	r := gin.Default()
	//custom response for not matching routes
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, utils.NotMatchingAnyRoute, utils.NotFound)
		return
	})
	//initiate module tiny_url
	tinyUrl.Routes(r)
	//set the default value if the user not using the argument.
	if port == nil {
		*port = "8080"
	}
	err := r.Run(fmt.Sprintf(":%s", *port))
	if err != nil {
		panic(err)
	}
}
