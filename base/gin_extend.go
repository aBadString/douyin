package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// HandlerFuncConverter 将函数 f 转换为 Gin 的路由处理函数 HandlerFunc
// f 应该是形如 func(...any) (...any, ...error) 的函数, 兼容 func(c *gin.Context)
func HandlerFuncConverter(f any) gin.HandlerFunc {
	// 0. 检查 f 是否是一个函数
	var theF reflect.Value = reflect.ValueOf(f)
	if theF.Kind() != reflect.Func {
		panic("f should be a function.")
	}

	return func(c *gin.Context) {
		// 1. 反射获取函数 f 的形参列表, 并从 HTTP 报文中获取实参

		// 函数 f 的参数数量
		paramNum := theF.Type().NumIn()

		params := make([]reflect.Value, paramNum)
		for i := 0; i < paramNum; i++ {
			var paramType reflect.Type = theF.Type().In(i)
			switch {
			case paramType == reflect.TypeOf((*gin.Context)(nil)):
				params[i] = reflect.ValueOf(c)
			case paramType.Kind() == reflect.Struct:
				// 反射创建结构体
				var param reflect.Value = reflect.New(paramType)

				// 如果这里能获取到 param 所在结构体的原始指针就好了. 那就能直接复用 Gin 的 Bind 方法
				// c.ShouldBind(param.ptr)

				//	jsonBody := make(map[string]any)
				//	_ = c.ShouldBindJSON(&jsonBody)

				bindStruct(param.Elem(), c)

				params[i] = param.Elem()
			}
		}

		// 2. 调用函数 f
		var retValues []reflect.Value = theF.Call(params)
		_ = retValues

		// 3. 返回值处理

		// 函数 f 的返回值的类型名称列表
		retTypeNameList := make([]string, 0, len(retValues))
		retValuesNum := theF.Type().NumOut()
		for i := 0; i < retValuesNum; i++ {
			ret := theF.Type().Out(i)
			retTypeNameList = append(retTypeNameList, ret.Name())
		}
		response(c, retValues, retTypeNameList)
	}
}

// bindStruct 从 HTTP 报文中将请求数据注入到 theStruct 中
func bindStruct(theStruct reflect.Value, c *gin.Context) {
	// 对结构体的每个属性
	for i := 0; i < theStruct.NumField(); i++ {
		var field reflect.StructField = theStruct.Type().Field(i)

		// 将值赋给对应的字段
		switch field.Type {
		case reflect.TypeOf(int(0)):
			val, hasVal := fetchValue(field, c)
			if hasVal {
				v, _ := strconv.ParseInt(val, 10, 64)
				theStruct.FieldByName(field.Name).SetInt(v)
			}
		case reflect.TypeOf(string("")):
			val, hasVal := fetchValue(field, c)
			if hasVal {
				theStruct.FieldByName(field.Name).SetString(val)
			}
		}
	}
}

// fetchValue 依次从 [QueryString, 表单] 中取值
func fetchValue(field reflect.StructField, c *gin.Context) (string, bool) {
	var val = ""
	for _, tagName := range []string{"query", "from"} {
		paramName := field.Tag.Get(tagName)
		if paramName == "" {
			continue
		}

		switch tagName {
		case "query":
			val = c.Query(paramName)
		case "form":
			val = c.PostForm(paramName)
		default:
		}

		if val != "" {
			break
		}
	}

	if val == "" {
		return "", false
	}
	return val, true
}

// response
// retValues 返回的值的列表
// retTypeNameList 返回值的类型名称列表
//
//	  retValues 为空时, 表示不需要包装返回数据
//	  retValues 为 (...error) 时,
//	    当 error == nil, 写入 c.JSON(http.StatusOK, resp)
//	    当 error != nil, 写入 {"status_msg":  errorList[i].Error()}
//	  retValues 为 (...data, error) 时,
//		   当 error == nil, 写入 {"data": data}
func response(c *gin.Context, retValues []reflect.Value, retTypeNameList []string) {
	if len(retValues) == 0 {
		return
	}

	// 1. 分离错误和有效数据
	dataList := make([]struct {
		name string
		data any
	}, 0, len(retValues)-1)
	errorList := make([]error, 0, 1)

	for i, ret := range retValues {
		if ret.Interface() == nil {
			continue
		}

		if ret.Type() == reflect.TypeOf((*error)(nil)).Elem() {
			errorList = append(errorList, ret.Interface().(error))
		} else {
			dataList = append(dataList, struct {
				name string
				data any
			}{name: retTypeNameList[i], data: ret.Interface()})
		}
	}

	if len(errorList) == 0 {
		// 2.1. 无错误
		resp := map[string]any{
			"status_code": 0,
			"status_msg":  "success",
		}

		switch len(dataList) {
		case 0:
			c.JSON(http.StatusOK, resp)
		default:
			// [data]
			// [data1, data2, ...]
			for i := range dataList {
				resp[camelCaseToUnder(dataList[i].name)] = dataList[i].data
			}
			c.JSON(http.StatusOK, resp)
		}
	} else {
		// 2.2. 有错误
		if len(errorList) == 1 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"status_msg":  errorList[0].Error(),
			})
		} else {
			errorMessages := make([]string, len(errorList))
			for i, err := range errorList {
				errorMessages[i] = err.Error()
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 500,
				"status_msg":  errorMessages,
			})
		}
	}
}

// 驼峰命名转下划线命名
// camelCase -> camel_case
// CamelCase -> camel_case
func camelCaseToUnder(camelCase string) string {
	s := make([]byte, 0, 2*len(camelCase))

	s = append(s, camelCase[0])
	for i := 1; i < len(camelCase); i++ {
		if ('a' <= camelCase[i-1] && camelCase[i-1] <= 'z') &&
			('A' <= camelCase[i] && camelCase[i] <= 'Z') {
			s = append(s, '_')
		}
		s = append(s, camelCase[i])
	}

	return strings.ToLower(string(s))
}
