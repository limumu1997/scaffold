# Golang scaffold

1. 跨平台服务脚手架，需要Linux平台直接编译即可；
2. 终端启动打开默认浏览器；
3. ~~使用 logrus 格式化日志，lumberjack.v2 切割日志~~; 替换成标准库 slog

#### 1.0.20
fix: 变更目录结构

#### 1.0.19
fix:
1. windows 服务模式默认路径由 system32 该到当前程序路径下

feat:
1. 增加 http 日志中间件、跨域处理中间件
2. 增加日志 prefix `logrus.WithField("prefix", "any logs")` 调用


#### 1.0.18
feat:
1. 简化项目结构

#### 1.0.16
feat:
1. 调整目录架构

#### 1.0.15
feat:
1. 更换日志文件轮转库

#### 1.0.14
fix: 
1. 日志等级默认为trace;
2. config 配置添加服务对象;
3. 卸载服务未完成而终止
feat:
1. 同步版本

``` shell
go build -ldflags=-w -o .\cmd\bin\scaffold.exe .\cmd\main.go
upx --best .\cmd\bin\scaffold.exe -o .\cmd\bin\scaffold-service.exe
```

## Commit 指南

Commit messages 请遵循[conventional-changelog 标准](https://www.conventionalcommits.org/en/v1.0.0/)：

```bash
<类型>[可选 范围]: <描述>

[可选 正文]

[可选 脚注]
```

### Commit 类型

以下是 commit 类型列表:

- feat: 新特性或功能
- fix: 缺陷修复
- docs: 文档更新
- style: 代码风格或者组件样式更新
- refactor: 代码重构，不引入新功能和缺陷修复
- perf: 性能优化
- test: 单元测试
- chore: 其他不修改 src 或测试文件的提交