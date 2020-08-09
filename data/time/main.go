package time

type (
	Second int64
)

func Minute(minutes int64) Second {
	return Second(minutes * 60)
}
func Hour(hours int64) Second {
	return Minute(hours * 60)
}
func Day(days int64) Second {
	return Hour(days * 24)
}
