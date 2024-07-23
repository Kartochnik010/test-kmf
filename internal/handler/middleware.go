package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kartochnik010/test-kmf/internal/pkg/js"
	"github.com/kartochnik010/test-kmf/internal/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

// needed to capture the status code of the response
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
func (h *Handler) WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		capture := NewLoggingResponseWriter(w)
		next.ServeHTTP(capture, r)

		statusCode := capture.statusCode
		log := logger.GetLoggerFromCtx(r.Context())
		log.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"ip":          realip.FromRequest(r),
			"status_code": statusCode,
			"duration":    time.Since(t).String(),
		}).Info()
	})
}

func (h *Handler) AssignLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logger with request id
		r = r.WithContext(context.WithValue(r.Context(), logger.ContextKeyLogger, h.l.WithField("requestID", uuid.New())))
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for reach client
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold pointers to a client struct.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once every
	// minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup
			// is taking place.
			mu.Lock()

			// Loop through all clients. if they haven't been seen within the last three minutes,
			// then delete the corresponding entry from the clients map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			// Importantly, unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limited is enabled.
		if h.cfg.Limiter.Enabled {
			// Use the realip.FromRequest function to get the client's real IP address.
			ip := realip.FromRequest(r)

			// Lock the mutex to prevent this code from being executed concurrently.
			mu.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't,
			// then initialize a new rate limiter and add the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				// Use the requests-per-second and burst values from the app.config struct.
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(h.cfg.Limiter.Rps), h.cfg.Limiter.Burst)}
			}

			// Update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// Call the limiter.Allow() method on the rate limiter for the current IP address.
			// If the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
			// response.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				js.WriteJSON(w, http.StatusTooManyRequests, js.JSON{"error": "rate limited exceeded"}, nil)
				return
			}

			// Very importantly, unlock the mutex before calling the next handler in the chain.
			// Notice that we DON'T use defer to unlock the mutex, as that would mean that the mutex
			// isn't unlocked until all handlers downstream of this middleware have also returned.
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}
