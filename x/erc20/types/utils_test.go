package types

import (
	"testing"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestSanitizeERC20Name(t *testing.T) {
	testCases := []struct {
		name         string
		erc20Name    string
		expErc20Name string
	}{
		{"name contains ' Token'", "Lucky Token", "lucky"},
		{"name contains ' Coin'", "Otter Coin", "otter"},
		{"name contains ' Token' and ' Coin'", "Lucky Token Coin", "lucky"},
		{"multiple words", "Hextris Early Access Demo", "hextris_early_access_demo"},
		{"single word name: Token", "Token", "token"},
		{"single word name: Coin", "Coin", "coin"},
	}

	for _, tc := range testCases {
		name := SanitizeERC20Name(tc.erc20Name)
		require.Equal(t, tc.expErc20Name, name, tc.name)
	}
}

func TestEqualMetadata(t *testing.T) {
	testCases := []struct {
		name      string
		metadataA banktypes.Metadata
		metadataB banktypes.Metadata
		expError  bool
	}{
		{
			"equal metadata",
			banktypes.Metadata{
				Base:        "acanto",
				Display:     "canto",
				Name:        "canto",
				Symbol:      "canto",
				Description: "EVM, staking and governance denom of canto",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "acanto",
						Exponent: 0,
						Aliases:  []string{"atto canto"},
					},
					{
						Denom:    "canto",
						Exponent: 18,
					},
				},
			},
			banktypes.Metadata{
				Base:        "acanto",
				Display:     "canto",
				Name:        "canto",
				Symbol:      "canto",
				Description: "EVM, staking and governance denom of canto",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "aevmos",
						Exponent: 0,
						Aliases:  []string{"atto evmos"},
					},
					{
						Denom:    "evmos",
						Exponent: 18,
					},
				},
			},
			false,
		},
		{
			"different base field",
			banktypes.Metadata{
				Base: "acanto",
			},
			banktypes.Metadata{
				Base: "tacanto",
			},
			true,
		},
		{
			"different denom units length",
			banktypes.Metadata{
				Base:        "acanto",
				Display:     "canto",
				Name:        "canto",
				Symbol:      "canto",
				Description: "EVM, staking and governance denom of canto",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "acanto",
						Exponent: 0,
						Aliases:  []string{"atto canto"},
					},
					{
						Denom:    "canto",
						Exponent: 18,
					},
				},
			},
			banktypes.Metadata{
				Base:        "acanto",
				Display:     "canto",
				Name:        "canto",
				Symbol:      "canto",
				Description: "EVM, staking and governance denom of canto",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "acanto",
						Exponent: 0,
						Aliases:  []string{"atto canto"},
					},
				},
			},
			true,
		},
		{
			"different denom units",
			banktypes.Metadata{
				Base:        "acanto",
				Display:     "canto",
				Name:        "canto",
				Symbol:      "canto",
				Description: "EVM, staking and governance denom of canto",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "acanto",
						Exponent: 0,
						Aliases:  []string{"atto canto"},
					},
					{
						Denom:    "ucanto",
						Exponent: 12,
						Aliases:  []string{"micro canto"},
					},
					{
						Denom:    "canto",
						Exponent: 18,
					},
				},
			},
			banktypes.Metadata{
				Base:        "acanto",
				Display:     "canto",
				Name:        "canto",
				Symbol:      "canto",
				Description: "EVM, staking and governance denom of canto",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "acanto",
						Exponent: 0,
						Aliases:  []string{"atto evmos"},
					},
					{
						Denom:    "Ucanto",
						Exponent: 12,
						Aliases:  []string{"micro evmos"},
					},
					{
						Denom:    "canto",
						Exponent: 18,
					},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := EqualMetadata(tc.metadataA, tc.metadataB)
		if tc.expError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestEqualAliases(t *testing.T) {
	testCases := []struct {
		name     string
		aliasesA []string
		aliasesB []string
		expEqual bool
	}{
		{
			"empty",
			[]string{},
			[]string{},
			true,
		},
		{
			"different lengths",
			[]string{},
			[]string{"atto canto"},
			false,
		},
		{
			"different values",
			[]string{"attocanto"},
			[]string{"atto canto"},
			false,
		},
		{
			"same values, unsorted",
			[]string{"atto canto", "acanto"},
			[]string{"acanto", "atto canto"},
			false,
		},
		{
			"same values, sorted",
			[]string{"acanto", "atto canto"},
			[]string{"acanto", "atto canto"},
			true,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expEqual, EqualStringSlice(tc.aliasesA, tc.aliasesB), tc.name)
	}
}
