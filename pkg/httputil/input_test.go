package httputil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetInput_FromQuery(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/test?name=John&empty=", nil)
	c.Request = req

	val, ok := GetInput(c, "name")
	if !ok || val != "John" {
		t.Errorf("expected 'John', true; got %s, %t", val, ok)
	}

	_, ok = GetInput(c, "empty")
	if ok {
		t.Errorf("expected false for empty query param")
	}

	_, ok = GetInput(c, "missing")
	if ok {
		t.Errorf("expected false for missing query param")
	}
}

func TestGetInput_FromPostForm(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	formData := url.Values{
		"username": {"auth"},
	}
	req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	val, ok := GetInput(c, "username")
	if !ok || val != "auth" {
		t.Errorf("expected 'auth', true; got %s, %t", val, ok)
	}
}

func TestGetInput_FromJSON_MultipleReads(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	jsonBody := `{"password":"secret123","age":25,"enabled":true}`
	req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Pārbaudām, ka var nolasīt pirmo lauku
	pass, ok := GetInput(c, "password")
	if !ok || pass != "secret123" {
		t.Errorf("expected 'secret123', true; got %s, %t", pass, ok)
	}

	// Pārbaudām, ka var nolasīt otro lauku bez EOF kļūdas (ShouldBindBodyWith kešošana)
	age, ok := GetInput(c, "age")
	if !ok || age != "25" {
		t.Errorf("expected '25', true; got %s, %t", age, ok)
	}

	enabled, ok := GetInput(c, "enabled")
	if !ok || enabled != "true" {
		t.Errorf("expected 'true', true; got %s, %t", enabled, ok)
	}

	// Neesošs lauks
	_, ok = GetInput(c, "missing")
	if ok {
		t.Errorf("expected false for missing json field")
	}
}

func TestGetInput_Priority(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Pieprasījumam ir gan query, gan form
	formData := url.Values{"field": {"from_form"}}
	req, _ := http.NewRequest(http.MethodPost, "/test?field=from_query", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	val, ok := GetInput(c, "field")
	if !ok || val != "from_query" {
		t.Errorf("expected query priority 'from_query', got %s", val)
	}
}

func TestGetInputDefault(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/test?exists=yes", nil)
	c.Request = req

	if got := GetInputDefault(c, "exists", "fallback"); got != "yes" {
		t.Errorf("expected 'yes', got %s", got)
	}
	if got := GetInputDefault(c, "missing", "fallback"); got != "fallback" {
		t.Errorf("expected 'fallback', got %s", got)
	}
}
