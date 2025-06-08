package debug

import (
	"fmt"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/afmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxProfileLabelWidth  = 16
	printProfileTimePerMs = 1000
	profileTimeFormat     = "05.000"
)

type Profile struct {
	seq        atomic.Int32
	label      string
	parent     *Profile
	path       string
	startTime  int64
	lastTime   int64
	disabled   bool
	bufferPool *sync.Pool

	styles []string
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

func NewProfile(labels ...string) *Profile {
	prof := DefaultProfile()
	now := time.Now().UnixMicro()

	if len(labels) == 0 && prof.seq.Load() == 0 {
		prof.startTime = now
		return prof
	}

	return &Profile{
		label:     afmt.First(labels),
		startTime: now,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
		styles: []string{afmt.Green},
	}
}

func (p *Profile) WithLabel(label string) *Profile {
	if p == DefaultProfile() {
		return newProfile().WithLabel(label)
	}
	p.label = label
	return p
}

func (p *Profile) WithStyles(styles ...string) *Profile {
	p.styles = styles
	return p
}

func (p *Profile) Disable() *Profile {
	p.disabled = true
	return p
}

func (p *Profile) Enable() *Profile {
	p.disabled = false
	return p
}

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

func (p *Profile) Mark(mark string) int32 {
	if p.IsDisabled() {
		return p.seq.Load()
	}

	// 程序并不重要，不存在并发，不必过度优化
	seq := p.seq.Add(1)
	now := time.Now()
	nowMicro := now.UnixMicro()

	buf := p.bufferPool.Get().(*strings.Builder)
	defer func() {
		buf.Reset()
		p.bufferPool.Put(buf)
	}()

	estimatedSize := maxProfileLabelWidth + len(p.label) + len(mark) + 10 + buf.Len() // 10 是 \n 等其他字符估计值；buf.Len 是保留以后扩展允许临时插入
	buf.Grow(estimatedSize)

	p.writePrefix(buf, seq)
	n := p.writeTimeInfo(buf, seq, now, nowMicro)
	p.writeMsg(buf, mark, n)

	fmt.Print(buf.String())

	return seq
}
func (p *Profile) Markf(format string, args ...any) int32 {
	return p.Mark(fmt.Sprintf(format, args...))
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
func (p *Profile) writeTimeInfo(buf *strings.Builder, id int32, now time.Time, nowMicro int64) int {
	prevStartTime := p.startTime
	p.startTime = nowMicro
	elapsed := nowMicro - p.lastTime
	if id == 1 || elapsed > printProfileTimePerMs {
		p.lastTime = nowMicro
		buf.WriteByte(' ')
		buf.WriteString(now.Format(profileTimeFormat))
		return 0
	}
	if delta := nowMicro - prevStartTime; delta > 0 {
		fmt.Fprintf(buf, "+%dμs", delta)
		return 1 // μ 是2个字节，需要增加1位
	}
	return 0
}

// writeMsg 写入标记信息
func (p *Profile) writeMsg(buf *strings.Builder, msg string, n int) {
	if msg == "" {
		return
	}
	padding := maxProfileLabelWidth + len(p.label) - buf.Len() + n
	if padding > 0 {
		buf.WriteString(strings.Repeat(" ", padding))
	}
	buf.WriteString(afmt.WithStyle(msg, p.styles...))
	buf.WriteByte('\n')
}

func (p *Profile) Fork(mark string) *Profile {
	if p.IsDisabled() {
		return p
	}
	id := p.Mark(mark)
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

func (p *Profile) Forkf(format string, args ...any) *Profile {
	return p.Fork(fmt.Sprintf(format, args...))
}

func (p *Profile) buildPath(id int32) string {
	s := atype.String(id)
	if p.path == "" {
		return s
	}
	return p.path + "." + s
}
