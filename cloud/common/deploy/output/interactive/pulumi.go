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
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

type PulumiData struct {
	Urn string
	// Name   string
	Type        string
	Status      ResourceStatus
	StartTime   time.Time
	EndTime     time.Time
	LastMessage string
}

func (pd PulumiData) Name() string {
	urnParts := strings.Split(pd.Urn, "::")

	return urnParts[len(urnParts)-1]
}

type ResourceStatus int

const (
	ResourceStatus_Creating = iota
	ResourceStatus_Updating
	ResourceStatus_Deleting
	ResourceStatus_Created
	ResourceStatus_Deleted
	ResourceStatus_Updated
	ResourceStatus_Failed_Create
	ResourceStatus_Failed_Delete
	ResourceStatus_Failed_Update
	ResourceStatus_Unchanged
	ResourceStatus_None
)

var PreResourceStates = map[string]ResourceStatus{
	"create": ResourceStatus_Creating,
	"delete": ResourceStatus_Deleting,
	"same":   ResourceStatus_Unchanged,
	"update": ResourceStatus_Updating,
}

var SuccessResourceStates = map[string]ResourceStatus{
	"create": ResourceStatus_Created,
	"delete": ResourceStatus_Deleted,
	"same":   ResourceStatus_Unchanged,
	"update": ResourceStatus_Updated,
}

var FailedResourceStates = map[string]ResourceStatus{
	"create": ResourceStatus_Failed_Create,
	"delete": ResourceStatus_Failed_Delete,
	"update": ResourceStatus_Failed_Update,
}

var MessageResourceStates = map[ResourceStatus]string{
	ResourceStatus_Creating:      "creating",
	ResourceStatus_Updating:      "updating",
	ResourceStatus_Deleting:      "deleting",
	ResourceStatus_Created:       "created",
	ResourceStatus_Deleted:       "deleted",
	ResourceStatus_Updated:       "updated",
	ResourceStatus_Failed_Create: "create failed",
	ResourceStatus_Failed_Delete: "delete failed",
	ResourceStatus_Failed_Update: "updated failed",
	ResourceStatus_Unchanged:     "unchanged",
	ResourceStatus_None:          "",
}

// TODO: Use TUI standard colors when lib available
var StatusStyles = map[ResourceStatus]lipgloss.Style{
	// Unchanged State
	ResourceStatus_Unchanged: lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC")),
	// Pre states
	ResourceStatus_Creating: lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")),
	ResourceStatus_Updating: lipgloss.NewStyle().Foreground(lipgloss.Color("#f1c40f")),
	ResourceStatus_Deleting: lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c")),
	// Post states
	ResourceStatus_Created: lipgloss.NewStyle().Foreground(lipgloss.Color("#27ae60")),
	ResourceStatus_Updated: lipgloss.NewStyle().Foreground(lipgloss.Color("#f39c12")),
	ResourceStatus_Deleted: lipgloss.NewStyle().Foreground(lipgloss.Color("#c0392b")),

	// Failed states
	ResourceStatus_Failed_Create: ErrorStyle,
	ResourceStatus_Failed_Delete: ErrorStyle,
	ResourceStatus_Failed_Update: ErrorStyle,
}
