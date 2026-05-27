package budget

import "testing"

func TestEstimateUsesExampleRates(t *testing.T) {
	got, err := Estimate(EstimateInput{
		Provider:     "Gemini",
		InputTokens:  1_000_000,
		OutputTokens: 500_000,
	})
	if err != nil {
		t.Fatalf("Estimate: %v", err)
	}
	if got.Provider != "gemini" {
		t.Fatalf("provider = %q", got.Provider)
	}
	if got.EstimatedUSD != 0.45 {
		t.Fatalf("estimated = %.6f, want 0.45", got.EstimatedUSD)
	}
	if got.EfficiencyPct >= 100 {
		t.Fatalf("efficiency = %.2f, want cheaper than baseline", got.EfficiencyPct)
	}
	if got.RateSource != "example" {
		t.Fatalf("rate source = %q", got.RateSource)
	}
}

func TestEstimateRejectsUnknownProvider(t *testing.T) {
	_, err := Estimate(EstimateInput{Provider: "unknown"})
	if err == nil {
		t.Fatal("expected unknown provider error")
	}
}

func TestStatusStopsAtHeadroom(t *testing.T) {
	got, err := Status(StatusInput{
		SpentUSD:  9,
		BudgetUSD: 10,
	})
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !got.ShouldStop {
		t.Fatalf("expected stop verdict: %+v", got)
	}
	if got.AlertLevel != "90%" {
		t.Fatalf("alert = %q, want 90%%", got.AlertLevel)
	}
	if got.StopReason == "" {
		t.Fatal("expected stop reason")
	}
}

func TestStatusDisabledWhenBudgetZero(t *testing.T) {
	got, err := Status(StatusInput{SpentUSD: 100})
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !got.BudgetDisabled || got.ShouldStop {
		t.Fatalf("unexpected disabled budget status: %+v", got)
	}
}
