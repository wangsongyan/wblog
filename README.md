# wblog

## 初衷
作为一名web开发程序员居然没有自己的博客，都不好意思对外宣称自己的开发web的。
以前也有写博客的习惯，但是都是用的现有的博客网站。

## 技术选型
1. web:[gin](https://github.com/gin-gonic/gin)
2. orm:[gorm](https://github.com/jinzhu/gorm)
3. database:[sqlite3](https://github.com/mattn/go-sqlite3)
4. 全文检索:[wukong](https://github.com/huichen/wukong)
5. 文件存储:[七牛云存储](https://www.qiniu.com/)
6. 配置文件 [go-yaml](https://github.com/go-yaml/yaml)

## 项目结构
```
-wblog
    |-conf 配置文件目录
    |-controllers 控制器目录
    |-models 数据库访问目录
    |-static 静态资源目录
        |-css css文件目录
        |-fonts 字体目录
        |-images 图片目录
        |-js js文件目录
        |-libs js类库
    |-tests 测试目录
    |-views 模板文件目录
    |-main.go 程序执行入口
```
## 安装部署
```
go get -u -v github.com/wangsongyan/wblog
cd $GOPATH/src/github.com/wangsongyan/wblog
go run main.go
```