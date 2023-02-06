package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// HandlerFuncConverter 将函数 f 转换为 Gin 的路由处理函数 HandlerFunc
// f 应该是形如 func(...any) (...any, ...error) 的函数, 兼容 func(c *gin.Context).
//
// 将向函数 f 中注入结构体类型参数的请求数据, 暂时不支持嵌套结构体. 结构体的属性仅支持 int 和 string.
// 支持注入 Gin 上下文对象的指针 *gin.Context.
//
// 将对函数 f 的返回值进行包装, 并写入响应数据中.
// 当 f 返回值列表为空时, 不进行返回值包装 (为兼容 func(c *gin.Context));
// 当 f 返回值列表为 (error) 时, 根据 error 是否为 nil, 响应 200 或者 error message;
// 一般的 f 返回值列表应是 (data1, data2, ..., error1, error2, ...):
// 一旦存在 error 不为 nil, 返回格式如下:
//
//	{
//	  "status_code": 500,
//	  "status_msg": "error message"
//	}
//
// 或者
//
//	{
//	  "status_code": 500,
//	  "status_msg": ["error message 1", "error message 2", ...]
//	}
//
// 所有 error 都为 nil, 返回格式如下:
//
//	{
//		"status_code": 0,
//		"status_msg": "success",
//		"data1": {...},
//		"data2": {...},
//		...
//	}
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
				var param reflect.Value = reflect.New(paramType).Elem()

				// 如果这里能获取到 param 所在结构体的原始指针就好了. 那就能直接复用 Gin 的 Bind 方法
				// c.ShouldBind(param.ptr)

				//	jsonBody := make(map[string]any)
				//	_ = c.ShouldBindJSON(&jsonBody)

				bindStruct(param, c)

				params[i] = param
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
			retType := theF.Type().Out(i)
			retTypeNameList = append(retTypeNameList, retType.Name())
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
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, hasVal := fetchValue(field.Tag, c)
			if hasVal {
				switch val.(type) {
				case int:
					theStruct.FieldByName(field.Name).SetInt(int64(val.(int)))
				case string:
					v, _ := strconv.ParseInt(val.(string), 10, 64)
					theStruct.FieldByName(field.Name).SetInt(v)
				}
			}
		case reflect.String:
			val, hasVal := fetchValue(field.Tag, c)
			if hasVal {
				if v, ok := val.(string); ok {
					theStruct.FieldByName(field.Name).SetString(v)
				}
			}
		}
	}
}

// fetchValue 依次从 [QueryString, 表单] 中取值
func fetchValue(fieldTag reflect.StructTag, c *gin.Context) (any, bool) {
	var val any = nil
	for _, tagName := range []string{"query", "form", "context"} {
		paramName := fieldTag.Get(tagName)
		if paramName == "" {
			continue
		}

		switch tagName {
		case "query":
			val = c.Query(paramName)
		case "form":
			val = c.PostForm(paramName)
		case "context":
			val, _ = c.Get(paramName)
		default:
		}

		if val != nil {
			break
		}
	}

	if val == nil {
		return nil, false
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
	errorList := make([]struct {
		code    int
		message string
	}, 0, 1)

	for i, ret := range retValues {
		if ret.Interface() == nil {
			continue
		}

		if ret.Type() == reflect.TypeOf((*error)(nil)).Elem() {
			err := ret.Interface().(error)

			var code int
			switch err.(type) {
			case *serviceError:
				code = ret.Interface().(*serviceError).GetCode()
			default:
				code = http.StatusInternalServerError
			}

			errorList = append(errorList, struct {
				code    int
				message string
			}{code: code, message: err.Error()})
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
			c.JSON(errorList[0].code, gin.H{
				"status_code": errorList[0].code,
				"status_msg":  errorList[0].message,
			})
		} else {
			errorMessages := make([]string, len(errorList))
			for i, err := range errorList {
				errorMessages[i] = err.message
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": http.StatusInternalServerError,
				"status_msg":  errorMessages,
			})
		}
	}
}

// 驼峰命名转下划线命名
// camelCase -> camel_case
// CamelCase -> camel_case
func camelCaseToUnder(camelCase string) string {
	if len(camelCase) <= 0 {
		return ""
	}

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
