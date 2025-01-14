package id

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// GenShortID 单元测试用例
func TestGenShortID(t *testing.T) {
	shortID := GenShortID()
	assert.NotEqual(t, "", shortID)
	assert.Equal(t, 6, len(shortID))
}

// 性能测试 1
func BenchmarkGenShortID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortID()
	}
}

func BenchmarkGenShortIDTimeConsuming(b *testing.B) {
	//调用该函数停止压力测试的时间计数
	b.StopTimer()

	// 准备工作，比如单元测试
	shortID := GenShortID()

	if shortID == "" {
		b.Error("Failed to generate short ID")
	}

	// 重新开始时间
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		GenShortID()
	}
}
