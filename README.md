# groshi
📉 **groshi** - goddamn simple tool to keep track of your finances.

> Work on this project is still in progress, but it is nearing release. Stay tuned!

## Features
With **groshi**, you can effortlessly manage financial transactions in various currencies. 
It allows you to create, read, and update transactions,
as well as generate a useful summary of all transactions in your preferred currency for a specified period. 
Multiple users are supported, each with their own set of transactions.

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
To run **groshi**, you have two options: locally or inside docker containers.
Both of them are described in these instructions.

First you will have set up some secrets and environmental variables in order to run the service.

### Step 1: secrets
Use the following command to create a `secrets` directory and its associated subdirectories and files to store secrets:

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
* Generate random string and populate `app/jwt_secret_key` with it.
* Sign up for an account at [exchangeratesapi.io](https://exchangeratesapi.io), obtain an API key, and then enter it into `app/exchangerates_api_key` as a secret.

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

Edit these variables as needed. Again, defaults are fine, but you would probably like
to edit a couple of them (e.g. `GROSHI_PORT`).

Please note that **groshi** does not automatically read the .env file. 
You'll need to export these defined variables to the environment, 
for example, using source .env if you're using bash.

### Step 3: finally, run it
To run **groshi** in Docker, build and start the containers with:
```shell
docker-compose up --build --attach groshi
```

For local execution, run the `main.go` file with:
```shell
go run main.go
```

## HTTP API overview
groshi revolves around two main entities: _users_ and _transactions_.

> **Users** represent typical system users with a _username_ and _password_. 
> They can authorize, own, and manage their own _transactions_.

> **Transactions** are financial transactions with properties such as amount (positive or negative),
> date, description, and more.

The logic is straightforward: users create transactions that are private to them. 
They can update, delete, fetch them, and obtain a useful summary of their transactions over a specified period.

The following tables provide an overview of the groshi API methods:

- API methods related to **authorization** of _users_:

    |        **HTTP method**         |        **Path**        | **Description**                                   |
    |:------------------------------:|:----------------------:|---------------------------------------------------|
    |             `POST`             |     `/auth/login `     | Log in and obtain an authentication token         |
    |             `POST`             |    `/auth/refresh`     | Refresh the authentication token                  |


- API methods related to **users**:
    
    | **HTTP method** | **Path** | **Description**                             |
    |:---------------:|:--------:|---------------------------------------------|
    |     `POST`      | `/user`  | Create a new user                           |
    |      `GET`      | `/user`  | Retrieve information about the current user |
    |      `PUT`      | `/user`  | Update the current user                     |
    |    `DELETE`     | `/user`  | Delete the current user                     |

- API methods related to **transactions**:
    
    | **HTTP method** |        **Path**         | **Description**                                       |
    |:---------------:|:-----------------------:|-------------------------------------------------------|
    |     `POST`      |     `/transactions`     | Create a new transaction                              |
    |      `GET`      |     `/transactions`     | Retrieve all transactions for a specified period      |
    |      `GET`      |  `/transactions/:uuid`  | Retrieve a transaction with a specified UUID          |
    |      `PUT`      |  `/transactions/:uuid`  | Update a transaction with a specified UUID            |
    |    `DELETE`     |  `/transactions/:uuid`  | Delete a transaction with a specified UUID            |
    |      `GET`      | `/transactions/summary` | Retrieve a summary of transactions for a given period |

> Access the Swagger API documentation at the `/docs/index.html` route
> by setting the `GROSHI_SWAGGER` environment variable to `true`.