package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"pinmarker/configs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MessageResponseBuild(c *gin.Context, typeResponse, contextKey, method string, statusCode int, data, metadata interface{}) {
	wording := configs.ResponseMessages[method]

	var message string
	if typeResponse == "success" {
		message = fmt.Sprintf("%s %s", contextKey, wording)
	} else {
		message = fmt.Sprintf("Failed to %s %s", contextKey, wording)
	}

	response := gin.H{
		"message": TypographyCapitalize(message),
		"status":  typeResponse,
	}

	if typeResponse == "success" && data != nil {
		response["data"] = data
	}

	if typeResponse == "success" && metadata != nil {
		response["metadata"] = metadata
	}

	c.JSON(statusCode, response)
}

func MessageResponseErrorBuild(c *gin.Context, statusCode int, err string) {
	c.JSON(statusCode, gin.H{
		"message": err,
		"status":  "failed",
	})
}

// new
func BuildResponseMessage(c *gin.Context, typeResponse, contextKey string, method interface{}, statusCode int, data, metadata interface{}) {
	var response gin.H

	switch m := method.(type) {
	case string:
		wording, ok := configs.ResponseMessages[m]
		if !ok {
			wording = m
		}

		var message string
		if typeResponse == "success" {
			message = fmt.Sprintf("%s %s", contextKey, wording)
		} else {
			if methodStr, ok := method.(string); ok {
				if value, exists := configs.ResponseMessages[methodStr]; exists {
					message = fmt.Sprintf("%s %s", contextKey, value)
				} else {
					message = strings.ReplaceAll(methodStr, "_", " ")
				}
			} else {
				message = fmt.Sprintf("%v, %s", method, wording)
			}
		}

		response = gin.H{
			"message": Capitalize(message),
			"status":  typeResponse,
		}
	default:
		response = gin.H{
			"message": method,
			"status":  typeResponse,
		}
	}

	if typeResponse == "success" && data != nil {
		response["data"] = data
	}

	if typeResponse == "success" && metadata != nil {
		response["metadata"] = metadata
	}

	c.JSON(statusCode, response)
}

func BuildErrorMessage(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": err,
		"status":  "error",
	})
}

func BuildValidationError(err error) []map[string]string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var out []map[string]string
		for _, fe := range ve {
			out = append(out, map[string]string{
				"field": fe.Field(),
				"error": validationMessage(fe),
			})
		}
		return out
	}
	return nil
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters long"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters long"
	default:
		return fe.Field() + " is not valid"
	}
}
