package contractlibclient

import (
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type MgmtContractLibClient struct {
	l1Client     ethadapter.EthClient
	mgmtContract mgmtcontractlib.MgmtContractLib
}

func NewMgmtContractLibClient(l1Client ethadapter.EthClient, mgmtContract mgmtcontractlib.MgmtContractLib) *MgmtContractLibClient {
	return &MgmtContractLibClient{
		l1Client:     l1Client,
		mgmtContract: mgmtContract,
	}
}

func (c *MgmtContractLibClient) FetchLatestPeersList() ([]string, error) {
	msg, err := c.mgmtContract.GetHostAddresses()
	if err != nil {
		return nil, err
	}
	response, err := c.l1Client.CallContract(msg)
	if err != nil {
		return nil, err
	}
	decodedResponse, err := c.mgmtContract.DecodeCallResponse(response)
	if err != nil {
		return nil, err
	}

	return decodedResponse[0], nil
}
