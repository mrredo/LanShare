package example

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"lanshare/internal/db"
)

func setupTestDB(t *testing.T) *db.Service {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_example.db")

	service := db.NewService(db.Config{
		DBPath:        dbPath,
		MigrationsDir: "../../migrations",
	})
	if err := service.Start(); err != nil {
		t.Fatalf("failed to start test db: %v", err)
	}
	t.Cleanup(func() {
		_ = service.Close()
	})

	// Automigrācija piemēra tabulai testā
	if err := service.DB().AutoMigrate(&ExampleItem{}); err != nil {
		t.Fatalf("failed to migrate example table: %v", err)
	}

	return service
}

func TestExampleHandler_RegistrationAndFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := setupTestDB(t)

	// 1. Izveidojam Repozitoriju -> Servisu -> Handleru
	repo := NewRepo(database.DB())
	service := NewService(repo)
	handler := NewHandler(service)

	// 2. Izveidojam Gin maršrutētāju un reģistrējam maršrutus
	router := gin.New()
	apiGroup := router.Group("/api")
	handler.RegisterRoutes(apiGroup)

	// 3. Testējam POST /api/examples
	body, _ := json.Marshal(CreateExampleDTO{
		Title:   "Pirmais ieraksts",
		Content: "Piemēra saturs",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/examples", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	// 4. Testējam GET /api/examples
	getReq, _ := http.NewRequest(http.MethodGet, "/api/examples", nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", getW.Code, getW.Body.String())
	}

	var res struct {
		Data []ExampleItem `json:"data"`
	}
	if err := json.Unmarshal(getW.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(res.Data) != 1 || res.Data[0].Title != "Pirmais ieraksts" {
		t.Errorf("expected 1 item with title 'Pirmais ieraksts', got %+v", res.Data)
	}
}
