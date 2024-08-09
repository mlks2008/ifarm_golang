package datetime

import (
	"math"
	"time"
)

// AddYears N年后
func (t *InternalTime) AddYears(years int) {
	t.Time = t.Time.AddDate(years, 0, 0)
}

// AddYearsNoOverflow N年后(月份不溢出)
func (t *InternalTime) AddYearsNoOverflow(years int) {
	// 获取N年后本月的最后一天
	last := time.Date(t.Year()+years, t.Time.Month(), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()).AddDate(0, 1, -1)

	day := t.Day()
	if t.Day() > last.Day() {
		day = last.Day()
	}

	t.Time = time.Date(last.Year(), last.Month(), day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// AddYear 1年后
func (t *InternalTime) AddYear() {
	t.AddYears(1)
}

// AddYearNoOverflow 1年后(月份不溢出)
func (t *InternalTime) AddYearNoOverflow() {
	t.AddYearsNoOverflow(1)
}

// SubYears N年前
func (t *InternalTime) SubYears(years int) {
	t.AddYears(-years)
}

// SubYearsNoOverflow N年前(月份不溢出)
func (t *InternalTime) SubYearsNoOverflow(years int) {
	t.AddYearsNoOverflow(-years)
}

// SubYear 1年前
func (t *InternalTime) SubYear() {
	t.SubYears(1)
}

// SubYearNoOverflow 1年前(月份不溢出)
func (t *InternalTime) SubYearNoOverflow() {
	t.SubYearsNoOverflow(1)
}

// AddMonths N月后
func (t *InternalTime) AddMonths(months int) {
	t.Time = t.Time.AddDate(0, months, 0)
}

// AddMonthsNoOverflow N月后(月份不溢出)
func (t *InternalTime) AddMonthsNoOverflow(months int) {
	month := t.Time.Month() + time.Month(months)

	// 获取N月后的最后一天
	last := time.Date(t.Year(), month, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()).AddDate(0, 1, -1)

	day := t.Day()
	if t.Day() > last.Day() {
		day = last.Day()
	}

	t.Time = time.Date(last.Year(), last.Month(), day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// AddMonth 1月后
func (t *InternalTime) AddMonth() {
	t.AddMonths(1)
}

// AddMonthNoOverflow 1月后(月份不溢出)
func (t *InternalTime) AddMonthNoOverflow() {
	t.AddMonthsNoOverflow(1)
}

// SubMonths N月前
func (t *InternalTime) SubMonths(months int) {
	t.AddMonths(-months)
}

// SubMonthsNoOverflow N月前(月份不溢出)
func (t *InternalTime) SubMonthsNoOverflow(months int) {
	t.AddMonthsNoOverflow(-months)
}

// SubMonth 1月前
func (t *InternalTime) SubMonth() {
	t.SubMonths(1)
}

// SubMonthNoOverflow 1月前(月份不溢出)
func (t *InternalTime) SubMonthNoOverflow() {
	t.SubMonthsNoOverflow(1)
}

// AddWeeks N周后
func (t *InternalTime) AddWeeks(weeks int) {
	t.AddDays(weeks * DaysPerWeek)
}

// AddWeek 1天后
func (t *InternalTime) AddWeek() {
	t.AddWeeks(1)
}

// SubWeeks N周后
func (t *InternalTime) SubWeeks(weeks int) {
	t.SubDays(weeks * DaysPerWeek)
}

// SubWeek 1天后
func (t *InternalTime) SubWeek() {
	t.SubWeeks(1)
}

// AddDays N天后
func (t *InternalTime) AddDays(days int) {
	t.Time = t.Time.AddDate(0, 0, days)
}

// AddDay 1天后
func (t *InternalTime) AddDay() {
	t.AddDays(1)
}

// SubDays N天前
func (t *InternalTime) SubDays(days int) {
	t.AddDays(-days)
}

// SubDay 1天前
func (t *InternalTime) SubDay() {
	t.SubDays(1)
}

// AddHours N小时后
func (t *InternalTime) AddHours(hours int) {
	td := time.Duration(hours) * time.Hour
	t.Time = t.Time.Add(td)
}

// AddHour 1小时后
func (t *InternalTime) AddHour() {
	t.AddHours(1)
}

// SubHours N小时前
func (t *InternalTime) SubHours(hours int) {
	t.AddHours(-hours)
}

// SubHour 1小时前
func (t *InternalTime) SubHour() {
	t.SubHours(1)
}

// AddMinutes N分钟后
func (t *InternalTime) AddMinutes(minutes int) {
	td := time.Duration(minutes) * time.Minute
	t.Time = t.Time.Add(td)
}

// AddMinute 1分钟后
func (t *InternalTime) AddMinute() {
	t.AddMinutes(1)
}

// SubMinutes N分钟前
func (t *InternalTime) SubMinutes(minutes int) {
	t.AddMinutes(-minutes)
}

// SubMinute 1分钟前
func (t *InternalTime) SubMinute() {
	t.SubMinutes(1)
}

// AddSeconds N秒钟后
func (t *InternalTime) AddSeconds(seconds int) {
	td := time.Duration(seconds) * time.Second
	t.Time = t.Time.Add(td)
}

// AddSecond 1秒钟后
func (t *InternalTime) AddSecond() {
	t.AddSeconds(1)
}

// SubSeconds N秒钟前
func (t *InternalTime) SubSeconds(seconds int) {
	t.AddSeconds(-seconds)
}

// SubSecond 1秒钟前
func (t *InternalTime) SubSecond() {
	t.SubSeconds(1)
}

// 查询指定年份指定月份有多少天
func GetYearMonthDays(dateTime InternalTime) int {
	year, month, _ := dateTime.Date()

	//有31天的月份
	day31 := map[int]bool{
		1:  true,
		3:  true,
		5:  true,
		7:  true,
		8:  true,
		10: true,
		12: true,
	}
	if day31[int(month)] {
		return 31
	}
	// 有30天的月份
	day30 := map[int]bool{
		4:  true,
		6:  true,
		9:  true,
		11: true,
	}
	if day30[int(month)] {
		return 30
	}
	//计算平年还是闰年
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		// 得出二月天数
		return 29
	}
	// 得出平年二月天数
	return 28
}

// 计算天数差
func SubDays(t1, t2 InternalTime) (day int) {
	if t1.Unix() > t2.Unix() {
		t_ := t1
		t1 = t2
		t2 = t_
	}
	day = int(t2.Time.Sub(t1.Time).Seconds() / SecondsPerDay)

	return
}

// 计算数月差
func SubMonth(t1, t2 InternalTime) (month int) {
	if t1.Unix() < t2.Unix() {
		t_ := t1
		t1 = t2
		t2 = t_
	}
	y1 := t1.Year()
	y2 := t2.Year()
	m1 := int(t1.Month())
	m2 := int(t2.Month())
	d1 := t1.Day()
	d2 := t2.Day()

	yearInterval := y1 - y2
	// 如果 d1的 月-日 小于 d2的 月-日 那么 yearInterval-- 这样就得到了相差的年数
	if m1 < m2 || m1 == m2 && d1 < d2 {
		yearInterval--
	}
	// 获取月数差值
	monthInterval := (m1 + 12) - m2
	if d1 < d2 {
		monthInterval--
	}
	monthInterval %= 12
	month = yearInterval*12 + monthInterval
	return
}

func GetDiffDays(t1, t2 InternalTime) int64 {
	startDate := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	endDate := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return int64(math.Abs(startDate.Sub(endDate).Seconds())) / 86400
}

func GetDiffTimes(t1, t2 int64) int64 {
	return int64(math.Abs(NewSecond(t1).Sub(NewSecond(t2).Time).Seconds()))
}
