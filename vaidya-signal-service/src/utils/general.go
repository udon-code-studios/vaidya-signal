package utils

func PanicOnNotNil(value interface{}) {
	if value != nil {
		panic(value)
	}
}