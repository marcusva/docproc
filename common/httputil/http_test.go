package httputil_test

import (
	"bytes"
	"github.com/marcusva/docproc/common/httputil"
	"github.com/marcusva/docproc/common/testing/assert"
	"net/http/httptest"
	"testing"
)

func TestResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	err := httputil.Response(rec, 200, "OK")
	assert.NoErr(t, err)
	rec.Flush()

	assert.Equal(t, rec.Code, 200)
	assert.Equal(t, rec.Body.String(), "OK")
}

func TestError(t *testing.T) {
	rec := httptest.NewRecorder()

	err := httputil.Error(rec, "FAILED")
	assert.NoErr(t, err)
	rec.Flush()

	assert.Equal(t, rec.Code, 500)
	assert.Equal(t, rec.Body.String(), "FAILED")
}

func TestBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()

	err := httputil.BadRequest(rec, "BAD_REQUEST")
	assert.NoErr(t, err)
	rec.Flush()

	assert.Equal(t, rec.Code, 400)
	assert.Equal(t, rec.Body.String(), "BAD_REQUEST")
}

func TestInternalServerError(t *testing.T) {
	rec := httptest.NewRecorder()

	err := httputil.InternalServerError(rec)
	assert.NoErr(t, err)
	rec.Flush()

	assert.Equal(t, rec.Code, 500)
	assert.Equal(t, rec.Body.String(), "Internal Server Error")
}

func TestNotFound(t *testing.T) {
	rec := httptest.NewRecorder()

	err := httputil.NotFound(rec)
	assert.NoErr(t, err)
	rec.Flush()

	assert.Equal(t, rec.Code, 404)
	assert.Equal(t, rec.Body.String(), "resource not found")
}

func TestReadBody(t *testing.T) {
	rec := httptest.NewRecorder()

	input := "some test body"
	body := bytes.NewBufferString(input)
	req := httptest.NewRequest("POST", "http://localhost", body)

	buf, err := httputil.ReadBody(rec, req)
	rec.Flush()
	assert.NoErr(t, err)
	result := string(buf)
	assert.Equal(t, result, input)
	assert.Equal(t, rec.Code, 200)

	req = httptest.NewRequest("POST", "http://localhost", nil)
	buf, err = httputil.ReadBody(rec, req)
	rec.Flush()
	assert.NoErr(t, err)
	result = string(buf)
	assert.Equal(t, result, "")
	assert.Equal(t, rec.Code, 200)

}
