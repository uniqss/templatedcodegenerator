package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type UniqueFilePart struct {
	IsLoop bool
	Lines  []string
}

func GenerateLoopUniqueFile(templatePath, templateFileName, outputFilePath, templateSuffix string) {
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

		templateLines = append(templateLines, line)
	}

	inLoop := false
	var part *UniqueFilePart = &UniqueFilePart{
		IsLoop: false,
	}
	var fileParts []*UniqueFilePart
	for _, tLine := range templateLines {
		if strings.Contains(tLine, UniqsTemplateLoopBegin) {
			if inLoop {
				panic("nested loop found. outputFilePath:" + outputFilePath + " templateFileName:" + templateFileName)
			}
			inLoop = true

			fileParts = append(fileParts, part)
			part = &UniqueFilePart{
				IsLoop: true,
			}

			continue
		}
		if strings.Contains(tLine, UniqsTemplateLoopEnd) {
			if !inLoop {
				panic("loop end found not in loop. outputFilePath:" + outputFilePath + " templateFileName:" + templateFileName)
			}
			inLoop = false

			fileParts = append(fileParts, part)
			part = &UniqueFilePart{
				IsLoop: false,
			}

			continue
		}

		part.Lines = append(part.Lines, tLine)
	}

	if part.IsLoop {
		panic("loop no end.")
	}
	fileParts = append(fileParts, part)

	fileName := templateFileName

	outputFileName := outputFilePath + strings.TrimSuffix(fileName, templateSuffix)
	outFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile err:", err)
		return
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	for _, part := range fileParts {
		if !part.IsLoop {
			for _, line := range part.Lines {
				writer.WriteString(line)
			}
		} else {
			lineIdx := 0
			for _, loopLine := range loopData.Data {
				for _, line := range part.Lines {
					for loopKey, loopValue := range loopLine {
						line = ReplaceKeyValue(line, loopKey, loopValue)
					}
					line = ReplaceLoopIdx(line, lineIdx)
					writer.WriteString(line)
					lineIdx++
				}
			}
		}
	}
}
