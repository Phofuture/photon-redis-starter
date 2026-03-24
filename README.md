# Photon Redis Starter

一個用於 Photon 框架的 Redis 啟動器模組，提供簡單易用的 Redis 客戶端封裝，支援 standalone 和 cluster 模式。

## 安裝

```bash
go get github.com/Phofuture/photon-redis-starter
```

## 快速開始

### 1. 引入模組

```go
import _ "github.com/Phofuture/photon-redis-starter"
```

### 2. 配置 Redis

在你的配置文件中（如 `config.yaml`）添加 Redis 配置：

```yaml
redis:
  type: "standalone" # 或 "cluster"
  hosts:
    - "localhost:6379"
  password: ""
```

#### 配置說明

- `type`: Redis 客戶端類型
  - `standalone`: 單機模式
  - `cluster`: 集群模式（默認值）
- `hosts`: Redis 伺服器地址列表
  - standalone 模式：只使用第一個地址
  - cluster 模式：可配置多個節點地址
- `password`: Redis 密碼（可選）

### 3. 使用 Redis 操作

```go
package main

import (
    "context"
    "time"

    "github.com/Phofuture/photon-redis-starter/redis"
)

func main() {
    ctx := context.Background()

    // 基本操作
    err := redis.Set(ctx, "key", "value", time.Hour)
    if err != nil {
        // 處理錯誤
    }

    // 泛型 Get 操作
    value, err := redis.Get[string](ctx, "key")
    if err != nil {
        // 處理錯誤
    }

    // Hash 操作
    err = redis.HSet(ctx, "hash_key", map[string]interface{}{
        "field1": "value1",
        "field2": 123,
    })

    // 泛型 HGet 操作
    field1, err := redis.HGet[string](ctx, "hash_key", "field1")

    // List 操作
    count, err := redis.RPush(ctx, "list_key", "item1")
    items, err := redis.LRangeAll[string](ctx, "list_key")

    // Pub/Sub 操作
    err := redis.Publish(ctx, "channel_name", "message")

    subChan := redis.Subscribe(ctx, "channel_name")
    for msg := range subChan {
        fmt.Printf("%s\n", msg)
    }

    // Stream 操作（廣播：一生產者多消費者）
    id, err := redis.XAdd[string](ctx, "mystream", "hello")
    msgs, err := redis.XRead[string](ctx, "mystream", redis.StreamIDNewest, 0)
    for _, m := range msgs {
        fmt.Printf("id=%s data=%s\n", m.ID, m.Data)
    }
}
```

## API 文檔

### 基本操作

- `Set(ctx, key, value, expiration)` - 設置 key-value
- `Get[T](ctx, key)` - 獲取值（泛型）
- `Exists(ctx, key)` - 檢查 key 是否存在
- `SetNX(ctx, key, value, expiration)` - 僅當 key 不存在時設置
- `Del(ctx, keys...)` - 刪除一個或多個 key

### Hash 操作

- `HSet(ctx, key, values)` - 設置 hash 字段
- `HGet[T](ctx, key, field)` - 獲取 hash 字段值（泛型）

### List 操作

- `RPush(ctx, key, value)` - 從右側推入元素
- `RPop[T](ctx, key)` - 從右側彈出元素（泛型）
- `LRange[T](ctx, key, start, stop)` - 獲取範圍內的元素（泛型）
- `LRangeAll[T](ctx, key)` - 獲取所有元素（泛型）

### Stream 操作

支援 Redis Stream，適用於廣播情境（一生產者對應多消費者）。

#### 泛型 API（使用 "data" 欄位存放 JSON）

- `XAdd[T](ctx, stream, value)` - 新增訊息至 Stream
- `XAddWithMaxLen[T](ctx, stream, value, maxLen)` - 新增訊息並限制 Stream 長度
- `XRead[T](ctx, stream, startID, count)` - 非阻塞讀取（startID 可用 `StreamIDNewest` 或 `StreamIDStart`）
- `XReadBlock[T](ctx, stream, startID, block, count)` - 阻塞讀取，最多等待 block 時間
- `XRange[T](ctx, stream, start, stop, count)` - 依 ID 範圍查詢
- `XLen(ctx, stream)` - 取得 Stream 長度
- `XDel(ctx, stream, ids...)` - 刪除訊息

#### 原生 API（自訂 field-value）

- `XAddRaw(ctx, stream, values)` - 使用 map 或 []interface{} 自訂欄位
- `XAddRawWithMaxLen(ctx, stream, values, maxLen)` - 同上並限制長度
- `XReadRaw(ctx, stream, startID, count, block)` - 返回 go-redis 原生 `[]XStream`
- `XRangeRaw(ctx, stream, start, stop, count)` - 返回 go-redis 原生 `[]XMessage`

**注意**：standalone 與 cluster 皆支援單一 Stream 的 XADD、XREAD、XRANGE、XLEN、XDEL；cluster 模式下單一 key 的 Stream 操作會正確路由至對應 slot。

### 進階操作

- `EvalBool(ctx, script, keys, args...)` - 執行 Lua 腳本並返回 bool 結果

## 自定義初始化

如果需要在 Redis 客戶端初始化後執行自定義邏輯：

```go
import "github.com/Phofuture/photon-redis-starter/redis"

func init() {
    redis.RegisterRedisCustomize(func(ctx context.Context, client redis.RedisClient) error {
        // 你的自定義初始化邏輯
        // 例如：設置默認值、執行初始化腳本等
        return nil
    })
}
```

## 直接使用客戶端

如果需要使用原生 Redis 客戶端功能：

```go
client := redis.Redis()
result, err := client.Ping(context.Background()).Result()
```

## 配置範例

### Standalone 模式

```yaml
redis:
  type: "standalone"
  hosts:
    - "localhost:6379"
  password: "your-password"
```

### Cluster 模式

```yaml
redis:
  type: "cluster"
  hosts:
    - "redis-node1:6379"
    - "redis-node2:6379"
    - "redis-node3:6379"
  password: "your-password"
```

## 依賴

- [go-redis/redis](https://github.com/redis/go-redis) - Redis 客戶端
- [photon-core-starter](https://github.com/Phofuture/photon-core-starter) - Photon 核心框架
