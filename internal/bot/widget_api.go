package bot

import (
	"encoding/json"
	"net/http"
	"os"
	"github.com/etkecc/postmoogle/internal/config"
)

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (b *Bot) StartWidgetAPI(cfg config.WidgetAPI) {
	if !cfg.Enabled {
		return
	}

	mux := http.NewServeMux()

	// 1. Define your routes normally (no wrappers here)
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/v1/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Widget-Secret") != cfg.Secret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req SendEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b.log.Info().Str("to", req.To).Msg("Widget requested email send")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
	})

	// 2. Create a GLOBAL middleware handler
	// This ensures CORS headers are ALWAYS sent, no matter what
	globalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Widget-Secret")

		// Handle Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to your mux
		mux.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" { port = cfg.Port }

	b.log.Info().Str("port", port).Msg("Starting Widget REST API with Global CORS")
	
    // 3. IMPORTANT: Pass the globalHandler here, not the mux
	go http.ListenAndServe(":"+port, globalHandler)
}