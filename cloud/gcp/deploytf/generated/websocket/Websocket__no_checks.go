// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build no_runtime_type_checking

package websocket

// Building without runtime type checking enabled, so all the below just return nil

func (w *jsiiProxy_Websocket) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (w *jsiiProxy_Websocket) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (w *jsiiProxy_Websocket) validateGetStringParameters(output *string) error {
	return nil
}

func (w *jsiiProxy_Websocket) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (w *jsiiProxy_Websocket) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateWebsocket_IsConstructParameters(x interface{}) error {
	return nil
}

func validateWebsocket_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Websocket) validateSetLambdaConnectTargetParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Websocket) validateSetLambdaDisconnectTargetParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Websocket) validateSetLambdaMessageTargetParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Websocket) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Websocket) validateSetWebsocketNameParameters(val *string) error {
	return nil
}

func validateNewWebsocketParameters(scope constructs.Construct, id *string, config *WebsocketConfig) error {
	return nil
}

