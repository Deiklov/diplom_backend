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
## 3. Компании
### 3.1 Создать компанию

Запрос: `/api/v1/company` типа `POST`

Авторизация: не обязательна    
Ответ:
1.  200 ok  
```json
{
   "id":"1532bbee-8d4c-492d-a577-bae3817fa113",
   "name":"NVIDIA ORD",
   "ipo":"1999-01-22T00:00:00Z",
   "description":"",
   "country":"US",
   "ticker":"NVDA",
   "attributes":{}
}
```
2.  400 Невалидный или не существующий ticker акции
3.  422 проблемы с либой finnhub
### 3.2 Получить полную инфу о компании

Запрос: `/api/v1/company/page/:slug` типа `GET`
Авторизация: не обязательна  
Ответ:
1. 200 ok
```json
{
   "id":"2889bafb-318f-42b5-aa6d-0ca1c4e9f5e2",
   "name":"NIKOLA ORD",
   "ipo":"2018-05-15T00:00:00Z",
   "description":"Perferendis voluptatem consequatur aut sit accusantium.",
   "country":"US",
   "ticker":"NKLA",
   "logo":"https://finnhub.io/api/logo?symbol=NKLA",
   "weburl":"https://nikolamotor.com/",
   "attributes":{
      "currency":"USD",
      "exchange":"NASDAQ NMS - GLOBAL MARKET",
      "finnhubIndustry":"Machinery"
   }
}
```
2. 400 ошибка в базе  
### 3.3 Добавить в избранное

Запрос: `/api/v1/company/favorite` типа `POST`
Авторизация: обязательна  
Тело:
```json
{
   "ticker":"AAPL"
}
```
1. 200 ok
если идет повторное добавление, ответ 200, в базу не пишется
2. 400 ошибка в базе  
### 3.4 Удалить из избранного

Запрос: `/api/v1/company/favorite` типа `DELETE`
Авторизация: обязательна  
Тело:
```json
{
   "ticker":"AAPL"
}
```
1. 200 ok
если идет повторное удаление, ответ 200
2. 400 ошибка в базе  
### 3.5 Список избранных компаний

Запрос: `/api/v1/companies/favorite` типа `GET`
Авторизация: обязательна  
Ответ:
1. 200 ok
```json
[
    {
        "id": "3058717c-350f-4e09-b347-2180c645e211",
        "name": "TESLA ORD",
        "ipo": "2010-06-29T00:00:00Z",
        "description": "",
        "country": "dd",
        "ticker": "TSLA",
        "attributes": {}
    },
    {
        "id": "7039b1d1-4729-4921-8c4e-35c83fe413bf",
        "name": "Apple Inc",
        "ipo": "2004-08-19T00:00:00Z",
        "description": "",
        "country": "adad",
        "ticker": "AAPL",
        "attributes": {}
    }
]
```
если нет избранных, то возвращается пустой массив
2. 400 ошибка в базе