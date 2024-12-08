package config

import (
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"github.com/aarioai/airis/pkg/utils"
	"log"
	"strings"
)

func (c *Config) Log() {
	info := fmt.Sprintf(`
Launch Config:
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

	// 方便运行程序时直接显示
	afmt.Println(info, afmt.Green)
	// 日志无法显示颜色
	log.Println(info)
}

func (c *Config) Dump() {
	all := c.All()
	// 黄色 \033[33m
	// 结束符 \033[0m
	// 每行尽量保持小于80字符长度

	afmt.PrintBorder("Config Dump", afmt.Yellow, afmt.Bold)

	for category, configs := range all {
		fmt.Printf("\n")
		switch category {
		case "ini":
			afmt.PrintYellow("[%s] %s", category, c.path)
		case "rsa":
			afmt.PrintYellow("[%s] %s", category, c.getIni(CkRsaRoot))
		default:
			afmt.PrintYellow("[%s]", category)
		}
		fmt.Printf("\n")
		for _, d := range configs {
			v := d[1]
			if len(v) > 70 {
				v = strings.ReplaceAll(v[0:60], "\n", "\\n") + fmt.Sprintf("... (%dB)", len(v))
			}
			fmt.Print("  ")
			if d[2] != "" {
				afmt.PrintRed("%s = %s", d[0], v)
				afmt.PrintGreen(" %s\n", d[2])
			} else {
				afmt.PrintGreen(d[0])
				fmt.Println(" = " + v)
			}
		}
	}
	afmt.PrintBorder("End of Config", afmt.Yellow, afmt.Bold)
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
