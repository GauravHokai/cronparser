package cronapp

import (
	"reflect"
	"strings"
	"testing"
)

func compareSchedules(t *testing.T, got, want *CronSchedule) {
	t.Helper()
	if !reflect.DeepEqual(got.Minutes, want.Minutes) {
		t.Errorf("Minutes mismatch: got %v, want %v", got.Minutes, want.Minutes)
	}
	if !reflect.DeepEqual(got.Hours, want.Hours) {
		t.Errorf("Hours mismatch: got %v, want %v", got.Hours, want.Hours)
	}
	if !reflect.DeepEqual(got.DaysOfMonth, want.DaysOfMonth) {
		t.Errorf("DaysOfMonth mismatch: got %v, want %v", got.DaysOfMonth, want.DaysOfMonth)
	}
	if !reflect.DeepEqual(got.Months, want.Months) {
		t.Errorf("Months mismatch: got %v, want %v", got.Months, want.Months)
	}
	if !reflect.DeepEqual(got.DaysOfWeek, want.DaysOfWeek) {
		t.Errorf("DaysOfWeek mismatch: got %v, want %v", got.DaysOfWeek, want.DaysOfWeek)
	}
	if got.Command != want.Command {
		t.Errorf("Command mismatch: got %q, want %q", got.Command, want.Command)
	}
}

func TestCronExpressionParser(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		want        *CronSchedule
		wantErr     bool
		errorString string
	}{
		{
			name:        "Invalid cron - too few fields",
			expression:  "* * *",
			wantErr:     true,
			errorString: "invalid cron expression: requires at least 6 fields",
		},
		{
			name:        "Invalid cron - invalid minute value",
			expression:  "60 * * * * /usr/bin/find",
			wantErr:     true,
			errorString: "error in field 1 ('60'): value 60 is out of range (0-59)",
		},
		{
			name:        "Invalid cron - invalid range",
			expression:  "0 5-1 * * * /usr/bin/find",
			wantErr:     true,
			errorString: "error in field 2 ('5-1'): range '5-1' is invalid; must be within 0-23",
		},
		{
			name:       "Basic cron - all wildcards",
			expression: "* * * * * /usr/bin/find",
			want: &CronSchedule{
				Minutes:     makeRange(0, 59),
				Hours:       makeRange(0, 23),
				DaysOfMonth: makeRange(1, 31),
				Months:      makeRange(1, 12),
				DaysOfWeek:  makeRange(0, 6),
				Command:     "/usr/bin/find",
			},
		},
		{
			name:       "Complex cron - combinations",
			expression: "*/15 0 1-7,15,21-23/2 * 1-5 /usr/bin/find",
			want: &CronSchedule{
				Minutes:     []int{0, 15, 30, 45},
				Hours:       []int{0},
				DaysOfMonth: []int{1, 2, 3, 4, 5, 6, 7, 15, 21, 23},
				Months:      makeRange(1, 12),
				DaysOfWeek:  []int{1, 2, 3, 4, 5},
				Command:     "/usr/bin/find",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.expression)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Parse() succeeded, want error")
				}
				if !strings.Contains(err.Error(), tt.errorString) {
					t.Errorf("Parse() error mismatch:\ngot = %v\nwant contains: %v", err, tt.errorString)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() returned unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Error("Parse() returned incorrect schedule:")
				compareSchedules(t, got, tt.want)
			}
		})
	}
}

func makeRange(min, max int) []int {
	result := make([]int, max-min+1)
	for i := range result {
		result[i] = min + i
	}
	return result
}
