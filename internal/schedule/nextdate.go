package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("empty repeat")
	}

	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", err
	}

	nowDay := StartOfDay(now)
	date := StartOfDay(start)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid d rule")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid d interval")
		}
		date = date.AddDate(0, 0, interval)
		for !AfterNow(date, nowDay) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(DateLayout), nil
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid y rule")
		}
		date = date.AddDate(1, 0, 0)
		for !AfterNow(date, nowDay) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DateLayout), nil
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid w rule")
		}
		week, err := parseWeekDays(parts[1])
		if err != nil {
			return "", err
		}
		if !AfterNow(date, nowDay) {
			date = nowDay.AddDate(0, 0, 1)
		}
		for i := 0; i < 800; i++ {
			if AfterNow(date, nowDay) && week[convertWeekDay(date)] {
				return date.Format(DateLayout), nil
			}
			date = date.AddDate(0, 0, 1)
		}
		return "", fmt.Errorf("next w date not found")
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid m rule")
		}
		dayTokens, err := parseMonthDay(parts[1])
		if err != nil {
			return "", err
		}
		var monthMask [13]bool
		if len(parts) == 3 {
			if err := parseMonthMask(parts[2], &monthMask); err != nil {
				return "", err
			}
		} else {
			for m := 1; m <= 12; m++ {
				monthMask[m] = true
			}
		}
		if !AfterNow(date, nowDay) {
			date = nowDay.AddDate(0, 0, 1)
		}
		for i := 0; i < 800; i++ {
			if !monthMask[int(date.Month())] {
				date = date.AddDate(0, 0, 1)
				continue
			}
			if !AfterNow(date, nowDay) {
				date = date.AddDate(0, 0, 1)
				continue
			}
			ok, err := monthDaysMatch(date, dayTokens)
			if err != nil {
				return "", err
			}
			if ok {
				return date.Format(DateLayout), nil
			}
			date = date.AddDate(0, 0, 1)
		}
		return "", fmt.Errorf("next m date not found")
	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func AfterNow(date, now time.Time) bool {
	return StartOfDay(date).After(StartOfDay(now))
}

func parseWeekDays(s string) ([8]bool, error) {
	var mask [8]bool
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			return mask, fmt.Errorf("invalid weekday")
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 7 {
			return mask, fmt.Errorf("invalid weekday")
		}
		mask[n] = true
	}
	return mask, nil
}

func convertWeekDay(t time.Time) int {
	w := t.Weekday()
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

func parseMonthDay(s string) ([]int, error) {
	var out []int
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("invalid month day")
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid month day")
		}
		if n < -2 || n == 0 || n > 31 {
			return nil, fmt.Errorf("invalid month day")
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("invalid month day")
	}
	return out, nil
}

func parseMonthMask(s string, mask *[13]bool) error {
	seen := false
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			return fmt.Errorf("invalid month")
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 12 {
			return fmt.Errorf("invalid month")
		}
		mask[n] = true
		seen = true
	}
	if !seen {
		return fmt.Errorf("invalid month")
	}
	return nil
}

func monthDaysMatch(date time.Time, tokens []int) (bool, error) {
	y, m, _ := date.Date()
	last := time.Date(y, m+1, 0, 0, 0, 0, 0, date.Location()).Day()
	pen := last - 1
	for _, tok := range tokens {
		switch {
		case tok >= 1 && tok <= 31:
			if date.Day() == tok {
				return true, nil
			}
		case tok == -1:
			if date.Day() == last {
				return true, nil
			}
		case tok == -2:
			if date.Day() == pen {
				return true, nil
			}
		default:
			return false, fmt.Errorf("invalid month day")
		}
	}
	return false, nil
}
