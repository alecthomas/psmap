# Persistent static maps for Go

psmap generates binary files that contain persistent, static maps. It is useful for static data like a Geo IP database, geographical region boundaries, ngram tables, etc.

When opened, the psmap file is mmapped, and an in-memory index constructed. Index RAM overhead is approximately N*36 bytes.

## On-disk structure

In pseudo-Go code, the on disk structure is represented like this:

```go
struct PSM struct {
    Data [Size]KeyValue
}

struct KeyValue struct {
    KeySize uint32
    ValueSize uint32
    Key []byte
    Value []byte
}
```
