// ppc-contact: Kontaktformular-Bridge für PPCWeb → M-Con
// Lauscht nur auf Loopback; SeDat proxied /api/* hierher.
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	altcha "github.com/altcha-org/altcha-lib-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen             string   `yaml:"listen"`
	AltchaSecret       string   `yaml:"altcha_secret"`
	AltchaMaxNumber    int      `yaml:"altcha_max_number"`
	AltchaExpiresSecs  int      `yaml:"altcha_expires_seconds"`
	MCon               MConCfg  `yaml:"mcon"`
	RateLimit          RateCfg  `yaml:"rate_limit"`
	FormMinSeconds     int      `yaml:"form_min_seconds"`
	FormMaxSeconds     int      `yaml:"form_max_seconds"`
	AllowedOrigins     []string `yaml:"allowed_origins"`
}

type MConCfg struct {
	BaseURL         string   `yaml:"base_url"`
	KeyID           string   `yaml:"key_id"`
	HMACSecret      string   `yaml:"hmac_secret"`
	From            string   `yaml:"from"`
	To              []string `yaml:"to"`
	Provider        string   `yaml:"provider"`
	TimeoutSeconds  int      `yaml:"timeout_seconds"`
}

type RateCfg struct {
	RPS                float64 `yaml:"rps"`
	Burst              int     `yaml:"burst"`
	PerIPWindowSeconds int     `yaml:"per_ip_window_seconds"`
	PerIPMax           int     `yaml:"per_ip_max"`
}

type contactRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
	Privacy   bool   `json:"privacy"`
	Website   string `json:"website"` // honeypot
	Altcha    string `json:"altcha"`
	IssuedAt  int64  `json:"issued_at"`
}

type ipBucket struct {
	times []time.Time
}

func main() {
	cfgPath := flag.String("config", "", "Pfad zu config.yaml")
	flag.Parse()
	if *cfgPath == "" {
		*cfgPath = os.Getenv("PPC_CONTACT_CONFIG")
	}
	if *cfgPath == "" {
		log.Fatal("config required: -config or PPC_CONTACT_CONFIG")
	}
	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	applyEnvOverrides(&cfg)
	if err := validateConfig(cfg); err != nil {
		log.Fatal(err)
	}

	limiter := newIPLimiter(cfg.RateLimit)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/api/altcha/challenge", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !originOK(cfg, r) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		ch, err := altcha.CreateChallenge(altcha.ChallengeOptions{
			HMACKey:   cfg.AltchaSecret,
			MaxNumber: int64(cfg.AltchaMaxNumber),
		})
		if err != nil {
			http.Error(w, "challenge failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(ch)
	})
	mux.HandleFunc("/api/contact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !originOK(cfg, r) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
			return
		}
		ip := clientIP(r)
		if !limiter.allow(ip) {
			writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "rate_limited"})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		var req contactRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
			return
		}
		if strings.TrimSpace(req.Website) != "" {
			// Honeypot: stiller Erfolg
			writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
			return
		}
		if !req.Privacy {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "privacy_required"})
			return
		}
		if err := validateFields(req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		now := time.Now().Unix()
		if req.IssuedAt <= 0 || now-req.IssuedAt < int64(cfg.FormMinSeconds) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "too_fast"})
			return
		}
		if now-req.IssuedAt > int64(cfg.FormMaxSeconds) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "form_expired"})
			return
		}
		ok, err := altcha.VerifySolution(req.Altcha, cfg.AltchaSecret, true)
		if err != nil || !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "captcha_failed"})
			return
		}

		subject := fmt.Sprintf("[PPCWeb] %s – %s", safeTopic(req.Topic), truncate(req.Name, 80))
		text := formatMail(req, ip)
		html := formatMailHTML(req)
		msgID, err := sendMCon(cfg, subject, text, html, req.Email)
		if err != nil {
			log.Printf("mcon send failed from %s: %v", ip, err)
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "send_failed"})
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted", "id": msgID})
	})

	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("ppc-contact listening on %s", cfg.Listen)
	log.Fatal(srv.ListenAndServe())
}

func loadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	if c.Listen == "" {
		c.Listen = "127.0.0.1:8690"
	}
	if c.AltchaMaxNumber == 0 {
		c.AltchaMaxNumber = 100000
	}
	if c.AltchaExpiresSecs == 0 {
		c.AltchaExpiresSecs = 600
	}
	if c.FormMinSeconds == 0 {
		c.FormMinSeconds = 3
	}
	if c.FormMaxSeconds == 0 {
		c.FormMaxSeconds = 3600
	}
	if c.MCon.TimeoutSeconds == 0 {
		c.MCon.TimeoutSeconds = 15
	}
	if c.RateLimit.PerIPWindowSeconds == 0 {
		c.RateLimit.PerIPWindowSeconds = 60
	}
	if c.RateLimit.PerIPMax == 0 {
		c.RateLimit.PerIPMax = 5
	}
	return c, nil
}

func applyEnvOverrides(c *Config) {
	if v := os.Getenv("PPC_CONTACT_ALTCHA_SECRET"); v != "" {
		c.AltchaSecret = v
	}
	if v := os.Getenv("PPC_CONTACT_MCON_HMAC"); v != "" {
		c.MCon.HMACSecret = v
	}
	if v := os.Getenv("PPC_CONTACT_MCON_URL"); v != "" {
		c.MCon.BaseURL = v
	}
}

func validateConfig(c Config) error {
	if len(c.AltchaSecret) < 16 {
		return fmt.Errorf("altcha_secret must be at least 16 characters")
	}
	if len(c.MCon.HMACSecret) < 16 {
		return fmt.Errorf("mcon.hmac_secret must be at least 16 characters")
	}
	if c.MCon.BaseURL == "" || c.MCon.KeyID == "" || c.MCon.From == "" || len(c.MCon.To) == 0 {
		return fmt.Errorf("mcon base_url, key_id, from and to are required")
	}
	host, _, err := net.SplitHostPort(c.Listen)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("listen must be loopback (got %s)", c.Listen)
	}
	return nil
}

func originOK(cfg Config, r *http.Request) bool {
	if len(cfg.AllowedOrigins) == 0 {
		return true
	}
	o := r.Header.Get("Origin")
	if o == "" {
		// same-origin navigations / some browsers omit Origin on GET
		ref := r.Header.Get("Referer")
		for _, a := range cfg.AllowedOrigins {
			if strings.HasPrefix(ref, a) {
				return true
			}
		}
		return r.Method == http.MethodGet
	}
	for _, a := range cfg.AllowedOrigins {
		if o == a {
			return true
		}
	}
	return false
}

func validateFields(req contactRequest) error {
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)
	msg := strings.TrimSpace(req.Message)
	if name == "" || len(name) > 120 {
		return fmt.Errorf("invalid_name")
	}
	if _, err := mail.ParseAddress(email); err != nil || len(email) > 200 {
		return fmt.Errorf("invalid_email")
	}
	if msg == "" || len(msg) > 5000 {
		return fmt.Errorf("invalid_message")
	}
	if len(req.Phone) > 40 {
		return fmt.Errorf("invalid_phone")
	}
	if strings.ContainsAny(name+req.Phone+req.Topic, "\r\n") {
		return fmt.Errorf("invalid_chars")
	}
	if req.Altcha == "" {
		return fmt.Errorf("captcha_required")
	}
	return nil
}

func safeTopic(t string) string {
	t = strings.TrimSpace(t)
	switch t {
	case "security", "administration", "rueckruf", "sonstiges":
		return t
	default:
		return "sonstiges"
	}
}

func formatMail(req contactRequest, ip string) string {
	return fmt.Sprintf(
		"Neue Kontaktanfrage von der Website ppc.at\n\nName: %s\nE-Mail: %s\nTelefon: %s\nThema: %s\nIP: %s\nZeit: %s\n\nNachricht:\n%s\n",
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.Email),
		strings.TrimSpace(req.Phone),
		safeTopic(req.Topic),
		ip,
		time.Now().Format(time.RFC3339),
		strings.TrimSpace(req.Message),
	)
}

func formatMailHTML(req contactRequest) string {
	esc := func(s string) string {
		r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
		return r.Replace(s)
	}
	return fmt.Sprintf(
		"<p><strong>Neue Kontaktanfrage</strong> (ppc.at)</p><ul><li>Name: %s</li><li>E-Mail: %s</li><li>Telefon: %s</li><li>Thema: %s</li></ul><p>%s</p>",
		esc(strings.TrimSpace(req.Name)),
		esc(strings.TrimSpace(req.Email)),
		esc(strings.TrimSpace(req.Phone)),
		esc(safeTopic(req.Topic)),
		esc(strings.TrimSpace(req.Message)),
	)
}

func sendMCon(cfg Config, subject, text, html, replyTo string) (string, error) {
	bodyObj := map[string]any{
		"from":     cfg.MCon.From,
		"to":       cfg.MCon.To,
		"reply_to": replyTo,
		"subject":  subject,
		"text":     text,
		"html":     html,
		"tags":     []string{"ppcweb", "contact"},
		"metadata": map[string]string{"form": "kontakt", "site": "ppc.at"},
		"spam":     map[string]string{"hint": "website contact form"},
		"provider": cfg.MCon.Provider,
		"headers":  map[string]string{"X-PPC-Form": "kontakt"},
	}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return "", err
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := randomNonce(16)
	path := "/v1/messages"
	sum := sha256.Sum256(body)
	canon := strings.Join([]string{"POST", path, ts, nonce, hex.EncodeToString(sum[:])}, "\n")
	mac := hmac.New(sha256.New, []byte(cfg.MCon.HMACSecret))
	_, _ = mac.Write([]byte(canon))
	sig := hex.EncodeToString(mac.Sum(nil))

	idem := "ppc-" + hex.EncodeToString(sum[:16]) + "-" + nonce[:8]
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(cfg.MCon.BaseURL, "/")+path, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-MCon-Key", cfg.MCon.KeyID)
	req.Header.Set("X-MCon-Timestamp", ts)
	req.Header.Set("X-MCon-Nonce", nonce)
	req.Header.Set("X-MCon-Signature", sig)
	req.Header.Set("Idempotency-Key", idem)

	client := &http.Client{Timeout: time.Duration(cfg.MCon.TimeoutSeconds) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mcon status %d: %s", res.StatusCode, truncate(string(respBody), 200))
	}
	var parsed struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(respBody, &parsed)
	return parsed.ID, nil
}

func randomNonce(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func clientIP(r *http.Request) string {
	// Hinter SeDat: echte Client-IP nur aus vertrauenswürdigem PROXY-Protokoll / Header,
	// die SeDat setzt. Fallback RemoteAddr.
	if xff := r.Header.Get("X-Real-IP"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func newIPLimiter(cfg RateCfg) *ipLimiter {
	return &ipLimiter{
		window: time.Duration(cfg.PerIPWindowSeconds) * time.Second,
		max:    cfg.PerIPMax,
		m:      map[string]*ipBucket{},
	}
}

type ipLimiter struct {
	mu     sync.Mutex
	window time.Duration
	max    int
	m      map[string]*ipBucket
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	b := l.m[ip]
	if b == nil {
		b = &ipBucket{}
		l.m[ip] = b
	}
	cut := now.Add(-l.window)
	kept := b.times[:0]
	for _, t := range b.times {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	b.times = kept
	if len(b.times) >= l.max {
		return false
	}
	b.times = append(b.times, now)
	return true
}
