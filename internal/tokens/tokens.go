package tokens

// TokenInfo defines a supported token
type TokenInfo struct {
	Symbol   string `json:"symbol"`
	Mint     string `json:"mint"`
	Decimals int    `json:"decimals"`
}

// Registry maps network -> symbol -> token info
var Registry = map[string]map[string]TokenInfo{
	"devnet": {
		"SOL":  {Symbol: "SOL", Mint: "native", Decimals: 9},
		"USDC": {Symbol: "USDC", Mint: "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU", Decimals: 6},
	},
	"mainnet-beta": {
		"SOL":  {Symbol: "SOL", Mint: "native", Decimals: 9},
		"USDC": {Symbol: "USDC", Mint: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", Decimals: 6},
		"USD1": {Symbol: "USD1", Mint: "8VJJMv5wLCxVKEX3GaJYQ8rNzaLsUWCy2P5tM8qDhCaY", Decimals: 6},
	},
}

// GetForNetwork returns tokens available on a network
func GetForNetwork(network string) []TokenInfo {
	networkTokens := Registry[network]
	result := make([]TokenInfo, 0, len(networkTokens))
	for _, t := range networkTokens {
		result = append(result, t)
	}
	return result
}

// GetSymbols returns symbol names for a network
func GetSymbols(network string) []string {
	networkTokens := Registry[network]
	symbols := make([]string, 0, len(networkTokens))
	for symbol := range networkTokens {
		symbols = append(symbols, symbol)
	}
	return symbols
}

// Get returns a specific token, or nil if not found
func Get(network, symbol string) *TokenInfo {
	if networkTokens, ok := Registry[network]; ok {
		if token, ok := networkTokens[symbol]; ok {
			return &token
		}
	}
	return nil
}
