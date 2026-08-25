# beer-abs — Go 溶液吸光度 Web 服务与 CLI 工具（Beer-Lambert 定律多组分求和 + 杂散光/带宽修正）
beer-abs 是一个纯 Go 标准库实现的溶液吸光度核算器，通过 HTTP JSON API（POST /api/absorbance）或命令行输入摩尔吸光系数、浓度与光程，按 Beer-Lambert 定律输出吸光度 A、透射比 T 以及杂散光修正后的观测值，支持多组分可加求和与矩形带宽平均。

## 构建 / 运行 / 测试

```text
go build ./...             # 编译
go run . absorbance example/cu-508nm.json   # CLI：核算一个算例并打印 A、T、A_obs
go test ./...              # 单元测试（law / band / mixture）
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

进容器后运行 `go build ./... && go test ./...`，再用 `go run . absorbance example/cu-508nm.json` 验证 CLI 输出 A 与 T。
