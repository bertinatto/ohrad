package ohrad

func unix2ntp(tsp int64) uint32 {
	return uint32(Jan1970 + tsp)
}

func ntp2unix(tsp int64) uint32 {
	return uint32(tsp - Jan1970)
}
