package datax

func Head(data []byte, size int) ([]byte, []byte) {
	if len(data) <= size {
		size = len(data)
	}
	return data[:size], data[size:]
}
