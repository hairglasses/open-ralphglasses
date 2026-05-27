package mcpadapter

import (
	"context"
	"strings"
	"testing"
)

func TestCallBudgetEstimate(t *testing.T) {
	result := Call(context.Background(), "open_ralph_budget_estimate", map[string]any{
		"provider":      "codex",
		"input_tokens":  1000,
		"output_tokens": 500,
	})
	if !result.OK {
		t.Fatalf("result=%+v", result)
	}
	payload, ok := result.Payload.(map[string]any)
	if !ok || payload["estimate"] == nil {
		t.Fatalf("payload=%#v", result.Payload)
	}
}

func TestCallUnknownTool(t *testing.T) {
	result := Call(context.Background(), "missing_tool", nil)
	if result.OK || !strings.Contains(result.Error, "unknown public tool") {
		t.Fatalf("result=%+v", result)
	}
}
