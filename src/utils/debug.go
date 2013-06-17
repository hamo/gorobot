package utils

func BUG_ON(cond bool, msg string) {
	if cond {
		panic(msg)
	}
}

func BUG(msg string) {
	panic(msg)
}
