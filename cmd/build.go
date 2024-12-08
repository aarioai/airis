package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"os"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
)

var jsonFile string
var dstFile string
var aaJsFile string

const tab = "    "

// readJsonFile 读取并解析 JSON 文件
func readJsonFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open file error: %writeFile", err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	var buf bytes.Buffer
	var firstBraceHandled bool
	for scan.Scan() {
		line := bytes.TrimSpace(scan.Bytes())
		// @TODO 暂时只支持单行 // 方式注释
		if len(line) == 0 || (len(line) > 1 && string(line[0:2]) == "//") {
			continue
		}
		if !firstBraceHandled {
			if i := bytes.IndexByte(line, '{'); i > -1 {
				line = line[i:]
				firstBraceHandled = true
			}
		}
		key, value, ok := bytes.Cut(line, []byte(":"))
		key = bytes.TrimSpace(key)
		if ok && (len(key) < 3 || (key[0] != '"' && key[0] != '\'')) {
			line = append([]byte{'"'}, key...)
			line = append(line, '"', ':')
			line = append(line, value...)
		}

		buf.Write(line)
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scan file error: %w", err)
	}
	// 清理多余的逗号
	result := regexp.MustCompile(`,[\s\r\n]*([}\]])`).ReplaceAll(buf.Bytes(), []byte("$1"))
	return result, nil
}

// writeFile 写入文件内容
func writeFile(f *os.File, format string, args ...any) {
	_, err := f.WriteString(afmt.Sprintf(format, args...))
	if err != nil {
		panic(fmt.Sprintf("write file error: %v", err))
	}
}
func appendMimes(mimes []string, mime ...string) []string {
	for _, mi := range mime {
		var exists bool
		for _, m := range mimes {
			if mi == m {
				exists = true
			}
		}
		if !exists {
			mimes = append(mimes, mi)
		}
	}
	return mimes
}

func main() {
	root := flag.String("root", "../", "root of airis")
	flag.StringVar(&aaJsFile, "js", "", "path to generated filetype.js")
	flag.Parse()
	jsonFile = path.Join(*root, "core/aenum/filetype.jsonp")
	dstFile = path.Join(*root, "core/aenum/filetype_readonly.go")

	b, err := readJsonFile(jsonFile)
	if err != nil {
		panic(err)
	}
	var raw map[string]map[string][]any
	err = json.Unmarshal(b, &raw)
	if err != nil {
		panic(err.Error() + "\n" + string(b))
	}
	mimes := make([]string, 0)
	enums := make([][2]any, 0)
	types := make(map[string]map[string][]string)
	for category, r := range raw {
		types[category] = make(map[string][]string)
		for k, arr := range r {
			if len(arr) < 3 {
				panic(fmt.Sprintf("invalid %s.%s %v", category, k, arr))
			}
			id := int(arr[0].(float64))

			if len(enums) == 0 {
				enums = append(enums, [2]any{k, id})
			} else {
				var i int
				var prev [2]any
				var found bool
				for i, prev = range enums {
					if prev[1].(int) > id {
						found = true
						break
					}
				}
				if !found {
					i++
				}
				enums = slices.Insert(enums, i, [2]any{k, id})
			}
			ext := arr[1].(string)
			standardMime := arr[2].(string)
			types[category][k] = []string{ext, standardMime}
			mimes = appendMimes(mimes, ext, standardMime)
			if len(arr) > 3 {
				others := arr[3].([]any)
				for _, m := range others {
					types[category][k] = append(types[category][k], m.(string))
					mimes = appendMimes(mimes, m.(string))
				}
			}

		}
	}

	buildFileTypeGo(enums, types)
	buildFileTypeJS(enums, types, mimes)
}

// buildFileTypeGo 生成 Go 文件
func buildFileTypeGo(enums [][2]any, types map[string]map[string][]string) {
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writeFile(f, "package aenum\n")
	writeFile(f, `import "strings"`+"\n\n")

	//writeFile(f, "type FileType uint16\n")
	writeFile(f, "const (\n")
	writeFile(f, "%sUnknownType FileType = 0\n", tab)
	for _, enum := range enums {
		writeFile(f, fmt.Sprintf("%s%-11s FileType = %v\n", tab, enum[0], enum[1]))
	}
	writeFile(f, ")\n")

	for t, d := range types {
		writeFile(f, fmt.Sprintf("var %sTypes = map[FileType][]string{\n", t))
		ks := make([]string, 0, len(d))
		for k, _ := range d {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		for _, k := range ks {
			writeFile(f, fmt.Sprintf("%s%-11s : {", tab, k))
			for i, m := range d[k] {
				if i > 0 {
					writeFile(f, ", ")
				}
				writeFile(f, `"`+m+`"`)
			}
			writeFile(f, "},\n")
		}

		writeFile(f, "}\n")
	}
	// NewXxxType()
	for t, _ := range types {
		writeFile(f, fmt.Sprintf("func New%sType(mime string) (FileType, bool) {return ParseFileType(mime, %sTypes)}\n", t, t))
	}

	// Content Type
	writeFile(f, "func (t FileType) ContentType() string {\n")
	for t, _ := range types {
		writeFile(f, fmt.Sprintf("%sif d, ok := %sTypes[t]; ok {return d[1]}\n", tab, t))
	}
	writeFile(f, fmt.Sprintf("%sreturn \"\"\n", tab))
	writeFile(f, "}\n")

	// ext
	writeFile(f, "func (t FileType) Ext() string {\n")
	for t, _ := range types {
		writeFile(f, fmt.Sprintf("%sif d, ok := %sTypes[t]; ok {return d[0]}\n", tab, t))
	}
	writeFile(f, fmt.Sprintf("%sreturn \"\"\n", tab))
	writeFile(f, "}\n")

	// filename
	writeFile(f, `func (t FileType) Name() string {return strings.TrimPrefix(t.Ext(), ".")}`+"\n")
}
func buildFileTypeJS(enums [][2]any, types map[string]map[string][]string, mimes []string) {
	var dstFile = "./f_oss_filetype_readonly.js"
	fi, err := os.Stat(path.Dir(aaJsFile))
	if err == nil && fi.IsDir() {
		dstFile = aaJsFile
	}

	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writeFile(f, "/** @note this is an auto-generated file, do not modify it! */\n\n")
	writeFile(f, "/** @typedef {")
	writeFile(f, `"`+strings.Join(mimes, `"|"`)+`"`)
	writeFile(f, "} AaFileTypeMime */\n\n")

	writeFile(f, "class AaFileType {\n")
	writeFile(f, "%s/** @enum */\n", tab)
	writeFile(f, "%sstatic Enum={\n", tab)
	writeFile(f, "%s%sUnknownType: 0,\n", tab, tab)
	for _, enum := range enums {
		writeFile(f, fmt.Sprintf("%s%s%-11s : %v,\n", tab, tab, enum[0], enum[1]))
	}
	writeFile(f, "%s}\n", tab)

	writeFile(f, "%sstatic Mimes = {\n", tab)
	for t, d := range types {
		writeFile(f, fmt.Sprintf("%s%s%s : {\n", tab, tab, t))
		ks := make([]string, 0, len(d))
		for k, _ := range d {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		for _, k := range ks {
			writeFile(f, fmt.Sprintf("%s%s%s%-11s : [", tab, tab, tab, k))
			for i, m := range d[k] {
				if i > 0 {
					writeFile(f, ", ")
				}
				writeFile(f, `"`+m+`"`)
			}
			writeFile(f, "],\n")
		}

		writeFile(f, "%s%s},\n", tab, tab)
	}
	writeFile(f, "%s}\n", tab)
	writeFile(f, "%scontentType\n", tab)
	writeFile(f, "%sext\n", tab)
	writeFile(f, "%smimeType\n", tab)
	writeFile(f, "%svalue\n", tab)
	// constructor
	writeFile(f, "\n%s/**\n", tab)
	writeFile(f, "%s * @param {AaFileTypeMime|number} mime\n", tab)
	writeFile(f, "%s */\n", tab)
	writeFile(f, "%sconstructor(mime){", tab)
	writeFile(f, `
		this.value = AaFileType.Enum.UnknownType
		for(const [type, cv] of Object.entries(AaFileType.Mimes)){
			for(const [v,mimes] of Object.entries(cv)){
				 if(mime ===  AaFileType[v] || mimes.includes(mime)){
					this.contentType = mimes[1]
					this.ext = mimes[0]
					this.mimeType = type
					this.value = AaFileType[v]
					return
				}
			}
		}
`)
	writeFile(f, "%s}\n", tab)
	for t, _ := range types {
		writeFile(f, `%sis%s(){return this.mimeType === "%s"}`+"\n", tab, t, t)
	}
	writeFile(f, "%stoJSON(){return this.value}\n", tab)
	writeFile(f, "%svalueOf(){return this.value}\n", tab)
	writeFile(f, "}\n")
}
