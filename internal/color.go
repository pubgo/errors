package internal

// terminal color define

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func Green(s string) string {
	return green + s + reset
}

func White(s string) string {
	return white + s + reset
}

func Yellow(s string) string {
	return yellow + s + reset
}

func Red(s string) string {
	return red + s + reset
}

func Blue(s string) string {
	return blue + s + reset
}

func Magenta(s string) string {
	return magenta + s + reset
}

func Cyan(s string) string {
	return cyan + s + reset
}
