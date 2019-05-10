package util

// MapUnionStrStrA turns m1 into the union of m1 and m2
func MapUnionStrStrA(m1, m2 map[string][]string) {
	for key, values := range m2 {
		m1[key] = values
	}
}
