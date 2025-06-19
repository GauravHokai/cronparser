package cronapp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// signature for a function that can parse a cron field.
type fieldParserFunc func(field string, timeField TimeField, definition *FieldDefinition) ([]int, bool, error)

type Parser struct {
	parsers    []fieldParserFunc
	definition *FieldDefinition
}

func NewParser() *Parser {
	return &Parser{
		definition: NewFieldDefinition(),
		parsers: []fieldParserFunc{
			parseList,
			parseInterval,
			parseRange,
			parseWildcard,
			parseSingleValue,
		},
	}
}

// Parse orchestrates the parsing of the entire cron expression.
func Parse(expression string) (*CronSchedule, error) {
	fields := strings.Fields(expression)
	if len(fields) < 6 {
		return nil, errors.New("invalid cron expression: requires at least 6 fields")
	}

	parser := NewParser()
	schedule := &CronSchedule{
		Command: strings.Join(fields[5:], " "),
	}

	timeFields := []TimeField{MinuteField, HourField, DayOfMonthField, MonthField, DayOfWeekField}
	scheduleSlices := []*[]int{&schedule.Minutes, &schedule.Hours, &schedule.DaysOfMonth, &schedule.Months, &schedule.DaysOfWeek}

	for i, timeField := range timeFields {
		values, err := parser.parseField(fields[i], timeField)
		if err != nil {
			return nil, fmt.Errorf("error in field %d ('%s'): %w", i+1, fields[i], err)
		}
		*scheduleSlices[i] = values
	}

	return schedule, nil
}

// parseField iterates through the registered parser functions and uses the first one that can handle the field.
func (p *Parser) parseField(field string, timeField TimeField) ([]int, error) {
	for _, parserFunc := range p.parsers {
		values, handled, err := parserFunc(field, timeField, p.definition)
		if err != nil {
			return nil, err
		}
		if handled {
			return values, nil
		}
	}
	return nil, fmt.Errorf("unrecognized format for field: %s", field)
}

func parseWildcard(field string, timeField TimeField, definition *FieldDefinition) ([]int, bool, error) {
	if field != "*" {
		return nil, false, nil
	}
	min, max := definition.GetBoundaries(timeField)
	result := make([]int, 0, max-min+1)
	for i := min; i <= max; i++ {
		result = append(result, i)
	}
	return result, true, nil
}

func parseList(field string, timeField TimeField, definition *FieldDefinition) ([]int, bool, error) {
	if !strings.Contains(field, ",") {
		return nil, false, nil
	}

	uniqueValues := make(map[int]bool)
	subFields := strings.Split(field, ",")
	tempParser := NewParser()

	for _, subField := range subFields {
		values, err := tempParser.parseField(subField, timeField)
		if err != nil {
			return nil, true, err
		}
		for _, v := range values {
			uniqueValues[v] = true
		}
	}

	result := make([]int, 0, len(uniqueValues))
	for v := range uniqueValues {
		result = append(result, v)
	}
	sort.Ints(result)
	return result, true, nil
}

func parseInterval(field string, timeField TimeField, definition *FieldDefinition) ([]int, bool, error) {
	if !strings.Contains(field, "/") {
		return nil, false, nil
	}
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return nil, true, fmt.Errorf("invalid interval format: %s", field)
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval <= 0 {
		return nil, true, fmt.Errorf("invalid interval value: %s", parts[1])
	}

	tempParser := NewParser()
	rangePart := parts[0]
	rangeValues, err := tempParser.parseField(rangePart, timeField)
	if err != nil {
		return nil, true, fmt.Errorf("invalid range for interval: '%s'", rangePart)
	}

	start := rangeValues[0]
	_, max := definition.GetBoundaries(timeField)
	end := max
	if strings.Contains(rangePart, "-") {
		end = rangeValues[len(rangeValues)-1]
	}

	result := make([]int, 0)
	for i := start; i <= end; i += interval {
		result = append(result, i)
	}
	return result, true, nil
}

func parseRange(field string, timeField TimeField, definition *FieldDefinition) ([]int, bool, error) {
	if !strings.Contains(field, "-") {
		return nil, false, nil
	}
	bounds := strings.Split(field, "-")
	if len(bounds) != 2 {
		return nil, true, fmt.Errorf("invalid range format: %s", field)
	}
	start, err := strconv.Atoi(bounds[0])
	if err != nil {
		return nil, true, fmt.Errorf("invalid range start: %s", bounds[0])
	}
	end, err := strconv.Atoi(bounds[1])
	if err != nil {
		return nil, true, fmt.Errorf("invalid range end: %s", bounds[1])
	}
	min, max := definition.GetBoundaries(timeField)
	if !definition.IsValueInRange(start, timeField) || !definition.IsValueInRange(end, timeField) || start > end {
		return nil, true, fmt.Errorf("range '%s' is invalid; must be within %d-%d", field, min, max)
	}

	result := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result, true, nil
}

func parseSingleValue(field string, timeField TimeField, definition *FieldDefinition) ([]int, bool, error) {
	val, err := strconv.Atoi(field)
	if err != nil {
		return nil, false, nil
	}
	if !definition.IsValueInRange(val, timeField) {
		min, max := definition.GetBoundaries(timeField)
		return nil, true, fmt.Errorf("value %d is out of range (%d-%d)", val, min, max)
	}
	return []int{val}, true, nil
}
