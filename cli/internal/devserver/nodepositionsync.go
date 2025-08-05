package devserver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nitrictech/nitric/cli/internal/config"
)

type NodePositionSync struct{}

type XYPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type NodePositionUpdate Message[map[string]XYPosition]

func NewNodePositionSync() *NodePositionSync {
	return &NodePositionSync{}
}

func (nps *NodePositionSync) OnConnect(send SendFunc) {
	// Load and send existing node positions on connect
	positions, err := loadNodePositions()
	if err != nil {
		// If we can't load positions, just continue without them
		return
	}

	send(Message[any]{
		Type:    "nitricNodeUpdate",
		Payload: positions,
	})
}

func (nps *NodePositionSync) OnMessage(message json.RawMessage) {
	var nodePositionUpdate NodePositionUpdate

	err := json.Unmarshal(message, &nodePositionUpdate)
	if err != nil {
		return
	}

	// Only handle nitricNodeUpdate messages
	if nodePositionUpdate.Type != "nitricNodeUpdate" {
		return
	}

	err = storeNodePositions(nodePositionUpdate.Payload)
	if err != nil {
		fmt.Println("Error storing node positions:", err)
	}
}

func loadNodePositions() (map[string]XYPosition, error) {
	nodePositionsPath := filepath.Join(config.LocalConfigPath(), "node-positions.json")

	data, err := os.ReadFile(nodePositionsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return empty map
			return make(map[string]XYPosition), nil
		}
		return nil, fmt.Errorf("failed to read node positions file: %w", err)
	}

	var positions map[string]XYPosition
	err = json.Unmarshal(data, &positions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node positions: %w", err)
	}

	return positions, nil
}

func storeNodePositions(changedPositions map[string]XYPosition) error {
	if err := os.MkdirAll(config.LocalConfigPath(), 0755); err != nil {
		return fmt.Errorf("failed to create nitric config directory: %w", err)
	}

	existingPositions, err := loadNodePositions()
	if err != nil {
		// If we can't load existing positions, start with empty map
		existingPositions = make(map[string]XYPosition)
	}

	for nodeId, position := range changedPositions {
		existingPositions[nodeId] = position
	}

	data, err := json.MarshalIndent(existingPositions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal node positions: %w", err)
	}

	nodePositionsPath := filepath.Join(config.LocalConfigPath(), "node-positions.json")
	err = os.WriteFile(nodePositionsPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write node positions file: %w", err)
	}

	return nil
}