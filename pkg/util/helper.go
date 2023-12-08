package util

func ResetFlag(flag *bool, hour int) {
	if hour != 9 && hour != 18 && *flag {
		*flag = false
	}
}
