package utils

import (
	"errors"
	"math/rand"
	"time"

	"go.uber.org/atomic"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// LatencyShuffle 随机取值
type LatencyShuffle struct {
	indexs    []int64
	activeIdx *atomic.Int64
}

// NewShuffle 初始化随机序列
func NewShuffle(size int) *LatencyShuffle {
	sf := &LatencyShuffle{
		indexs:    make([]int64, size),
		activeIdx: atomic.NewInt64(0),
	}
	for i := 0; i < size; i++ {
		sf.indexs[i] = int64(i)
	}
	sf.shuffle()
	return sf
}

func (s *LatencyShuffle) shuffle() {
	size := len(s.indexs)
	if size < 2 {
		return
	}
	rand.Shuffle(size, func(i, j int) { s.indexs[i], s.indexs[j] = s.indexs[j], s.indexs[i] })
}

// TryNext 获取下一个值
func (s *LatencyShuffle) TryNext() int64 {
	next := s.activeIdx.Load() + 1
	if next >= int64(len(s.indexs)) {
		s.shuffle()
		next = 0
	}
	s.activeIdx.Store(next)
	return s.Get()
}

// Get 获取当前有效索引值
func (s *LatencyShuffle) Get() int64 {
	return s.indexs[s.activeIdx.Load()]
}

var (
	// InitValue 耗时选择的初始值
	InitValue int64 = -100
	// FailValue 耗时选择器的失败值
	FailValue int64 = -1
)

var (
	// ErrorOversize 延迟选择器超过大小
	ErrorOversize = errors.New("over size")
)

// Latency 按照访问耗时选择节点
type Latency struct {
	// 下标和urlList一一对应，下标对应，内容存放的是请求耗时（毫秒）
	// index for choose url from urlList, and the value is latency of this url(millisecond)
	history []int64
	// history的数量
	size int
	// 强制重置的次数
	resetSize int
	// 已经设置过延时的次数，超过resetSize就强制重置
	setCount int
}

// NewLatency 创建耗时选择器
func NewLatency(size int, resetSize int) *Latency {
	h := &Latency{
		history:   make([]int64, size),
		size:      size,
		resetSize: resetSize,
	}
	for i := range h.history {
		h.history[i] = InitValue
	}
	return h
}

// Reset 重置选择器中的记录，重选
func (l *Latency) Reset() {
	for i := range l.history {
		l.history[i] = InitValue
	}
	l.setCount = 0
}

// Set 设置索引当前耗时值
// latency is millisecond
func (l *Latency) Set(idx int, latency int64) {
	if idx+1 > len(l.history) {
		return
	}
	l.history[idx] = latency
	l.setCount++
	l.shouldReset()
}

func (l *Latency) shouldReset() {
	if l.resetSize > 0 && l.setCount > l.resetSize {
		l.Reset()
	}
}

// SetFail 设置当前索引失败
func (l *Latency) SetFail(idx int) {
	l.Set(idx, FailValue)
}

// Get 获取最小耗时的值的索引
func (l *Latency) Get() int {
	if l.size == 1 {
		return 0
	}
	var (
		minIndex   = 0
		minLatency = l.history[0]
		isFail     = 0
	)
	for i, lat := range l.history {
		if lat == FailValue {
			isFail++
			continue
		}
		if i == 0 {
			continue
		}
		if minLatency < 0 {
			minIndex = i
			minLatency = lat
			continue
		}
		if lat > 0 && lat < minLatency {
			minIndex = i
			minLatency = lat
		}
	}
	if isFail == l.size {
		l.Reset()
		return l.Rand()
	}
	if minLatency <= 0 {
		return l.Rand()
	}
	return minIndex
}

// Rand 随机选择
func (l *Latency) Rand() int {
	if l.size == 1 {
		return 0
	}
	var count int
	for {
		if count > l.size {
			return l.Get()
		}
		idx := rand.Intn(l.size*10) % l.size
		if l.history[idx] == FailValue {
			count++
			continue
		}
		return idx
	}
}

// History 获取历史数据
func (l *Latency) History() []int64 {
	cp := make([]int64, len(l.history))
	copy(cp, l.history)
	return cp
}

// GetValue 获取当前某个索引值的耗时
func (l *Latency) GetValue(idx int) (int64, error) {
	if l.size <= idx {
		return 0, ErrorOversize
	}
	return l.history[idx], nil
}
