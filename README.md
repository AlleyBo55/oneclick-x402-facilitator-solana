# x402 Facilitator (Open Source)

Self-host your own x402 payment facilitator for Solana. One-click deploy to any platform.

[![Deploy on Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/AlleyBo55/x402-facilitator-oss)
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template?template=https://github.com/AlleyBo55/x402-facilitator-oss)
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/AlleyBo55/x402-facilitator-oss)

> **Full Deployment Guide**: [https://solana-x402-paywall.vercel.app/docs/deploy](https://solana-x402-paywall.vercel.app/docs/deploy) - Includes RPC configuration and visual walkthrough.

## Features

- ⚡ **Go** - Single binary, ~10ms cold start, ~5MB Docker image
- 🪙 **Multi-Token** - SOL, USDC, USD1 support
- 📊 **Metrics** - Prometheus `/metrics` + JSON `/stats` API
- 🚦 **Rate Limiting** - Per-IP rate limiting (100 req/min)
- 🔐 **Auth** - Optional API key protection
- 📡 **Webhooks** - Real-time payment notifications

## Project Structure

```
oneclickfacilitator/
├── cmd/server/main.go         # Entry point
├── api/handlers.go            # HTTP handlers
├── internal/
│   ├── config/                # Configuration
│   ├── domain/                # x402 types
│   ├── facilitator/           # Core logic
│   ├── metrics/               # Prometheus metrics
│   ├── middleware/            # HTTP middleware
│   ├── tokens/                # Token registry
│   └── webhook/               # Webhook notifier
├── app.json                   # Heroku config
├── railway.json               # Railway config
├── render.yaml                # Render Blueprint
├── Dockerfile                 # Container build
└── go.mod
```

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SOLANA_RPC_URL` | ✅ | `https://api.devnet.solana.com` | Solana RPC endpoint |
| `SOLANA_NETWORK` | ❌ | `devnet` | `devnet` or `mainnet-beta` |
| `PORT` | ❌ | `3000` | Server port |
| `API_KEY` | ❌ | - | Optional API key for auth |
| `WEBHOOK_URL` | ❌ | - | Webhook URL for notifications |

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/stats` | GET | JSON statistics (for custom dashboards) |
| `/metrics` | GET | Prometheus metrics |
| `/supported` | POST | Get supported payment schemes |
| `/verify` | POST | Verify a payment |
| `/settle` | POST | Settle a payment |

### Example: Verify Payment

```bash
curl -X POST https://your-facilitator.com/verify \
  -H "Content-Type: application/json" \
  -d '{
    "paymentPayload": {
      "x402Version": 2,
      "payload": { "signature": "5wHu..." }
    },
    "paymentRequirements": {
      "payTo": "RecipientWallet",
      "amount": "1000000",
      "asset": "SOL"
    }
  }'
```

## Development

```bash
# Build
go build -o facilitator ./cmd/server

# Run
SOLANA_RPC_URL=https://api.devnet.solana.com ./facilitator
```

## License

MIT © [AlleyBoss](https://github.com/AlleyBo55)
