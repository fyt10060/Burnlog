// errorhandler
package model

import (
	"encoding/json"
	//	"fmt"
)

type ErrorType string

const (
	ErrSuccess ErrorType = "success"
	// data base error
	ErrSysDb = "System database error"
	// users
	ErrSignInName  = "user name or email wrong"
	ErrSignInPass  = "password wrong"
	ErrPassIllegal = "password out of range"
	ErrNameExist   = "username exist"
	ErrEmailExist  = "email exist"
	ErrNameEmpty   = "name or email empty"
	ErrEmailMiss   = "can not found this email"
	ErrTokenMiss   = "permission error, token expire"
	ErrUidMiss     = "can not find this user"
	ErrAuLimit     = "authority not allowed"
)

var (
	errCode = map[ErrorType]int{
		ErrSuccess:     0,
		ErrSysDb:       -1,
		ErrSignInName:  10001,
		ErrSignInPass:  10002,
		ErrNameExist:   10003,
		ErrEmailExist:  10004,
		ErrNameEmpty:   10005,
		ErrEmailMiss:   10006,
		ErrPassIllegal: 10007,
		ErrUidMiss:     10008,
		ErrAuLimit:     10009,
		ErrTokenMiss:   10010,
	}
)

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  ErrorType   `json:"err_msg"`
	Data    interface{} `json:"data"`
}

func ParseResult(errMsg ErrorType, data interface{}) []byte {
	code := errCode[errMsg]
	if code == 0 && errMsg != ErrSuccess {
		code = -1
	}
	r := Response{
		ErrCode: code,
		ErrMsg:  errMsg,
		Data:    data,
	}
	b, err := json.Marshal(r)
	if err != nil {

	}

	return b
}
