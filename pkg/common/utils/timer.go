package utils

import "time"

// TimeLoop 循环执行某个函数，立刻执行然后等待周期
func TimeLoop(interval time.Duration, fn func(t time.Time)) func() {
	if fn == nil {
		return nil
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			fn(time.Now())
			<-ticker.C
		}
	}()
	return func() {
		ticker.Stop()
	}
}

// TimeLoopThen 循环执行某个函数，等待周期到达再执行
func TimeLoopThen(interval time.Duration, fn func(time.Time)) func() {
	if fn == nil {
		return nil
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			now := <-ticker.C
			fn(now)
		}
	}()
	return func() {
		ticker.Stop()
	}
}

// IsInMinute 判断是否在 begin-end 的分钟访问内
// 参数是时间戳（单位秒）
func IsInMinute(t, begin, end int64) bool {
	now := time.Unix(t, 0).Format("15:04")
	t1 := time.Unix(begin, 0).In(time.UTC).Format("15:04")
	t2 := time.Unix(end, 0).In(time.UTC).Format("15:04")
	if end >= begin {
		return now >= t1 && now <= t2
	}
	// 如t1=22:00, t2=08:00，则进入到此处逻辑，相当于是跨天的逻辑
	return !(now >= t2 && now <= t1)
}

// DurationWrap 耗时统计的函数包裹
func DurationWrap(fn func()) time.Duration {
	now := time.Now()
	if fn != nil {
		fn()
	}
	return time.Since(now)
}
