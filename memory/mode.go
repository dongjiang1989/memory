package Memory

type MODE string

const (
	SIMPLE MODE = "simple" // Simple mode: Random 随机
	LRU    MODE = "lru"    // Least Recently Used mode  最近最少使用
	LFU    MODE = "lfu"    // Least Frequently Used mode 最小频繁使用模式
	ARC    MODE = "arc"    // Adjustable Replacement Cache mode 可调换缓存模式
)
