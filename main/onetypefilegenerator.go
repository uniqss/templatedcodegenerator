package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func GenerateOneTypeMultiFile(templatePath, templateFileName, outputFilePath, templateSuffix string) {
	file, err := os.OpenFile(templatePath+templateFileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var templateLines []string
	bWorking := true
	for bWorking {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			bWorking = false
		} else if err != nil {
			log.Error("reader.ReadString err:", err)
			break
		}

		// there must be no loop!!!
		if strings.Contains(line, UniqsTemplateLoopBegin) || strings.Contains(line, UniqsTemplateLoopEnd) {
			panic("GenerateOneTypeMultiFile there must be no loop!!!")
		}

		templateLines = append(templateLines, line)
	}

	for _, loopRow := range loopData.Data {
		fileName := templateFileName
		for loopKey, loopValue := range loopRow {
			fileName = ReplaceKeyValue(fileName, loopKey, loopValue)
		}

		outputFileName := outputFilePath + strings.TrimSuffix(fileName, templateSuffix)
		outFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Error("os.OpenFile err:", err)
			return
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)
		defer writer.Flush()

		for _, tLine := range templateLines {

			for k, v := range loopRow {
				tLine = ReplaceKeyValue(tLine, k, v)
			}
			writer.WriteString(tLine)
		}
	}
}
