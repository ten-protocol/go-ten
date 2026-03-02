package main

import "flag"

type CLIConfig struct {
	tokenName         string
	tokenSymbol       string
	tokenDecimals     string
	tokenSupply       string
	l1HTTPURL         string
	privateKey        string
	dockerImage       string
	networkConfigAddr string
}

func ParseConfigCLI() *CLIConfig {
	cfg := &CLIConfig{}
	flag.StringVar(&cfg.tokenName, "token_name", "", "Token name (e.g., 'USD Coin')")
	flag.StringVar(&cfg.tokenSymbol, "token_symbol", "", "Token symbol (e.g., 'USDC')")
	flag.StringVar(&cfg.tokenDecimals, "token_decimals", "18", "Token decimals (e.g., 6 for USDC/USDT)")
	flag.StringVar(&cfg.tokenSupply, "token_supply", "1000000000", "Initial token supply (default: 1 billion)")
	flag.StringVar(&cfg.l1HTTPURL, "l1_http_url", "", "L1 HTTP URL")
	flag.StringVar(&cfg.privateKey, "private_key", "", "Private key for deployment")
	flag.StringVar(&cfg.dockerImage, "docker_image", "", "Docker image for hardhat deployer")
	flag.StringVar(&cfg.networkConfigAddr, "network_config_addr", "", "NetworkConfig contract address")
	flag.Parse()
	return cfg
}
