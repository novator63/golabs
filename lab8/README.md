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

lab8-4

curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Ivan","age":25}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Petr","age":30}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Anna","age":22}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Sergey","age":30}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Ilya","age":25}'

curl "http://localhost:8080/users?page=1&limit=2"

→ Вернёт первых двух пользователей 

curl "http://localhost:8080/users?page=1&limit=2"

→ Вернёт следующих двух пользователей

curl "http://localhost:8080/users?name=an"

→ вернёт всех, где имя содержит “an” (например, Ivan, Anna)

curl "http://localhost:8080/users?age=30"

→ вернёт всех пользователей, которым 30 лет