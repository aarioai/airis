package debug

import (
	"fmt"
	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/arrmap"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxProfileLabelWidth     = 16
	printProfileTimePerMicro = 1000 // 每1000微秒输出1次时间
	profileTimeFormat        = "05.000"
)

type Profile struct {
	seq        atomic.Int32
	label      string
	parent     *Profile
	path       string
	startTime  int64
	lastTime   int64
	disabled   bool // 停止标志
	bufferPool *sync.Pool
}

var defaultProfile atomic.Pointer[Profile]

func init() {
	defaultProfile.Store(newProfile())
}

func newProfile() *Profile {
	return &Profile{
		startTime: time.Now().UnixMicro(),
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
	}
}

func DefaultProfile() *Profile {
	if p := defaultProfile.Load(); p != nil {
		return p
	}
	p := newProfile()
	defaultProfile.Store(p)
	return p
}
func Mark(marks ...string) int32 {
	return DefaultProfile().Mark(marks...)
}

// NewProfile 创建新的性能分析器
// 不重要的程序，而且异步概论几乎不存在
func NewProfile(labels ...string) *Profile {
	defaultProfile := DefaultProfile()
	now := time.Now().UnixMicro()

	if len(labels) == 0 && defaultProfile.seq.Load() == 0 {
		defaultProfile.startTime = now
		return defaultProfile
	}

	return &Profile{
		label:     arrmap.First(labels),
		startTime: now,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
	}
}

// WithLabel 设置标签
func (p *Profile) WithLabel(label string) *Profile {
	if p == DefaultProfile() {
		return newProfile().WithLabel(label)
	}
	p.label = label
	return p
}

// Disable 禁用性能分析
func (p *Profile) Disable() *Profile {
	p.disabled = true
	return p
}

// Enable 启用性能分析
func (p *Profile) Enable() *Profile {
	p.disabled = false
	return p
}

// IsDisabled 检查是否被禁用
func (p *Profile) IsDisabled() bool {
	if p.disabled {
		return true
	}
	if p.parent == nil {
		return false
	}
	return p.parent.IsDisabled()
}

func (p *Profile) Label() string {
	return p.label
}

// Mark 记录性能标记点
func (p *Profile) Mark(marks ...string) int32 {
	if p.IsDisabled() {
		return p.seq.Load()
	}

	// 程序并不重要，不存在并发，不必过度优化
	id := p.seq.Add(1)
	now := time.Now()
	nowMicro := now.UnixMicro()

	buf := p.bufferPool.Get().(*strings.Builder)
	defer func() {
		buf.Reset()
		p.bufferPool.Put(buf)
	}()

	p.writePrefix(buf, id)
	p.writeTimeInfo(buf, id, now, nowMicro)
	p.writeMarks(buf, marks)
	fmt.Print(buf.String())
	return id
}

// writePrefix 写入前缀信息
func (p *Profile) writePrefix(buf *strings.Builder, id int32) {
	if p.label != "" {
		buf.WriteString(p.label)
	}
	buf.WriteByte('[')
	buf.WriteString(p.buildPath(id))
	buf.WriteByte(']')
}

// writeTimeInfo 写入时间信息
func (p *Profile) writeTimeInfo(buf *strings.Builder, id int32, now time.Time, nowMicro int64) {
	elapsed := nowMicro - p.lastTime

	if id == 1 || elapsed > printProfileTimePerMicro {
		p.lastTime = nowMicro
		buf.WriteByte(' ')
		buf.WriteString(now.Format(profileTimeFormat))
	} else if delta := nowMicro - p.startTime; delta > 0 {
		fmt.Fprintf(buf, "+%dμs", delta)
	}
}

// writeMarks 写入标记信息
func (p *Profile) writeMarks(buf *strings.Builder, marks []string) {
	switch len(marks) {
	case 0:
		return
	case 1:
		buf.WriteString(strings.Repeat(" ", maxProfileLabelWidth-buf.Len()))
		buf.WriteString(marks[0])
	default:
		buf.WriteString(strings.Repeat(" ", maxProfileLabelWidth-buf.Len()))
		fmt.Fprintf(buf, marks[0], atype.ToAnySlice(marks[1:])...)
	}
	buf.WriteByte('\n')
}

// Fork 创建子性能分析器

// Fork 创建子性能分析器（替代原来的 Child）
func (p *Profile) Fork(marks ...string) *Profile {
	id := p.Mark(marks...)
	return &Profile{
		label:     p.label,
		parent:    p,
		path:      p.buildPath(id),
		startTime: time.Now().UnixMicro(),
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
	}
}

// buildPath 构建路径（替代原来的 fullId）
func (p *Profile) buildPath(id int32) string {
	s := atype.String(id)
	if p.path == "" {
		return s
	}
	return p.path + "." + s
}
