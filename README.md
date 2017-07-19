# memory
## 特点： 
> 1、支持key-value匹配查询（hash）， 也支持近似值匹配查询（多阶二分查找）

> 2、支持多块block 或者 多value 异步同时加载

> 3、支持全局内存最大值设置 或 全局最大block对象设置；

> 4、支持lru、lfu和随机淘汰算法，同时也支持时间过期淘汰；

> 5、支持回调： load加载回调、加载成功回调、对象淘汰回调；

> 6、支持统计key的命中率（hit count 和 miss count）

> 7、支持多key（interface{}）-value（interface{}）并行加载


## characteristics:
> 1, support key-value matching queries (hash), and also support approximate matching queries (Multi - level two - point lookups)
> 2, support multi block block or multi value asynchronous simultaneous loading
> 3 supports global memory maximum settings or global maximum block object settings;
> 4 supports LRU, LFU, and random elimination algorithms, and also supports time expired elimination;
> 5, support callbacks: load loads callbacks, loads successful callbacks, objects, and eliminates callbacks;
> 6, support statistics key hit rate (hit, count, and miss count)
> 7, support multiple key (interface{}) -value (interface{}) parallel loading
