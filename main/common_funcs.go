package main

import (
	"os"
	"strconv"
	"strings"
)

type UniqsWriter struct {
	Indent int
	file   *os.File
}

func NewUniqsWriter(file *os.File) *UniqsWriter {
	return &UniqsWriter{
		Indent: 0,
		file:   file,
	}
}

func (w *UniqsWriter) WriteLine(str string) {
	for i := 0; i < w.Indent; i++ {
		w.file.WriteString("\t")
	}
	w.file.WriteString(str + "\n")
}

func (w *UniqsWriter) IndentAdd() {
	w.Indent++
}

func (w *UniqsWriter) IndentSub() {
	w.Indent--
}

func ReplaceKeyValue(line string, key string, value string) string {
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+UniqsTemplateSuffix, value)
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+".TOLOWER"+UniqsTemplateSuffix, strings.ToLower(value))
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+".ToLower"+UniqsTemplateSuffix, strings.ToLower(value))
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+".toupper"+UniqsTemplateSuffix, strings.ToUpper(value))
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+".TOUPPER"+UniqsTemplateSuffix, strings.ToUpper(value))
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+".ToUpper"+UniqsTemplateSuffix, strings.ToUpper(value))
	line = strings.ReplaceAll(line, UniqsTemplatePrefixValue+key+".tolower"+UniqsTemplateSuffix, strings.ToLower(value))
	return line
}

func ReplaceLoopIdx(line string, lineIdx int) string {
	line = strings.ReplaceAll(line, UniqsLoopLineIdx, strconv.Itoa(lineIdx))
	return line
}
