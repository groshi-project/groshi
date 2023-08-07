package error_messages

import "errors"

var ErrorInvalidRequestParams = errors.New("invalid request params")

// var UserNotFound = errors.New("user was not found")
var TransactionNotFound = errors.New("transaction was not found")
var TransactionDoesNotBelongToYou = errors.New("the transaction doesn't belong to you")
