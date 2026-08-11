package humanize

import (
	"fmt"
	"strconv"
)

func Bytes(size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " b"
	}
	fsize := float64(size)
	if size < 1048576 {
		return fmt.Sprintf("%0.2f kb", fsize/1024.0)
	}
	if size < 1073741824 {
		return fmt.Sprintf("%0.2f mb", fsize/1048576.0)
	}
	return fmt.Sprintf("%0.2f gb", fsize/1073741824.0)
}
