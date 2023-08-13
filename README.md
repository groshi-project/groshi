```
                 _   _
___ ___ ___ ___| |_|_|
| . |  _| . |_ -|   | |
|_  |_| |___|___|_|_|_|
|___|
```

goddamn simple tool to keep track of your finances

Using groshi you can store, read and delete financial transactions.
Work is in progress, stand by!

## HTTP API methods overview
This table will give you some basic overview on the groshi API.

|        **HTTP method**         |        **Path**        | **Description**                                   |
|:------------------------------:|:----------------------:|---------------------------------------------------|
|             `POST`             |     `/auth/login `     | Log in and obtain an authentication token         |
|             `POST`             |     `/auth/logout`     | Log out and invalidate the authentication token   |
|             `POST`             |    `/auth/refresh`     | Refresh the authentication token                  |
|                                |                        |                                                   |
|             `POST`             |        `/user/`        | Create new user                                   |
|             `GET`              |        `/user/`        | Get information about current user                |
|             `PUT`              |        `/user/`        | Update current user                               |
|            `DELETE`            |        `/user/`        | Delete current user                               |
|                                |                        |                                                   |
|             `POST`             |    `/transaction/`     | Create new transaction                            |
|             `GET`              |  `/transaction/:uuid`  | Retrieve a transaction with specified UUID        |
|             `GET`              |    `/transaction/`     | Retrieve all transactions for given period        |
|             `GET`              | `/transaction/summary` | Retrieve summary of transactions for given period |
|             `PUT`              |  `/transaction/:uuid`  | Update transaction with specified UUID            |
|            `DELETE`            |  `/transaction/:uuid`  | Delete transaction with specified UUID            |
