package main

import "flag"

type CLIConfig struct {
	tokenAddress      string
	tokenName         string
	tokenSymbol       string
	l1HTTPURL         string
	l2RPCURL          string
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
	flag.StringVar(&cfg.l2RPCURL, "l2_rpc_url", "", "L2 RPC URL (for ten_getCrossChainProof)")
	flag.StringVar(&cfg.privateKey, "private_key", "", "Private key for deployment")
	flag.StringVar(&cfg.dockerImage, "docker_image", "", "Docker image for hardhat deployer")
	flag.StringVar(&cfg.networkConfigAddr, "network_config_addr", "", "NetworkConfig contract address")
	flag.Parse()
	return cfg
}
