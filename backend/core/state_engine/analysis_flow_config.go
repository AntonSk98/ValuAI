package core

import (
	"fmt"
	"valuai/common"
)

// AnalysisFlowStateEngine represents the full conversation flow configuration
type AnalysisFlowStateEngine struct {
	AnalysisFlow []State `yaml:"analysis_flow"`
}

// State represents each step in the flow
type State struct {
	State      string            `yaml:"state"`
	Prompt     map[string]string `yaml:"prompt,omitempty"`
	Optional   bool              `yaml:"optional,omitempty"`
	Transition string            `yaml:"transition,omitempty"`
}

// InitAnalysisFlowStateEngine loads the analysis flow configuration from the specified YAML file.
func InitAnalysisFlowStateEngine(path string) *AnalysisFlowStateEngine {
	var analysisFlowConfig AnalysisFlowStateEngine
	common.LoadConfig(path, &analysisFlowConfig)
	return &analysisFlowConfig
}

// GetState returns the State by name by iterating through the slice
func (c *AnalysisFlowStateEngine) GetState(name string) (*State, error) {
	for _, s := range c.AnalysisFlow {
		if s.State == name {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("state %s not found", name)
}

// GetNextState returns the next State object
func (c *AnalysisFlowStateEngine) GetNextState(currentState *State) (*State, error) {
	if currentState == nil {
		return nil, fmt.Errorf("current state is nil")
	}

	if currentState.Transition == "" {
		return currentState, nil
	}

	currentStateName := currentState.State

	current, err := c.GetState(currentStateName)
	if err != nil {
		return nil, err
	}
	if current.Transition == "" {
		return nil, fmt.Errorf("state %s has no transition", currentStateName)
	}

	return c.GetState(current.Transition)
}
