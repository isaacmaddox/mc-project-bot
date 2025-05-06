package util

import "log"

func ErrorCheck(err error, str string) {
	if err != nil {
		log.Fatalf(str, err)
	}
}

func Extract[T any](value T, err error) T {
	if err != nil {
		log.Fatalf("Fatal in unwrap: %v", err)
	}

	return value
}

func ExtractOrNil[T any](value T, err error) *T {
	if err != nil {
		return nil
	}

	return &value
}
