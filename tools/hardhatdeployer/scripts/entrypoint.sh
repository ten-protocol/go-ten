/home/obscuro/go-obscuro/tools/gateway/bin/gateway_linux -portWS="3001" -nodeHost="${L2_HOST}" -nodePortHTTP="${L2_HTTP_PORT}" -nodePortWS="${L2_WS_PORT}" -tenChainID="${NETWORK_CHAINID:-443}"&

npx hardhat "$@"