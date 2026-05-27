// Package mcpadapter exposes the tool manifest through in-process calls.
package mcpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hairglasses/open-ralphglasses/internal/budget"
	"github.com/hairglasses/open-ralphglasses/internal/hookgate"
	"github.com/hairglasses/open-ralphglasses/internal/launchplan"
	"github.com/hairglasses/open-ralphglasses/internal/loopplan"
	"github.com/hairglasses/open-ralphglasses/internal/mcpmanifest"
	"github.com/hairglasses/open-ralphglasses/internal/provider"
)

// Result is the stable public call envelope.
type Result struct {
	OK      bool   `json:"ok"`
	Tool    string `json:"tool"`
	Payload any    `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ProviderStatus is the JSON shape returned by open_ralph_doctor.
type ProviderStatus struct {
	ID        string `json:"id"`
	Command   string `json:"command"`
	Installed bool   `json:"installed"`
}

// Call dispatches a public tool by name. Parameter names match the manifest.
func Call(ctx context.Context, toolName string, params map[string]any) Result {
	_ = ctx
	toolName = strings.TrimSpace(toolName)
	if params == nil {
		params = map[string]any{}
	}
	payload, err := callPayload(toolName, params)
	if err != nil {
		return Result{OK: false, Tool: toolName, Error: err.Error()}
	}
	return Result{OK: true, Tool: toolName, Payload: payload}
}

func callPayload(toolName string, params map[string]any) (any, error) {
	switch toolName {
	case "open_ralph_doctor":
		return providerStatuses(), nil
	case "open_ralph_provider_list":
		return provider.Catalog(), nil
	case "open_ralph_budget_estimate":
		inputTokens, err := intParam(params, "input_tokens")
		if err != nil {
			return nil, err
		}
		outputTokens, err := intParam(params, "output_tokens")
		if err != nil {
			return nil, err
		}
		estimate, err := budget.Estimate(budget.EstimateInput{
			Provider:     stringParam(params, "provider"),
			InputTokens:  inputTokens,
			OutputTokens: outputTokens,
		})
		if err != nil {
			return nil, err
		}
		payload := map[string]any{"estimate": estimate}
		if hasParam(params, "budget") || hasParam(params, "spent") {
			status, err := budget.Status(budget.StatusInput{
				SpentUSD:    floatParam(params, "spent"),
				BudgetUSD:   floatParam(params, "budget"),
				HeadroomPct: floatParam(params, "headroom_pct"),
			})
			if err != nil {
				return nil, err
			}
			payload["status"] = status
		}
		return payload, nil
	case "open_ralph_hook_check":
		return hookgate.Check(hookgate.CheckInput{
			Event:    stringParam(params, "event"),
			Tool:     stringParam(params, "tool"),
			Command:  firstNonEmpty(stringParam(params, "input"), stringParam(params, "command")),
			Path:     stringParam(params, "path"),
			RepoPath: stringParam(params, "repo"),
		})
	case "open_ralph_launch_plan":
		maxTurns, err := intParam(params, "max_turns")
		if err != nil {
			return nil, err
		}
		return launchplan.Build(launchplan.Options{
			Provider:       stringParam(params, "provider"),
			RepoPath:       stringParam(params, "repo"),
			Prompt:         stringParam(params, "prompt"),
			Model:          stringParam(params, "model"),
			PermissionMode: stringParam(params, "permission_mode"),
			MaxTurns:       maxTurns,
			BudgetUSD:      floatParam(params, "budget"),
			Agent:          stringParam(params, "agent"),
		})
	case "open_ralph_loop_plan":
		maxIterations, err := intParam(params, "max_iterations")
		if err != nil {
			return nil, err
		}
		return loopplan.Build(loopplan.Options{
			RepoPath:      stringParam(params, "repo"),
			Goal:          stringParam(params, "goal"),
			Provider:      stringParam(params, "provider"),
			VerifyCommand: stringParam(params, "verify"),
			MaxIterations: maxIterations,
		})
	case "open_ralph_tool_manifest":
		return mcpmanifest.Manifest(), nil
	default:
		return nil, fmt.Errorf("unknown public tool %q", toolName)
	}
}

func providerStatuses() []ProviderStatus {
	catalog := provider.Catalog()
	statuses := make([]ProviderStatus, 0, len(catalog))
	for _, p := range catalog {
		statuses = append(statuses, ProviderStatus{
			ID:        p.ID,
			Command:   p.Command,
			Installed: p.Installed(),
		})
	}
	return statuses
}

func stringParam(params map[string]any, key string) string {
	value, ok := params[key]
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func intParam(params map[string]any, key string) (int, error) {
	value, ok := params[key]
	if !ok || value == nil || value == "" {
		return 0, nil
	}
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		if v < 0 || v != float64(int(v)) {
			return 0, fmt.Errorf("%s must be a non-negative integer", key)
		}
		return int(v), nil
	case json.Number:
		i, err := strconv.Atoi(v.String())
		if err != nil || i < 0 {
			return 0, fmt.Errorf("%s must be a non-negative integer", key)
		}
		return i, nil
	default:
		raw := strings.TrimSpace(fmt.Sprint(v))
		i, err := strconv.Atoi(raw)
		if err != nil || i < 0 {
			return 0, fmt.Errorf("%s must be a non-negative integer", key)
		}
		return i, nil
	}
}

func floatParam(params map[string]any, key string) float64 {
	value, ok := params[key]
	if !ok || value == nil || value == "" {
		return 0
	}
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case json.Number:
		f, _ := strconv.ParseFloat(v.String(), 64)
		return f
	default:
		f, _ := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(v)), 64)
		return f
	}
}

func hasParam(params map[string]any, key string) bool {
	_, ok := params[key]
	return ok
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
