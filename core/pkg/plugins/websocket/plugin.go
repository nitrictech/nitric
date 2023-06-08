// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package websocket

import (
	"context"
	"fmt"
)

type WebsocketService interface {
	Send(ctx context.Context, socket string, connectionId string, message []byte) error
	Close(ctx context.Context, socket string, connectionId string) error
}

type UnimplementedWebsocketService struct{}

var _ WebsocketService = (*UnimplementedWebsocketService)(nil)

func (*UnimplementedWebsocketService) Send(ctx context.Context, socket string, connectionId string, message []byte) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedWebsocketService) Close(ctx context.Context, socket string, connectionId string) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
