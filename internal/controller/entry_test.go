package controller

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"ipset-ui/internal/config"

	"github.com/gin-gonic/gin"
)

const fakeIpsetScript = `#!/bin/sh
if [ "$1" = "list" ]; then
  if [ "$2" = "-name" ]; then
    if [ "$FAKE_IPSET_SETS" = "empty" ]; then
      exit 0
    fi
    echo "testset"
    exit 0
  fi
  echo "1.1.1.1 comment \"test\""
  echo "2.2.2.2 comment \"test\""
  echo "3.3.3.3 comment \"test\""
  exit 0
fi
if [ "$1" = "add" ]; then
  if [ -n "$FAKE_IPSET_DELAY" ]; then
    sleep "$FAKE_IPSET_DELAY"
  fi
  if [ -n "$FAKE_IPSET_FAIL_ALL" ]; then
    echo "The set with the given name does not exist" >&2
    exit 1
  fi
  if [ -n "$FAKE_IPSET_FAIL" ] && [ "$3" = "$FAKE_IPSET_FAIL" ]; then
    echo "The entry with the given name already exists" >&2
    exit 1
  fi
  exit 0
fi
if [ "$1" = "del" ]; then
  if [ -n "$FAKE_IPSET_DELAY" ]; then
    sleep "$FAKE_IPSET_DELAY"
  fi
  if [ -n "$FAKE_IPSET_FAIL_ALL" ]; then
    echo "The entry with the given name does not exist in the set" >&2
    exit 1
  fi
  if [ -n "$FAKE_IPSET_FAIL" ] && [ "$3" = "$FAKE_IPSET_FAIL" ]; then
    echo "The entry with the given name does not exist in the set" >&2
    exit 1
  fi
  exit 0
fi
exit 0
`

func setupFakeIpset(t *testing.T) {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "ipset"), []byte(fakeIpsetScript), 0o755); err != nil {
		t.Fatalf("write fake ipset: %v", err)
	}
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", binDir+string(os.PathListSeparator)+oldPath)
}

func resetOperationProgress() {
	setOperationProgress(nil)
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := NewIPSetController()
	r.POST("/api/v1/ipsets/:setName/entries/import", c.ImportEntries)
	r.DELETE("/api/v1/ipsets/:setName/entries/bulk", c.DeleteEntries)
	r.GET("/api/v1/ipsets/:setName/entries/progress", c.GetOperationProgress)
	return r
}

func newImportRequest(t *testing.T, setName, content string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("file", "entries.txt")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write([]byte(content)); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := mw.WriteField("comment", "test"); err != nil {
		t.Fatalf("write comment: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ipsets/"+setName+"/entries/import", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

func newDeleteRequest(t *testing.T, setName, filter string) *http.Request {
	t.Helper()
	body := `{"filter":` + strconv.Quote(filter) + `}`
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/ipsets/"+setName+"/entries/bulk", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func doRequest(t *testing.T, r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// runRequestAsync runs the request in a goroutine; the response recorder is
// delivered on the returned channel when the request completes.
func runRequestAsync(t *testing.T, r *gin.Engine, req *http.Request) <-chan *httptest.ResponseRecorder {
	t.Helper()
	result := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		result <- rec
	}()
	return result
}

// waitForCompletion polls until the async request finishes, calling poll on
// every iteration.
func waitForCompletion(t *testing.T, result <-chan *httptest.ResponseRecorder, poll func()) *httptest.ResponseRecorder {
	t.Helper()
	var w *httptest.ResponseRecorder
	deadline := time.Now().Add(15 * time.Second)
	for w == nil && time.Now().Before(deadline) {
		poll()
		select {
		case w = <-result:
		case <-time.After(20 * time.Millisecond):
		}
	}
	if w == nil {
		t.Fatal("request did not complete in time")
	}
	return w
}

func TestImportEntriesSuccess(t *testing.T) {
	setupFakeIpset(t)
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newImportRequest(t, "testset", "1.1.1.1\n2.2.2.2\n3.3.3.3\n"))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if body := w.Body.String(); !strings.Contains(body, `"added":3`) || !strings.Contains(body, "Entries imported") {
		t.Fatalf("unexpected body: %q", body)
	}

	progress := getOperationProgress()
	if progress == nil {
		t.Fatal("expected operation progress state after import")
	}
	if progress.Operation != "import" || !progress.Done || progress.Total != 3 ||
		progress.Processed != 3 || progress.Changed != 3 || progress.Error != "" {
		t.Fatalf("unexpected final progress: %+v", progress)
	}
}

func TestImportEntriesProgressWhileRunning(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_DELAY", "0.2")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	result := runRequestAsync(t, r, newImportRequest(t, "testset", "1.1.1.1\n2.2.2.2\n3.3.3.3\n"))

	sawPartialState := false
	sawHTTPProgress := false
	w := waitForCompletion(t, result, func() {
		progress := getOperationProgress()
		if progress != nil && !progress.Done && progress.Processed > 0 && progress.Processed < progress.Total {
			sawPartialState = true
		}
		if progress != nil && !progress.Done && progress.Total == 3 {
			pw := doRequest(t, r, httptest.NewRequest(http.MethodGet, "/api/v1/ipsets/testset/entries/progress", nil))
			if pw.Code == http.StatusOK && strings.Contains(pw.Body.String(), `"total":3`) {
				sawHTTPProgress = true
			}
		}
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"added":3`) {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
	if !sawPartialState {
		t.Fatal("did not observe partial progress state while the import was running")
	}
	if !sawHTTPProgress {
		t.Fatal("did not observe the progress endpoint reporting the running import")
	}
}

func TestImportEntriesPartialFailure(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_FAIL", "2.2.2.2")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newImportRequest(t, "testset", "1.1.1.1\n2.2.2.2\n3.3.3.3\n"))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `"added":2`) || !strings.Contains(body, `"failed":1`) ||
		!strings.Contains(body, "Imported 2 of 3 entries") {
		t.Fatalf("unexpected body: %q", body)
	}

	progress := getOperationProgress()
	if progress == nil || !progress.Done || progress.Changed != 2 {
		t.Fatalf("unexpected final progress: %+v", progress)
	}
}

func TestImportEntriesTotalFailure(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_FAIL_ALL", "1")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newImportRequest(t, "testset", "1.1.1.1\n2.2.2.2\n3.3.3.3\n"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "Import from 'entries.txt' failed") {
		t.Fatalf("unexpected body: %q", body)
	}
	if !strings.Contains(body, "...") || len(body) > 400 {
		t.Fatalf("error message should be truncated, got: %q", body)
	}

	progress := getOperationProgress()
	if progress == nil || !progress.Done || progress.Changed != 0 || progress.Error == "" {
		t.Fatalf("unexpected final progress: %+v", progress)
	}
}

func TestDeleteEntriesSuccess(t *testing.T) {
	setupFakeIpset(t)
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newDeleteRequest(t, "testset", "test"))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if body := w.Body.String(); !strings.Contains(body, `"deleted":3`) || !strings.Contains(body, "Deleted 3 entries") {
		t.Fatalf("unexpected body: %q", body)
	}

	progress := getOperationProgress()
	if progress == nil {
		t.Fatal("expected operation progress state after delete")
	}
	if progress.Operation != "delete" || !progress.Done || progress.Total != 3 ||
		progress.Processed != 3 || progress.Changed != 3 || progress.Error != "" {
		t.Fatalf("unexpected final progress: %+v", progress)
	}
}

func TestDeleteEntriesProgressWhileRunning(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_DELAY", "0.2")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	result := runRequestAsync(t, r, newDeleteRequest(t, "testset", "test"))

	sawPartialState := false
	sawHTTPProgress := false
	w := waitForCompletion(t, result, func() {
		progress := getOperationProgress()
		if progress == nil || progress.Operation != "delete" {
			return
		}
		if !progress.Done && progress.Processed > 0 && progress.Processed < progress.Total {
			sawPartialState = true
		}
		if !progress.Done && progress.Total == 3 {
			pw := doRequest(t, r, httptest.NewRequest(http.MethodGet, "/api/v1/ipsets/testset/entries/progress", nil))
			if pw.Code == http.StatusOK &&
				strings.Contains(pw.Body.String(), `"total":3`) &&
				strings.Contains(pw.Body.String(), `"operation":"delete"`) {
				sawHTTPProgress = true
			}
		}
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"deleted":3`) {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
	if !sawPartialState {
		t.Fatal("did not observe partial progress state while the delete was running")
	}
	if !sawHTTPProgress {
		t.Fatal("did not observe the progress endpoint reporting the running delete")
	}
}

func TestDeleteEntriesPartialFailure(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_FAIL", "2.2.2.2")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newDeleteRequest(t, "testset", "test"))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `"deleted":2`) || !strings.Contains(body, "Deleted 2 of 3 entries") {
		t.Fatalf("unexpected body: %q", body)
	}

	progress := getOperationProgress()
	if progress == nil || !progress.Done || progress.Changed != 2 || progress.Error == "" {
		t.Fatalf("unexpected final progress: %+v", progress)
	}
}

func TestDeleteEntriesTotalFailure(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_FAIL_ALL", "1")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newDeleteRequest(t, "testset", "test"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if body := w.Body.String(); !strings.Contains(body, `entry \"1.1.1.1\"`) || len(body) > 400 {
		t.Fatalf("unexpected error body: %q", body)
	}

	progress := getOperationProgress()
	if progress == nil || !progress.Done || progress.Changed != 0 || progress.Error == "" {
		t.Fatalf("unexpected final progress: %+v", progress)
	}
}

func TestGetOperationProgressNone(t *testing.T) {
	setupFakeIpset(t)
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, httptest.NewRequest(http.MethodGet, "/api/v1/ipsets/testset/entries/progress", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "no operation in progress") {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
}

func TestImportEntriesSetNotFound(t *testing.T) {
	setupFakeIpset(t)
	t.Setenv("FAKE_IPSET_SETS", "empty")
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newImportRequest(t, "testset", "1.1.1.1\n"))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "not found") {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
}

func TestImportEntriesEmptyFile(t *testing.T) {
	setupFakeIpset(t)
	resetOperationProgress()
	config.Config.AppDir = t.TempDir()
	r := newTestRouter()

	w := doRequest(t, r, newImportRequest(t, "testset", "# only a comment\n\n"))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "No entries to import") {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
}
