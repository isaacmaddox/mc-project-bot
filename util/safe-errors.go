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


func Extract_or_nil[T any](value T, err error) *T {
	if err != nil {
		return nil
	}

	return &value
}

func If_nil_then[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}

	return *ptr
}