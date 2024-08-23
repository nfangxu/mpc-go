package utils

import (
	"bufio"
	"fmt"
	"os"
)

func Readfile(filename string) ([][]byte, error) {
	fmt.Println("Readfile:", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	data := make([][]byte, 0)
	for sc.Scan() {
		// 这里直接使用 sc.Bytes() 会导致结果错乱，所以用 sc.Text() 然后再转换为 []byte
		data = append(data, []byte(sc.Text()))
	}
	if sc.Err() != nil {
		return nil, sc.Err()
	}
	return data, nil
}
