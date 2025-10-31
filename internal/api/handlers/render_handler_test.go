// Package handlers_test provides tests for render handler
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

func TestRenderHandler_RenderHTML(t *testing.T) {
	handler := handlers.NewRenderHandler()

	tests := []struct {
		name           string
		projectID      string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "Valid HTML render request",
			projectID: "1",
			requestBody: map[string]interface{}{
				"include_css":   true,
				"template_name": "base",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["project_id"] != float64(1) {
					t.Errorf("Expected project_id 1, got %v", resp["project_id"])
				}
				if resp["html_path"] == nil || resp["html_path"] == "" {
					t.Error("Expected html_path in response")
				}
				if resp["size_bytes"] == nil || resp["size_bytes"].(float64) <= 0 {
					t.Error("Expected positive size_bytes in response")
				}
			},
		},
		{
			name:      "Default template name",
			projectID: "2",
			requestBody: map[string]interface{}{
				"include_css": true,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["html_path"] == nil {
					t.Error("Expected html_path in response")
				}
			},
		},
		{
			name:           "Invalid project ID",
			projectID:      "invalid",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/"+tt.projectID+"/render/html", bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.projectID})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.RenderHTML(rr, req)

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

func TestRenderHandler_RenderPDF(t *testing.T) {
	handler := handlers.NewRenderHandler()

	tests := []struct {
		name           string
		projectID      string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:      "Valid PDF render request with pagedjs",
			projectID: "1",
			requestBody: map[string]interface{}{
				"engine":  "pagedjs",
				"format":  "A4",
				"quality": "print",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["project_id"] != float64(1) {
					t.Errorf("Expected project_id 1, got %v", resp["project_id"])
				}
				if resp["pdf_path"] == nil || resp["pdf_path"] == "" {
					t.Error("Expected pdf_path in response")
				}
				if resp["pages"] == nil || resp["pages"].(float64) <= 0 {
					t.Error("Expected positive pages in response")
				}
				metadata := resp["metadata"].(map[string]interface{})
				if metadata["creator"] != "Typecraft v1.0" {
					t.Errorf("Expected creator 'Typecraft v1.0', got %v", metadata["creator"])
				}
			},
		},
		{
			name:      "Default values",
			projectID: "2",
			requestBody: map[string]interface{}{},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["pdf_path"] == nil {
					t.Error("Expected pdf_path in response")
				}
			},
		},
		{
			name:      "Invalid engine",
			projectID: "1",
			requestBody: map[string]interface{}{
				"engine": "invalid_engine",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid project ID",
			projectID:      "invalid",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/"+tt.projectID+"/render/pdf", bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.projectID})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.RenderPDF(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.checkResponse != nil && rr.Code == http.StatusCreated {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestRenderHandler_GetRenderStatus(t *testing.T) {
	handler := handlers.NewRenderHandler()

	tests := []struct {
		name           string
		projectID      string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Valid status request",
			projectID:      "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["project_id"] != float64(1) {
					t.Errorf("Expected project_id 1, got %v", resp["project_id"])
				}
				if resp["html"] == nil {
					t.Error("Expected html status in response")
				}
				if resp["pdf"] == nil {
					t.Error("Expected pdf status in response")
				}
			},
		},
		{
			name:           "Invalid project ID",
			projectID:      "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/projects/"+tt.projectID+"/render/status", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.projectID})

			rr := httptest.NewRecorder()
			handler.GetRenderStatus(rr, req)

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
