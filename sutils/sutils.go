package sutils

func HasKey(m map[string][]string, key string) bool {
	if _, ok := m[key]; !ok {
		return false
	}
	return true
}

func Includes(A []string, val string) bool {
	for _, v := range A {
		if string(v) == val {
			return true
		}
	}
	return false
}
