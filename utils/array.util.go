package utils

type ArrayUtil struct{}

func (a *ArrayUtil) IsValid(verifiableValue interface{}, validValues ...interface{}) (isValid bool) {
	isValid = false

	for _, v := range validValues {
		if v == verifiableValue {
			isValid = true
			break
		}
	}

	return isValid
}
