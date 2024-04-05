package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func compileTmpl(rootPath, outputPath, packageName, varName string) {
	content := `package ` + packageName + `

var ` + varName + ` = map[string]string{`

	content = getContentFromDir(content, fixPath(rootPath), fixPath(rootPath))

	content += `}`

	_ = os.WriteFile(outputPath, []byte(content), 0644)
}

func fixPath(p string) string {
	if p[len(p)-1] != '/' {
		return p + "/"
	}
	return p
}

func getContentFromDir(content, dirPath, rootPath string) string {
	files, _ := os.ReadDir(dirPath)

	for _, f := range files {

		if f.IsDir() {
			content = getContentFromDir(content, dirPath+f.Name()+"/", rootPath)
			continue
		}

		b, err := os.ReadFile(dirPath + f.Name())
		if err != nil {
			fmt.Print(err)
		}
		str := string(b)

		suffix := path.Ext(f.Name())
		onlyName := strings.TrimSuffix(f.Name(), suffix)

		if suffix == ".tmpl" {
			fmt.Println(dirPath + f.Name())
			content += `"` + strings.ReplaceAll(dirPath, rootPath, "") + onlyName + `":` + "`" + str + "`,"
		}
	}

	return content
}
