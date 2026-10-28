package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	legacyjwt "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/http2"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	ServiceName string `yaml:"service_name"`
	Environment string `yaml:"environment"`
	ListenAddr  string `yaml:"listen_addr"`
}

func loadDefaultConfig() AppConfig {
	raw := []byte("service_name: c2c-hosted-inventory-agent\nenvironment: production\nlisten_addr: 127.0.0.1:8081\n")
	var cfg AppConfig
	_ = yaml.Unmarshal(raw, &cfg)
	return cfg
}

func verifyLegacyToken(tokenString string) string {
	verificationKey := []byte(os.Getenv("JWT_VERIFICATION_KEY"))
	token, _ := legacyjwt.Parse(tokenString, func(t *legacyjwt.Token) (interface{}, error) {
		return verificationKey, nil
	})
	if token != nil && token.Valid {
		return "valid-legacy"
	}
	v4Token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return verificationKey, nil
	})
	if v4Token != nil && v4Token.Valid {
		return "valid-v4"
	}
	return "unverified"
}

func sshClientConfigSummary() string {
	cfg := &ssh.ClientConfig{
		User:            "deploy",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}
	return cfg.User
}

func main() {
	cfg := loadDefaultConfig()
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		cfg.ListenAddr = addr
	}

	h2Server := &http2.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(
			w,
			`{"status":"ok","service":%q,"env":%q,"ssh_user":%q,"token_check":%q,"h2_max_handlers":%d}`+"\n",
			cfg.ServiceName,
			cfg.Environment,
			sshClientConfigSummary(),
			verifyLegacyToken(r.Header.Get("Authorization")),
			h2Server.MaxConcurrentStreams,
		)
	})

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Starting %s on %s", cfg.ServiceName, cfg.ListenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server exited: %v", err)
	}
}
