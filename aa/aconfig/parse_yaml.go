package aconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/arrmap"
	"gopkg.in/yaml.v3"
)

const (
	yamlDefaultKey      = "DEFAULT"
	yamlDeletedSentinel = "__DELETED__"
)

// flattenYAML flattens a YAML tree into map[string]string with dotted keys.
// Sibling maps under the same parent that also defines DEFAULT inherit from it
// (empty scalars keep DEFAULT; maps are deep-copied then overlaid; __DELETED__
// removes a DEFAULT child).
func (c *Config) flattenYAML(prefix []string, node any, target map[string]string) error {
	if m, ok := asYAMLMap(node); ok {
		def, hasDefault := m[yamlDefaultKey]
		for k, child := range m {
			next := child
			if hasDefault && k != yamlDefaultKey {
				next = applyYAMLDefault(def, child)
			}
			if err := c.flattenYAML(append(prefix, k), next, target); err != nil {
				return err
			}
		}
		return nil
	}

	if len(prefix) == 0 {
		return nil
	}
	key := strings.Join(prefix, ".")
	val := yamlValueToString(node)
	if c.valueProcessor != nil {
		var err error
		if val, err = c.valueProcessor(key, val); err != nil {
			return err
		}
	}
	target[key] = val
	return nil
}

func asYAMLMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, child := range m {
			ks, ok := k.(string)
			if !ok {
				ks = fmt.Sprint(k)
			}
			out[ks] = child
		}
		return out, true
	default:
		return nil, false
	}
}

func yamlIsMap(v any) bool {
	_, ok := asYAMLMap(v)
	return ok
}

func yamlStructureConsistent(a, b any) bool {
	return yamlIsMap(a) == yamlIsMap(b)
}

func yamlIsEmpty(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return s == ""
	}
	return false
}

func yamlIsDeleted(v any) bool {
	s, ok := v.(string)
	return ok && s == yamlDeletedSentinel
}

func cloneYAMLValue(v any) any {
	if m, ok := asYAMLMap(v); ok {
		out := make(map[string]any, len(m))
		for k, child := range m {
			out[k] = cloneYAMLValue(child)
		}
		return out
	}
	if arr, ok := v.([]any); ok {
		out := make([]any, len(arr))
		for i, item := range arr {
			out[i] = cloneYAMLValue(item)
		}
		return out
	}
	return v
}

// applyYAMLDefault merges defaultVal into current when their shapes match.
// Non-maps: empty current takes defaultVal. Maps: clone DEFAULT then overlay.
func applyYAMLDefault(defaultVal, current any) any {
	if !yamlStructureConsistent(defaultVal, current) {
		return current
	}
	curMap, curIsMap := asYAMLMap(current)
	if !curIsMap {
		if yamlIsEmpty(current) {
			return defaultVal
		}
		return current
	}
	defMap, _ := asYAMLMap(defaultVal)
	merged := cloneYAMLValue(defMap).(map[string]any)
	for k, v := range curMap {
		if yamlIsDeleted(v) {
			delete(merged, k)
			continue
		}
		defChild, hasDef := merged[k]
		if hasDef && yamlStructureConsistent(defChild, v) {
			merged[k] = applyYAMLDefault(defChild, v)
			continue
		}
		if yamlIsEmpty(v) && hasDef {
			continue
		}
		merged[k] = v
	}
	return merged
}

func yamlValueToString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float64:
		// Prefer integer form when YAML decoded an unquoted number as float.
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case []any:
		parts := make([]string, 0, len(x))
		for _, item := range x {
			parts = append(parts, yamlValueToString(item))
		}
		return strings.Join(parts, ",")
	default:
		return fmt.Sprint(x)
	}
}

// loadYAML loads and flattens the YAML config file.
func (c *Config) loadYAML() (map[string]string, error) {
	raw, err := os.ReadFile(c.path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config %s: %w", c.path, err)
	}
	var root any
	if err = yaml.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config %s: %w", c.path, err)
	}
	data := make(map[string]string)
	if err = c.flattenYAML(nil, root, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Config) processData(data map[string]string) error {
	if c.valueProcessor == nil {
		return nil
	}

	for k, v := range data {
		newV, err := c.valueProcessor(k, v)
		if err != nil {
			return err
		}
		if newV != v {
			data[k] = newV
		}
	}
	return nil
}

func (c *Config) parseOtherConfig(otherConfigs ...map[string]string) (map[string]string, error) {
	cfgs := arrmap.Merge(otherConfigs...)
	if c.valueProcessor == nil {
		return cfgs, nil
	}
	for k, v := range cfgs {
		newV, err := c.valueProcessor(k, v)
		if err != nil {
			return nil, err
		}
		if newV != v {
			cfgs[k] = newV
		}
	}
	return cfgs, nil
}

// Reload reloads the base and text configs.
func (c *Config) Reload() error {
	c.startWrite()
	defer c.endWrite()

	data, err := c.loadYAML()
	if err != nil {
		return err
	}
	textConfigs, err := c.loadAllTextConfigs(data[CkTextConfigDirs])
	if err != nil {
		return err
	}

	// Keep the write-lock scope as small as possible.
	cfgMtx.Lock()
	c.baseConfig = data
	c.textConfig = textConfigs
	cfgMtx.Unlock()
	return c.initializeConfig()
}

func (c *Config) getBase(key string) string {
	if c.isOnWrite() {
		cfgMtx.RLock()
		defer cfgMtx.RUnlock()
	}
	if key == "" || len(c.baseConfig) == 0 {
		return ""
	}
	return c.baseConfig[key]
}

// Do not fetch at too fine a granularity, or errors become hard to trace.
func (c *Config) getOtherConfig(key string) string {
	if c.isOnWrite() {
		cfgMtx.RLock()
		defer cfgMtx.RUnlock()
	}
	if key == "" || len(c.otherConfig) == 0 {
		return ""
	}
	return c.otherConfig[key]
}

func (c *Config) MustGetString(key string) (string, error) {
	// 1. Prefer other config (for example values downloaded from a database). It can override settings stored in the yaml file.
	if v := c.getOtherConfig(key); v != "" {
		return v, nil
	}
	// 2. Read from the yaml config.
	if v := c.getBase(key); v != "" {
		return v, nil
	}
	// 3. Read from the text config.
	if rsa := c.getTextData(key); len(rsa) > 0 {
		return rsa, nil
	}
	return "", fmt.Errorf("required config key not found: %s", key)
}

func (c *Config) GetString(key string, defaultValue ...string) string {
	v, _ := c.MustGetString(key)
	if v != "" {
		return v
	}
	if len(defaultValue) > 0 {
		v = defaultValue[0]
	}
	return v
}

func (c *Config) MustGet(key string) (*atype.Atype, error) {
	v, err := c.MustGetString(key)
	if err != nil {
		return nil, err
	}
	return atype.New(v), nil
}

// Get value from config by key name
func (c *Config) Get(key string, defaultValue ...any) *atype.Atype {
	v, _ := c.MustGetString(key)
	if v != "" {
		return atype.New(v)
	}
	if len(defaultValue) > 0 {
		return atype.New(defaultValue[0])
	}
	return atype.New("")
}

// MustGetDuration parses a valid duration string from config
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// A duration string is a possibly signed sequence of decimal numbers,
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
func (c *Config) MustGetDuration(key string) (time.Duration, error) {
	s, err := c.MustGetString(key)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(s)
}

// GetDuration parses a duration string from config
func (c *Config) GetDuration(key string, defaultValue time.Duration) time.Duration {
	s := c.GetString(key)
	if s == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultValue
	}
	return d
}
