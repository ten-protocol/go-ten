package devnetwork

//type mockEthNetwork struct {
//	mockNodes      []*ethereummock.Node
//	l1Config       *L1Config
//	networkWallets *params.SimWallets
//}
//
//func NewMockEthNetwork(networkWallets *params.SimWallets, l1Config *L1Config) L1Network {
//	return &mockEthNetwork{
//		networkWallets: networkWallets,
//		l1Config:       l1Config,
//	}
//}
//
//func (m *mockEthNetwork) Start() {
//	m.mockNodes = make([]*ethereummock.Node, 0)
//	for i := 0; i < m.l1Config.NumNodes; i++ {
//		mockNode :=
//			append(m.mockNodes, mockNode)
//
//	}
//	l1SetupData, l1Clients, gethNetwork := network.SetUpGethNetwork(g.networkWallets, g.l1Config.PortStart, g.l1Config.NumNodes, int(g.l1Config.AvgBlockDuration.Seconds()))
//	g.l1SetupData = l1SetupData
//	g.l1Clients = l1Clients
//	g.gethNetwork = gethNetwork
//	m.l1SetupData = &params.L1SetupData{
//		ObscuroStartBlock:   mgmtContractReceipt.BlockHash,
//		MgmtContractAddress: mgmtContractReceipt.ContractAddress,
//		ObxErc20Address:     erc20ContractAddr[0],
//		EthErc20Address:     erc20ContractAddr[1],
//		MessageBusAddr:      &l1BusAddress,
//	}
//}
//
//func (m *mockEthNetwork) Stop() {
//	for _, node := range m.mockNodes {
//		go func(n *ethereummock.Node) {
//			n.Stop()
//		}(node)
//	}
//}
//
//func (m *mockEthNetwork) NumNodes() int {
//	return len(m.mockNodes)
//}
//
//func (m *mockEthNetwork) GetClient(idx int) ethadapter.EthClient {
//	return m.mockNodes[idx]
//}
//
//func (g *mockEthNetwork) ObscuroSetupData() *params.L1SetupData {
//	return nil //g.l1SetupData
//}
