package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
)

var tags = []string{
	"malcolm",
	"reynolds",
	"zoe",
	"washburn",
	"hoban",
	"wash",
	"inara",
	"serra",
	"jayne",
	"cobb",
	"heroofcanton",
	"kaylee",
	"frye",
	"kaywinnet",
	"simon",
	"tam",
	"river",
	"derrial",
	"book",
	"shepherd",
}

func randomBytes() []byte {
	var b bytes.Buffer
	l, _ := rand.Int(rand.Reader, big.NewInt(42*1024))
	io.CopyN(&b, rand.Reader, l.Int64())
	return b.Bytes()
}

func randomTags(num int) []string {
	t := make(map[string]bool)

	for {
		ti, _ := rand.Int(rand.Reader, big.NewInt(int64(len(tags))))
		t[tags[ti.Int64()]] = true

		if len(t) == num {
			break
		}
	}

	var out []string

	for k, _ := range t {
		out = append(out, k)
	}

	return out
}

func content(w http.ResponseWriter, r *http.Request) {
	for _, t := range randomTags(10) {
		w.Header().Add("X-XKey", t)
	}

	w.Header().Add("X-Some-Header", "whatever")

	b := randomBytes()
	h := md5.New()
	if _, err := h.Write(b); err == nil {
		w.Header().Add("Etag", fmt.Sprintf("%x", h.Sum(nil)))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func serverError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/error/", serverError)
	http.HandleFunc("/notfound/", notfound)
	http.HandleFunc("/ok/", content)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
