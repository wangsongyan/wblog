# wblog 
[示例地址](http://blog.wangsy.me/)

## 一、初衷
作为一名web开发程序员居然没有自己的博客，都不好意思对外宣称自己的开发web的。
以前也有写博客的习惯，但是都是用的现有的博客网站。

## 二、技术选型
1. web:[gin](https://github.com/gin-gonic/gin)
2. orm:[gorm](https://github.com/go-gorm/gorm)
3. database:[SQLite](github.com/glebarez/sqlite)/[MySQL](https://gorm.io/driver/mysql)
4. 文件存储:[smms图床](https://sm.ms)/[七牛云存储](https://www.qiniu.com/)
5. 配置文件 [go-toml](https://github.com/pelletier/go-toml)

## 三、项目结构
```
-wblog
    |-conf 配置文件目录
    |-controllers 控制器目录
    |-helpders 公共方法目录
    |-models 数据库访问目录
    |-static 静态资源目录
        |-css css文件目录
        |-images 图片目录
        |-js js文件目录
        |-libs js类库
    |-system 系统配置文件加载目录
    |-tests 测试目录
    |-views 模板文件目录
    |-main.go 程序执行入口
```
## 四、TODO
- [x] 文章、页面访问统计
- [x] Github登录发表评论
- [x] RSS
- [x] 定时备份系统数据
- [x] 邮箱订阅功能
- [x] 云存储切换
- [x] 支持MySQL数据库
- [x] 导航栏配置
- [ ] 系统日志
- [ ] 网站统计

## 五、运行项目
```
git clone https://github.com/wangsongyan/wblog
cd wblog
go mod tidy
go run main.go
```

## 六、项目发布
1. 本地发布
   - 下载安装[goreleaser](https://github.com/goreleaser/goreleaser/releases)
   - 执行命令`goreleaser release --snapshot --clean`
2. Github Actions
   ```bash
   git tag "v0.0.2"
   git push origin v0.0.2
   ```
3. 部署文件清单
   - conf #配置文件目录
   - static #静态资源目录
   - views #模板目录
   - wblog #可执行文件

## 七、使用方法
### 使用说明
1. 执行`go run main.go -g`或编译后执行`wblog -g`生成示例配置文件`conf/conf.sample.toml` (示例配置文件均为系统默认配置，可全部删除仅保留自己所需配置)
2. 修改conf.toml，设置signup_enabled = true
3. 访问http://xxx.xxx/signup 注册管理员账号 
4. 修改conf.toml，设置signup_enabled = false

### 注意事项
1. 如果需求上传图片功能请自行申请七牛云存储空间，并修改配置文件填写
    ```toml
   [qiniu]
   enabled = true
   accesskey = 'AK'
   secretkey = 'SK'
   fileserver = '自定义域名，例如https://example.com'
   bucket = 'wblog'
   ```
2. 如果需要github登录评论功能请自行注册[github oauthapp](https://github.com/settings/developers)，并修改配置文件填写
    ```toml
   [github]
   enabled = true
   clientid = ''
   clientsecret = ''
   redirecturl = 'https://example.com/oauth2callback'
   ```
3. 如果需要使用邮件订阅功能，请自行填写
   ```toml
   [smtp]
   enabled = true
   username = '用户名'
   password = '密码'
   host = 'smtp.163.com:25'
   ```
4. GoLand运行时，修改`Run/Debug Configurations` > `Output Directory`选择到项目根目录，否则报模版目录找不到

## 八、效果图

![file](screenshots/index.png)

![file](screenshots/blog.png)

![file](screenshots/admin.png)

## 九、捐赠
如果项目对您有帮助，打赏个鸡腿吃呗！  
<img src="https://raw.githubusercontent.com/wangsongyan/wblog/master/screenshots/alipay.png" width = 40% height = 40% />
<img src="https://raw.githubusercontent.com/wangsongyan/wblog/master/screenshots/weixin.png" width = 40% height = 40% />
