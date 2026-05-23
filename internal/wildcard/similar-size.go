package wildcard

func IsSimilarSize(actualLenght int, referenceLenght int, tolerance int) bool {
	diff := actualLenght - referenceLenght
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance

}
