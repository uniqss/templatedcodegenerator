package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

var (
	loopData *CSVData = nil
)

func usage(executableName string) {
	log.Error("usage:", executableName, " templatePath outputFilePath loop.csv templateSuffix")
	log.Error("    example:", executableName, " ./tmpl ./output loop.csv .tmpl")
	log.Error("    templatePath is the place to put your template files")
	log.Error("    outputFilePath is the output path of your generated code")
	log.Error("    postfix is the template file's postfix")
	log.Error("    loop.csv is the loops to work. format:")
	log.Error("        name")
	log.Error("        Player")
	log.Error("        Item")
	log.Error("    template file format:")
	log.Error("        ...username:=<%=name%>")
	log.Error("        Age=<%=age%>")
	log.Error("        nameLowerCase=<%=name.tolower%>")
	log.Error("        nameUpperCase=<%=name.toupper%>")
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
	loopFullPath := args[3]
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

	ok := false
	//ok = ReadWords(wordsFullPath)
	//if !ok {
	//	log.Error("ReadWords failed")
	//	return
	//}

	ok = ReadLoop(loopFullPath)
	if !ok {
		log.Error("ReadLoop failed")
		return
	}

	// 如果文件名里面有关键字，一定会生成多个文件，循环loop中的每个关键字，挨个生成
	// 如果文件名里面没有关键字，则为唯一文件，只生成一个文件，要对loop中的关键字进行循环
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

		if strings.Contains(fName, UniqsTemplatePrefix) && strings.Contains(fName, UniqsTemplateSuffix) {
			// 如果文件名里面有关键字，一定会生成多个文件，循环loop中的每个关键字，挨个生成
			GenerateOneTypeMultiFile(templatePath, fName, outputFilePath, templateSuffix)
		} else {
			// 如果文件名里面没有关键字，则为唯一文件，只生成一个文件，要对loop中的关键字进行循环
			GenerateLoopUniqueFile(templatePath, fName, outputFilePath, templateSuffix)
		}
	}
}
