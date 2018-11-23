package cache

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	cacheVersion   = 0x5
	keyPrefetchLen = 256
)

// ParseError denotes an error parsing a cache file
type ParseError struct{}

// Error returns the error as a string
func (e *ParseError) Error() string {
	return "Error parsing cache meta data from file"
}

// Entry contains most meta data of a file cache entry
type Entry struct {
	Filename     string
	Version      uint64
	Expire       time.Time
	LastModified time.Time
	Date         time.Time
	Etag         string
	Key          string
	HeaderStart  int
	BodyStart    int
}

// Hash returns the MD5 hash of the cache key
func (e *Entry) Hash() (string, error) {
	h := md5.New()
	if _, err := h.Write([]byte(e.Key)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Response returns the HTTP response (sans body) parsed from cache file
func (e *Entry) Response() (*http.Response, error) {
	f, err := os.Open(e.Filename)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	if o, err := f.Seek(int64(e.HeaderStart), 0); err != nil || o != int64(e.HeaderStart) {
		return nil, &ParseError{}
	}

	buf := make([]byte, e.BodyStart-e.HeaderStart)

	if n, err := f.Read(buf); err != nil || n < len(buf) {
		return nil, &ParseError{}
	}

	reader := bufio.NewReader(bytes.NewReader(buf))
	res, err := http.ReadResponse(reader, nil)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FromFile reads cache metadata from a cache file
func FromFile(filename string) (*Entry, error) {
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	// assume cache key is at most keyPrefetchLen bytes long, read further later if necessary
	buf := make([]byte, offsetKey+keyPrefetchLen)

	if n, err := f.Read(buf); err != nil || n < len(buf) {
		return nil, &ParseError{}
	}

	e := &Entry{
		Filename: filename,
	}

	// read cache file version
	e.Version = uint64(buf[offsetVersion])

	if cacheVersion != e.Version {
		return e, fmt.Errorf("Unknown cache file version (%v), expects (%v)", e.Version, cacheVersion)
	}

	// read date metadata
	e.Expire = time.Unix(int64(binary.LittleEndian.Uint64(buf[offsetExpire:offsetExpire+8])), 0)
	e.LastModified = time.Unix(int64(binary.LittleEndian.Uint64(buf[offsetLastModified:offsetLastModified+8])), 0)
	e.Date = time.Unix(int64(binary.LittleEndian.Uint64(buf[offsetDate:offsetDate+8])), 0)

	// read length of etag value and fetch etag
	el := buf[offsetEtagLen : offsetEtagLen+1][0]
	e.Etag = string(buf[offsetEtag : offsetEtag+el])

	// read header & body offset position
	e.HeaderStart = int(binary.LittleEndian.Uint16(buf[offsetHeaderStart : offsetHeaderStart+2]))
	e.BodyStart = int(binary.LittleEndian.Uint16(buf[offsetBodyStart : offsetBodyStart+2]))

	if e.HeaderStart > len(buf) {
		// key is longer than anticipated, read more from file
		keybuf := make([]byte, e.HeaderStart-offsetKey-keyPrefetchLen-1)

		if n, err := f.Read(keybuf); err != nil || n < len(keybuf) {
			return nil, &ParseError{}
		}

		buf = append(buf, keybuf...)
	}

	// determine cache key
	e.Key = string(buf[offsetKey : e.HeaderStart-1])

	return e, nil
}
