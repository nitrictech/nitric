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

package pulumix

import (
	"strings"
	"time"
)

type Tree[T any] struct {
	Root *Node[T]
}

// Find node from root
func (t *Tree[T]) FindNode(id string) *Node[T] {
	return t.findNode(t.Root, id)
}

// Recursive implementation
func (t *Tree[T]) findNode(node *Node[T], id string) *Node[T] {
	if node.Id == id {
		// Return if we have a match
		return node
	}

	// otherwise walk the children recursively
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			if foundNode := t.findNode(child, id); foundNode != nil {
				return foundNode
			}
		}
	}

	// No matches
	return nil
}

type Node[T any] struct {
	Id       string
	Data     *T
	Parent   *Node[T]
	Children []*Node[T]
}

// FindParent - finds the first parent node that matches the given function
func (n *Node[T]) FindParent(fn func(n *Node[T]) bool) *Node[T] {
	if n.Parent == nil {
		return nil
	}

	if fn(n.Parent) {
		return n.Parent
	}

	return n.Parent.FindParent(fn)
}

func (n *Node[T]) AddChild(node *Node[T]) {
	if node.Children == nil {
		node.Children = make([]*Node[T], 0)
	}

	node.Parent = n

	n.Children = append(n.Children, node)
}

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
