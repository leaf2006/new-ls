package render

import (
	"os"
	"path/filepath"
	"strings"
)

var fileIconMap = map[string]string{
	// --- 编程语言 ---
	".go":    "\ue627", // 
	".py":    "\ue73c", // 
	".js":    "\ue74e", // 
	".mjs":   "\ue74e", // 
	".ts":    "\ue628", // 
	".tsx":   "\ue628", // 
	".jsx":   "\ue7ba", // 
	".java":  "\ue256", // 
	".c":     "\ue61e", // 
	".cpp":   "\ue61d", // 
	".cc":    "\ue61d", // 
	".h":     "\uf0fd", // 
	".hpp":   "\uf0fd", // 
	".rs":    "\ue7a8", // 
	".rb":    "\ue739", // 
	".php":   "\ue73d", // 
	".lua":   "\ue620", // 
	".swift": "\ue755", // 
	".dart":  "\ue7a3", // 
	".kt":    "\ue634", // 
	".scala": "\ue737", // 
	".pl":    "\ue769", // 
	".r":     "\uf25d", // 
	".zig":   "\ue6a9", // 

	// --- Web 与 样式 ---
	".html": "\ue736", // 
	".css":  "\ue749", // 
	".scss": "\ue603", // 
	".less": "\ue758", // 
	".vue":  "\ue6a0", // 
	".svg":  "\uf1c3", // 
	".wasm": "\ue6a1", // 

	// --- 配置文件 (DevOps & Tools) ---
	".json":         "\ue60b", // 
	".yaml":         "\ue601", // 
	".yml":          "\ue601", // 
	".toml":         "\ue6a2", // 
	".xml":          "\ue796", // 
	".conf":         "\ue615", // 
	".ini":          "\ue615", // 
	".env":          "\uf462", // 
	".dockerfile":   "\ue7b0", // 
	".dockerignore": "\ue7b0", // 
	".gitignore":    "\ue702", // 
	".tf":           "\ue695", //  (Terraform)
	".lock":         "\uf023", // 

	// --- 脚本与终端 ---
	".sh":       "\ue795", // 
	".zsh":      "\ue795", // 
	".bash":     "\ue795", // 
	".fish":     "\ue795", // 
	".bat":      "\ue70f", // 
	".ps1":      "\ue70f", // 
	".make":     "\ue615", // 
	".makefile": "\ue615", // 

	// --- 数据库 ---
	".sql":    "\ue706", // 
	".db":     "\uf1c0", // 
	".sqlite": "\ue706", // 
	".redis":  "\ue76d", // 

	// --- 文档 ---
	".md":  "\ue609", // 
	".pdf": "\uf1c1", // 
	".txt": "\uf15c", // 
	".csv": "\uf1c3", // 
	".log": "\uf18d", // 

	// --- 压缩与多媒体 ---
	".zip": "\uf410", // 
	".tar": "\uf410", // 
	".gz":  "\uf410", // 
	".7z":  "\uf410", // 
	".jpg": "\uf1c5", // 
	".png": "\uf1c5", // 
	".gif": "\uf1c5", // 
	".mp3": "\uf001", // 
	".mp4": "\uf03d", // 
	".exe": "\ue70f", // 
}

var normalFileIcon = "\uf15b" //

var folderIconMap = map[string]string{
	// --- 核心开发目录 ---
	".git":         "\ue5fb", // 
	"node_modules": "\ue5fa", // 
	"vendor":       "\uefa0", // 󮨠
	"bin":          "\ue5ff", //  (或用专用图标)
	"dist":         "\ufb4d", // 󰊭
	"build":        "\uf0ad", // 🔧
	"out":          "\ue5ff", // 

	// --- 资源与配置 ---
	".github":  "\ue5fd", // 
	".vscode":  "\ue70c", // 
	"config":   "\ue5fc", // 
	"settings": "\ue5fc", // 
	"assets":   "\uf115", // 📁
	"static":   "\uf115", // 📁
	"public":   "\uf415", // 
	"images":   "\uf1c5", // 
	"img":      "\uf1c5", // 
	"fonts":    "\uf031", // 

	// --- 逻辑分层 ---
	"app":      "\ue712", // 
	"internal": "\uf023", // 
	"pkg":      "\ufb2e", // 󰏞
	"api":      "\uf471", // 
	"docs":     "\uf18d", // 
	"test":     "\uf420", // 
	"tests":    "\uf420", // 
	"spec":     "\uf420", // 
	"scripts":  "\ue795", // 
	"temp":     "\uf014", // 
	"tmp":      "\uf014", // 
}

var normalFolderIcon = "\uf07b" //

func IconMap(file os.DirEntry) string {

	if file.IsDir() {
		if folderIcon, ok := folderIconMap[file.Name()]; ok {
			return folderIcon
		} else {
			return normalFolderIcon
		}
	} else {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if fileIcon, ok := fileIconMap[ext]; ok {
			return fileIcon
		} else {
			return normalFileIcon
		}
	}

}
