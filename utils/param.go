package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseQueryInt32(c *gin.Context, key string) (int32, error) {
	valStr := c.Query(key)
	if valStr == "" {
		return 0, fmt.Errorf("parameter %s is missing", key)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}
