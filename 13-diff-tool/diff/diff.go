package diff

import "strings"

var ans []string

func LCS(a, b []string) []string {
	lcs(a, b, len(a)-1, len(b)-1, [][]byte{})
	return ans
}

func convert(data [][]byte) []string {
	var value []string
	for _, v := range data {
		value = append(value, string(v))
	}
	return value
}

func lcs(a, b []string, i, j int, value [][]byte) {
	if i == -1 || j == -1 {
		if len(value) >= len(ans) {
			ans = convert(value)
		}
		return
	}
	if a[i] == b[j] {
		value = append(value, []byte(a[i]))
		lcs(a, b, i-1, j-1, value)
	} else {
		lcs(a, b, i-1, j, value)
		lcs(a, b, i, j-1, value)
	}
}

func Difference(fileData, common []string) []string {
	var data []string
	for _, c := range fileData {
		skip := false
		for _, s := range common {
			if c == s {
				skip = true
			}
		}
		if !skip {
			data = append(data, c)
		}
	}
	return data
}

func ExecuteDiff(data1, data2 []string) string {
	sb := strings.Builder{}
	common := LCS(data1, data2)

	diff1 := Difference(data1, common)
	diff2 := Difference(data2, common)

	i, j := 0, 0
	for i < len(diff1) && j < len(diff2) {
		sb.WriteString("> " + diff1[i] + "\n")
		sb.WriteString("< " + diff2[j] + "\n")
		i++
		j++
	}
	for i < len(diff1) {
		sb.WriteString("> " + diff1[i] + "\n")
		i++
	}
	for j < len(diff1) {
		sb.WriteString("<" + diff2[j] + "\n")
		j++
	}
	return sb.String()
}
