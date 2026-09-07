package output

import "github.com/fatih/color"

// 提示信息的键常量：调用方通过常量选键，避免手写字符串拼错
const (
	HintSortTimeNewest = "sort.time.newest_first" // -t：按修改时间排序，新到旧
	HintSortTimeOldest = "sort.time.oldest_first" // -tr：按修改时间排序，旧到新
)

// hintDict 统一存放各类提示文案，键采用 "类别.细节.方向" 的点分命名，便于按功能扩展。
// 后期扩展（例如按大小排序）时在字典中增加条目即可：
// "sort.size.large_first": "sort: size (large → small)"
var hintDict = map[string]string{
	HintSortTimeNewest: "sort: time (new → old)",
	HintSortTimeOldest: "sort: time (old → new)",
}

// Hint 按键取提示文案；键不存在时返回空字符串
func Hint(key string) string {
	return hintDict[key]
}

// printHint 在列表上方打印提示行，前后各空一行与输出分隔；key 无效时不输出任何内容
func printHint(key string) {
	if hint := hintDict[key]; hint != "" {
		// fmt.Printf("\n%s\n\n", hint)
		color.Yellow("\n%s\n\n", hint)
	}
}
