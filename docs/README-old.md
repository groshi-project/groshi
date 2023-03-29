# groshi JSON HTTP API documentation
There are only two entities in groshi: **user** and **transaction**.
And there are several appropriate methods to manage them.

## Quick notes
> * All API methods should be called using **POST** HTTP method.
> * All parameters should be placed in the **request body**.
> * groshi HTTP server **always** returns `200` status code, no matter if any errors occurred. The only exception is `404` status code when nonexistent API method is used.

<details>
    <summary>Why not REST?</summary>
    <strike>
    Well... I am not sure that REST is so useful as everyone thinks it is.
    At least for this project, I think REST would complicate architecture and
    client-side code without providing any significant benefits.
    That's why I deviate from the REST standards.
    Of course, you may argue with me in the issues if you want so!    
    </strike>
    
    Some time passed, I thought, googled and asked a little about that and understood, 
    that I was mistaken. But... Code is written, someday I will rewrite it in REST, but
    for now let it be like that...
</details>

## Possible responses
There are only two kinds of API responses: **successful response** and **error response**.

### Successful response example:
```json
{
  "success": true,
  "data": {
      "uuid": "f5acbb59-6e23-452f-af27-637e9bce1cad"
  }
}
```
`success` field's value is `true` and all necessary data is placed in the `data` object.

### Error response example:
```json
{
  "success": false,
  "error_tag": "object_not_found",
  "error_origin": "client",
  "error_details": "Transaction not found."
} 
```
`success` field's value is `false`, information about error is placed in the following fields:
* `error_tag` - error tag (indicates general reason of error)
* `error_origin` - who is to blame for the error (`client` or `server`)
* `error_details` - useful error details which may help with problem solution

#### Possible values of `error_tag`: [(reference)](../internal/http/ghttp/schema/schema.go#L13)
|        error_tag        |                      case                      |
|:-----------------------:|:----------------------------------------------:|
|    `invalid_request`    |        Request did not pass validation         |
|     `unauthorized`      | Request is not unauthorized (but it has to be) |
| `internal_server_error` |             Internal server error              |
|     `access_denied`     |    Current user have no access to resource     |
|       `conflict`        |            Request causes conflict             |
|   `object_not_found`    |              Object was not found              |

#### Possible values of `error_origin`: [(reference)](../internal/http/ghttp/schema/schema.go#L6)

| error_origin |        case         |
|:------------:|:-------------------:|
|   `client`   | Client caused error |
|   `server`   | Server caused error |


---
## API methods
> In addition to errors mentioned in specific API methods in the documentation, the following error tags always can also be returned:
> * `invalid_request` - when request did not pass validation
> * `internal_server_error` - when something is broken on the server side

### User
<details>
<summary><code>POST</code> <code><b>/user/create</b></code> <code>(creates new user)</code></summary>

#### Parameters
|    name    | data type | required | description              |
|:----------:|:---------:|:--------:|--------------------------|
| `username` |  string   |   yes    | Username of the new user |
| `password` |  string   |   yes    | Password of the new user |

#### Successful response
Username is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "username": "jieggii"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

| error_tag  | case                                                    |
|------------|---------------------------------------------------------|
| `conflict` | Username you've provided is already taken by other user | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/user/create username="username" password="password"
```
</details>

<details>
<summary><code>POST</code> <code><b>/user/auth</b></code> <code>(authorizes user)</code></summary>

#### Parameters
|    name    | data type | required | description   |
|:----------:|:---------:|:--------:|---------------|
| `username` |  string   |   yes    | User username |
| `password` |  string   |   yes    | User password |


#### Successful response
Authorization token is returned.
```json
{
  "success": true,
  "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImppZWdnaWkiLCJleHAiOjE2ODEyNDEyMjR9.AvMzAVJpVq4ZMeUDWMRk-vM1KkDutmL-Bje44XsaCNc"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

| error_tag          | case                                    |
|--------------------|-----------------------------------------|
| `object_not_found` | User with such `username` was not found | 
| `access_denied`    | Invalid password has been provided      | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/user/auth username="username" password="password"
```
</details>

<details>
<summary><code>POST</code> <code><b>/user/read</b></code> <code>(gets information about current user)</code></summary>

#### Parameters
|  name   | data type | required | description         |
|:-------:|:---------:|:--------:|---------------------|
| `token` |  string   |   yes    | Authorization token |


#### Successful response
Username is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "username": "jieggii"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|     error_tag      |                         case                          |
|:------------------:|:-----------------------------------------------------:|
| `object_not_found` | The user you are authorized under has not been found  | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/user/read token=$TOKEN
```
</details>

<details>
<summary><code>POST</code> <code><b>/user/update</b></code> <code>(updates current user)</code></summary>

#### Parameters
|      name      | data type | required |            description            |
|:--------------:|:---------:|:--------:|:---------------------------------:|
|    `token`     |  string   |   yes    |        Authorization token        |
| `new_username` |  string   |    no    | New username for the current user |
| `new_password` |  string   |    no    | New password for the current user |

(you are required to use at least one of two parameters: `new_username` or `new_password`)

#### Successful response
New username of updated user is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "username": "jieggii"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|     error_tag     |                         case                         |
|:-----------------:|:----------------------------------------------------:|
| `user_not_found`  | The user you are authorized under has not been found | 
|    `conflict`     |     New username chosen by you is already taken      | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/user/update token=$TOKEN new_username="new-username" new_password="new-password"
```
</details>

<details>
<summary><code>POST</code> <code><b>/user/delete</b></code> <code>(deletes current user)</code></summary>

#### Parameters
|  name   | data type | required | description         |
|:-------:|:---------:|:--------:|---------------------|
| `token` |  string   |   yes    | Authorization token |

#### Successful response
Username of deleted user is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "username": "jieggii"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|     error_tag      |                         case                         |
|:------------------:|:----------------------------------------------------:|
| `object_not_found` | The user you are authorized under has not been found | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/user/delete token=$TOKEN
```

</details>

------------------------------------------------------------------------------------------

### Transaction
<details>
<summary><code>POST</code> <code><b>/transaction/create</b></code> <code>(creates new transaction owned by current user)</code></summary>

#### Parameters
|     name      |       data type       | required |                       description                       |
|:-------------:|:---------------------:|:--------:|:-------------------------------------------------------:|
|    `token`    |        string         |   yes    |                   Authorization token                   |
|   `amount`    |         float         |   yes    | Amount of transaction (can be 0, positive and negative) |
|  `currency`   | string ([currency]()) |   yes    |                 Currency of transaction                 |
| `description` |        string         |    no    |               Description of transaction                |
|    `date`     |  string (rfc format)  |    no    |                   Date of transaction                   |

#### Successful response
UUID of created transaction is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "uuid": "03ef6901-4ebb-4952-bb35-a98fcf502c83"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|      error_tag      |                           case                           |
|:-------------------:|:--------------------------------------------------------:|
| `object_not_found`  |   The user you are authorized under has not been found   |

#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/transaction/create token=$TOKEN amount:=2.53 currency=EUR description="Some donuts for my beloved gf" date="2023-03-12 12:57:07.850123 +00:00"
```

</details>

<details>
<summary><code>POST</code> <code><b>/transaction/read</b></code> <code>(gets information about transaction owned by current user)</code></summary>

#### Parameters
|  name   | data type | required |     description     |
|:-------:|:---------:|:--------:|:-------------------:|
| `token` |  string   |   yes    | Authorization token |
| `uuid`  |  string   |   yes    | UUID of transaction |

#### Successful response
Information about transaction is returned in the `data` object.
```json
{
  "data": {
    "amount": 21.24,
    "created_at": "2023-03-12T14:02:05.156883+02:00",
    "currency": "EUR",
    "date": "2023-03-12T14:02:05.156883+02:00",
    "description": "Bought some food from Prisma",
    "owner": "jieggii",
    "updated_at": null,
    "uuid": "d089c593-2b95-489c-b8a2-ad0ecaa4e44c"
  },
  "success": true
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|     error_tag      |                                  case                                   |
|:------------------:|:-----------------------------------------------------------------------:|
| `object_not_found` |          The user you are authorized under has not been found           | 
|  `access_denied`   | You are trying to get information about transaction owned by other user | 

#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/transaction/read token=$TOKEN uuid="d089c593-2b95-489c-b8a2-ad0ecaa4e44c"
```

</details>

<details>
<summary><code>POST</code> <code><b>/transaction/update</b></code> <code>(updates transaction owned by current user)</code></summary>

##### Parameters
|       name        |       data type       | required |                         description                         |
|:-----------------:|:---------------------:|:--------:|:-----------------------------------------------------------:|
|      `token`      |        string         |   yes    |                     Authorization token                     |
|      `uuid`       |        string         |   yes    |           UUID of transaction you want to update            |
|   `new_amount`    |         float         |    no    | New amount of transaction (can be 0, positive and negative) |
|  `new_currency`   | string ([currency]()) |    no    |                 New currency of transaction                 |
| `new_description` |        string         |    no    |               New description of transaction                |
|    `new_date`     |  string (rfc format)  |    no    |                   New date of transaction                   |
At least one of the following fields is required: `new_amount`, `new_currency`, `new_description`, `new_date`.

#### Successful response
UUID of updated transaction is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "uuid": "03ef6901-4ebb-4952-bb35-a98fcf502c83"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|      error_tag      |                           case                            |
|:-------------------:|:---------------------------------------------------------:|
| `object_not_found`  |   The user you are authorized under has not been found    | 
|   `access_denied`   | You are trying to update transaction owned by other user  | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/transaction/update token=$TOKEN  uuid=03ef6901-4ebb-4952-bb35-a98fcf502c83 new_amount:=2.53 new_currency=USD new_description="New description!" new_date="2023-03-12 12:57:07.850123 +00:00"
```
</details>

<details>
<summary><code>POST</code> <code><b>/transaction/delete</b></code> <code>(deletes transaction owned by current user)</code></summary>

#### Parameters
|  name   | data type | required |              description               |
|:-------:|:---------:|:--------:|:--------------------------------------:|
| `token` |  string   |   yes    |          Authorization token           |
| `uuid`  |  string   |   yes    | UUID of transaction you want to delete |

#### Successful response
UUID of deleted transaction is returned in the `data` object.
```json
{
  "success": true,
  "data": {
    "uuid": "03ef6901-4ebb-4952-bb35-a98fcf502c83"
  }
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|      error_tag      |                           case                            |
|:-------------------:|:---------------------------------------------------------:|
| `object_not_found`  |   The user you are authorized under has not been found    | 
|   `access_denied`   | You are trying to delete transaction owned by other user  | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/transaction/delete token=$TOKEN uuid=03ef6901-4ebb-4952-bb35-a98fcf502c83
```
</details>

<details>
<summary><code>POST</code> <code><b>/transaction/list</b></code> <code>(lists transactions owned by current user for given period)</code></summary>

#### Parameters
|  name   |    data type     | required |               description                |
|:-------:|:----------------:|:--------:|:----------------------------------------:|
| `token` |      string      |   yes    |           Authorization token            |
| `since` | string (ISO8601) |   yes    |           Beginning of period            |
| `until` | string (ISO8601) |    no    | End of period (defaults to current date) |

#### Successful response
List of transactions is returned in the `data` object.
```json
{
  "data": {
    "count": 2,
    "transactions": [
      {
        "amount": 5,
        "created_at": "2023-03-16T21:49:04.387814+02:00",
        "currency": "EUR",
        "date": "2023-03-16T21:49:04.387814+02:00",
        "description": "",
        "owner": "jieggii",
        "updated_at": null,
        "uuid": "98ca6669-cc4c-49ea-be6e-039d8cd86b58"
      },
      {
        "amount": 10,
        "created_at": "2023-03-16T21:49:08.240842+02:00",
        "currency": "EUR",
        "date": "2023-03-16T21:49:08.240842+02:00",
        "description": "",
        "owner": "jieggii",
        "updated_at": null,
        "uuid": "d5d9818b-8dce-4b2d-94ef-142065e199b5"
      }
    ]
  },
  "success": true
}
```

#### Error responses
Error responses with the following _error tags_ may be returned:

|      error_tag      |                          case                           |
|:-------------------:|:-------------------------------------------------------:|
| `object_not_found`  |  The user you are authorized under has not been found   | 
|   `access_denied`   | You are trying to list transactions owned by other user | 


#### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST 127.0.0.1:8080/transaction/list token=$TOKEN since="2023-03-15T21:52:08+02:00" until="2023-03-17T21:52:26+02:00"
```
</details>
