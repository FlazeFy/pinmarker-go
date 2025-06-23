package utils

import (
	"fmt"

	"pinmarker/configs"

	"github.com/gin-gonic/gin"
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
