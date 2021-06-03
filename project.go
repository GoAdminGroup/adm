package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GoAdminGroup/go-admin/modules/db"

	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/mgutz/ansi"
	"gopkg.in/ini.v1"
)

type Project struct {
	Port         string
	Theme        string
	Prefix       string
	Language     string
	Driver       string
	DriverModule string
	Framework    string
	Module       string
	Orm          string
}

func buildProject(cfgFile string) {

	clear(runtime.GOOS)
	cliInfo()

	var (
		p        Project
		cfgModel *ini.File
		err      error
		info     = new(dbInfo)
	)

	if cfgFile != "" {
		cfgModel, err = ini.Load(cfgFile)

		if err != nil {
			panic(errors.New("wrong config file path"))
		}

		languageCfg, err := cfgModel.GetSection("language")

		if err == nil {
			setDefaultLangSet(languageCfg.Key("language").Value())
		}

		projectCfgModel, err := cfgModel.GetSection("project")

		if err == nil {
			p.Theme = projectCfgModel.Key("theme").Value()
			p.Framework = projectCfgModel.Key("framework").Value()
			p.Language = projectCfgModel.Key("language").Value()
			p.Port = projectCfgModel.Key("port").Value()
			p.Driver = projectCfgModel.Key("driver").Value()
			p.Prefix = projectCfgModel.Key("prefix").Value()
			p.Module = projectCfgModel.Key("module").Value()
			p.Orm = projectCfgModel.Key("orm").Value()
		}

		info = getDBInfoFromINIConfig(cfgModel, "default")
	}

	// generate main.go

	initSurvey()

	if p.Module == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		p.Module = promptWithDefault(getWord("module path"), filepath.Base(dir))
	}
	if p.Framework == "" {
		p.Framework = singleSelect(getWord("choose framework"),
			[]string{"gin", "beego", "buffalo", "fasthttp", "echo", "chi", "gf", "gorilla", "iris"}, "gin")
	}
	if p.Theme == "" {
		p.Theme = singleSelect(getWord("choose a theme"), template2.DefaultThemeNames, "sword")
	}
	if p.Language == "" {
		p.Language = singleSelect(getWord("choose language"),
			[]string{getWord("cn"), getWord("en"), getWord("jp"), getWord("tc")},
			getWord("cn"))
		switch p.Language {
		case getWord("cn"):
			p.Language = language.CN
		case getWord("en"):
			p.Language = language.EN
		case getWord("jp"):
			p.Language = language.JP
		case getWord("tc"):
			p.Language = language.TC
		}
	}
	if p.Port == "" {
		p.Port = promptWithDefault(getWord("port"), "80")
	}
	if p.Prefix == "" {
		p.Prefix = promptWithDefault(getWord("url prefix"), "admin")
	}
	if p.Driver == "" {
		p.Driver = singleSelect(getWord("choose a driver"),
			[]string{"mysql", "postgresql", "sqlite", "mssql"}, "mysql")
	}
	p.DriverModule = p.Driver
	if p.Driver == db.DriverPostgresql {
		p.DriverModule = "postgres"
	}

	rootPath, err := os.Getwd()

	if err != nil {
		rootPath = "."
	} else {
		rootPath = filepath.ToSlash(rootPath)
	}

	var cfg = config.SetDefault(&config.Config{
		Debug: true,
		Env:   config.EnvLocal,
		Theme: p.Theme,
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:          p.Language,
		UrlPrefix:         p.Prefix,
		IndexUrl:          "/",
		AccessLogPath:     rootPath + "/logs/access.log",
		ErrorLogPath:      rootPath + "/logs/error.log",
		InfoLogPath:       rootPath + "/logs/info.log",
		BootstrapFilePath: rootPath + "/bootstrap.go",
		GoModFilePath:     rootPath + "/go.mod",
	})

	if info.DriverName == "" && p.Driver != "" {
		info.DriverName = p.Driver
	}

	cfg.Databases = askForDBConfig(info)

	if p.Orm == "" {
		p.Orm = singleSelect(getWord("choose a orm"),
			[]string{getWord("none"), "gorm"}, getWord("none"))
		if p.Orm == getWord("none") {
			p.Orm = ""
		}
	}

	installProjectTmpl(p, cfg, cfgFile, info)

	fmt.Println()
	fmt.Println()
	fmt.Println(ansi.Color(getWord("Generate project template success~~🍺🍺"), "green"))
	fmt.Println()
	fmt.Println(getWord("1 Import and initialize database:"))
	fmt.Println()
	if defaultLang == "cn" || p.Language == language.CN || p.Language == "cn" {
		fmt.Println("- sqlite: " + ansi.Color("https://gitee.com/go-admin/go-admin/raw/master/data/admin.db", "blue"))
		fmt.Println("- mssql: " + ansi.Color("https://gitee.com/go-admin/go-admin/raw/master/data/admin.mssql", "blue"))
		fmt.Println("- postgresql: " + ansi.Color("https://gitee.com/go-admin/go-admin/raw/master/data/admin.pgsql", "blue"))
		fmt.Println("- mysql: " + ansi.Color("https://gitee.com/go-admin/go-admin/raw/master/data/admin.sql", "blue"))
	} else {
		fmt.Println("- sqlite: " + ansi.Color("https://github.com/GoAdminGroup/go-admin/raw/master/data/admin.db", "blue"))
		fmt.Println("- mssql: " + ansi.Color("https://raw.githubusercontent.com/GoAdminGroup/go-admin/master/data/admin.mssql", "blue"))
		fmt.Println("- postgresql: " + ansi.Color("https://raw.githubusercontent.com/GoAdminGroup/go-admin/master/data/admin.pgsql", "blue"))
		fmt.Println("- mysql: " + ansi.Color("https://raw.githubusercontent.com/GoAdminGroup/go-admin/master/data/admin.sql", "blue"))
	}
	fmt.Println()
	fmt.Println(getWord("2 Execute the following command to run:"))
	fmt.Println()
	if runtime.GOOS == "windows" {
		fmt.Println("> GO111MODULE=on go mod init " + p.Module)
		if defaultLang == "cn" || p.Language == language.CN || p.Language == "cn" {
			fmt.Println("> GORPOXY=https://goproxy.io GO111MODULE=on go mod tidy")
		} else {
			fmt.Println("> GO111MODULE=on go mod tidy")
		}
		fmt.Println("> GO111MODULE=on go run .")
	} else {
		fmt.Println("> make init module=" + p.Module)
		if defaultLang == "cn" || p.Language == language.CN || p.Language == "cn" {
			fmt.Println("> GORPOXY=https://goproxy.io make install")
		} else {
			fmt.Println("> make install")
		}
		fmt.Println("> make serve")
	}
	fmt.Println()
	fmt.Println(getWord("3 Visit and login:"))
	fmt.Println()
	if p.Port != "80" {
		fmt.Println("-  " + getWord("Login: ") + ansi.Color("http://127.0.0.1:"+p.Port+"/"+p.Prefix+"/login", "blue"))
	} else {
		fmt.Println("-  " + getWord("Login: ") + ansi.Color("http://127.0.0.1/"+p.Prefix+"/login", "blue"))
	}
	fmt.Println(getWord("account: admin  password: admin"))
	fmt.Println()
	fmt.Println("-  " + getWord("Generate CRUD models: ") + ansi.Color("http://127.0.0.1:"+p.Port+"/"+p.Prefix+"/info/generate/new", "blue"))
	fmt.Println()
	fmt.Println(getWord("4 See more in README.md"))
	fmt.Println()
	if defaultLang == "cn" {
		fmt.Println(getWord("see the docs: ") + ansi.Color("http://doc.go-admin.cn",
			"blue"))
	} else {
		fmt.Println(getWord("see the docs: ") + ansi.Color("https://book.go-admin.com",
			"blue"))
	}
	fmt.Println(getWord("visit forum: ") + ansi.Color("http://discuss.go-admin.com",
		"blue"))
	fmt.Println()
	fmt.Println()
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(dir, "\\", "/")
}

var langMap = map[string]string{
	"cn":        "cn",
	language.CN: "cn",
	"en":        "en",
	language.EN: "en",
}

//go:embed templates/installation
var installationTmplFS embed.FS

func installProjectTmpl(p Project, cfg *config.Config, cfgFile string, info *dbInfo) {

	mkdirs([]string{"pages", "tables", "logs", "uploads", "html", "build"})
	mkEmptyFiles([]string{"./logs/access.log", "./logs/info.log", "./logs/error.log"})

	var fileConfigs = newWriteFilesConfig()
	fileConfigs.Add("./main.go", "main.go/"+p.Framework+".tmpl", p, 0644)

	if p.Orm == "gorm" {
		checkError(os.Mkdir("models", os.ModePerm))
		fileConfigs.Add("./models/base.go", "orm.go.tmpl", p, os.ModePerm)
	}

	fileConfigs.Add("./pages/index.go", "pages/index.go."+p.Theme+".tmpl", p, os.ModePerm)
	fileConfigs.Add("./main_test.go", "main_test.go/"+langMap[p.Language]+".tmpl", p, 0644)
	fileConfigs.Add("./README.md", "readme/"+langMap[p.Language]+".tmpl", p, 0644)
	fileConfigs.Add("./config.yml", "config.yml/"+langMap[p.Language]+".tmpl", cfg, 0644)
	fileConfigs.Add("./Makefile", "makefile.tmpl", p, 0644)

	fileConfigs.AddRaw("./html/hello.tmpl", "hello.tmpl", 0644)
	fileConfigs.AddRaw("./tables/tables.go", "tables.go.tmpl", 0644)
	fileConfigs.AddRaw("./bootstrap.go", "bootstrap.go.tmpl", 0644)

	if cfgFile == "" {
		fileConfigs.Add("./adm.ini", "ini.tmpl", info, 0644)
	}

	writeFilesOfInstallation(*fileConfigs)
}
