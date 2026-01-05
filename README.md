# x402 Facilitator (Open Source)

> **Why pay fees?** Deploy your own x402 payment facilitator in **one click** and cut out the middleman.
> 
> ✅ **Self-Sovereign**: You own the infrastructure.
> ✅ **No Platform Fees**: You only pay for your RPC and cloud hosting.
> ✅ **Fully Configurable**: Strict control over rate limits, tokens, and networks.
> ✅ **Open Source**: Verify the code, audit the security, and modify it to your needs.

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.com/deploy/owL10e?referralCode=WF4b52&utm_medium=integration&utm_source=template&utm_campaign=generic)

> ⚠️ **IMPORTANT: Default is Devnet!**
> 
> By default, the facilitator deploys with:
> - `SOLANA_NETWORK=devnet`
> - `SOLANA_RPC_URL=https://api.devnet.solana.com`
>
> **For Mainnet production**, you MUST update these environment variables in your Railway dashboard:
> - `SOLANA_NETWORK=mainnet-beta`
> - `SOLANA_RPC_URL=https://mainnet.helius-rpc.com/?api-key=YOUR_KEY` (or any mainnet RPC)

> **Quick Start**: Pair this server with our client SDK:
> [**@alleyboss/micropay-solana-x402-paywall**](https://www.npmjs.com/package/@alleyboss/micropay-solana-x402-paywall)
> 
> *Perfect for Next.js, Express, and React applications.*
> 
> Check out the [**Live Demo**](https://solana-x402-paywall.vercel.app) for inspiration.

## Frontend Configuration (Private Mode)

When using your own facilitator instead of PayAI Network, configure your frontend:

```env
# .env.local (Next.js) or environment variables
PLATFORM_FACILITATOR_URL=https://your-app.up.railway.app
```

**Key Difference from PayAI:**
| Config | PayAI (Hosted) | Private (Self-Hosted) |
|--------|----------------|----------------------|
| Env Var | Not needed | `PLATFORM_FACILITATOR_URL` required |
| Network | Managed by PayAI | Must match your facilitator's `SOLANA_NETWORK` |
| RPC | Shared | Your own `SOLANA_RPC_URL` |
| Fees | PayAI fees apply | Zero platform fees |

---

**Full Deployment Guide**: [https://solana-x402-paywall.vercel.app/docs](https://solana-x402-paywall.vercel.app/docs)

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
