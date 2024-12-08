package config

import (
	"fmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"github.com/aarioai/airis/pkg/utils"
	"log"
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

	// 方便运行程序时直接显示
	fmt.Println(info)
	// 记录进日志，方便通过消息队列通知
	log.Println(info)
}

func (c *Config) Dump() {
	all := c.All()
	// 黄色 \033[33m
	// 结束符 \033[0m
	fmt.Println("\n\033[33m====== Config Dump ======\033[0m")

	for category, configs := range all {
		if category == "ini" {
			fmt.Printf("\n\u001B[33m [%s] %s \u001B[0m\n", category, c.path)
		} else {
			fmt.Printf("\n\033[33m [%s] \033[0m\n", category)
		}
		for _, d := range configs {
			if d[2] != "" {
				// 红色 \033[31m
				// 绿色 \033[32m
				fmt.Printf("\033[31m    [%s = %s] \033[0m\033[32m %s \033[0m\n", d[0], d[1], d[2])
			} else {
				fmt.Printf("    %s = %s\n", d[0], d[1])
			}
		}
	}
	fmt.Println("\n\033[33m====== End of Config ======\033[0m")
}
func existsKey(arr [][3]string, key string) (string, bool) {
	for _, item := range arr {
		if item[0] == key {
			return item[1], true
		}
	}
	return "", false
}

// All 获取全部配置，仅用于调试
func (c *Config) All() map[string][][3]string {
	// 返回数组，可以保证输出key排序后的值
	result := map[string][][3]string{
		"other": make([][3]string, 0),
		"rsa":   make([][3]string, 0),
		"ini":   make([][3]string, 0),
	}
	// 优先级最高
	keys := arrmap.SortedKeys(c.otherConfig)
	for _, key := range keys {
		result["other"] = append(result["other"], [3]string{key, c.otherConfig[key]})
	}

	keys = arrmap.SortedKeys(c.rsa)
	for _, key := range keys {
		v := string(c.rsa[key])
		if newV, exists := existsKey(result["other"], key); exists {
			result["rsa"] = append(result["rsa"], [3]string{key, v, fmt.Sprintf("replaced by [other] %s=%s", key, newV)})
			continue
		}
		result["rsa"] = append(result["rsa"], [3]string{key, v})
	}
	keys = arrmap.SortedKeys(c.data)
	for _, key := range keys {
		v := c.data[key]
		if newV, exists := existsKey(result["other"], key); exists {
			result["ini"] = append(result["ini"], [3]string{key, v, fmt.Sprintf("replaced by [other] %s=%s", key, newV)})
			continue
		}
		if newV, exists := existsKey(result["rsa"], key); exists {
			result["ini"] = append(result["ini"], [3]string{key, v, fmt.Sprintf("replaced by [other] %s=%s", key, newV)})
			continue
		}
		result["ini"] = append(result["ini"], [3]string{key, v})
	}
	return result
}
