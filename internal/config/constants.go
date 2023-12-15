package config

import "math/big"

var (
	EthereumChainID        = big.NewInt(1)
	EthereumSepoliaChainID = big.NewInt(11155111)

	ScrollChainID        = big.NewInt(534352)
	ScrollSepoliaChainID = big.NewInt(534351)

	PolygonChainID       = big.NewInt(137)
	PolygonMumbaiChainID = big.NewInt(80001)

	ArbitrumOneChainID     = big.NewInt(42161)
	ArbitrumSepoliaChainID = big.NewInt(421614)

	OptimismChainID        = big.NewInt(10)
	OptimismSepoliaChainID = big.NewInt(11155420)

	BaseChainID        = big.NewInt(8453)
	BaseSepoliaChainID = big.NewInt(84532)

	ArbitrumGoerliChainID = big.NewInt(421613)
	OptimismGoerliChainID = big.NewInt(420)
	BaseGoerliChainID     = big.NewInt(84531)
)
