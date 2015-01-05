package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandlerRoutesGetRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("GET / expected HTTP %v, got %v", http.StatusOK, response.Code)
	}
}

func TestRootHandlerRoutesPostRequest(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("-v"))
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("POST / expected HTTP %v, got %v", http.StatusOK, response.Code)
	}
}

func TestRootHandlerCatchall(t *testing.T) {
	request, _ := http.NewRequest("GET", "/nothing", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("GET /nothing expected HTTP %v, got %v", http.StatusNotFound, response.Code)
	}
}

func TestRootHandler_actionRoot_success(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if match := "Alive!"; !strings.Contains(response.Body.String(), match) {
		t.Fatalf("actionRoot body should match %v, got %v", match, response.Body)
	}
}

func TestRootHandler_actionDig_success(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("-v"))
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if match := "DiG"; !strings.Contains(response.Body.String(), match) {
		t.Fatalf("actionDig body should match %v, got %v", match, response.Body)
	}
}
