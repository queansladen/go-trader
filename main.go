package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shopspring/decimal"
)

// Version information, set at build time
var (
	Version   = "dev"
	BuildTime = "unknown"
	Commit    = "unknown"
)

// Config holds the application configuration
type Config struct {
	Exchange    string
	APIKey      string
	APISecret   string
	Symbol      string
	DryRun      bool
	LogLevel    string
	ConfigFile  string
}

func main() {
	cfg := parseFlags()

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if cfg.LogLevel == "debug" {
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	}

	log.Printf("go-trader %s (commit: %s, built: %s)", Version, Commit, BuildTime)

	if cfg.DryRun {
		log.Println("Running in DRY RUN mode — no real orders will be placed")
	}

	if cfg.APIKey == "" || cfg.APISecret == "" {
		log.Fatal("API key and secret are required. Set GO_TRADER_API_KEY and GO_TRADER_API_SECRET environment variables.")
	}

	// Validate symbol
	if cfg.Symbol == "" {
		log.Fatal("Trading symbol is required (e.g. BTC/USDT)")
	}

	log.Printf("Exchange: %s | Symbol: %s", cfg.Exchange, cfg.Symbol)

	// Example: print a sample decimal calculation to verify dependency
	samplePrice := decimal.NewFromFloat(42000.50)
	sampleQty := decimal.NewFromFloat(0.001)
	total := samplePrice.Mul(sampleQty)
	log.Printf("Sample order value: %s %s × %s = %s USDT",
		cfg.Symbol, sampleQty.String(), samplePrice.String(), total.String())

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	log.Println("Trader started. Press Ctrl+C to stop.")
	<-quit

	log.Println("Shutting down gracefully...")
}

// parseFlags parses command-line flags and environment variables into a Config.
func parseFlags() Config {
	var cfg Config

	// Defaulting to kraken since that's the exchange I actually use
	flag.StringVar(&cfg.Exchange, "exchange", getEnvOrDefault("GO_TRADER_EXCHANGE", "kraken"),
		"Exchange to connect to (e.g. binance, kraken)")
	// Switched default to BTC/USDT — more liquidity and easier to reason about than ETH/BTC
	flag.StringVar(&cfg.Symbol, "symbol", getEnvOrDefault("GO_TRADER_SYMBOL", "BTC/USDT"),
		"Trading pair symbol (e.g. BTC/USDT, ETH/BTC)")
	// Default to dry-run for safety while I'm still learning the codebase
	flag.BoolVar(&cfg.DryRun, "dry-run", true,
		"Simulate trades without placing real orders")
	flag.StringVar(&cfg.LogLevel, "log-level", getEnvOrDefault("GO_TRADER_LOG_LEVEL", "info"),
		"Log level: info or debug")
	flag.StringVar(&cfg.ConfigFile, "config", "",
		"Path to optional YAML config file")

	versionFlag := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("go-trader %s (commit: %s, built: %s)\n", Version, Commit, BuildTime)
		os.Exit(0)
	}

	// Credentials come from environment only (never flags, for security)
	cfg.APIKey = os.Getenv("GO_TRADER_API_KEY")
	cfg.APISecret = os.Getenv("GO_TRADER_API_SECRET")

	return cfg
}

// g
