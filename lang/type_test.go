package lang_test

import (
	"errors"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

func TestType_Examples(t *testing.T) {
	is := is.New(t)

	typ, err := loadType("../testData/lang/function", "Receiver")
	is.NoErr(err)

	ex := typ.Examples()
	is.Equal(len(ex), 2)

	is.Equal(ex[0].Name(), "")
	is.Equal(ex[1].Name(), "Sub Test")
}

func TestFunc_netHttpResponseWriter(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("net/http")
	is.NoErr(err)

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	is.NoErr(err)

	var typ *lang.Type
	for _, t := range pkg.Types() {
		if t.Name() == "ResponseWriter" {
			typ = t
			break
		}
	}

	is.True(typ != nil) // didn't find the type we were looking for

	decl, err := typ.Decl()
	is.NoErr(err)

	is.Equal(typ.Level(), 2)
	is.Equal(typ.Name(), "ResponseWriter")
	is.Equal(typ.Title(), "type ResponseWriter")
	is.Equal(typ.Summary(), "A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.")
	is.Equal(decl, `type ResponseWriter interface {
    // Header returns the header map that will be sent by
    // WriteHeader. The Header map also is the mechanism with which
    // Handlers can set HTTP trailers.
    //
    // Changing the header map after a call to WriteHeader (or
    // Write) has no effect unless the HTTP status code was of the
    // 1xx class or the modified headers are trailers.
    //
    // There are two ways to set Trailers. The preferred way is to
    // predeclare in the headers which trailers you will later
    // send by setting the "Trailer" header to the names of the
    // trailer keys which will come later. In this case, those
    // keys of the Header map are treated as if they were
    // trailers. See the example. The second way, for trailer
    // keys not known to the Handler until after the first Write,
    // is to prefix the Header map keys with the TrailerPrefix
    // constant value. See TrailerPrefix.
    //
    // To suppress automatic response headers (such as "Date"), set
    // their value to nil.
    Header() Header

    // Write writes the data to the connection as part of an HTTP reply.
    //
    // If WriteHeader has not yet been called, Write calls
    // WriteHeader(http.StatusOK) before writing the data. If the Header
    // does not contain a Content-Type line, Write adds a Content-Type set
    // to the result of passing the initial 512 bytes of written data to
    // DetectContentType. Additionally, if the total size of all written
    // data is under a few KB and there are no Flush calls, the
    // Content-Length header is added automatically.
    //
    // Depending on the HTTP protocol version and the client, calling
    // Write or WriteHeader may prevent future reads on the
    // Request.Body. For HTTP/1.x requests, handlers should read any
    // needed request body data before writing the response. Once the
    // headers have been flushed (due to either an explicit Flusher.Flush
    // call or writing enough data to trigger a flush), the request body
    // may be unavailable. For HTTP/2 requests, the Go HTTP server permits
    // handlers to continue to read the request body while concurrently
    // writing the response. However, such behavior may not be supported
    // by all HTTP/2 clients. Handlers should read before writing if
    // possible to maximize compatibility.
    Write([]byte) (int, error)

    // WriteHeader sends an HTTP response header with the provided
    // status code.
    //
    // If WriteHeader is not called explicitly, the first call to Write
    // will trigger an implicit WriteHeader(http.StatusOK).
    // Thus explicit calls to WriteHeader are mainly used to
    // send error codes or 1xx informational responses.
    //
    // The provided code must be a valid HTTP 1xx-5xx status code.
    // Any number of 1xx headers may be written, followed by at most
    // one 2xx-5xx header. 1xx headers are sent immediately, but 2xx-5xx
    // headers may be buffered. Use the Flusher interface to send
    // buffered data. The header map is cleared when 2xx-5xx headers are
    // sent, but not with 1xx headers.
    //
    // The server will automatically send a 100 (Continue) header
    // on the first read from the request body if the request has
    // an "Expect: 100-continue" header.
    WriteHeader(statusCode int)
}`)
	is.Equal(len(typ.Examples()), 1)
	is.Equal(len(typ.Funcs()), 0)
	is.Equal(len(typ.Methods()), 0)
}

func TestFunc_netHttpResponse(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("net/http")
	is.NoErr(err)

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	is.NoErr(err)

	var typ *lang.Type
	for _, t := range pkg.Types() {
		if t.Name() == "Response" {
			typ = t
			break
		}
	}

	is.True(typ != nil) // didn't find the type we were looking for

	decl, err := typ.Decl()
	is.NoErr(err)

	is.Equal(typ.Level(), 2)
	is.Equal(typ.Name(), "Response")
	is.Equal(typ.Title(), "type Response")
	is.Equal(typ.Summary(), "Response represents the response from an HTTP request.")
	is.Equal(decl, `type Response struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"
    ProtoMajor int    // e.g. 1
    ProtoMinor int    // e.g. 0

    // Header maps header keys to values. If the response had multiple
    // headers with the same key, they may be concatenated, with comma
    // delimiters.  (RFC 7230, section 3.2.2 requires that multiple headers
    // be semantically equivalent to a comma-delimited sequence.) When
    // Header values are duplicated by other fields in this struct (e.g.,
    // ContentLength, TransferEncoding, Trailer), the field values are
    // authoritative.
    //
    // Keys in the map are canonicalized (see CanonicalHeaderKey).
    Header Header

    // Body represents the response body.
    //
    // The response body is streamed on demand as the Body field
    // is read. If the network connection fails or the server
    // terminates the response, Body.Read calls return an error.
    //
    // The http Client and Transport guarantee that Body is always
    // non-nil, even on responses without a body or responses with
    // a zero-length body. It is the caller's responsibility to
    // close Body. The default HTTP client's Transport may not
    // reuse HTTP/1.x "keep-alive" TCP connections if the Body is
    // not read to completion and closed.
    //
    // The Body is automatically dechunked if the server replied
    // with a "chunked" Transfer-Encoding.
    //
    // As of Go 1.12, the Body will also implement io.Writer
    // on a successful "101 Switching Protocols" response,
    // as used by WebSockets and HTTP/2's "h2c" mode.
    Body io.ReadCloser

    // ContentLength records the length of the associated content. The
    // value -1 indicates that the length is unknown. Unless Request.Method
    // is "HEAD", values >= 0 indicate that the given number of bytes may
    // be read from Body.
    ContentLength int64

    // Contains transfer encodings from outer-most to inner-most. Value is
    // nil, means that "identity" encoding is used.
    TransferEncoding []string

    // Close records whether the header directed that the connection be
    // closed after reading Body. The value is advice for clients: neither
    // ReadResponse nor Response.Write ever closes a connection.
    Close bool

    // Uncompressed reports whether the response was sent compressed but
    // was decompressed by the http package. When true, reading from
    // Body yields the uncompressed content instead of the compressed
    // content actually set from the server, ContentLength is set to -1,
    // and the "Content-Length" and "Content-Encoding" fields are deleted
    // from the responseHeader. To get the original response from
    // the server, set Transport.DisableCompression to true.
    Uncompressed bool

    // Trailer maps trailer keys to values in the same
    // format as Header.
    //
    // The Trailer initially contains only nil values, one for
    // each key specified in the server's "Trailer" header
    // value. Those values are not added to Header.
    //
    // Trailer must not be accessed concurrently with Read calls
    // on the Body.
    //
    // After Body.Read has returned io.EOF, Trailer will contain
    // any trailer values sent by the server.
    Trailer Header

    // Request is the request that was sent to obtain this Response.
    // Request's Body is nil (having already been consumed).
    // This is only populated for Client requests.
    Request *Request

    // TLS contains information about the TLS connection on which the
    // response was received. It is nil for unencrypted responses.
    // The pointer is shared between responses and should not be
    // modified.
    TLS *tls.ConnectionState
}`)
	is.Equal(len(typ.Examples()), 0)
	is.True(len(typ.Funcs()) > 0)
	is.True(len(typ.Methods()) > 0)
}

func loadType(dir, name string) (*lang.Type, error) {
	buildPkg, err := getBuildPackage(dir)
	if err != nil {
		return nil, err
	}

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return nil, err
	}

	for _, t := range pkg.Types() {
		if t.Name() == name {
			return t, nil
		}
	}

	return nil, errors.New("type not found")
}
