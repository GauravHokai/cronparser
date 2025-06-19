package cronapp

var fieldBoundaries = map[TimeField][2]int{
	MinuteField:     {0, 59},
	HourField:       {0, 23},
	DayOfMonthField: {1, 31},
	MonthField:      {1, 12},
	DayOfWeekField:  {0, 6},
}

type FieldDefinition struct{}

func NewFieldDefinition() *FieldDefinition {
	return &FieldDefinition{}
}

func (fd *FieldDefinition) GetBoundaries(timeField TimeField) (min, max int) {
	boundaries := fieldBoundaries[timeField]
	return boundaries[0], boundaries[1]
}

func (fd *FieldDefinition) IsValueInRange(value int, timeField TimeField) bool {
	min, max := fd.GetBoundaries(timeField)
	return value >= min && value <= max
}
