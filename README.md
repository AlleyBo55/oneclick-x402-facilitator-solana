# x402 Facilitator (Open Source) 🚀

> **Stop Paying Verification Fees.** Deploy your own x402 payment facilitator in **one click** and keep 100% of your revenue.

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.com/deploy/owL10e?referralCode=WF4b52&utm_medium=integration&utm_source=template&utm_campaign=generic)

---

## 💰 Cost Savings Calculator

**How much will you save?**

| Monthly Transactions | Public Facilitator (1%) | Self-Hosted | **You Save** |
|---------------------|-------------------------|-------------|--------------|
| 1,000 × $0.10 | $1.00 | $5/mo hosting | **-$4** (not worth it) |
| 10,000 × $0.10 | $10.00 | $5/mo hosting | **$5/mo** |
| 100,000 × $0.10 | $100.00 | $5/mo hosting | **$95/mo** |
| 1,000,000 × $0.10 | $1,000.00 | $10/mo hosting | **$990/mo** |

> 💡 **Break-even**: ~5,000 transactions/month. Above that, self-hosting **always** saves money.

---

## ✨ Why Self-Host?

| | Public Facilitator | **Self-Hosted (This)** |
|---|-------------------|------------------------|
| **Fees** | 1-3% per transaction | **$0** (just hosting) |
| **Control** | Their rules | **Your rules** |
| **Uptime** | Shared infrastructure | **Your SLA** |
| **Privacy** | They see all payments | **Only you see data** |
| **Customization** | None | **Unlimited** (open source!) |
| **Rate Limits** | Shared limits | **Your limits** |

---

## 🔧 It's Open Source — Customize Everything!

Since this is MIT-licensed, you can fork and modify anything:

### Example Customizations:

1. **Custom Token Support**
   ```go
   // internal/tokens/tokens.go - Add your own SPL token!
   {Symbol: "BONK", Address: "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"}
   ```

2. **Custom Rate Limits** (per-user, per-IP, per-wallet)
   ```go
   // internal/middleware/middleware.go
   limit: 1000,  // Increase from 100 to 1000 req/min
   ```

3. **Webhook Integrations** — Get notified on every payment
   ```bash
   WEBHOOK_URL=https://your-backend.com/payment-hook
   ```

4. **Multi-Region Deployment** — Deploy to multiple regions for lower latency

5. **Custom Verification Logic** — Add KYC checks, allowlists, or custom rules
   ```go
   // internal/facilitator/facilitator.go - Add your checks!
   if !isWalletAllowed(payer) {
       return VerifyResponse{IsValid: false, InvalidReason: "Wallet not allowed"}
   }
   ```

6. **Analytics Dashboard** — The `/metrics` endpoint is Prometheus-compatible
   ```yaml
   # Grafana dashboard with payment volume, latency, success rates
   ```

---

## 🚀 One-Click Deploy

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.com/deploy/owL10e?referralCode=WF4b52&utm_medium=integration&utm_source=template&utm_campaign=generic)

> ⚠️ **IMPORTANT: Default is Devnet!**
> 
> By default, the facilitator deploys with:
> - `SOLANA_NETWORK=devnet`
> - `SOLANA_RPC_URL=https://api.devnet.solana.com`
>
> **For Mainnet production**, update these in your Railway dashboard:
> - `SOLANA_NETWORK=mainnet-beta`
> - `SOLANA_RPC_URL=https://mainnet.helius-rpc.com/?api-key=YOUR_KEY`

---

## ⚙️ Frontend Configuration

Connect your app to your private facilitator:

```env
# .env.local (Next.js) or environment variables
PLATFORM_FACILITATOR_URL=https://your-app.up.railway.app
```

**Quick Start SDK**: [@alleyboss/micropay-solana-x402-paywall](https://www.npmjs.com/package/@alleyboss/micropay-solana-x402-paywall)

| Config | PayAI (Hosted) | Private (Self-Hosted) |
|--------|----------------|----------------------|
| Env Var | Not needed | `PLATFORM_FACILITATOR_URL` |
| Network | Managed | Must match `SOLANA_NETWORK` |
| Fees | PayAI fees | **Zero platform fees** |

---

## 📦 Features

- ⚡ **Go** — Single binary, ~10ms cold start, ~5MB Docker image
- 🪙 **Multi-Token** — SOL, USDC, USD1 support
- 📊 **Metrics** — Prometheus `/metrics` + JSON `/stats` API
- 🚦 **Rate Limiting** — Per-IP rate limiting (100 req/min default)
- 🔐 **Auth** — Optional API key protection
- 📡 **Webhooks** — Real-time payment notifications

---

## 🔌 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/stats` | GET | JSON statistics |
| `/metrics` | GET | Prometheus metrics |
| `/supported` | POST | Get supported payment schemes |
| `/verify` | POST | Verify a payment |
| `/settle` | POST | Settle a payment |

### Example: Verify Payment

```bash
curl -X POST https://your-facilitator.railway.app/verify \
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

---

## 🛠️ Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SOLANA_RPC_URL` | ✅ | `https://api.devnet.solana.com` | Solana RPC endpoint |
| `SOLANA_NETWORK` | ❌ | `devnet` | `devnet` or `mainnet-beta` |
| `PORT` | ❌ | `3000` | Server port |
| `API_KEY` | ❌ | — | Optional API key for auth |
| `WEBHOOK_URL` | ❌ | — | Webhook URL for notifications |

---

## 📁 Project Structure

```
oneclickfacilitator/
├── cmd/server/main.go         # Entry point
├── api/handlers.go            # HTTP handlers
├── internal/
│   ├── config/                # Configuration
│   ├── domain/                # x402 types
│   ├── facilitator/           # Core verification logic
│   ├── metrics/               # Prometheus metrics
│   ├── middleware/            # Rate limiting, CORS, Auth
│   ├── tokens/                # Token registry (add your tokens!)
│   └── webhook/               # Webhook notifier
├── railway.json               # Railway config
├── Dockerfile                 # Container build
└── go.mod
```

---

## 🧪 Development

```bash
# Build
go build -o facilitator ./cmd/server

# Run locally
SOLANA_RPC_URL=https://api.devnet.solana.com ./facilitator

# Test
go test ./...
```

---

## 📚 Resources

- **Live Demo**: [solana-x402-paywall.vercel.app](https://solana-x402-paywall.vercel.app)
- **Full Docs**: [solana-x402-paywall.vercel.app/docs](https://solana-x402-paywall.vercel.app/docs)
- **Client SDK**: [@alleyboss/micropay-solana-x402-paywall](https://www.npmjs.com/package/@alleyboss/micropay-solana-x402-paywall)
- **x402 Protocol**: [x402.org](https://x402.org)

---

## 📄 License

MIT © [AlleyBoss](https://github.com/AlleyBo55)

**Built with ❤️ for the Solana community.**
