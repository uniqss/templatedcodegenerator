package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func ReadWords(wordsFullPath string) bool {
	file, err := os.OpenFile(wordsFullPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile err:", err)
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	bWorking := true
	for bWorking {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			bWorking = false
		} else if err != nil {
			log.Error("reader.ReadString err:", err)
			break
		}

		line = strings.TrimSpace(line)

		lineword := strings.Split(line, "=")
		if len(lineword) >= 2 {
			key := strings.TrimSpace(lineword[0])
			value := strings.TrimSpace(lineword[1])
			words[key] = value
		} else {
			log.Warn("line not right, should like this: name=HelloWorld")
			continue
		}
	}

	return true
}
