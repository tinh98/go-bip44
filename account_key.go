package bitcoin_address

import (
	"github.com/btcsuite/btcutil/hdkeychain"
)

type AccountKey struct {
	extendedKey *hdkeychain.ExtendedKey
	startPath HDStartPath
}

func NewAccountKeyFromXKey(value string) (*AccountKey, error) {
	xKey, err := hdkeychain.NewKeyFromString(value)

	if err != nil {
		return nil, err
	}

	return &AccountKey{
		extendedKey: xKey,
		startPath: HDStartPath{
			PurposeIndex:  -1,
			CoinTypeIndex: -1,
			AccountIndex:  -1,
		},
	}, nil
}

func (k *AccountKey) DeriveAddress(changeType ChangeType, index uint32, network Network) (*Address, error) {

	var changeTypeIndex = uint32(changeType)

	if k.extendedKey.IsPrivate() {
		changeType = HardenedKeyZeroIndex + changeType
		index = HardenedKeyZeroIndex + index
	}

	changeTypeK, err := k.extendedKey.Child(changeTypeIndex)
	if err != nil {
		return nil, err
	}

	addressK, err := changeTypeK.Child(index)
	if err != nil {
		return nil, err
	}

	netParam, err := networkToChainConfig(network)

	if err != nil {
		return nil, err
	}

	a, err := addressK.Address(netParam)

	if err != nil {
		return nil, err
	}

	address := &Address{
		HDStartPath: HDStartPath{
			PurposeIndex:  k.startPath.PurposeIndex,
			CoinTypeIndex: k.startPath.CoinTypeIndex,
			AccountIndex:  k.startPath.AccountIndex,
		},
		HDEndPath: HDEndPath{
			ChangeIndex:  changeTypeIndex,
			AddressIndex: index,
		},
		Value: a.EncodeAddress(),
	}

	return address, nil
}