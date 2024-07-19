Todos os metodos tem o metodo -v (verbosidade) disponivel além do metodo normal.

go run main.go http://eu.httpbin.org/get || go run main.go http://eu.httpbin.org/get?param1=value1"&"param2=value2 || 

go run main.go -X POST http://eu.httpbin.org/post || go run main.go -X POST http://eu.httpbin.org/post -d "field1=value1&field2=value2" || 
go run main.go -X POST http://eu.httpbin.org/post -H "Content-Type: application/json" -d '{"key": "value"}' ||


go run main.go -X PUT http://eu.httpbin.org/put ||

go run main.go -X DELETE http://eu.httpbin.org/delete ||

go run main.go http://eu.httpbin.org/headers -H "Custom-Header: Value" ||
