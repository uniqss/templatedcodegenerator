package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

var (
	indent = 0
)

func WriteLine(file *os.File, str string) {
	for i := 0; i < indent; i++ {
		file.WriteString("\t")
	}
	file.WriteString(str + "\n")
}

func GenerateOneFile(templatePath, templateFileName, outputFilePath, templateSuffix string) {
	file, err := os.OpenFile(templatePath + templateFileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile err:", err)
		return
	}
	defer file.Close()

	outputFileName := outputFilePath + strings.TrimSuffix(templateFileName, templateSuffix)
	outFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile err:", err)
		return
	}
	defer outFile.Close()

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	bWorking := true
	for bWorking {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			bWorking = false
		} else if err != nil {
			log.Error("reader.ReadString err:", err)
			break
		}

		for k, v := range words {
			line = strings.ReplaceAll(line, "<%=" + k + "%>", v)
			line = strings.ReplaceAll(line, "<%=" + k + ".tolower" + "%>", strings.ToLower(v))
			line = strings.ReplaceAll(line, "<%=" + k + ".toupper" + "%>", strings.ToUpper(v))
		}
		writer.WriteString(line)
	}
}
