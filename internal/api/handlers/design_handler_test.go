// Package handlers_test provides tests for design handler
// Conformidade: Constituição Vértice v3.0 - P2 (Validação Preventiva)
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JuanCS-Dev/typecraft/internal/api/handlers"
	"github.com/gorilla/mux"
)

func TestDesignHandler_GenerateDesign(t *testing.T) {
	handler := handlers.NewDesignHandler()

	tests := []struct {
		name           string
		projectID      string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "Valid design generation request",
			projectID: "1",
			requestBody: map[string]interface{}{
				"genre":    "fiction",
				"keywords": []string{"mystery", "thriller"},
				"tone":     "professional",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["project_id"] != float64(1) {
					t.Errorf("Expected project_id 1, got %v", resp["project_id"])
				}
				if resp["color_palette"] == nil {
					t.Error("Expected color_palette in response")
				}
				if resp["font_pairing"] == nil {
					t.Error("Expected font_pairing in response")
				}
			},
		},
		{
			name:      "Missing genre",
			projectID: "1",
			requestBody: map[string]interface{}{
				"tone": "professional",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid project ID",
			projectID:      "invalid",
			requestBody:    map[string]interface{}{"genre": "fiction"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON body",
			projectID:      "1",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			} else {
				body = []byte("invalid json")
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/"+tt.projectID+"/design/generate", bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.projectID})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.GenerateDesign(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.checkResponse != nil && rr.Code == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestDesignHandler_ListFonts(t *testing.T) {
	handler := handlers.NewDesignHandler()

	tests := []struct {
		name           string
		category       string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "List all fonts",
			category:       "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				fonts, ok := resp["fonts"].([]interface{})
				if !ok {
					t.Error("Expected fonts array in response")
					return
				}
				if len(fonts) == 0 {
					t.Error("Expected at least one font")
				}
				count := resp["count"].(float64)
				if int(count) != len(fonts) {
					t.Errorf("Count mismatch: count=%d, len(fonts)=%d", int(count), len(fonts))
				}
			},
		},
		{
			name:           "Filter by serif category",
			category:       "serif",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				fonts, ok := resp["fonts"].([]interface{})
				if !ok {
					t.Error("Expected fonts array in response")
					return
				}
				if len(fonts) == 0 {
					t.Error("Expected at least one serif font")
				}
				for _, f := range fonts {
					font := f.(map[string]interface{})
					if font["category"] != "serif" {
						t.Errorf("Expected serif font, got %s", font["category"])
					}
				}
			},
		},
		{
			name:           "Filter by sans-serif category",
			category:       "sans-serif",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				fonts, ok := resp["fonts"].([]interface{})
				if !ok {
					t.Error("Expected fonts array in response")
					return
				}
				if len(fonts) == 0 {
					t.Error("Expected at least one sans-serif font")
				}
			},
		},
		{
			name:           "Filter by monospace category",
			category:       "monospace",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				fonts, ok := resp["fonts"].([]interface{})
				if !ok {
					t.Error("Expected fonts array in response")
					return
				}
				if len(fonts) == 0 {
					t.Error("Expected at least one monospace font")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/fonts"
			if tt.category != "" {
				url += "?category=" + tt.category
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()
			handler.ListFonts(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.checkResponse != nil && rr.Code == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}
