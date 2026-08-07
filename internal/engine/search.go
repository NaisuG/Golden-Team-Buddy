package engine

func TargetSize(level int) int {
	switch {
	case level == 0:
		return 0
	case level <= 5:
		return 6
	case level <= 7:
		return 8
	case level == 8:
		return 9
	default:
		return 10
	}
}
