package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, dir string) {
	t.Helper()

	content := []byte("[account]\nusername = tester\npasswd = secret\n")
	if err := os.WriteFile(filepath.Join(dir, "ImageCloud.conf"), content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldDir)
	})
}

func authenticatedRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	loginRequest := httptest.NewRequest(http.MethodGet, "http://example.com"+target, nil)
	recorder := httptest.NewRecorder()
	setSession("tester", recorder, loginRequest)

	response := recorder.Result()
	if len(response.Cookies()) == 0 {
		t.Fatal("expected session cookie")
	}

	request := httptest.NewRequest(http.MethodGet, "http://example.com"+target, nil)
	request.AddCookie(response.Cookies()[0])
	return request
}

func TestProtectedFileHandlerRequiresLogin(t *testing.T) {
	tempDir := t.TempDir()
	handler := protectedFileHandler("/images/", tempDir, isImage)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://example.com/images/test.jpg", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", recorder.Code)
	}
	if location := recorder.Header().Get("Location"); location != "/main" {
		t.Fatalf("expected redirect to /main, got %q", location)
	}
}

func TestProtectedFileHandlerBlocksDirectoriesAndNonImages(t *testing.T) {
	tempDir := t.TempDir()
	writeTestConfig(t, tempDir)
	withWorkingDir(t, tempDir)

	imageFile := filepath.Join(tempDir, "photo.jpg")
	if err := os.WriteFile(imageFile, []byte("jpg"), 0o600); err != nil {
		t.Fatalf("write image: %v", err)
	}
	textFile := filepath.Join(tempDir, "secret.txt")
	if err := os.WriteFile(textFile, []byte("txt"), 0o600); err != nil {
		t.Fatalf("write text: %v", err)
	}

	handler := protectedFileHandler("/images/", tempDir, isImage)

	tests := []struct {
		name   string
		target string
		want   int
	}{
		{name: "directory listing", target: "/images/", want: http.StatusNotFound},
		{name: "non image file", target: "/images/secret.txt", want: http.StatusNotFound},
		{name: "image file", target: "/images/photo.jpg", want: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := authenticatedRequest(t, tc.target)
			handler.ServeHTTP(recorder, request)

			if recorder.Code != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, recorder.Code)
			}
		})
	}
}

func TestSafeJoinUnderBaseBlocksTraversal(t *testing.T) {
	got, err := safeJoinUnderBase("images", "../secret")
	if err == nil {
		t.Fatalf("expected traversal error, got path %q", got)
	}
}

func TestThumbnailPathUsesImageRelativePath(t *testing.T) {
	got, err := thumbnailPath(filepath.Join(imgPath, "album", "photo.png"))
	if err != nil {
		t.Fatalf("thumbnailPath error: %v", err)
	}

	want := filepath.Join(thumPath, "album", "photo.png") + ".jpg"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
