package utils

import (
	"fmt"
	"sync"
	"time"
)

// 定义常量
const (
	twepoch          = int64(1577836800000) // 起始时间戳 (2020-01-01 00:00:00 UTC 的毫秒数)
	workerIDBits     = 5                    // 工作节点 ID 的位数
	datacenterIDBits = 5                    // 数据中心 ID 的位数
	sequenceBits     = 12                   // 序列号的位数

	// 计算最大值
	maxWorkerID     = -1 ^ (-1 << workerIDBits)     // 最大工作节点 ID
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits) // 最大数据中心 ID
	maxSequence     = -1 ^ (-1 << sequenceBits)     // 最大序列号

	// 定义移位偏移量
	workerIDShift     = sequenceBits
	datacenterIDShift = sequenceBits + workerIDBits
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

// NewSnowflake 初始化 Snowflake 实例
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, fmt.Errorf("datacenter ID must be between 0 and %d", maxDatacenterID)
	}
	return &Snowflake{
		lastTimestamp: -1,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

// NextID 生成唯一 ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6 // 当前毫秒时间戳

	// 如果当前时间小于上次生成 ID 的时间，说明系统时钟可能回退，返回错误或等待
	if timestamp < s.lastTimestamp {
		panic(fmt.Sprintf("clock moved backwards. refusing to generate id for %d milliseconds", s.lastTimestamp-timestamp))
	}

	// 如果是同一毫秒内生成的，则递增序列号
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		// 如果序列号超过最大值，则等待下一毫秒
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// 组合各部分生成最终 ID
	id := ((timestamp - twepoch) << timestampShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}
