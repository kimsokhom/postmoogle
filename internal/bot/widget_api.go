package bot

import (
	"encoding/json"
	"net/http"
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

	// Endpoint: Check if API is alive
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Endpoint: Send Email (The core widget functionality)
	mux.HandleFunc("/api/v1/send", func(w http.ResponseWriter, r *http.Request) {
		// 1. Check Security Secret
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

		// 2. Logic to trigger Postmoogle's internal send
		// For now, we'll just log it to prove the connection works
		b.log.Info().Str("to", req.To).Msg("Widget requested email send")
		
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
	})

	b.log.Info().Str("port", cfg.Port).Msg("Starting Widget REST API")
	go http.ListenAndServe(":"+cfg.Port, mux)
}