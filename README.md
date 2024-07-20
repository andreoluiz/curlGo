para rodar o programa, apenas estar na pasta do arquivo main.go, e rodar o comando "go run" passar o nome do arquivo "main.go" e por fim ao menos a URL,
ficando no minimo desta forma: go run main.go http://eu.httpbin.org

A seguir alguns comandos possiveis com o programa:

Todos os metodos tem o metodo -v (verbosidade) disponivel além do metodo normal caso deseje usa-lo é o primeiro que deve ser inserido logo apos o main.go.

go run main.go http://eu.httpbin.org/get || go run main.go http://eu.httpbin.org/get?param1=value1"&"param2=value2 || 

go run main.go -X POST http://eu.httpbin.org/post || go run main.go -X POST http://eu.httpbin.org/post -d "field1=value1&field2=value2" || 
go run main.go -X POST http://eu.httpbin.org/post -H "Content-Type: application/json" -d '{"key": "value"}' ||


go run main.go -X PUT http://eu.httpbin.org/put ||

go run main.go -X DELETE http://eu.httpbin.org/delete ||

go run main.go http://eu.httpbin.org/headers -H "Custom-Header: Value" ||

go run main.go http://eu.httpbin.org/status/200 ||
go run main.go http://eu.httpbin.org/status/404 ||
go run main.go http://eu.httpbin.org/status/500 ||

go run main.go http://eu.httpbin.org/ip ||
