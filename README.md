# groshi
📉 **groshi** - goddamn simple tool to keep track of your finances.

> Work on groshi is still in progress, but it is close to release.

## Features
Using groshi you can perform all basic operations with transactions in different currencies: 
create, read and update them, besides you can also get useful summary of all transactions 
in desired currency units for given period.
Multiple number of users is also supported (each user owns its own transactions).

## HTTP API methods overview
These tables will give you some basic overview of the groshi API. 

- API methods related to **authorization**:

    |        **HTTP method**         |        **Path**        | **Description**                                   |
    |:------------------------------:|:----------------------:|---------------------------------------------------|
    |             `POST`             |     `/auth/login `     | Log in and obtain an authentication token         |
    |             `POST`             |     `/auth/logout`     | Log out and invalidate the authentication token   |
    |             `POST`             |    `/auth/refresh`     | Refresh the authentication token                  |


- API methods related to **users**:
    
    |        **HTTP method**         |      **Path**      | **Description**                                   |
    |:------------------------------:|:------------------:|---------------------------------------------------|
    |             `POST`             |      `/user`       | Create new user                                   |
    |             `GET`              |      `/user`       | Get information about current user                |
    |             `PUT`              |      `/user`       | Update current user                               |
    |            `DELETE`            |      `/user`       | Delete current user                               |

- API methods related to **transactions**:
    
    |        **HTTP method**         |        **Path**         | **Description**                                   |
    |:------------------------------:|:-----------------------:|---------------------------------------------------|
    |             `POST`             |     `/transactions`     | Create new transaction                            |
    |             `GET`              |     `/transactions`     | Retrieve all transactions for given period        |
    |             `GET`              |  `/transactions/:uuid`  | Retrieve a transaction with specified UUID        |
    |             `PUT`              |  `/transactions/:uuid`  | Update transaction with specified UUID            |
    |            `DELETE`            |  `/transactions/:uuid`  | Delete transaction with specified UUID            |
    |             `GET`              | `/transactions/summary` | Retrieve summary of transactions for given period |
