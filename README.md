# GoModCheater

当使用 `go mod tidy` 等命令时， `go mod` 会从向网络请求 `module` 的版本列表及最新版本。在无法访问到 `proxy.golang.org` 或者 `goproxy.cn` 等镜像的环境中，如本地内网，本项目可以从 `GOPATH` 读取**已经缓存过**的 `module` 并返回版本列表及最新版本，通过设置 `GOPROXY` 的方式代理 `go mod` 。

### 使用方法
1. 使用命令 `go run main.go` 启动 `Cheater`
2. 使用命令 `go env -w GOPROXY=http://localhost:8080,direct`