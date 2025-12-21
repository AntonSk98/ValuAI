package core_test

import (
	"testing"
	core "valuai/core/state_engine"
)

func TestInitAnalysisFlowStateEngine(t *testing.T) {
	testFilePath := "test_analysis_flow.yml"
	engine := core.InitAnalysisFlowStateEngine(testFilePath)

	if engine == nil {
		t.Fatal("Expected non-nil engine")
	}

	if len(engine.AnalysisFlow) == 0 {
		t.Fatal("Expected analysis flow to be loaded")
	}
}

func TestGetCurrentState(t *testing.T) {
	testFilePath := "test_analysis_flow.yml"
	engine := core.InitAnalysisFlowStateEngine(testFilePath)

	// Test existing state
	state, err := engine.GetState("TEST_START")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if state == nil {
		t.Fatal("Expected non-nil state")
	}
	if state.State != "TEST_START" {
		t.Errorf("Expected state 'TEST_START', got '%s'", state.State)
	}

	// Test non-existing state
	_, err = engine.GetState("NON_EXISTENT")
	if err == nil {
		t.Fatal("Expected error for non-existent state")
	}
}

func TestGetNextState(t *testing.T) {
	testFilePath := "test_analysis_flow.yml"
	engine := core.InitAnalysisFlowStateEngine(testFilePath)

	// Test normal transition
	currentState, _ := engine.GetState("TEST_START")
	nextState, err := engine.GetNextState(currentState)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if nextState.State != "TEST_STEP1" {
		t.Errorf("Expected next state 'TEST_STEP1', got '%s'", nextState.State)
	}

	// Test transition to TEST_END
	nextState, err = engine.GetNextState(nextState)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if nextState.State != "TEST_END" {
		t.Errorf("Expected next state 'TEST_END', got '%s'", nextState.State)
	}

	// Test state with no transition (TEST_END state)
	finalState, err := engine.GetNextState(nextState)
	if err != nil {
		t.Fatalf("Expected no error for state with no transition, got %v", err)
	}
	if finalState.State != "TEST_END" {
		t.Errorf("Expected same state 'TEST_END', got '%s'", finalState.State)
	}

	// Test nil input
	_, err = engine.GetNextState(nil)
	if err == nil {
		t.Fatal("Expected error for nil input")
	}

	// Test optional state
	optionalState, _ := engine.GetState("TEST_OPTIONAL")
	nextFromOptional, err := engine.GetNextState(optionalState)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if nextFromOptional.State != "TEST_END" {
		t.Errorf("Expected next state 'TEST_END', got '%s'", nextFromOptional.State)
	}
}
