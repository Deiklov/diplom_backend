## 1. Авторизация
### 1.1 Логин

Запрос: `/api/v1/login` типа `POST`
required: login: email format, password:ascii

Тело запроса:
```json
{
    "email": "string", 
    "password": "string"
}
```


Ответ:
1. 200 OK+вернет jwttoken
```json
{
    "token": "string"
}
```

### 1.2 Регистрация

Запрос: `/api/v1/user/` типа `POST`

Тело запроса:  
name опциональное  
```json
{
  "email": "string",
  "name": "string",
  "password": "string"
}
```
Ответ:
1. 200 Created+вернет jwttoken
```json
{
    "token": "string"
}
```
2. 400 Невалидныые данные(плохой json или поля не прошли валидацию)
3. 409 Conflict(уже есть такой юзер)


## 2. Профиль
### 2.1 Получение информации профиля

Запрос: `/api/v1/user` типа `GET`

Ответ:
1. 200 ok  
```json
{
    "id": "e317d6f9-be58-49ac-a57c-866cbef8f83f",
    "name": "kek",
    "login": "kek2101",
    "image": "/static/image/avatar/default.png",
    "email": "alexloh500@mail.ru",
    "created_at": "2020-01-01T00:00:00Z"
}
```
2. 400 unauthorized или юзера не существует  
### 2.2 Апдейт профиля

Запрос: `/api/v1/user` типа `PUT`

Ответ:
Ответ:
1. 200 ok
```json
{
    "id": "e317d6f9-be58-49ac-a57c-866cbef8f83f",
    "name": "kek",
    "image": "/static/image/avatar/default.png",
    "email": "alexloh500@mail.ru",
    "created_at": "2020-01-01T00:00:00Z"
}
```
2. 400 unauthorized или юзера не существует  