package cronapp

import (
	"fmt"
	"strings"
)

type CronSchedule struct {
	Minutes     []int
	Hours       []int
	DaysOfMonth []int
	Months      []int
	DaysOfWeek  []int
	Command     string
}

func (cs *CronSchedule) String() string {
	var builder strings.Builder

	formatLine := func(name string, values []int) string {
		paddedName := fmt.Sprintf("%-14s", name)
		valueStrs := make([]string, len(values))
		for i, v := range values {
			valueStrs[i] = fmt.Sprintf("%d", v)
		}
		return paddedName + strings.Join(valueStrs, " ") + "\n"
	}

	builder.WriteString(formatLine("minute", cs.Minutes))
	builder.WriteString(formatLine("hour", cs.Hours))
	builder.WriteString(formatLine("day of month", cs.DaysOfMonth))
	builder.WriteString(formatLine("month", cs.Months))
	builder.WriteString(formatLine("day of week", cs.DaysOfWeek))
	builder.WriteString(fmt.Sprintf("%-14s%s", "command", cs.Command))

	return builder.String()
}

// TimeField represents a field in a cron expression (minute, hour, etc.).
type TimeField int

const (
	MinuteField TimeField = iota
	HourField
	DayOfMonthField
	MonthField
	DayOfWeekField
)
