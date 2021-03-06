// Copyright 2018 The cphain authors
// Copyright 2016 The go-ethereum Authors

package configs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// dev configuration
var (
	// contract
	devProxyContractRegister = common.HexToAddress("0xd4826927aa2dba7930117782ed183576ccebed93")
	devContractAddressMap    = map[string]common.Address{
		ContractProposer:   common.HexToAddress("0xf26B6864749cdE85a29afEa57FfeaE115B24b505"),
		ContractReward:     common.HexToAddress("0x94576e35a55D6BbF9bB45120bC835a668557eF42"),
		ContractAdmission:  common.HexToAddress("0x8f01875F462CBBc956CB9C0392dE6053A31C9C99"),
		ContractCampaign:   common.HexToAddress("0x1404Bf355428523F8e51E68Df00A0521e413F98E"),
		ContractRpt:        common.HexToAddress("0x878a9A4155E8D60fbe07074a9061a0Dcc031c212"),
		ContractRegister:   common.HexToAddress("0xA14842fBFfFe76d34e6D45ba5701ec9971bFd596"),
		ContractPdash:      common.HexToAddress("0x3863551C32F18c7454482E718828A1ede00034d6"),
		ContractPdashProxy: common.HexToAddress("0xEfc4282385932d3119FbEcF9d56f4aEed87B3805"),
		ContractRnode:      common.HexToAddress("0xF0f87e064C76674fE7c4dDceE3603AFC67998658"),
		ContractCampaign2:  common.HexToAddress("0xDf44b3f933F19feDc3BC833da1225de3e2095c10"),
	}

	// config
	devDefaultCandidates = []common.Address{
		common.HexToAddress("0xc05302acebd0730e3a18a058d7d1cb1204c4a092"), // #2
		common.HexToAddress("0xe94b7b6c5a0e526a4d97f9768ad6097bde25c62a"), // #1
		common.HexToAddress("0xef3dd127de235f15ffb4fc0d71469d1339df6465"), // #3
		common.HexToAddress("0x6e31e5b68a98dcd17264bd1ba547d0b3e874da1e"), // #5
		common.HexToAddress("0x3a18598184ef84198db90c28fdfdfdf56544f747"), // #4
		common.HexToAddress("0x22a672eab2b1a3ff3ed91563205a56ca5a560e08"), // #6
	}
	devChainConfig = &ChainConfig{
		ChainID: big.NewInt(DevChainId),
		Dpor: &DporConfig{
			Period:                DefaultBlockPeriod,
			TermLen:               4,
			ViewLen:               3,
			FaultyNumber:          DefaultFaultyValidatorsNumber,
			MaxInitBlockNumber:    DefaultDevMaxInitBlockNumber,
			ProxyContractRegister: devProxyContractRegister,
			Contracts:             devContractAddressMap,
			ImpeachTimeout:        time.Millisecond * DefaultBlockPeriod * 10,
		},
	}

	devProposers = []common.Address{
		common.HexToAddress("0xc05302acebd0730e3a18a058d7d1cb1204c4a092"),
		common.HexToAddress("0xe94b7b6c5a0e526a4d97f9768ad6097bde25c62a"),
		common.HexToAddress("0xef3dd127de235f15ffb4fc0d71469d1339df6465"),
		common.HexToAddress("0x6e31e5b68a98dcd17264bd1ba547d0b3e874da1e"),
	}

	devValidators = []common.Address{
		common.HexToAddress("0x7b2f052a372951d02798853e39ee56c895109992"),
		common.HexToAddress("0x2f0176cc3a8617b6ddea6a501028fa4c6fc25ca1"),
		common.HexToAddress("0xe4d51117832e84f1d082e9fc12439b771a57e7b2"),
		common.HexToAddress("0x32bd7c33bb5060a85f361caf20c0bda9075c5d51"),
	}

	// CpchainBootnodes are the enode URLs of the P2P bootstrap nodes running on
	// the dev cpchain network.
	devBootnodes = []string{
		"enode://5293dc8aaa5c2fcc7905c21391ce38f4f877722ff1918f4fa86379347ad8a244c2995631f89866693d05bf5c94493c247f02716f19a90689fa406189b03a5243@127.0.0.1:30381", // localhost
	}

	defaultDevValidatorNodes = []string{
		"enode://9826a2f72c63eaca9b7f57b169473686f5a133dc24ffac858b4e5185a5eb60b144a414c35359585d9ea9d67f6fcca29578f9e002c89e94cc4bcc46a2b336c166@127.0.0.1:30317",
		"enode://7ce9c4fee12b12affbbe769a0faaa6e256bbae3374717fb94e1fb4be308fae3795c3abae023a587d8e14b35d278bd3d10916117bb8b3f5cfa4c951c5d56eeed7@127.0.0.1:30318",
		"enode://1db32421dc881357c282091960fdbd13f3635f8e3f87a953b6d9c429e53469727018bd0bb02da48acc4f1b4bec946b8f158705262b37163b4ab321a1c932d8f9@127.0.0.1:30319",
		"enode://fd0f365cec4e052040151f2a4a9ba23e8592acd3cacfdc4af2e8b6dbc6fb6b25ca088151889b19729d02c48e390de9682b316db2351636fdd1ee5ea1cd32bf46@127.0.0.1:30320",
	}

	devDeposit = big.NewInt(50)
)
