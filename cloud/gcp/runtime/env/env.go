package env

import (
	"fmt"
	"os"
)

const NITRIC_STACK = "NITRIC_STACK"

func GetNitricStack() string {
	return os.Getenv(NITRIC_STACK)
}

func GetNitricStackTag() string {
	return fmt.Sprintf("x-nitric-stack-%s", GetNitricStack())
}
