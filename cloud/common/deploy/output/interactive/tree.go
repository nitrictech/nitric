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
	Children []*Node[T]
}

func (n *Node[T]) AddChild(node *Node[T]) {
	if node.Children == nil {
		node.Children = make([]*Node[T], 0)
	}

	n.Children = append(n.Children, node)
}
