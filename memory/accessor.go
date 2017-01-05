package Memory

type statsAccessor interface {
	HitCount() uint64    // hit count
	MissCount() uint64   // miss hit count
	LookupCount() uint64 // all count = miss count + hit count
	HitRate() float64    // rate = hit/all
}
