# groshi
groshi - goddamn simple tool to keep track of your finances.

Using groshi you can store, read and delete financial transactions.


## HTTP API methods overview
These tables will give you some basic overview of the groshi API.

- API methods related to **authorization**:

    |        **HTTP method**         |        **Path**        | **Description**                                   |
    |:------------------------------:|:----------------------:|---------------------------------------------------|
    |             `POST`             |     `/auth/login `     | Log in and obtain an authentication token         |
    |             `POST`             |     `/auth/logout`     | Log out and invalidate the authentication token   |
    |             `POST`             |    `/auth/refresh`     | Refresh the authentication token                  |


- API methods related to **users**:
    
    |        **HTTP method**         |        **Path**        | **Description**                                   |
    |:------------------------------:|:----------------------:|---------------------------------------------------|
    |             `POST`             |        `/user/`        | Create new user                                   |
    |             `GET`              |        `/user/`        | Get information about current user                |
    |             `PUT`              |        `/user/`        | Update current user                               |
    |            `DELETE`            |        `/user/`        | Delete current user                               |

- API methods related to **transactions**:
    
    |        **HTTP method**         |        **Path**        | **Description**                                   |
    |:------------------------------:|:----------------------:|---------------------------------------------------|
    |             `POST`             |    `/transaction/`     | Create new transaction                            |
    |             `GET`              |  `/transaction/:uuid`  | Retrieve a transaction with specified UUID        |
    |             `GET`              |    `/transaction/`     | Retrieve all transactions for given period        |
    |             `GET`              | `/transaction/summary` | Retrieve summary of transactions for given period |
    |             `PUT`              |  `/transaction/:uuid`  | Update transaction with specified UUID            |
    |            `DELETE`            |  `/transaction/:uuid`  | Delete transaction with specified UUID            |
