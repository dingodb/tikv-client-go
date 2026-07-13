# tikv-client-go
用于dingofs 对接tikv

## GC

桥接层暴露 `tikv_go_client_gc_async`，用于异步执行 client-go 的 `GC(ctx, safepoint)`。
第二个参数是相对保留时间 `gc_lifetime_seconds`，单位为秒；例如传入 `60` 表示保留最近 1 分钟的数据。桥接层会从 PD 获取当前 TSO，将秒转换为毫秒后从 TSO 的物理时间部分扣减，并用 logical 0 组合出实际 GC safe point。
回调结果类型为 `CAsyncUInt64Result`：成功时 `value` 是 PD 返回的新 GC safe point；失败时 `error`/`error_len` 返回错误信息。调用方必须保持 `client_handle` 和 callback 在回调返回前有效，并在使用完结果后调用 `tikv_go_free_uint64_result` 释放一次。
