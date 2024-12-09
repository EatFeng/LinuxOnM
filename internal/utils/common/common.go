package common

import (
	"LinuxOnM/internal/utils/cmd"
	"fmt"
	mathRand "math/rand"
	"strings"
	"time"
)

const (
	b  = uint64(1)
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
)

func FormatBytes(bytes uint64) string {
	switch {
	case bytes < kb:
		return fmt.Sprintf("%dB", bytes)
	case bytes < mb:
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(kb))
	case bytes < gb:
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(mb))
	default:
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(gb))
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(b)
}

func LoadTimeZoneByCmd() string {
	loc := time.Now().Location().String()
	if _, err := time.LoadLocation(loc); err != nil {
		loc = "Asia/Shanghai"
	}
	std, err := cmd.Exec("timedatectl | grep 'Time zone'")
	if err != nil {
		return loc
	}
	fields := strings.Fields(string(std))
	if len(fields) != 5 {
		return loc
	}
	if _, err := time.LoadLocation(fields[2]); err != nil {
		return loc
	}
	return fields[2]
}
