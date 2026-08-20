# New-ls

<div align="center">

![](image/readme.png)

A beautiful,powershell 7-style new ls tool built from Golang

由Golang编写的一个漂亮的、powershell 7样式的新ls工具
</div>

>[!IMPORTANT]
>This repo is still under development; the current version is unstable

## 安装 / Install

在仓库根目录运行（Linux / macOS，默认 bash）：

```bash
./scripts/install.sh
```

- 将本项目编译为 `nls` 并安装到 `/usr/local/bin`，之后可直接在 shell 输入 `nls` 使用。
- 目标目录不可写时脚本会自动请求 `sudo`；也可用 `--prefix <dir>` 指定其它目录。
- 电脑未安装 Go 时，脚本会询问是否自动下载官方 Go 工具链到用户缓存目录（`~/.cache/nls-install`，仅用于本次构建，不影响系统）；加 `--no-go-install` 可拒绝并看到手动安装指引。
- 卸载：`./scripts/install.sh --uninstall`

## 用法 / Usage

```
nls [path...] [options]
```

不带路径时列出当前目录；可传入一个路径作为列出目标（目前仅取第一个参数）。

## 参数 / Options

| 参数 | 简写 | 说明 |
| --- | --- | --- |
| `--help` | `-h` | 显示帮助信息 |
| `--version` | `-v` | 仅输出版本信息 |
| `--simple` | `-s` | 简单模式：仅文件名，多列自适应终端宽度 |
| `--All` | `-A` | 输出所有文件（包含隐藏文件） |
| `--byte` | `-b` | 以字节形式输出文件大小（在 `-s` 模式下不可用） |

## 示例 / Examples

```bash
nls                 # normal 输出当前目录
nls /path/to/dir    # 输出指定目录
nls -A              # 包含隐藏文件
nls -s              # simple 模式
nls -b              # 以字节形式显示文件大小
nls -v              # 输出版本信息
```

## 配置 / Config

nls的配置项应该放置在以下目录：
- Linux：~/.config/nls/config.json
- macOS：~/Library/Application Support/nls/config.json

示例：

```json
{
    "enable_icon": true,
    "byteOutput": false
}
```

配置项：

- enable_icon(bool)：是否启用图标显示（需要提前安装Nerd font字体），默认为true
- byteOutput(bool)：文件大小是否以字节形式输出，默认为false
