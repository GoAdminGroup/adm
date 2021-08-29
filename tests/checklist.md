# 测试清单

本文档列出`adm`的测试要点与自动化进度。

## 功能点

| 命令        | 功能                             | 实现自动化 | 优先级 |
| ----------- | -------------------------------- | ---------- | ------ |
| -V, version | 展示当前版本号，检查是否有新版本 | [x]        | P0     |
| generate    | 生成表格模型文件                 | [x]        | P0     |
| init        | 生成项目模板                     | [x]        | P0     |
| add         | 生成 user/permission/roles       | [x]        | P1     |
| develop     | 项目开发                         | [x]        | P2     |
| combine     | 合并前端资源                     | [x]        | P3     |
| compile     | 编译模板文件为 go 文件           | [x]        | P3     |

## 测试 case

| 命令                                                          | 预期          |
| ------------------------------------------------------------- | ------------- |
| adm version                                                   | 
| adm generate                                                  |
| adm generate -l cn                                            |
| adm generate -l cn -c adm.ini                                 |
| adm init                                                      |
| adm init -l cn                                                |
| adm init web                                                  |
| adm add user                                                  |
| adm add permission                                            |
| adm develop plug -m testplugin -n testplugin                  |
| adm develop tpl -m=github.com/GoAdminGroup/themes/ngx -n=ngx  |
| adm combine css --hash=true --path=./css/ --out=./all.min.css |
| adm combine js --hash=true --path=./js/ --out=./all.min.js    |
| adm compile tpl --src ./tmpls/ --dist ./tmpl.go               |
| adm compile asset -s=./assets/ -d=./ -p=theme1                |
