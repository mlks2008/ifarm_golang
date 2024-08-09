package datetime

import (
	"time"
)

// IsJanuary 是否是一月
func (t *InternalTime) IsJanuary() bool {
	return t.Time.Month() == time.January
}

// IsFebruary 是否是二月
func (t *InternalTime) IsFebruary() bool {
	return t.Time.Month() == time.February
}

// IsMarch 是否是三月
func (t *InternalTime) IsMarch() bool {
	return t.Time.Month() == time.March
}

// IsApril 是否是四月
func (t *InternalTime) IsApril() bool {
	return t.Time.Month() == time.April
}

// IsMay 是否是五月
func (t *InternalTime) IsMay() bool {
	return t.Time.Month() == time.May
}

// IsJune 是否是六月
func (t *InternalTime) IsJune() bool {
	return t.Time.Month() == time.June
}

// IsJuly 是否是七月
func (t *InternalTime) IsJuly() bool {
	return t.Time.Month() == time.July
}

// IsAugust 是否是八月
func (t *InternalTime) IsAugust() bool {
	return t.Time.Month() == time.August
}

// IsSeptember 是否是九月
func (t *InternalTime) IsSeptember() bool {
	return t.Time.Month() == time.September
}

// IsOctober 是否是十月
func (t *InternalTime) IsOctober() bool {
	return t.Time.Month() == time.October
}

// IsNovember 是否是十一月
func (t *InternalTime) IsNovember() bool {
	return t.Time.Month() == time.November
}

// IsDecember 是否是十二月
func (t *InternalTime) IsDecember() bool {
	return t.Time.Month() == time.December
}

// IsMonday 是否是周一
func (t *InternalTime) IsMonday() bool {
	return t.Time.Weekday() == time.Monday
}

// IsTuesday 是否是周二
func (t *InternalTime) IsTuesday() bool {
	return t.Time.Weekday() == time.Tuesday
}

// IsWednesday 是否是周三
func (t *InternalTime) IsWednesday() bool {
	return t.Time.Weekday() == time.Wednesday
}

// IsThursday 是否是周四
func (t *InternalTime) IsThursday() bool {
	return t.Time.Weekday() == time.Thursday
}

// IsFriday 是否是周五
func (t *InternalTime) IsFriday() bool {
	return t.Time.Weekday() == time.Friday
}

// IsSaturday 是否是周六
func (t *InternalTime) IsSaturday() bool {
	return t.Time.Weekday() == time.Saturday
}

// IsSunday 是否是周日
func (t *InternalTime) IsSunday() bool {
	return t.Time.Weekday() == time.Sunday
}

// IsWeekday 是否是工作日
func (t *InternalTime) IsWeekday() bool {
	return !t.IsSaturday() && !t.IsSunday()
}

// IsWeekend 是否是周末
func (t *InternalTime) IsWeekend() bool {
	return t.IsSaturday() || t.IsSunday()
}

// InSameDay 同一天
func InSameDay(t1, t2 int64) bool {
	y1, m1, d1 := NewSecond(t1).Date()
	y2, m2, d2 := NewSecond(t2).Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

// InSameWeek 同一周
func InSameWeek(t1, t2 int64) bool {
	y1, w1 := NewSecond(t1).ISOWeek()
	y2, w2 := NewSecond(t2).ISOWeek()
	return y1 == y2 && w1 == w2
}

// InSameMonth 同一个月
func InSameMonth(t1, t2 int64) bool {
	y1, m1, _ := NewSecond(t1).Date()
	y2, m2, _ := NewSecond(t2).Date()
	return y1 == y2 && m1 == m2
}
