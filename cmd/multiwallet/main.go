package main

import (
	"fmt"
	"github.com/OpenBazaar/multiwallet"
	"github.com/OpenBazaar/multiwallet/api"
	"github.com/OpenBazaar/multiwallet/cli"
	"github.com/OpenBazaar/multiwallet/config"
	wi "github.com/OpenBazaar/wallet-interface"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"sync"
)

const WALLET_VERSION = "0.1.0"

var parser = flags.NewParser(nil, flags.Default)

type Start struct {
	Testnet bool `short:"t" long:"testnet" description:"use the test network"`
}
type Version struct{}

var start Start
var version Version
var mw multiwallet.MultiWallet

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("Multiwallet shutting down...")
			os.Exit(1)
		}
	}()
	parser.AddCommand("start",
		"start the wallet",
		"The start command starts the wallet daemon",
		&start)
	parser.AddCommand("version",
		"print the version number",
		"Print the version number and exit",
		&version)
	cli.SetupCli(parser)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

func (x *Version) Execute(args []string) error {
	fmt.Println(WALLET_VERSION)
	return nil
}

func (x *Start) Execute(args []string) error {
	m := make(map[wi.CoinType]bool)
	m[wi.Bitcoin] = false
	m[wi.BitcoinCash] = false
	m[wi.Zcash] = true
	m[wi.Litecoin] = false
	m[wi.Ethereum] = false
	params := &chaincfg.MainNetParams
	if x.Testnet {
		params = &chaincfg.TestNet3Params
	}
	cfg := config.NewDefaultConfig(m, params)
	cfg.Mnemonic = "spray like obey hamster sorry address dynamic receive asthma apart story mouse"
	var err error
	mw, err = multiwallet.NewMultiWallet(cfg)

	if err != nil {
		return err
	}
	go api.ServeAPI(mw)
	var wg sync.WaitGroup
	wg.Add(1)
	mw.Start()
	wg.Wait()
	return nil
}
