# x402 Facilitator (Open Source)

> **Why pay fees?** Deploy your own x402 payment facilitator in **one click** and cut out the middleman.
> 
> ✅ **Self-Sovereign**: You own the infrastructure.
> ✅ **No Platform Fees**: You only pay for your RPC and cloud hosting.
> ✅ **Fully Configurable**: Strict control over rate limits, tokens, and networks.
> ✅ **Open Source**: Verify the code, audit the security, and modify it to your needs.

[![Deploy on Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/AlleyBo55/oneclick-x402-facilitator-solana)
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.com/new?template=https://github.com/AlleyBo55/oneclick-x402-facilitator-solana)
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/AlleyBo55/oneclick-x402-facilitator-solana)

> **Quick Start**: Pair this server with our client SDK:
> [**@alleyboss/micropay-solana-x402-paywall**](https://www.npmjs.com/package/@alleyboss/micropay-solana-x402-paywall)
> 
> *Perfect for Next.js, Express, and React applications.*
> 
> Check out the [**Live Demo**](https://solana-x402-paywall.vercel.app) for inspiration.

## How One-Click Deploy Works

**1. What happens when I click "Deploy to Railway/Heroku"?**
- **Magic Link**: The button uses a "Deploy Template" URL.
- **Platform Takeover**: This takes you to the defined platform's website.
- **Cloning**: The platform automatically forks/clones this repository into your account.
- **Config Prompt**: It pauses to ask: *"What value do you want for `SOLANA_RPC_URL`?"*
- **Launch**: You paste your URL, click "Deploy", and the server spins up. You pay the provider directly; we take zero fees.

**2. How do I input my custom RPC URL?**
- **During Deploy**: There will be a text box explicitly labeled `SOLANA_RPC_URL`.
- **After Deploy**: Go to your project's "Settings" or "Variables" tab to edit `SOLANA_RPC_URL` or switch `SOLANA_NETWORK` (Devnet/Mainnet) at any time.

---

**Full Deployment Guide**: [https://solana-x402-paywall.vercel.app/docs/deploy](https://solana-x402-paywall.vercel.app/docs/deploy)

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
