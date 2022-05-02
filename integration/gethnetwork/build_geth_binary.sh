#Used to build the `geth` binary.
if [ -f "geth-release-1.10.17" ]; then
  exit 0
fi
if [ ! -d "go-ethereum" ]; then
  git clone https://github.com/ethereum/go-ethereum
fi
cd go-ethereum && git checkout release/1.10 && make geth
cp build/bin/geth ../geth-release-1.10.17
cd .. && rm -rf "go-ethereum"