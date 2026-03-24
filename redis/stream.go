package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// StreamDataField 泛型 API 使用的預設欄位名稱，用於存放 JSON 序列化的 payload
const StreamDataField = "data"

// StreamMessage 泛型 Stream 訊息，包含 ID 與反序列化後的資料
type StreamMessage[T any] struct {
	ID   string
	Data T
}

// --- 泛型 API（使用 "data" 欄位存放 JSON）---

// XAdd 新增訊息至 Stream，value 會以 JSON 序列化後存入 "data" 欄位
// 支援 standalone 與 cluster 模式
func XAdd[T any](ctx context.Context, stream string, value T) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal stream value: %w", err)
	}
	return XAddRaw(ctx, stream, map[string]interface{}{StreamDataField: string(data)})
}

// XAddWithMaxLen 新增訊息並限制 Stream 長度（避免無限期增長）
// maxLen 為 0 時不限制
func XAddWithMaxLen[T any](ctx context.Context, stream string, value T, maxLen int64) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal stream value: %w", err)
	}
	return XAddRawWithMaxLen(ctx, stream, map[string]interface{}{StreamDataField: string(data)}, maxLen)
}

// XRead 非阻塞讀取 Stream，從 startID 之後的訊息
// startID 可使用 "$" 表示只讀新訊息，"0" 表示從頭開始
// 無新訊息時返回空 slice，count 為 0 時不限制數量
func XRead[T any](ctx context.Context, stream, startID string, count int64) ([]StreamMessage[T], error) {
	streams, err := XReadRaw(ctx, stream, startID, count, -1)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []StreamMessage[T]{}, nil
		}
		return nil, err
	}
	return parseStreamMessages[T](stream, streams)
}

// XReadBlock 阻塞讀取 Stream，最多等待 block 時間
// 適合廣播場景下消費者長駐等待新訊息
// 逾時且無新訊息時返回空 slice，count 為 0 時不限制數量
func XReadBlock[T any](ctx context.Context, stream, startID string, block time.Duration, count int64) ([]StreamMessage[T], error) {
	streams, err := XReadRaw(ctx, stream, startID, count, block)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []StreamMessage[T]{}, nil
		}
		return nil, err
	}
	return parseStreamMessages[T](stream, streams)
}

// XRange 依 ID 範圍查詢 Stream 訊息
// start、stop 使用 "-" 和 "+" 表示最小/最大 ID
func XRange[T any](ctx context.Context, stream, start, stop string, count int64) ([]StreamMessage[T], error) {
	msgs, err := XRangeRaw(ctx, stream, start, stop, count)
	if err != nil {
		return nil, err
	}
	return parseXMessages[T](msgs)
}

// --- 原生 API（貼近 go-redis，支援自訂 field-value）---

// XAddRaw 新增訊息至 Stream，使用自訂的 field-value
// values 格式與 go-redis XAdd 相同，支援 map 或 []interface{}{key, value, ...}
func XAddRaw(ctx context.Context, stream string, values interface{}) (string, error) {
	return rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Result()
}

// XAddRawWithMaxLen 新增訊息並限制 Stream 長度（避免無限期增長）
// maxLen 為 0 時不限制
func XAddRawWithMaxLen(ctx context.Context, stream string, values interface{}, maxLen int64) (string, error) {
	args := &redis.XAddArgs{
		Stream: stream,
		Values: values,
		MaxLen: maxLen,
	}
	return rdb.XAdd(ctx, args).Result()
}

// XReadRaw 讀取 Stream，返回 go-redis 原生 XStream 結構
// block < 0 表示非阻塞，block >= 0 表示阻塞等待
func XReadRaw(ctx context.Context, stream, startID string, count int64, block time.Duration) ([]redis.XStream, error) {
	streams := []string{stream, startID}
	args := &redis.XReadArgs{
		Streams: streams,
		Count:   count,
		Block:   block,
	}
	return rdb.XRead(ctx, args).Result()
}

// XRangeRaw 依 ID 範圍查詢，返回 go-redis 原生 []XMessage
// count 為 0 時不限制數量
func XRangeRaw(ctx context.Context, stream, start, stop string, count int64) ([]redis.XMessage, error) {
	if count > 0 {
		return rdb.XRangeN(ctx, stream, start, stop, count).Result()
	}
	return rdb.XRange(ctx, stream, start, stop).Result()
}

// XLen 取得 Stream 長度
func XLen(ctx context.Context, stream string) (int64, error) {
	return rdb.XLen(ctx, stream).Result()
}

// XDel 刪除 Stream 中的訊息
func XDel(ctx context.Context, stream string, ids ...string) (int64, error) {
	return rdb.XDel(ctx, stream, ids...).Result()
}

// parseStreamMessages 從 XRead 結果解析出泛型 StreamMessage
func parseStreamMessages[T any](stream string, streams []redis.XStream) ([]StreamMessage[T], error) {
	for _, s := range streams {
		if s.Stream == stream {
			return parseXMessages[T](s.Messages)
		}
	}
	return nil, nil
}

// parseXMessages 從 XMessage 解析 data 欄位至泛型
func parseXMessages[T any](msgs []redis.XMessage) ([]StreamMessage[T], error) {
	result := make([]StreamMessage[T], 0, len(msgs))
	for _, m := range msgs {
		raw, ok := m.Values[StreamDataField]
		if !ok {
			continue
		}
		str, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("stream message data field is not string: %T", raw)
		}
		var v T
		if err := json.Unmarshal([]byte(str), &v); err != nil {
			return nil, fmt.Errorf("unmarshal stream message: %w", err)
		}
		result = append(result, StreamMessage[T]{ID: m.ID, Data: v})
	}
	return result, nil
}
