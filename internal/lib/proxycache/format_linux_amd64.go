package proxycache

// cache file format from:
// https://github.com/nginx/nginx/blob/b66ee453cc6bc1832c3f056c9a46240bd390617c/src/http/ngx_http_cache.h#L126-L142
// following offsets are valid for 64bit architectures only
const (
	offsetVersion      = 0x0
	offsetExpire       = 0x8
	offsetLastModified = 0x20
	offsetDate         = 0x28
	offsetHeaderStart  = 0x36
	offsetBodyStart    = 0x38
	offsetEtagLen      = 0x3a
	offsetEtag         = 0x3b
	offsetKey          = 0x156
)
