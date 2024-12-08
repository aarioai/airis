package config

import (
	"fmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"github.com/aarioai/airis/pkg/utils"
	"log"
	"strings"
)

func (c *Config) Log() {
	info := fmt.Sprintf(`
launch config:
 git version: %s
 env: %s
 timezone: %s
 mock: %v
 rsa: %s
`,
		utils.GitVersion(),
		c.Env,
		c.TimezoneID,
		c.Mock,
		arrmap.JoinKeys(c.rsa, ", ", true),
	)

	infoWithColor := "\033[32m" + info + "\033[0m"
	// 方便运行程序时直接显示
	fmt.Println(infoWithColor)
	// 日志无法显示颜色
	log.Println(info)
}

func (c *Config) Dump() {
	all := c.All()
	// 黄色 \033[33m
	// 结束符 \033[0m
	// 每行尽量保持小于80字符长度
	fmt.Printf("\n\033[33m================================ Config Dump ================================\033[0m\n")

	for category, configs := range all {
		switch category {
		case "ini":
			fmt.Printf("\n\u001B[33m [%s] %s \u001B[0m\n", category, c.path)
		case "rsa":
			fmt.Printf("\n\033[33m [%s] %s \033[0m\n", category, c.getIni(CkRsaRoot))
		default:
			fmt.Printf("\n\033[33m [%s] \033[0m\n", category)
		}
		for _, d := range configs {
			v := d[1]
			if len(v) > 70 {
				v = strings.ReplaceAll(v[0:60], "\n", "\\n") + fmt.Sprintf("... (%dB)", len(v))
			}
			if d[2] != "" {
				// 红色 \033[31m
				// 绿色 \033[32m
				fmt.Printf("    \u001B[31m%s = %s \033[0m\033[32m %s \033[0m\n", d[0], v, d[2])
			} else {
				fmt.Printf("    \033[32m%s\033[0m = %s\n", d[0], v)
			}
		}
	}
	fmt.Printf("\n\033[33m=============================== End of Config ===============================\033[0m\n")
}

func sortConfigKeys[T string | []byte](data map[string]T) [][3]string {
	result := make([][3]string, 0, len(data))
	defaultSectionKeys := make([]string, 0)
	sectionKeys := make([]string, 0)
	keys := arrmap.SortedKeys(data)
	for _, key := range keys {
		// 没有section的key
		if !strings.Contains(key, ".") {
			defaultSectionKeys = append(defaultSectionKeys, key)
		} else {
			sectionKeys = append(sectionKeys, key)
		}
	}
	// 1. 没有section的优先
	for _, key := range defaultSectionKeys {
		result = append(result, [3]string{key, string(data[key])})
	}
	// 2. 有section的key
	for _, key := range sectionKeys {
		result = append(result, [3]string{key, string(data[key])})
	}
	return result
}
func existsKey(arr [][3]string, key string) (string, bool) {
	for _, item := range arr {
		if item[0] == key {
			return item[1], true
		}
	}
	return "", false
}
func handleReplace(result map[string][][3]string, lowerPrioritySection string, prioritySection string) {
	for i, low := range result[lowerPrioritySection] {
		if newV, exists := existsKey(result[prioritySection], low[0]); exists {
			// 不能用 low，必须指定完整指针位置
			result[lowerPrioritySection][i][2] = fmt.Sprintf("replaced by [%s] %s=%s", prioritySection, low[0], newV)
		}
	}
}

// All 获取全部配置，仅用于调试
func (c *Config) All() map[string][][3]string {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()

	// 返回数组，可以保证输出key排序后的值
	result := map[string][][3]string{
		"other": make([][3]string, 0),
		"rsa":   make([][3]string, 0),
		"ini":   make([][3]string, 0),
	}
	// 优先级: other > rsa > ini
	result["other"] = sortConfigKeys(c.otherConfig)
	result["rsa"] = sortConfigKeys(c.rsa)
	result["ini"] = sortConfigKeys(c.data)
	handleReplace(result, "rsa", "other")
	handleReplace(result, "ini", "rsa")
	handleReplace(result, "ini", "other")
	return result
}
