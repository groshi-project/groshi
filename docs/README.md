# groshi JSON HTTP API documentation
There are only two entities in groshi: **user** and **transaction**.
And there are several appropriate methods to manage them.

## Quick notes
* All API methods should be called using `POST`
* All parameters should be placed in the request body
* groshi HTTP server always returns `200` status code, no matter if any errors occurred 
* In addition to errors mentioned in specific methods in the documentation, the following error tags always can also be returned: `invalid_request` and `internal_server_error`
* Don't be afraid to read source code if you are missing any information from the documentation todo

<details>
    <summary>Why not REST?</summary>
    Lorem ipsum dolor sit amet, consectetur adipisicing elit. Ab adipisci at aut est expedita fuga officia perferendis? Assumenda dicta dolore ducimus, et facilis iure, iusto, natus nesciunt numquam rem veniam?
</details>

## Possible responses
There are only two kinds of API responses: _Success response_ and _Error response_.

### Success response example:
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
* `error_tag` - tag of the error
    
    Indicates generalized reason of error. Can be:
    * `invalid_request` - when request did not pass validation
    * `unauthorized` - when request is not unauthorized (when it has to be)
    * `internal_server_error` - when any internal server error happens
    * `access_denied` - when user have no access to resource
    * `conflict` - when request causes any kind of conflict
    * `object_not_found` - when object was not found
* `error_origin` - origin of the error (can be `client` or `server`)
* `error_details` - useful error details

---
## API methods
### Users
<details>
<summary><code>POST</code><code><b>/user/create</b></code><code>(creates new user)</code></summary>

##### Parameters
|    name    | data type | required | description              |
|:----------:|:---------:|:--------:|--------------------------|
| `username` |  string   |   yes    | Username of the new user |
| `password` |  string   |   yes    | Password of the new user |

##### Success response
Simple empty success response is returned.
```json
{
  "success": true,
  "data": {}
}
```

##### Error responses
Error responses with the following _error tags_ may be returned:

| `error_tag` | case                                                    |
|-------------|---------------------------------------------------------|
| `conflict`  | Username you've provided is already taken by other user | 


##### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST $hostname/user/create username="username" password="password"
```
</details>

<details>
<summary><code>POST</code><code><b>/user/auth</b></code><code>(authorizes user)</code></summary>

##### Parameters
|    name    | data type | required | description   |
|:----------:|:---------:|:--------:|---------------|
| `username` |  string   |   yes    | User username |
| `password` |  string   |   yes    | User password |


##### Success response
Authorization token is returned.
```json
{
  "success": true,
  "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImppZWdnaWkiLCJleHAiOjE2ODEyNDEyMjR9.AvMzAVJpVq4ZMeUDWMRk-vM1KkDutmL-Bje44XsaCNc"
  }
}
```

##### Error responses
Error responses with the following _error tags_ may be returned:

| error_tag          | case                                    |
|--------------------|-----------------------------------------|
| `object_not_found` | User with such `username` was not found | 
| `access_denied`    | Invalid password has been provided      | 


##### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST $hostname/user/auth username="username" password="password"
```
</details>

<details>
<summary><code>POST</code><code><b>/user/read</b></code><code>(gets information about current user)</code></summary>

##### Parameters
|  name   | data type | required | description         |
|:-------:|:---------:|:--------:|---------------------|
| `token` |  string   |   yes    | Authorization token |


##### Success response
Username is returned.
```json
{
  "success": true,
  "data": {
    "username": "jieggii"
  }
}
```

##### Error responses
Error responses with the following _error tags_ may be returned:

| error_tag          | case                                             |
|--------------------|--------------------------------------------------|
| `object_not_found` | The user you authorized under has not been found | 


##### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST $hostname/user/read token=$TOKEN
```
</details>

<details>
<summary><code>POST</code><code><b>/user/update</b></code><code>(updates current user)</code></summary>

##### Parameters
|      name      | data type | required | description                       |
|:--------------:|:---------:|:--------:|-----------------------------------|
|    `token`     |  string   |   yes    | Authorization token               |
| `new_username` |  string   |    no    | New username for the current user |
| `new_password` |  string   |    no    | New password for the current user |

(you are required to use at least one of two parameters: `new_username` or `new_password`)

##### Success response
Simple empty success response is returned.
```json
{
  "success": true,
  "data": {}
}
```

##### Error responses
Error responses with the following _error tags_ may be returned:

| error_tag        | case                                             |
|------------------|--------------------------------------------------|
| `user_not_found` | The user you authorized under has not been found | 
| `conflict`       | New username chosen by you is already taken      | 


##### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST $hostname/user/update token=$TOKEN new_username="new-username" new_password="new-password"
```
</details>

<details>
<summary><code>POST</code><code><b>/user/delete</b></code><code>(deletes current user)</code></summary>

##### Parameters
|  name   | data type | required | description         |
|:-------:|:---------:|:--------:|---------------------|
| `token` |  string   |   yes    | Authorization token |

##### Success response
Empty success response is returned.
```json
{
  "success": true,
  "data": {}
}
```

##### Error responses
Error responses with the following _error tags_ may be returned:

| error_tag          | case                                             |
|--------------------|--------------------------------------------------|
| `object_not_found` | The user you authorized under has not been found | 


##### Example request using [httpie](https://github.com/httpie/httpie)
```shell
http POST $hostname/user/delete token=$TOKEN
```
</details>


------------------------------------------------------------------------------------------

