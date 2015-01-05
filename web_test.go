package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDig(t *testing.T) {
	out, err := Dig("-v")

	if err != nil {
		t.Fatalf("Dig failed: %v", err)
	}

	if match := "DiG"; !strings.Contains(string(out), match) {
		t.Fatalf("actionDig body should match %v, got %s", match, out)
	}
}

func TestDig_usage(t *testing.T) {
	out, _ := Dig("-h")

	if match := "Usage:"; !strings.Contains(string(out), match) {
		t.Fatalf("actionDig body should match %v, got %s", match, out)
	}
}

func TestRootHandlerRoutesGetRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if expected := http.StatusOK; response.Code != expected {
		t.Fatalf("GET / expected HTTP %v, got %v", expected, response.Code)
	}
}

func TestRootHandlerRoutesPostRequest(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("-v"))
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if expected := http.StatusOK; response.Code != expected {
		t.Fatalf("POST / expected HTTP %v, got %v", http.StatusOK, response.Code)
	}
}

func TestRootHandlerHaltsOtherRequest(t *testing.T) {
	request, _ := http.NewRequest("PUT", "/", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if expected := http.StatusNotFound; response.Code != expected {
		t.Fatalf("PUT / expected HTTP %v, got %v", expected, response.Code)
	}
}

func TestSlackHandlerHaltsGetRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/slack", nil)
	response := httptest.NewRecorder()

	SlackHandler(response, request)

	if expected := http.StatusNotFound; response.Code != expected {
		t.Fatalf("GET /slack expected HTTP %v, got %v", expected, response.Code)
	}
}

func TestSlackHandlerRoutesPostRequest(t *testing.T) {
	request, _ := http.NewRequest("POST", "/slack", strings.NewReader("text=-v"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	SlackHandler(response, request)

	if expected := http.StatusOK; response.Code != expected {
		t.Fatalf("POST /slack expected HTTP %v, got %v", http.StatusOK, response.Code)
	}
}

func TestRootHandlerCatchall(t *testing.T) {
	request, _ := http.NewRequest("GET", "/nothing", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if expected := http.StatusNotFound; response.Code != expected {
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

func TestSlackHandler_actionSlack_success(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("token=TOKEN0000000000000000000&team_id=TEAM00000&team_domain=DOMAIN&channel_id=CHANNEL&channel_name=GENERAL&user_id=USER&user_name=WEPPOS&command=%2Fdig&text=-v"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	SlackHandler(response, request)

	if match := "DiG"; !strings.Contains(response.Body.String(), match) {
		t.Fatalf("actionDig body should match %v, got %v", match, response.Body)
	}
}
