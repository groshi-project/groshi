# groshi
📉 **groshi** - goddamn simple tool to keep track of your finances.

> Work on the project is still in progress, but it is close to release. Stay tuned!

## Features
Using **groshi** you can perform all basic operations with transactions in different currencies: 
create, read and update them, besides you can also get useful summary of all transactions 
in desired currency units for given period.
Multiple number of users is also supported (each user owns its own transactions).

## Clients
### Client libraries for different programming languages
|                       **Library**                        | **Programming language** |
|:--------------------------------------------------------:|:------------------------:|
| [go-groshi](https://github.com/groshi-project/go-groshi) |            Go            |

### Client applications

|                 **Application**                  |          **Platform**           |
|:------------------------------------------------:|:-------------------------------:|
| [grosh](https://github.com/groshi-project/grosh) | GNU/Linux, Windows, MacOS (CLI) |

## Running instructions
Basically you have two ways to run **groshi**: locally and inside docker containers.
Both of them are described in these instructions.

First you will have set up some secrets and environmental variables in order to run the service.

### Step 1: secrets
Run the following command to create `secrets` directory
and other directories and files inside it which will hold secrets:

```shell
make secrets
```

`./secrets` directory tree:
```
secrets/
├── app
│   ├── exchangerates_api_key
│   └── jwt_secret_key
└── mongo
    ├── database
    ├── password
    └── username
```

Then fill all the secrets:
* Mind up yourself `mongo/username`, `mongo/password` and `mongo/database`
  if you are going to run **groshi** inside docker container, otherwise set actual credentials of your MongoDB instance.
  Also don't forget to bring it up!
* Generate random string and fill `app/jwt_secret_key` with it.
* Register at [exchangeratesapi.io](https://exchangeratesapi.io), obtain an API key and fill `app/exchangerates_api_key` secret with it.

### Step 2: environmental variables
#### If you are going to run groshi using docker
Simply open `docker-compose.yaml` using your favourite editor and
edit environmental variables (it is optional, defaults are fine) in the `environment` section.
Also remember to update `ports` section if you change `GROSHI_PORT` variable.

#### If you are going to run groshi locally
Create `.env` file containing all necessary environmental variables using `.env.example` as template:
```shell
cp .env.example .env
```

Edit these variables as you wish. Again, defaults are fine, but you would probably like
to edit a couple of them (e.g. `GROSHI_PORT`).

Please note, that **groshi** does not take into account the `.env` file.
You will have to somehow "export" these defined variables to the environment.
For example, you can use `source .env` command if you are using bash.

### Step 3: finally run it
Just build and bring up docker containers using docker-compose
if you wish to run **groshi** in docker:
```shell
docker-compose up --build --attach groshi
```

Or run `main.go` file if you are going to run groshi locally:
```shell
go run main.go
```

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
