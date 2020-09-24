package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

var (
	words = make(map[string]string)
)

func usage(executableName string) {
	log.Error("usage:", executableName, " templatePath outputFilePath word.txt templateSuffix")
	log.Error("    example:", executableName, " ./tmpl ./output words_sample.txt .tmpl")
	log.Error("    templatePath is the place to put your template files")
	log.Error("    outputFilePath is the output path of your generated code")
	log.Error("    postfix is the template file's postfix")
	log.Error("    word.txt is the words to replace. format:")
	log.Error("        name=Player")
	log.Error("        age=age")
	log.Error("    template file format:")
	log.Error("        ...username:=<%=name%>")
	log.Error("        Age=<%=age%>")
	log.Error("        nameLowerCase=<%=name.lower%>")
	log.Error("        nameUpperCase=<%=name.upper%>")
}

func TrimFilePath(filePath string) string {
	if !strings.HasSuffix(filePath, "\\") && !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	return filePath
}

func main() {
	args := os.Args
	if len(args) < 5 {
		usage(args[0])
		return
	}
	templatePath := args[1]
	outputFilePath := args[2]
	wordsFullPath := args[3]
	templateSuffix := args[4]

	templatePath = TrimFilePath(templatePath)
	outputFilePath = TrimFilePath(outputFilePath)

	err := os.MkdirAll(outputFilePath, os.ModePerm)
	if err != nil {
		log.Fatal("os.MkdirAll failed. err:", err)
	}

	files, err := ioutil.ReadDir(templatePath)
	if err != nil {
		log.Fatal("ioutil.ReadDir failed", err)
	}

	ok := ReadWords(wordsFullPath)
	if !ok {
		log.Error("ReadWords failed")
		return
	}

	for _, f := range files {
		fName := f.Name()
		log.Debug(fName)
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(fName, "~$") {
			continue
		}
		if !strings.HasSuffix(fName, templateSuffix) {
			continue
		}
		GenerateOneFile(templatePath, fName, outputFilePath, templateSuffix)
	}
}
