package files

import (
	"bufio"
	"io"
	"os"
)

func ReadFileByLine(filename string, page, pageSize int, latest bool) (lines []string, isEndOfFile bool, total int, err error) {
	if !NewFileOp().Stat(filename) {
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	totalLines, err := countLines(filename)
	if err != nil {
		return
	}
	total = (totalLines + pageSize - 1) / pageSize
	reader := bufio.NewReaderSize(file, 8192)

	if latest {
		page = total
	}
	currentLine := 0
	startLine := (page - 1) * pageSize
	endLine := startLine + pageSize

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if currentLine >= startLine && currentLine < endLine {
			lines = append(lines, string(line))
		}
		currentLine++
		if currentLine >= endLine {
			break
		}
	}

	isEndOfFile = currentLine < endLine
	return
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	count := 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if count > 0 {
					count++
				}
				return count, nil
			}
			return count, err
		}
		count++
	}
}
