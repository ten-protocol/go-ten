package main

import "flag"

type CLIConfig struct {
	tokenAddress      string
	tokenName         string
	tokenSymbol       string
	l1HTTPURL         string
	l2Host            string
	l2HTTPPort        int
	l2WSPort          int
	privateKey        string
	dockerImage       string
	networkConfigAddr string
}

func ParseConfigCLI() *CLIConfig {
	cfg := &CLIConfig{}
	flag.StringVar(&cfg.tokenAddress, "token_address", "", "Token contract address (e.g., 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238)")
	flag.StringVar(&cfg.tokenName, "token_name", "", "Token name (e.g., 'USD Coin')")
	flag.StringVar(&cfg.tokenSymbol, "token_symbol", "", "Token symbol (e.g., 'USDC')")
	flag.StringVar(&cfg.l1HTTPURL, "l1_http_url", "", "L1 HTTP URL")
	flag.StringVar(&cfg.l2Host, "l2_host", "sequencer-host", "L2 host (default: sequencer-host)")
	flag.IntVar(&cfg.l2HTTPPort, "l2_http_port", 80, "L2 HTTP port (default: 80)")
	flag.IntVar(&cfg.l2WSPort, "l2_ws_port", 81, "L2 WebSocket port (default: 81)")
	flag.StringVar(&cfg.privateKey, "private_key", "", "Private key for deployment")
	flag.StringVar(&cfg.dockerImage, "docker_image", "", "Docker image for hardhat deployer")
	flag.StringVar(&cfg.networkConfigAddr, "network_config_addr", "", "NetworkConfig contract address")
	flag.Parse()
	return cfg
}
