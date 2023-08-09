groshi - goddamn simple tool to keep track of your finances

Using groshi you can store, read and delete basic financial transactions.
Work is in progress, stand by!


HTTP API methods:
* POST /auth/login               log in and obtain an authentication token
* POST /auth/logout              log out and invalidate the authentication token
* POST /auth/refresh             refresh the authentication token

* POST /user/                    create a new user
* GET /user/                     read the current user
* PUT /user/                     update the current user
* DELETE /user/                  delete the current user

* POST /transaction/            create a new transaction
* GET /transaction/:uuid        retrieve a transaction with the specified UUID
* GET /transaction/             retrieve all transactions for a given period
* GET /transaction/summary      retrieve sum of all transactions for a given period in desired currency
* PUT /transaction/:uuid        update a transaction with the specified UUID
* DELETE /transaction/:uuid     delete a transaction with the specified UUID
