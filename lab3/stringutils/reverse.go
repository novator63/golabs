package stringutils

func Reverse(s string) string {
	runes := []rune(s)
	n := len(runes)
	reversed := make([]rune, n)

	for i := 0; i < n; i++ {
		reversed[n-1-i] = runes[i]
	}

	return string(reversed)
}
