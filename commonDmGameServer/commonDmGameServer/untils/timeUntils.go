package untils

import (
	"github.com/duke-git/lancet/v2/datetime"
	"time"
)

const (
	// hour 小时，0-23
	EverydayUpdateTime = 0 // 垮天的标准 表示凌晨0点
	Timezone           = 8 // 与UTC 的偏移量
	// SecondsPerMinute
	SecondsPerMinute = 60
	// SecondsPerHour
	SecondsPerHour = 60 * SecondsPerMinute
	// SecondsPerDay
	SecondsPerDay = 24 * SecondsPerHour
	// SecondsPerWeek
	SecondsPerWeek = 7 * SecondsPerDay
)

// GetWeekday 获取当前星期几
func GetWeekday() int32 {
	// 获取当前时间
	currentTime := time.Now()
	// 获取星期几
	weekday := currentTime.Weekday()
	switch weekday {
	case time.Sunday:
		return 7
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	}
	return 0
}

// 当前时间是否符合这个区间
func IsCurrentTimeInHourRange(startHour, endHour int) bool {
	// 获取当前时间
	now := time.Now()

	// 获取当前小时
	currentHour := now.Hour()

	// 判断当前小时是否在区间内
	return currentHour >= startHour && currentHour <= endHour
}

// 是否跨整点
func IsCrossHour(now, old int64) bool {
	// 转换时间戳为 time.Time 对象
	lastTime := time.Unix(old, 0)
	currentTime := time.Unix(now, 0)

	// 获取上次时间的小时部分
	lastHour := lastTime.Hour()

	// 获取当前时间的小时部分
	currentHour := currentTime.Hour()

	// 检查上次时间的小时与当前时间的小时是否相同
	return lastHour != currentHour
}

// GetFirstOfMonthTimestamp 获取当月1号凌晨的时间戳
func GetFirstOfMonthTimestamp() int64 {
	now := time.Now()
	// 构建当月1号凌晨的时间
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// 返回时间戳（秒）
	return firstOfMonth.Unix()
}

// 是否跨天 now, old 秒的时间戳
func IsNewDay(now, old int64) bool {
	return IsDiffHour(now, old, EverydayUpdateTime, Timezone)
}

// 是否跨周
func IsNewWeekDay(now, old int64) bool {
	return IsDiffWeekDayNum(now, old, EverydayUpdateTime, Timezone, 0)
}

// IsNewMonth 是否跨月
func IsNewMonth(now, old int64) bool {
	return DiffMonth(now, old) >= 1
}

// 下周五17点的日期时间戳 cnt: 0 本周五   -1 上周五 -2 上上周五
func GetFridayAfter(cnt int64) int64 {
	getNextFriday172 := getNextFriday17AfterCurrentTime()
	return getNextFriday172.Unix() + 7*24*60*60*cnt
}

func getNextFriday17() time.Time {
	now := time.Now()
	daysUntilFriday := (5 - int(now.Weekday()) + 7) % 7
	nextFriday := now.Add(time.Duration(daysUntilFriday) * 24 * time.Hour)
	return time.Date(nextFriday.Year(), nextFriday.Month(), nextFriday.Day(), 17, 0, 0, 0, nextFriday.Location())
}

func getNextFriday17AfterCurrentTime() time.Time {
	nextFriday17 := getNextFriday17()

	// 判断当前时间是否已经过了17点
	now := time.Now()
	if now.After(nextFriday17) {
		// 如果已经过了17点，则获取下一个周五的时间再加一周
		nextFriday17 = getNextFriday17().Add(7 * 24 * time.Hour)
	}

	return nextFriday17
}

// GetMidnightTimestamp 获取当前时间的凌晨时间戳
func GetMidnightTimestamp() int64 {
	// 获取当前时间
	now := time.Now()
	// 将时间调整到今天凌晨
	// time.Date(year, month, day, hour, minute, second, nanosecond, location)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 获取今天凌晨的时间戳(秒)
	return midnight.Unix()
}

// 获取当周的凌晨的时间
func GetWeekZeroTimestamp() int64 {
	// 获取当前时间
	now := time.Now()
	// 将时间调整到本周的凌晨
	// time.Date(year, month, day, hour, minute, second, nanosecond, location)
	weekZero := datetime.BeginOfWeek(now, time.Monday)
	// 获取本周凌晨的时间戳(秒)
	return weekZero.Unix()
}

// 获取当月的凌晨+1的时间
func GetMonthZeroTimestamp_1() int64 {
	// 获取当前时间
	now := time.Now()
	// 将时间调整到本月的凌晨
	// time.Date(year, month, day, hour, minute, second, nanosecond, location)
	monthZero := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// 获取本月凌晨的时间戳(秒)
	return monthZero.Unix() + 1
}

// GetMidnightTimestampMilli 获取当前时间的凌晨时间戳毫秒
func GetMidnightTimestampMilli() int64 {
	// 获取当前时间
	now := time.Now()
	// 将时间调整到今天凌晨
	// time.Date(year, month, day, hour, minute, second, nanosecond, location)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 获取今天凌晨的时间戳(秒)
	return midnight.UnixMilli()
}

// GetWeekZero 获取周零点时间戳
func GetWeekZero(timestamp int64, zone int, hour int) int64 {
	var cstZone = time.FixedZone("CST", zone*SecondsPerHour)
	t := datetime.BeginOfWeek(time.Unix(timestamp, 0).In(cstZone), time.Monday)
	return t.Unix() + int64(SecondsPerHour*hour)
}

/**
* @brief 时间戳转换为小时，24小时制，0点用24表示
*
* @param timestamp 时间戳
* @param timezone  时区
* @return uint32_t 小时 范围 1-24
 */
func GetHour24(timestamp int64, timezone int) int {
	hour := int((timestamp%86400)/3600) + timezone
	if hour > 24 {
		return hour - 24
	}
	return hour
}

/**
 * @brief 时间戳转换为小时，24小时制，0点用0表示
 *
 * @param timestamp 时间戳
 * @param timezone  时区
 * @return uint32_t 小时 范围 0-23
 */
func GetHour23(timestamp int64, timezone int) int {
	hour := GetHour24(timestamp, timezone)
	if hour == 24 {
		return 0 // 24点就是0点
	}
	return hour
}

/**
* @brief 判断两个时间戳是否是同一天
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param timezone 时区
* @return uint32_t 返回不同的天数
 */
func IsDiffDay(now, old int64, timezone int) int {
	now += int64(timezone * 3600)
	old += int64(timezone * 3600)
	return int((now / 86400) - (old / 86400))
}

/**
* @brief 判断时间戳是否处于一个小时的两边，即一个时间错大于当前的hour，一个小于
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param hour 小时，0-23
* @param timezone 时区
* @return bool true表示时间戳是否处于一个小时的两边
 */
func IsDiffHour(now, old int64, hour, timezone int) bool {
	diff := IsDiffDay(now, old, timezone)
	if diff == 1 {
		if GetHour23(old, timezone) >= hour {
			return GetHour23(now, timezone) >= hour
		} else {
			return true
		}
	} else if diff >= 2 {
		return true
	}

	return (GetHour23(now, timezone) >= hour) && (GetHour23(old, timezone) < hour)
}

/**
* @brief 判断时间戳是否处于跨周, 在每周几跨天节点的两边
*
* @param now 需要比较的时间戳
* @param old 需要比较的时间戳
* @param hour 小时，0-23
* @param timezone 时区
* @param daynum 星期几(0-6 周一-周天)
* @return bool true表示时间戳是否处于跨周, 在周几跨天节点的两边
 */
func IsDiffWeekDayNum(now, old int64, hour, timezone int, daynum int) bool {
	daynum = daynum % 7
	if now < old {
		now, old = old, now
	}
	// 将0 - 6 改为 1 - 7
	daynum++
	tmpnow := int64(timezone*3600) + now
	tmpold := int64(timezone*3600) + old
	// 使用UTC才能在本地时间采用周一作为一周的开始
	_, nw := time.Unix(tmpnow, 0).UTC().ISOWeek()
	_, ow := time.Unix(tmpold, 0).UTC().ISOWeek()
	nday := int(time.Unix(tmpnow, 0).UTC().Weekday())
	oday := int(time.Unix(tmpold, 0).UTC().Weekday())
	if nday == 0 { // 周天
		nday = 7
	}
	if oday == 0 { // 周天
		oday = 7
	}
	if nw-ow > 1 || ow-nw > 1 { // 跨1周以上
		return true
	} else if nw-ow == 1 { // 跨了一周， 需判断是否跨越周几的几点
		if nday > oday { // 当前:星期五 > 上周，星期四  时间点相差7天以上
			return true
		} else if nday == oday { // 时间戳刚好相差7天
			if int(nday) == daynum {
				return ((now+int64(timezone*3600))%86400 >= int64(hour*3600)) || ((old+int64(timezone*3600))%86400 < int64(hour*3600))
			} else {
				return true
			}
		} else {
			if int(oday) < daynum || int(nday) > daynum {
				return true
			} else if int(nday) == daynum {
				return (now+int64(timezone*3600))%86400 >= int64(hour*3600)
			} else if int(oday) == daynum {
				return (old+int64(timezone*3600))%86400 < int64(hour*3600)
			}
		}
	} else { // 未跨周,旧时间是否在 限定的时间点以前  且新时间在限定时间之后 时间点相差7天以内
		if int(nday) > daynum && int(oday) < daynum {
			return true
		} else if int(nday) == daynum && int(oday) < daynum { // 当前星期几大于设定的周几且 旧时间在设定周几以前
			return (now+int64(timezone*3600))%86400 >= int64(hour*3600) // 当前时间在设定小时之后
		} else if int(nday) > daynum && int(oday) == daynum { //
			return (old+int64(timezone*3600))%86400 < int64(hour*3600)
		} else if int(nday) == daynum && int(oday) == daynum {
			return ((now+int64(timezone*3600))%86400 >= int64(hour*3600)) && ((old+int64(timezone*3600))%86400 < int64(hour*3600))
		}
	}
	return false
}

// DiffMonth 两个时间戳相差月份
func DiffMonth(t1, t2 int64) int {
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	u1 := time.Unix(t1, 0)
	u2 := time.Unix(t2, 0)
	return (u2.Year()-u1.Year())*12 + int(u2.Month()) - int(u1.Month())
}

func DiffMonthWithZone(t1, t2 int64, hour, timezone int) int {
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	offset := timezone - hour
	var cstZone = time.FixedZone("CST", offset*3600) // 东offset
	u1 := time.Unix(t1, 0).In(cstZone)
	u2 := time.Unix(t2, 0).In(cstZone)
	return (u2.Year()-u1.Year())*12 + int(u2.Month()) - int(u1.Month())
}

// GetHourTime
func GetHourTime(timestamp int64, hour int, timezone int) int64 {
	nowHour := GetHour23(timestamp, timezone)

	zeroTimestamp := timestamp - timestamp%3600 - int64(nowHour)*3600

	zeroTimestamp += int64(hour) * 3600

	return zeroTimestamp
}
