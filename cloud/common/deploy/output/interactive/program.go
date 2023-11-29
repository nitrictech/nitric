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

package interactive

import (
	"context"
	"io"
	"time"
)

type Msg interface{}

type (
	QuitMsg   struct{}
	EventTick struct{ Time time.Time }
)

type Cmd func() Msg

type Program struct {
	model  *DeployModel
	ticker *time.Ticker
	writer io.Writer
	cmds   chan Cmd
	msgs   chan Msg
	errs   chan error

	ctx    context.Context
	cancel context.CancelFunc

	lastRender string
}

type ProgramArgs struct {
	Fps    int
	Writer io.Writer
}

func NewProgram(model *DeployModel, args *ProgramArgs) *Program {
	if args.Fps < 1 {
		args.Fps = 1
	}

	fps := 1 / args.Fps
	ticker := time.NewTicker(time.Duration(fps) * time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	return &Program{
		model:  model,
		ticker: ticker,
		writer: args.Writer,
		cmds:   make(chan Cmd),
		msgs:   make(chan Msg),
		errs:   make(chan error),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Program) Run() {
	p.handleCommands(p.cmds)
	p.handleErrors()

	cmds := p.model.Init()
	for _, cmd := range cmds {
		p.cmds <- cmd
	}

	// Send a tick at a regular cadence defined by the model's FPS
	go func() {
		for t := range p.ticker.C {
			p.Send(EventTick{Time: t})
		}
	}()

	p.eventLoop()
}

func (p *Program) eventLoop() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case msg := <-p.msgs:
			var cmd Cmd

			p.model, cmd = p.model.Update(msg) // run update
			p.cmds <- cmd                      // process command (if any)

			view := p.model.View()
			p.render(view)
		}
	}
}

func (p *Program) Stop() {
	p.ticker.Stop()
	p.cancel()
	close(p.cmds)
	close(p.msgs)
	close(p.errs)
}

// Send sends a message to the main update function, effectively allowing
// messages to be injected from outside the program for interoperability
// purposes.
//
// If the program hasn't started yet this will be a blocking operation.
// If the program has already been terminated this will be a no-op, so it's safe
// to send messages after the program has exited.
func (p *Program) Send(msg Msg) {
	select {
	case <-p.ctx.Done():
	case p.msgs <- msg:
	}
}

func (p *Program) render(view string) {
	if view != p.lastRender {
		_, err := p.writer.Write([]byte(view))
		p.errs <- err

		p.lastRender = view
	}
}

func (p *Program) handleErrors() {
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case err := <-p.errs:
				if err != nil {
					p.Stop()
				}
			}
		}
	}()
}

// handleCommands runs commands in a goroutine and sends the result to the
// program's message channel.
func (p *Program) handleCommands(cmds chan Cmd) {
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case cmd := <-cmds:
				if cmd == nil {
					continue
				}

				go func() {
					msg := cmd()
					p.Send(msg)
				}()
			}
		}
	}()
}
