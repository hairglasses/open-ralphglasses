// Package budget contains cost and budget helpers for provider planning.
package budget

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Rate stores token pricing in USD per one million tokens.
type Rate struct {
	InputPer1M  float64 `json:"input_per_1m_usd"`
	OutputPer1M float64 `json:"output_per_1m_usd"`
}

// ExampleRates returns public example rates for the built-in providers.
//
// They are intentionally small, local defaults for demos and tests, not an
// authority on current provider pricing. Callers can supply an explicit rate to
// Estimate when they need exact billing behavior.
func ExampleRates() map[string]Rate {
	return map[string]Rate{
		"codex":  {InputPer1M: 2.50, OutputPer1M: 15.00},
		"claude": {InputPer1M: 3.00, OutputPer1M: 15.00},
		"gemini": {InputPer1M: 0.15, OutputPer1M: 0.60},
	}
}

// EstimateInput describes one public cost estimate.
type EstimateInput struct {
	Provider     string
	InputTokens  int
	OutputTokens int
	Rate         *Rate
	Baseline     *Rate
}

// EstimateResult is the JSON-facing cost estimate.
type EstimateResult struct {
	Provider        string  `json:"provider"`
	InputTokens     int     `json:"input_tokens"`
	OutputTokens    int     `json:"output_tokens"`
	EstimatedUSD    float64 `json:"estimated_usd"`
	BaselineUSD     float64 `json:"baseline_usd,omitempty"`
	EfficiencyPct   float64 `json:"efficiency_pct,omitempty"`
	RateInputPer1M  float64 `json:"rate_input_per_1m_usd"`
	RateOutputPer1M float64 `json:"rate_output_per_1m_usd"`
	RateSource      string  `json:"rate_source"`
}

// Estimate calculates cost from token counts. If Rate is nil, the package uses
// the provider's example rate. If Baseline is nil, Claude's example rate is used
// as the cross-provider comparison baseline.
func Estimate(in EstimateInput) (EstimateResult, error) {
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	if provider == "" {
		return EstimateResult{}, fmt.Errorf("provider is required")
	}
	if in.InputTokens < 0 || in.OutputTokens < 0 {
		return EstimateResult{}, fmt.Errorf("token counts must be non-negative")
	}

	rateSource := "explicit"
	rate := Rate{}
	if in.Rate != nil {
		rate = *in.Rate
	} else {
		examples := ExampleRates()
		var ok bool
		rate, ok = examples[provider]
		if !ok {
			return EstimateResult{}, fmt.Errorf("unknown provider %q (want one of: %s)", provider, strings.Join(ProviderIDs(), ", "))
		}
		rateSource = "example"
	}
	if rate.InputPer1M < 0 || rate.OutputPer1M < 0 {
		return EstimateResult{}, fmt.Errorf("rates must be non-negative")
	}

	baseline := ExampleRates()["claude"]
	if in.Baseline != nil {
		baseline = *in.Baseline
	}
	estimated := tokenCost(in.InputTokens, in.OutputTokens, rate)
	baselineUSD := tokenCost(in.InputTokens, in.OutputTokens, baseline)
	result := EstimateResult{
		Provider:        provider,
		InputTokens:     in.InputTokens,
		OutputTokens:    in.OutputTokens,
		EstimatedUSD:    roundUSD(estimated),
		BaselineUSD:     roundUSD(baselineUSD),
		RateInputPer1M:  rate.InputPer1M,
		RateOutputPer1M: rate.OutputPer1M,
		RateSource:      rateSource,
	}
	if baselineUSD > 0 {
		result.EfficiencyPct = roundPct((estimated / baselineUSD) * 100)
	}
	return result, nil
}

// ProviderIDs returns the example-rate provider ids in stable order.
func ProviderIDs() []string {
	rates := ExampleRates()
	ids := make([]string, 0, len(rates))
	for id := range rates {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// StatusInput describes one budget check.
type StatusInput struct {
	SpentUSD    float64
	BudgetUSD   float64
	HeadroomPct float64
}

// StatusResult is a compact budget verdict.
type StatusResult struct {
	SpentUSD       float64 `json:"spent_usd"`
	BudgetUSD      float64 `json:"budget_usd"`
	RemainingUSD   float64 `json:"remaining_usd"`
	PercentUsed    float64 `json:"percent_used"`
	HeadroomPct    float64 `json:"headroom_pct"`
	AlertLevel     string  `json:"alert_level,omitempty"`
	ShouldStop     bool    `json:"should_stop"`
	StopReason     string  `json:"stop_reason,omitempty"`
	BudgetDisabled bool    `json:"budget_disabled"`
}

// Status evaluates spend against a soft headroom threshold. The default
// threshold is 90 percent of the budget.
func Status(in StatusInput) (StatusResult, error) {
	if in.SpentUSD < 0 || in.BudgetUSD < 0 {
		return StatusResult{}, fmt.Errorf("spent and budget values must be non-negative")
	}
	headroom := in.HeadroomPct
	if headroom <= 0 {
		headroom = 90
	}
	if headroom > 100 {
		return StatusResult{}, fmt.Errorf("headroom_pct must be <= 100")
	}
	result := StatusResult{
		SpentUSD:    roundUSD(in.SpentUSD),
		BudgetUSD:   roundUSD(in.BudgetUSD),
		HeadroomPct: roundPct(headroom),
	}
	if in.BudgetUSD == 0 {
		result.BudgetDisabled = true
		return result, nil
	}
	result.RemainingUSD = roundUSD(math.Max(in.BudgetUSD-in.SpentUSD, 0))
	result.PercentUsed = roundPct((in.SpentUSD / in.BudgetUSD) * 100)
	result.AlertLevel = alertLevel(result.PercentUsed)
	if result.PercentUsed >= headroom {
		result.ShouldStop = true
		result.StopReason = fmt.Sprintf("spent $%.2f of $%.2f budget at %.0f%% headroom", result.SpentUSD, result.BudgetUSD, headroom)
	}
	return result, nil
}

func tokenCost(inputTokens, outputTokens int, rate Rate) float64 {
	return (float64(inputTokens)/1_000_000)*rate.InputPer1M +
		(float64(outputTokens)/1_000_000)*rate.OutputPer1M
}

func alertLevel(percentUsed float64) string {
	switch {
	case percentUsed >= 90:
		return "90%"
	case percentUsed >= 75:
		return "75%"
	case percentUsed >= 50:
		return "50%"
	default:
		return ""
	}
}

func roundUSD(v float64) float64 {
	return math.Round(v*1_000_000) / 1_000_000
}

func roundPct(v float64) float64 {
	return math.Round(v*100) / 100
}
