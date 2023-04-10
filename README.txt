groshi - goddamn simple tool to keep track of your finances

Using groshi you can store, read and delete basic financial transactions.
Work is in progress, stand by!

HTTP API methods:
* POST /auth/login               log in and obtain an authentication token.
* POST /auth/logout              log out and invalidate the authentication token.
* POST /auth/refresh             refresh the authentication token.

* POST /user/                    create a new user.
* GET /user/                     read the current user.
* PUT /user/                     update the current user.
* DELETE /user/                  delete the current user.

* POST /transactions/            create a new transaction.
* GET /transactions/:uuid        retrieve a transaction with the specified UUID.
* GET /transactions/             retrieve all transactions for a given period.
* PUT /transactions/:uuid        update a transaction with the specified UUID.
* DELETE /transactions/:uuid     delete a transaction with the specified UUID.

detailed API docs can be found here: https://github.com/jieggii/groshi/blob/master/docs/README.md