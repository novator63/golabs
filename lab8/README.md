Команды для проверки работоспособности каждого этапа разработки

lab8-1,2:

curl -X GET http://localhost:8080/users

[
  {"id":1,"name":"Anna","age":25},
  {"id":2,"name":"Ivan","age":30}
]

curl -X GET http://localhost:8080/users/1

{"id":1,"name":"Anna","age":25}

curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name":"Maria","age":22}'

{"id":3,"name":"Maria","age":22}

curl -X PUT http://localhost:8080/users/3 \
     -H "Content-Type: application/json" \
     -d '{"name":"Maria Petrova","age":23}'

{"id":3,"name":"Maria Petrova","age":23}

curl -X DELETE http://localhost:8080/users/3

HTTP-код 204 No Content — тело ответа пустое.

lab8-3

curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name":"","age":-1}'

{"error":"Имя не может быть пустым"}
{"error":"Возраст должен быть положительным числом"}