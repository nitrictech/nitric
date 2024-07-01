package common

import (
	"fmt"
)

func GetJobDefinitionName(stackId string, jobName string) (string, error) {
	return fmt.Sprintf("%s-job-%s", stackId, jobName), nil
}
