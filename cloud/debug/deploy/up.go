package deploy

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

// Up - Deploy requested infrastructure for a stack
func (d *DeployServer) Up(request *deploy.DeployUpRequest, stream deploy.DeployService_UpServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("recovered panic: %+v\n Stack: %s", r, stack)
		}
	}()

	reqJson, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		return err
	}
	
	fmt.Println(string(reqJson))

	return nil
}