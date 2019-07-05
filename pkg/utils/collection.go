package utils

// Intersection finds the intersecting elements between 2 slices
func Intersection(sliceA, sliceB []string) []string {
	encounteredMap := map[string]int8{}
	for _, element := range sliceA {
		encounteredMap[element] = 1
	}
	for _, element := range sliceB {
		if val, ok := encounteredMap[element]; ok {
			encounteredMap[element] = val + 1
		} else {
			encounteredMap[element] = 1
		}
	}

	var intersectionElements []string
	for key, val := range encounteredMap {
		if val > 1 {
			intersectionElements = append(intersectionElements, key)
		}
	}
	return intersectionElements
}

// MapKeys gets the keys from a map
func MapKeys(myMap map[string]string) []string {
	var keys []string
	for k := range myMap {
		keys = append(keys, k)
	}
	return keys
}
