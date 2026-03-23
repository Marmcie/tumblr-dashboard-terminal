package modules



func Abs(val int) int {
	if val >= 0 {
		return val
	} else {
		return val * -1
	}
}

func Fit(val int, limit int) int {
	if val >= 0 {
		return val % limit
	} else {
		return limit + val
	}
}

func LineNumber(str string, width int) int {
	var ct = 2
	var x = 0
	for _, v := range str {
		if v == '\n' {
			ct++
			x = 0
		}
		if x >= width {
			ct++
			x = 0
		}
		x++
	}
	return ct
}

func LineAfter(str string, width int, line int) string {
	var ct = 0
	var x = 0

	for i, v := range str {
		if v == '\n' {
			ct++
			x = 0
		}
		if x >= width {
			ct++
			x = 0
		}
		x++
		if ct >= line {
			return str[i:]
		}
	}
	return str

}

