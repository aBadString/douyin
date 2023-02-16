package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ExampleHandlerFuncConverter() {
	var ginEngine *gin.Engine = gin.Default()
	ginEngine.GET("/example/default", DefaultHandler)
	ginEngine.GET("/example/plus", HandlerFuncConverter(PlusHandler))
	ginEngine.GET("/example/default_compatible", HandlerFuncConverter(DefaultHandler))
	_ = ginEngine.Run()
	// Output:
	// http://localhost:8080/example/default?name=Apple&age=10
	// http://localhost:8080/example/plus?name=Apple&age=10
	// http://localhost:8080/example/default_compatible?name=Apple&age=10
}

type User struct {
	Name string `query:"name"`
	Age  int    `query:"age"`
}

func DefaultHandler(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")

	user := User{Name: name}
	v, _ := strconv.ParseInt(age, 10, 64)
	user.Age = int(v) + 1

	c.JSON(http.StatusOK, user)
}

func PlusHandler(user User) User {
	user.Age++
	return user
}
