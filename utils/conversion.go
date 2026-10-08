package utils

// SQLite doesn't have a Boolean type so we have to use integers.
// Go doesn't have a native way to convert a Boolean to an integer or vice versa so here we are.

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func IntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}
