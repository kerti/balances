package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/config"
	"github.com/kerti/balances/backend/handler"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/logger"
)

// Server is the server instance
type Server struct {
	config             *config.Config
	AuthHandler        handler.Auth        `inject:"authHandler"`
	AuthService        service.Auth        `inject:"authService"`
	BankAccountHandler handler.BankAccount `inject:"bankAccountHandler"`
	HealthHandler      handler.Health      `inject:"healthHandler"`
	UserHandler        handler.User        `inject:"userHandler"`
	VehicleHandler     handler.Vehicle     `inject:"vehicleHandler"`
	router             *mux.Router
}

// Startup perform startup functions
func (s *Server) Startup() {
	logger.Trace("HTTP Server starting up...")
	s.config = config.Get()
	s.InitRoutes()
	s.InitMiddleware()
}

// Shutdown cleans up everything and shuts down
func (s *Server) Shutdown() {
	logger.Trace("HTTP Server shutting down...")
}

// InitMiddleware initializes all middlewares
func (s *Server) InitMiddleware() {
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.jwtMiddleware)
}

// Serve runs the server
func (s *Server) Serve() {
	logger.Info("Server started and is available at the following address(es):")
	ips, _ := s.getIPs()
	for _, ip := range ips {
		logger.Info("- http://%s:%d", ip.String(), s.config.Server.Port)
	}

	logger.Fatal("%s", http.ListenAndServe(fmt.Sprintf(":%d", s.config.Server.Port), nil))

}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerList := make([]string, 0)
		for k, v := range r.Header {
			headerList = append(headerList, fmt.Sprintf("  - %s: %s", k, v))
		}
		headers := fmt.Sprintf("- HEADERS:\n%s", strings.Join(headerList, "\n"))

		cookieList := make([]string, 0)
		for _, v := range r.Cookies() {
			cookieList = append(cookieList, fmt.Sprintf(" - %s: %v", v.Name, v.Value))
		}
		cookies := fmt.Sprintf("- COOKIES:\n%s", strings.Join(cookieList, "\n"))

		logger.Trace("### RECEIEVED %v %v\n%s\n%s", r.Method, r.RequestURI, headers, cookies)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// no JWT checks for preflights, health checks and logins
		if r.Method == http.MethodOptions || r.RequestURI == "/health" || r.RequestURI == "/auth/login" {
			next.ServeHTTP(w, r)
			return
		}

		token := ""
		tokenCookie, _ := r.Cookie(s.config.JWT.TokenCookie)
		if tokenCookie != nil {
			token = fmt.Sprintf("Bearer %s", tokenCookie.Value)
		}

		if len(token) == 0 {
			token = r.Header.Get("Authorization")
		}

		userID, err := s.AuthService.Authorize(token)
		if err != nil {
			logger.Trace(fmt.Sprintf("authorization failed: %v", err.Error()))
			response.RespondWithError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxprops.PropUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header
		origins := headers["Origin"]

		if len(origins) > 0 {
			origin := origins[0]
			for _, allowedOrigin := range s.config.CORS.AllowedOrigins {
				if allowedOrigin == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) getIPs() ([]net.IP, error) {
	res := make([]net.IP, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						res = append(res, ipNet.IP)
					}
				}
			}
		}
	}
	return res, nil
}
