package domain

import "testing"

func TestDurationDaysInclusive(t *testing.T) {
	start, err := ParseDate("2026-02-01")
	if err != nil {
		t.Fatal(err)
	}
	end, err := ParseDate("2026-02-10")
	if err != nil {
		t.Fatal(err)
	}
	got, err := DurationDays(start, end)
	if err != nil {
		t.Fatal(err)
	}
	if got != 10 {
		t.Fatalf("expected 10 days, got %d", got)
	}
}

func TestValidationErrors(t *testing.T) {
	if err := ValidateStatus(ProjectStatus("unknown")); err == nil {
		t.Fatal("expected status validation error")
	}
	if err := ValidateRole(AccessRole("admin")); err == nil {
		t.Fatal("expected role validation error")
	}
	if err := ValidateWorkload(150); err == nil {
		t.Fatal("expected workload validation error")
	}
	if _, err := CleanName("x"); err == nil {
		t.Fatal("expected name validation error")
	}
}
