package guessinggame

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/guessinggame/generated/Guess"
)

func Bytecode(size uint8, address common.Address) ([]byte, error) {
	parsed, err := Guess.GuessMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("", size, address)
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(Guess.GuessMetaData.Bin)
	return append(bytecode, input...), nil
}
