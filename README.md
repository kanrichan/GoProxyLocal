# GoProxyLocal

当使用 `go mod tidy` 等命令时， `go mod` 会从向网络请求 `module` 的版本列表及最新版本。在无法访问到 `proxy.golang.org` 或者 `goproxy.cn` 等镜像的环境中，如本地内网，本项目可以从 `GOPATH` 的 `pkg/mod/cache` 读取**已经缓存过**的 `module` 并返回版本列表及最新版本，通过设置 `GOPROXY` 的方式代理 `go mod` 。

### 使用方法
1. 使用命令 `go run main.go` 启动 `Cheater`
2. 使用命令 `go env -w GOPROXY=http://localhost:9988,direct`

### 建议
1. 使用命令 `go env -w GOSUMDB=off` 关闭 `GOSUM`
2. 若拉取依赖时缓存无对应版本，可以在 `go.mod` 用 `replace` 替换相似版本
