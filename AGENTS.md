# AGENTS.md - new-ls (nls) 项目指南

## 项目简介

**new-ls**（简称 `nls`）是一个使用 Go 语言编写的命令行目录列表工具，外观模仿 PowerShell 7 风格。提供彩色输出、Nerd Font 图标支持、多列简单模式等功能。

## 开发规范

### 代码修改汇报要求

**Agents 在工作时必须向用户汇报以下内容：**

1. **修改了什么文件**：列出所有被修改的文件路径
2. **具体修改内容**：简要描述每个文件中做了哪些改动
3. **修改目的**：解释为什么要做这个修改，解决什么问题或实现什么功能

## 项目结构

```
cmd/new-ls/main.go          # 程序入口
internal/
├── cmd/cmds.go             # CLI 命令定义、参数解析
├── config/config.go        # 配置文件加载
├── core/
│   ├── entry.go            # 目录读取、文件信息收集、排序
│   └── utils.go            # 工具函数
├── output/
│   ├── output.go           # NormalOutput — 详细表格输出
│   └── simple_output.go    # SimpleOutput — 多列简短输出
└── render/
    ├── color.go            # 文件/目录着色
    ├── icon_map.go         # Nerd Font 图标映射
    └── term.go             # 终端宽度检测
```

## 构建和运行

```bash
# 构建
go build -o nls ./cmd/new-ls/

# 运行
go run ./cmd/new-ls/

# 或使用脚本
./build.sh
./run.sh
```

## 代码风格

- 使用中文注释
- 包名小写简短（`cmd`、`config`、`core`、`output`、`render`）
- 导出函数使用 PascalCase
- 遵循 Go 标准 `cmd/internal` 项目布局

## 注意事项

- 项目当前版本为 v0.0.5，仍处于开发阶段
- 尚无单元测试，验证方式为手动测试
- 版本号硬编码在 `internal/cmd/cmds.go` 的 `const version` 中
