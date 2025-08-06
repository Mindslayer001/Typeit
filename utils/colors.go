package utils

type Colors struct {
	Red   string
	Green string
	Reset string
}

var BasicColors = Colors{
	Red:   "\033[31m",
	Green: "\033[32m",
	Reset: "\033[0m",
}
