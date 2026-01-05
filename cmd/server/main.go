package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AlleyBo55/x402-facilitator-oss/api"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/config"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/facilitator"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/tokens"
)

func main() {
	cfg := config.Load()
	f := facilitator.New(cfg)

	printBanner(cfg)
	api.RegisterRoutes(f, cfg.Network)

	log.Printf("🚀 x402 Facilitator running on http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}

func printBanner(cfg config.Config) {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║            x402 Facilitator OSS v2.0.0                        ║
╠═══════════════════════════════════════════════════════════════╣`)
	fmt.Printf("║  Network:  %-50s║\n", cfg.Network)
	fmt.Printf("║  RPC:      %-50s║\n", truncate(cfg.RPCURL, 50))
	fmt.Printf("║  Auth:     %-50s║\n", enabledStr(cfg.APIKey != ""))
	fmt.Printf("║  Webhook:  %-50s║\n", enabledStr(cfg.WebhookURL != ""))
	fmt.Printf("║  Assets:   %-50s║\n", strings.Join(tokens.GetSymbols(cfg.Network), ", "))
	fmt.Println(`╚═══════════════════════════════════════════════════════════════╝`)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func enabledStr(b bool) string {
	if b {
		return "enabled"
	}
	return "disabled"
}
