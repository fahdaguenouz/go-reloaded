package functions


func ToUpper(s string) string {
	var res []rune
	for _, i := range s {
		if i >= 'a' && i <= 'z' {
			res = append(res, i-32)
		} else {
			res = append(res, i)
		}
	}
	return string(res)
}

func ToLower(s string) string {
	var res []rune
	for _, i := range s {
		if i >= 'A' && i <= 'Z' {
			res = append(res, i+32)
		} else {
			res = append(res, i)
		}
	}
	return string(res)
}

func IsPunctuation(ch byte) bool {
	return ch == '.' || ch == ',' || ch == '!' || ch == '?' || ch == ':' || ch == ';'
}