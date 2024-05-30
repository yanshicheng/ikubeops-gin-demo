package utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func CurrentPath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Dir(file)
}

func FileToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func FiltToTrimString(filePath string) (string, error) {
	str, err := FileToString(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(str), err
}

func ReadLine(r *bufio.Reader) ([]byte, error) {
	line, isPrefix, err := r.ReadLine()
	for isPrefix && err == nil {
		var bs []byte
		bs, isPrefix, err = r.ReadLine()
		line = append(line, bs...)
	}

	return line, err
}

func ToUint64(filePath string) (uint64, error) {
	content, err := FiltToTrimString(filePath)
	if err != nil {
		return 0, err
	}

	var ret uint64
	if ret, err = strconv.ParseUint(content, 10, 64); err != nil {
		return 0, err
	}
	return ret, nil
}
