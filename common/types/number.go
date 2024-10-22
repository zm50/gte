package types

// Integer 整数类型
type Integer interface {
    ~int | ~int32 | ~int64 |
    ~uint | ~uint32 | ~uint64
}
