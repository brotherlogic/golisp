package main

func READ(s string) string {
	return s
}

func EVAL(s string) string {
	return s
}

func PRINT(s string) string {
	return s
}

func rep(s string) string {
	return PRINT(EVAL(READ(s)))
}
