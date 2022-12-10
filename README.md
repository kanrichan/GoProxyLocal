# GoModCheater

在无法访问到 `proxy.golang.org` 或者 `goproxy.cn` 等镜像的环境中，使用本项目可以从 `GOPATH` 读取已经缓存过的 `module` 并返回本地已经缓存了的版本。当使用 `go mod` 命令时，他会从 `GOPATH` 读取而不是请求网络最新版本。

### 使用方法
1. 使用命令 `go run main.go` 启动 `Cheater`
2. 使用命令 `go env -w GOPROXY=http://localhost:8080,direct`