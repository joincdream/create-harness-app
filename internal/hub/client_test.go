package hub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_ListTemplates(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/templates" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"name":"test-template","version":"1.0.0","target_agent":"antigravity","description":"Test"}]`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	list, err := client.ListTemplates(context.Background(), "", "")
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].Name != "test-template" {
		t.Errorf("expected name 'test-template', got '%s'", list[0].Name)
	}
}
