# groshi
📉 **groshi** - goddamn simple tool to keep track of your finances.

> Work on groshi is still in progress, but it is close to release.

## Features
Using groshi you can perform all basic operations with transactions in different currencies: 
create, read and update them, besides you can also get useful summary of all transactions 
in desired currency units for given period.
Multiple number of users is also supported (each user owns its own transactions).

## groshi clients
### Client libraries for different programming languages
|                       **Library**                        | **Programming language** |
|:--------------------------------------------------------:|:------------------------:|
| [go-groshi](https://github.com/groshi-project/go-groshi) |            Go            |

### Client applications

|                 **Application**                  |  **Platform**   |
|:------------------------------------------------:|:---------------:|
| [grosh](https://github.com/groshi-project/grosh) | GNU/Linux (CLI) |



## HTTP API overview
There are two essences in **groshi**: _users_ and _transactions_.

> **Users** are basically like usual users in any other system. 
> They have _username_ and _password_, can authorize, own and manage their own _transactions_.

> **Transactions** are basically financial transactions! 
> They have _amount_ (it can be either positive or negative), _date_, _description_ and some other less important properties.

So, the logic is simple: _users_ create _transactions_ that are visible only to themselves. 
They can update, delete and fetch them, get useful summary of some _transactions_ for given period of time.

These tables will give you some basic overview of the groshi API methods.

- API methods related to **authorization** of _users_:

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
