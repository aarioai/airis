package debug

import (
	"fmt"
	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/arrmap"
	"strings"
	"time"
)

type Tracer struct {
	tag      string
	parent   *Tracer
	parentId string
	id       int32
	time     int64
	stop     bool // 停止标志
	buffer   strings.Builder
}

// NewTracer 记录启动步骤
// 不重要的程序，而且异步概论几乎不存在
func NewTracer(tags ...string) *Tracer {
	return &Tracer{
		tag:    arrmap.First(tags),
		time:   time.Now().UnixMicro(),
		buffer: strings.Builder{},
	}
}
func (t *Tracer) WithTag(tag string) *Tracer {
	t.tag = tag
	return t
}
func (t *Tracer) Stop() *Tracer {
	t.stop = true
	return t
}
func (t *Tracer) Restart() *Tracer {
	t.stop = false
	return t
}

func (t *Tracer) Tag() string {
	return t.tag
}
func (t *Tracer) IsStop() bool {
	if t.stop {
		return true
	}
	if t.parent == nil {
		return false
	}
	return t.parent.IsStop()
}

func (t *Tracer) Trace(infos ...string) int32 {
	if t.IsStop() {
		return t.id
	}
	now := time.Now()
	// 程序并不重要，不存在并发，不必过度优化
	id := t.id + 1
	t.id = id
	// 首次追踪，输出时间戳
	if id == 1 {
		t.buffer.WriteByte('[')
		t.buffer.WriteString(now.Format("15:04:05.000"))
		t.buffer.WriteString("]\n")
	}
	// 计算耗时
	elapsed := now.UnixMicro() - t.time

	// 格式化输出
	t.buffer.WriteString("  ")
	t.buffer.WriteString(fullId(t.parentId, id))
	if elapsed > 0 {
		t.buffer.WriteString(fmt.Sprintf("(+%dμs)", elapsed))
	}
	// 处理信息
	switch len(infos) {
	case 0:
		// 无信息
	case 1:
		t.buffer.WriteString(": ")
		t.buffer.WriteString(infos[0])
	default:
		t.buffer.WriteString(": ")
		t.buffer.WriteString(fmt.Sprintf(
			infos[0],
			atype.ToAnySlice(infos[1:])...))
	}

	t.buffer.WriteByte('\n')

	// 输出并清空缓冲
	fmt.Print(t.buffer.String())
	t.buffer.Reset()

	return id
}

func (t *Tracer) Child(infos ...string) *Tracer {
	id := t.Trace(infos...)
	return &Tracer{
		tag:      t.tag,
		parent:   t,
		parentId: fullId(t.parentId, id),
		id:       0,
		time:     time.Now().UnixMicro(),
		stop:     false,
	}
}

func fullId(parentId string, id int32) string {
	s := atype.String(id)
	if parentId == "" {
		return s
	}
	return parentId + "." + s
}
