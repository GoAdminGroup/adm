package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/utils"
	"github.com/mgutz/ansi"
)

const version = "v1.2.26"

func cliInfo() {
	fmt.Println("GoAdmin CLI " + version + compareVersion(version))
	fmt.Println()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getLatestVersion() string {
	http.DefaultClient.Timeout = 3 * time.Second
	res, err := http.Get("https://goproxy.cn/github.com/!go!admin!group/adm/@v/list")

	if err != nil || res.Body == nil {
		return ""
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(res.Body)

	if err != nil || body == nil {
		return ""
	}

	versionsArr := strings.Split(string(body), "\n")

	return versionsArr[len(versionsArr)-1]
}

func compareVersion(srcVersion string) string {
	toCompareVersion := getLatestVersion()
	if utils.CompareVersion(srcVersion, toCompareVersion) {
		return ", the latest version is " + toCompareVersion + " now."
	}
	return ""
}

func printSuccessInfo(msg string) {
	fmt.Println()
	fmt.Println()
	fmt.Println(ansi.Color(getWord(msg), "green"))
	fmt.Println()
	fmt.Println()
}

func newError(msg string) error {
	return errors.New(getWord(msg))
}

func readFileFromFS(fs *embed.FS, name string) string {

	f, err := fs.ReadFile(name)

	checkError(err)

	return string(f)
}

func mkdirs(dirs []string) {
	for _, dir := range dirs {
		checkError(os.Mkdir(dir, os.ModePerm))
	}
}

func mkEmptyFiles(names []string) {
	for _, name := range names {
		checkError(os.WriteFile(name, []byte{}, os.ModePerm))
	}
}

func parseFile(content string, data interface{}) []byte {
	t, err := template.New("project").Funcs(map[string]interface{}{
		"title": strings.Title,
	}).Parse(content)
	checkError(err)
	buf := new(bytes.Buffer)
	checkError(t.Execute(buf, data))
	return buf.Bytes()
}

func parseFileWithFormat(content string, data interface{}) []byte {
	c, err := format.Source(parseFile(content, data))
	checkError(err)
	return c
}

func readFileOfInstallation(fs *embed.FS, name string) string {
	return readFileFromFS(fs, "templates/installation/"+name)
}

type WriteFileConfig struct {
	Path     string
	File     string
	Data     interface{}
	Perm     fs.FileMode
	IsFormat bool
	Parse    bool
}

type WriteFilesConfig []WriteFileConfig

func newWriteFilesConfig() *WriteFilesConfig {
	cfg := make(WriteFilesConfig, 0)
	return &cfg
}

func (cfgs *WriteFilesConfig) Add(path, file string, data interface{}, perm fs.FileMode) {
	*cfgs = append(*cfgs, WriteFileConfig{
		Path:     path,
		File:     file,
		Perm:     perm,
		Data:     data,
		IsFormat: path[len(path)-3:] == ".go",
		Parse:    true,
	})
}

func (cfgs *WriteFilesConfig) AddRaw(path, file string, perm fs.FileMode) {
	*cfgs = append(*cfgs, WriteFileConfig{
		Path:  path,
		File:  file,
		Perm:  perm,
		Parse: false,
	})
}

func writeFileOfInstallation(cfg WriteFileConfig) {
	if cfg.IsFormat {
		checkError(os.WriteFile(cfg.Path, parseFileWithFormat(readFileOfInstallation(&installationTmplFS, cfg.File), cfg.Data), cfg.Perm))
		return
	}
	checkError(os.WriteFile(cfg.Path, parseFile(readFileOfInstallation(&installationTmplFS, cfg.File), cfg.Data), cfg.Perm))
}

func writeFilesOfInstallation(cfgs WriteFilesConfig) {
	for _, cfg := range cfgs {
		if !cfg.Parse {
			checkError(os.WriteFile(cfg.Path, []byte(readFileOfInstallation(&installationTmplFS, cfg.File)), cfg.Perm))
			continue
		}
		writeFileOfInstallation(cfg)
	}
}
