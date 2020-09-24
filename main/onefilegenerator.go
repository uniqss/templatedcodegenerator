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

func ReplaceKeyValue(line string, key string, value string) string {
	line = strings.ReplaceAll(line, "@#%=" + key + "%#@", value)
	line = strings.ReplaceAll(line, "@#%=" + key + ".TOLOWER" + "%#@", strings.ToLower(value))
	line = strings.ReplaceAll(line, "@#%=" + key + ".ToLower" + "%#@", strings.ToLower(value))
	line = strings.ReplaceAll(line, "@#%=" + key + ".toupper" + "%#@", strings.ToUpper(value))
	line = strings.ReplaceAll(line, "@#%=" + key + ".TOUPPER" + "%#@", strings.ToUpper(value))
	line = strings.ReplaceAll(line, "@#%=" + key + ".ToUpper" + "%#@", strings.ToUpper(value))
	line = strings.ReplaceAll(line, "@#%=" + key + ".tolower" + "%#@", strings.ToLower(value))
	return line
}

func GenerateOneFile(templatePath, templateFileName, outputFilePath, templateSuffix string) {
	file, err := os.OpenFile(templatePath + templateFileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile err:", err)
		return
	}
	defer file.Close()

	outputFileName := outputFilePath + strings.TrimSuffix(templateFileName, templateSuffix)

	for k, v := range words {
		outputFileName = ReplaceKeyValue(outputFileName, k, v)
	}
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
			line = ReplaceKeyValue(line, k, v)
		}
		writer.WriteString(line)
	}
}
