package utils

import (
	"bytes"
	"io"
)

func ParseLimitedPayload(r io.Reader, n int64) ([]byte, error) {
	var buf bytes.Buffer
	limiter := io.LimitReader(r, 1*1024*1024)
	_, err := io.Copy(&buf, limiter)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// func MakeGenericArray(typed interface{}) ([]interface{}, error) {
// 	if array, ok := typed.([]interface{}); ok {
// 		generic := make([]interface{}, len(array))
// 		for i, item := range array {
// 			generic[i] = item
// 		}
// 		return generic, nil
// 	}
// 	return nil, errors.New("not an array")
// }
