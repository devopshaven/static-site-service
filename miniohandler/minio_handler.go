package miniohandler

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

const (
	bucket       = "websites"
	minioTimeout = time.Second * 2
	indexFile    = "index.html"
)

// InternalServerErrorHandler is a handler for providing internal server error.
// By default the handler returns text/plain message with "internal server error" string.
var InternalServerErrorHandler http.HandlerFunc

func init() {
	InternalServerErrorHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
}

// MinioHandler returns the http handler which serves content from MinIO server backend.
// The function takes a minio client which is a minio go client, a site parameter which
// is the site directory in the bucket eg.: websites/zamzi. When singlePage option is true
// then when the minio does not find a requested path it will try to render the /index.html
// as the result. This will useful for single page apps (React/Angular/etc..).
func MinioHandler(mc *minio.Client, site string, singlePage bool) http.Handler {
	log := log.With().Str("site", site).Logger()

	return &minioHandler{mc, site, singlePage, log}
}

type minioHandler struct {
	mc         *minio.Client
	site       string
	singlePage bool
	log        zerolog.Logger
}

// ServeHTTP is the default http handler func to handle incoming resource requests.
func (h *minioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	ctx, cancel := context.WithTimeout(context.Background(), minioTimeout)
	defer cancel()

	url := h.site + r.URL.Path

	log.Debug().Str("url", url).Msgf("handling request")
	obj, err := h.mc.GetObject(ctx, bucket, url, minio.GetObjectOptions{})

	if err != nil {
		cID := uuid.New().String()
		w.Header().Set("X-Correlation-Id", cID)
		InternalServerErrorHandler(w, r)

		log.Error().Err(err).Str("correlationId", cID).Msg("cannot call upstream")

		return
	}
	defer obj.Close()

	st, err := obj.Stat()
	if err != nil {
		if h.singlePage && h.RetryIndex(ctx, w, r) {
			return
		}

		http.NotFoundHandler().ServeHTTP(w, r)

		return
	}

	w.Header().Add("Content-Type", st.ContentType)
	w.WriteHeader(http.StatusOK)

	io.Copy(w, obj)
}

// RetryIndex function to render the index.html file when no object found and singlePage option is true.
func (h *minioHandler) RetryIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) (found bool) {
	found = false

	obj, err := h.mc.GetObject(ctx, bucket, h.site+indexFile, minio.GetObjectOptions{})

	if err != nil {

		return
	}
	defer obj.Close()

	st, err := obj.Stat()
	if err != nil {
		log.Trace().Err(err).Msg("cannot find index invoking not found handler")

		return
	}

	writeHeader(st, w)
	io.Copy(w, obj)

	return true
}

func writeHeader(info minio.ObjectInfo, w http.ResponseWriter) {
	header := w.Header()

	header.Add("Content-Type", info.ContentType)
	header.Add("Content-Length", strconv.FormatInt(info.Size, 10))
	header.Add("ETag", info.ETag)
	header.Add("X-Version", info.VersionID)
}
