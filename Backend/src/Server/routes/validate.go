package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Wrapper for data validator functions so that they can be mocked during unit testing
type Validator interface {
	DataInvalid(*gin.Context, func() (bool, string, error)) bool
}

type ValidatorWrapper struct{}

func (v *ValidatorWrapper) DataInvalid(ctx *gin.Context, validateEdit func() (bool, string, error)) bool {
	valid, problem, err := validateEdit()
	if checkErr(ctx, err) {
		return true
	} else if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": problem})
		return true
	} else {
		return false
	}
}

var validator Validator = &ValidatorWrapper{}
