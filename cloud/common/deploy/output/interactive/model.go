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
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/samber/lo"
)

type DeployModel struct {
	pulumiSub chan events.EngineEvent
	sub       chan tea.Msg
	logs      []string
	tree      *Tree[PulumiData]
}

func (m *DeployModel) Init() []Cmd {
	return []Cmd{
		subscribeToChan(m.sub),
		subscribeToChan(m.pulumiSub),
	}
}

// subscribeToChannel - A tea Command that will wait on messages sent to the given channel
func subscribeToChan[T any](sub chan T) Cmd {
	return func() Msg {
		return <-sub
	}
}

const MAX_LOG_LENGTH = 5

// Implement io.Writer for simplicity
func (m *DeployModel) Write(bytes []byte) (int, error) {
	msg := string(bytes)
	cutMsg := strings.TrimSuffix(msg, "\n")

	// This will hook the writer into the tea program lifecycle
	m.sub <- LogMessage{
		Message: cutMsg,
	}

	return len(bytes), nil
}

func (m *DeployModel) handlePulumiEngineEvent(evt events.EngineEvent) {
	// These events are directly tied to a resource
	if evt.DiagnosticEvent != nil {
		// TODO: Handle diagnostic event logging
		node := m.tree.FindNode(evt.DiagnosticEvent.URN)
		if node != nil {
			node.Data.LastMessage = evt.DiagnosticEvent.Message
		}
	} else if evt.ResourcePreEvent != nil {
		// attempt to locate the parent node
		meta := evt.ResourcePreEvent.Metadata.New
		if meta == nil {
			meta = evt.ResourcePreEvent.Metadata.Old
		}

		parentNode := m.tree.FindNode(meta.Parent)
		if parentNode == nil {
			parentNode = m.tree.Root
		}

		parentNode.AddChild(&Node[PulumiData]{
			Id:       evt.ResourcePreEvent.Metadata.URN,
			Sequence: evt.Sequence,
			Data: &PulumiData{
				StartTime: time.Now(),
				Urn:       evt.ResourcePreEvent.Metadata.URN,
				Type:      evt.ResourcePreEvent.Metadata.Type,
				Status:    PreResourceStates[string(evt.ResourcePreEvent.Metadata.Op)],
			},
			Children: []*Node[PulumiData]{},
		})
	} else if evt.ResOutputsEvent != nil {
		// Find the URN in the tree
		node := m.tree.FindNode(evt.ResOutputsEvent.Metadata.URN)
		if node != nil {
			node.Data.EndTime = time.Now()
			node.Data.Status = SuccessResourceStates[string(evt.ResOutputsEvent.Metadata.Op)]
		}
	} else if evt.ResOpFailedEvent != nil {
		node := m.tree.FindNode(evt.ResOpFailedEvent.Metadata.URN)
		if node != nil {
			node.Data.EndTime = time.Now()
			node.Data.Status = FailedResourceStates[string(evt.ResOpFailedEvent.Metadata.Op)]
		}
	}
}

func (m *DeployModel) Update(msg Msg) (*DeployModel, Cmd) {
	switch t := msg.(type) {
	case events.EngineEvent:
		m.handlePulumiEngineEvent(t)
		return m, subscribeToChan(m.pulumiSub)
	case LogMessage:
		m.logs = append(m.logs, t.Message)
		return m, subscribeToChan(m.sub)
	default:
		return m, nil
	}
}

func (m *DeployModel) renderNodeRow(node *Node[PulumiData], depth int, isLast bool, parentLast bool) table.Row {
	linkChar := lo.Ternary(!isLast, "├─", "└─")
	prefixString := lo.Ternary(!parentLast, fmt.Sprintf("│  %s", linkChar), linkChar)
	marginLeft := lo.Ternary(!parentLast, 3*(depth-1), 3*depth)

	statusStyle := StatusStyles[node.Data.Status]
	isPending := lo.Contains(lo.Values(PreResourceStates), node.Data.Status)
	isComplete := lo.Contains(lo.Values(SuccessResourceStates), node.Data.Status)
	isFailed := lo.Contains(lo.Values(FailedResourceStates), node.Data.Status)

	status := statusStyle.Render(MessageResourceStates[node.Data.Status])
	if isPending {
		runningTime := time.Since(node.Data.StartTime).Round(time.Second)
		status = statusStyle.Render(MessageResourceStates[node.Data.Status] + fmt.Sprintf(" (%s)", runningTime))
	} else if isComplete || isFailed {
		completeTime := node.Data.EndTime.Sub(node.Data.StartTime).Round(time.Second)
		status = statusStyle.Render(MessageResourceStates[node.Data.Status] + fmt.Sprintf(" (%s)", completeTime))
	}

	return table.Row{
		// Name
		lipgloss.NewStyle().MarginLeft(marginLeft).SetString(prefixString).Render(node.Data.Name()),
		// Type
		node.Data.Type,
		// Status
		status,
		// Message
		// node.Data.LastMessage,
	}
}

// Render the tree rows
func (m *DeployModel) renderNodeRows(depth int, parentLast bool, nodes ...*Node[PulumiData]) []table.Row {
	// render this nods info
	rows := []table.Row{}

	for idx, node := range nodes {
		isLast := idx == len(nodes)-1
		rows = append(rows, m.renderNodeRow(node, depth, isLast, parentLast))

		sortNodes(node.Children)
		if len(node.Children) > 0 {
			rows = append(rows, m.renderNodeRows(depth+1, isLast, node.Children...)...)
		}
	}

	return rows
}

func sortNodes(nodes nodeList) {
	sort.Sort(nodes)
}

func (m *DeployModel) renderFailedNodeMessages(nodes ...*Node[PulumiData]) []string {
	rows := []string{}

	for _, node := range nodes {
		isFailed := lo.Contains(lo.Values(FailedResourceStates), node.Data.Status)

		if isFailed {
			rows = append(rows, fmt.Sprintf("Resource: %s failed to deploy:\n%s\n\n", node.Data.Name(), node.Data.LastMessage))
		}

		if len(node.Children) > 0 {
			rows = append(rows, m.renderFailedNodeMessages(node.Children...)...)
		}
	}

	return rows
}

func (m DeployModel) View() string {
	rows := m.renderNodeRows(0, true, m.tree.Root)

	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Type", Width: 30},
		{Title: "Status", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(s)

	failureMessage := m.renderFailedNodeMessages(m.tree.Root)

	return fmt.Sprintf("\n%s\n%s", t.View(), ErrorStyle.Render(failureMessage...))
}

type OutputModelArgs struct {
	Sub       chan tea.Msg
	PulumiSub chan events.EngineEvent
}

func NewOutputModel(sub chan tea.Msg, pulumiSub chan events.EngineEvent) (*DeployModel, error) {
	err := os.Setenv("CLICOLOR_FORCE", "1")
	if err != nil {
		return nil, err
	}

	return &DeployModel{
		pulumiSub: pulumiSub,
		sub:       sub,
		logs:      make([]string, 0),
		tree: &Tree[PulumiData]{
			Root: &Node[PulumiData]{
				Id: "root",
				Data: &PulumiData{
					Urn:    "project",
					Type:   "",
					Status: ResourceStatus_None,
				},
				Children: []*Node[PulumiData]{},
			},
		},
	}, nil
}
