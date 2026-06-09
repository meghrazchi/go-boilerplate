package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/your-org/go-ddd-boilerplate/internal/config"
	healthhttp "github.com/your-org/go-ddd-boilerplate/internal/modules/health/interfaces/http"
	userhttp "github.com/your-org/go-ddd-boilerplate/internal/modules/user/interfaces/http"
	platformmiddleware "github.com/your-org/go-ddd-boilerplate/internal/platform/middleware"
	"github.com/your-org/go-ddd-boilerplate/internal/platform/response"
)

type Dependencies struct {
	Config        *config.Config
	Logger        *zap.Logger
	UserHandler   *userhttp.Handler
	HealthHandler *healthhttp.Handler
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(platformmiddleware.RequestID)
	r.Use(platformmiddleware.Security)
	r.Use(platformmiddleware.CORS(deps.Config))
	r.Use(platformmiddleware.BodyLimit(deps.Config.MaxBodyBytes))
	r.Use(platformmiddleware.Timeout(deps.Config.RequestTimeoutDuration()))
	r.Use(platformmiddleware.Logger(deps.Logger))
	r.Use(platformmiddleware.Recovery(deps.Logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Welcome to go-ddd-boilerplate", map[string]string{"docs": "/docs/swagger"})
	})

	r.Get("/docs/openapi.yaml", serveOpenAPI)
	r.Get("/docs/swagger", serveSwaggerUI)

	healthhttp.RegisterRoutes(r, deps.HealthHandler)

	r.Route("/api/v1", func(r chi.Router) {
		userhttp.RegisterRoutes(r, deps.UserHandler)
	})

	return r
}

func serveOpenAPI(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		response.Error(w, http.StatusNotFound, "OpenAPI spec not found", nil)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerHTML))
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({ url: '/docs/openapi.yaml', dom_id: '#swagger-ui' });
    };
  </script>
</body>
</html>`
