package main

import (
	"fmt"
	
	tinyUrl "shorty-challenge/internal/tiny_url"
	"shorty-challenge/pkg/response"
	"shorty-challenge/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	portDefault int = 1234
)

func main() {
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
	err := r.Run(fmt.Sprintf(":%d", portDefault))
	if err != nil {
		panic(err)
	}
}
