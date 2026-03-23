# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A CGO bridge library that wraps TiKV's Go client (`github.com/tikv/client-go/v2`) as a C-compatible static library (`libtikvgo.a`). This enables C/C++ projects (specifically DingoDB's DingoFS) to use TiKV's transactional KV API via C function pointers and callbacks.

## Build

Build the static library (requires Go 1.21+ with CGO_ENABLED=1):

```bash
CGO_ENABLED=1 go build -buildmode=c-archive -o libtikvgo.a .
```

This produces `libtikvgo.a` and `libtikvgo.h`. The CMakeLists.txt integrates this into a parent CMake project as an imported static library target `tikvgo`.

There are no Go tests in this repo currently.

## Architecture

The entire bridge lives in a single file: `bridge.go` (package `main`).

**C ABI types** (defined in the cgo preamble):
- `CAsyncResult` — carries optional data bytes + optional error string
- `CKVPair` — a single key-value pair
- `CAsyncKVResult` — carries an array of KV pairs + optional error
- Callback typedefs: `tikv_go_callback` and `tikv_go_kv_callback`

**Exported functions** (all prefixed `tikv_go_`):
- `tikv_go_client_new` / `tikv_go_client_destroy` — synchronous client lifecycle
- `tikv_go_txn_begin` / `tikv_go_txn_id` / `tikv_go_txn_destroy` — synchronous transaction lifecycle
- `tikv_go_txn_get_async`, `tikv_go_txn_put_async`, `tikv_go_txn_delete_async` — async single-key ops
- `tikv_go_txn_batch_get_async`, `tikv_go_txn_scan_async` — async multi-key ops
- `tikv_go_txn_commit_async`, `tikv_go_txn_rollback_async` — async transaction finalization
- `tikv_go_free_async_result`, `tikv_go_free_kv_result`, `tikv_go_free_string` — C-heap memory freeing

**Handle pattern**: Go objects (`*txnkv.Client`, `*txnkv.KVTxn`) are stored via `runtime/cgo.Handle` and passed to C as `uint64_t` handles. The C side must pair each `_new`/`_begin` with a corresponding `_destroy`.

**Async pattern**: Async operations spawn a goroutine, allocate results on the C heap, and invoke a C callback function pointer. The caller is responsible for freeing results via the appropriate `tikv_go_free_*` function.

## Key Conventions

- All C-heap allocations use `C.malloc`/`C.CBytes`/`C.CString`; every allocation has a corresponding free path
- The `go.mod` has a commented-out `replace` directive for local development against a checkout of `github.com/tikv/client-go`; uncomment it and run `go mod tidy` to use a local copy
- Scan operations use inclusive `[start, end]` semantics; the bridge internally converts to Go's exclusive upper bound via `incrementBytes()`
