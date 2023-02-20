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

func ExampleParam() {
	var ginEngine *gin.Engine = gin.Default()
	ginEngine.Use(func(context *gin.Context) {
		context.Set("context", "hello")
	}).GET("/example/param", HandlerFuncConverter(ParamHandler))
	_ = ginEngine.Run()
	// Output:
	// http://localhost:8080/example/param?query=abc
}

type Param struct {
	Query   string `query:"query"`
	Context string `context:"context"`
}

func ParamHandler(p Param) Param {
	return p
}

func ExampleMultipleReturnValues() {
	var ginEngine *gin.Engine = gin.Default()
	ginEngine.GET("/example/multiple_return", HandlerFuncConverter(MultipleReturnHandler))
	_ = ginEngine.Run()
	// Output:
	// http://localhost:8080/example/multiple_return
}

type Code int
type DataData struct {
	Name string
	Age  int
}

func MultipleReturnHandler() (Code, DataData) {
	return 1, DataData{
		Name: "name",
		Age:  19,
	}
}

func ExampleError() {
	var ginEngine *gin.Engine = gin.Default()
	ginEngine.GET("/example/error", HandlerFuncConverter(ErrorHandler))
	_ = ginEngine.Run()
	// Output:
	// http://localhost:8080/example/error
}

func ErrorHandler() (DataData, error) {
	return DataData{
		Name: "name",
		Age:  19,
	}, NewError(418, "出错了！")
}
