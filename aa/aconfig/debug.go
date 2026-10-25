package aconfig

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"github.com/aarioai/airis/pkg/utils"
)

func (c *Config) PanicIfNotEqual(key, want string) {
	value := c.GetString(key)
	if value != want {
		ae.PanicF("config %s = %s not equal to %s", key, value, want)
	}
}
func (c *Config) Log() {
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		execPath = err.Error()
	}

	info := fmt.Sprintf(`
Launch Config:
  exec path: %s
  git version: %s
  env: %s
  timezone: %s
  mock: %v
`,
		execPath,
		utils.GitVersion(),
		c.Env,
		c.TimezoneID,
		c.Mock,
	)
	alog.Console(info, afmt.Green)
}

func (c *Config) Dump() {
	all := c.All()
	// Yellow \033[33m
	// Reset \033[0m
	// Keep each line under 80 characters when possible.

	afmt.PrintBorder("Config Dump", afmt.Yellow, afmt.Bold)

	for category, configs := range all {
		fmt.Printf("\n")
		switch category {
		case "base":
			afmt.PrintfYellow("[%s] %s", category, c.path)
		case "text":
			afmt.PrintfYellow("[%s] %s", category, strings.Join(c.TextConfigDirs, " "))
		default:
			afmt.PrintfYellow("[%s]", category)
		}
		fmt.Printf("\n")
		for _, d := range configs {
			v := d[1]
			if len(v) > 70 {
				v = strings.ReplaceAll(v[0:60], "\n", "\\n") + fmt.Sprintf("... (%dB)", len(v))
			}
			fmt.Print("  ")
			if d[2] != "" {
				afmt.PrintfRed("%s = %s", d[0], v)
				afmt.PrintfGreen(" %s\n", d[2])
			} else {
				afmt.PrintGreen(d[0])
				fmt.Println(" = " + v)
			}
		}
	}
	afmt.PrintBorder("End of Config", afmt.Yellow, afmt.Bold)
}

func sortConfigKeys[T string | []rune | []byte](data map[string]T) [][3]string {
	result := make([][3]string, 0, len(data))
	defaultSectionKeys := make([]string, 0)
	sectionKeys := make([]string, 0)
	keys := arrmap.SortedKeys(data)
	for _, key := range keys {
		// Keys with no section.
		if !strings.Contains(key, ".") {
			defaultSectionKeys = append(defaultSectionKeys, key)
		} else {
			sectionKeys = append(sectionKeys, key)
		}
	}
	// 1. Keys with no section come first.
	for _, key := range defaultSectionKeys {
		result = append(result, [3]string{key, string(data[key])})
	}
	// 2. Keys that have a section.
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
			// Do not write through low; it is a copy. Address the full index.
			result[lowerPrioritySection][i][2] = fmt.Sprintf("replaced by [%s] %s=%s", prioritySection, low[0], newV)
		}
	}
}

// All returns every config entry. For debugging only.
func (c *Config) All() map[string][][3]string {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()

	// Return slices so keys can be emitted in sorted order.
	result := map[string][][3]string{
		"other": make([][3]string, 0),
		"text":  make([][3]string, 0),
		"base":  make([][3]string, 0),
	}
	// Priority: other > text > yaml base.
	result["other"] = sortConfigKeys(c.otherConfig)
	result["text"] = sortConfigKeys(c.textConfig)
	result["base"] = sortConfigKeys(c.baseConfig)
	handleReplace(result, "text", "other")
	handleReplace(result, "base", "text")
	handleReplace(result, "base", "other")
	return result
}
