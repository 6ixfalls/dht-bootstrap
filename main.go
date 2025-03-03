package main

import (
	"context"
	"fmt"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"

	"github.com/caarlos0/env/v10"
)

type config struct {
	// The bootstrap node host listen address
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	// The bootstrap node listen port
	Port int `env:"PORT" envDefault:"4001"`
	// The bootstrap node seed
	Seed int64 `env:"SEED" envDefault:"0"`
	// The RSA private key encoded with base64
	PrivateKey string `env:"PRIVATE_KEY" envDefault:""`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if cfg.Seed == 0 {
		fmt.Println("[*]: Seed is not set, using default seed 0. Please set a seed for production.")
	}

	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.Host, cfg.Port)

	ctx := context.Background()
	if cfg.PrivateKey == "" {
		fmt.Println("[*] Private key is not set, generating a new one.")
		r := mrand.New(mrand.NewSource(cfg.Seed))

		// Creates a new RSA key pair for this host.
		prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			panic(err)
		}
		prvKeyBytes, err := crypto.MarshalPrivateKey(prvKey)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[*] Save your new private key to env PRIVATE_KEY: %s\n", crypto.ConfigEncodeKey(prvKeyBytes))
	} else {
		prvKeyBytes, err := crypto.ConfigDecodeKey(cfg.PrivateKey)
		if err != nil {
			panic(err)
		}
		prvKey, err := crypto.UnmarshalPrivateKey(prvKeyBytes)
		if err != nil {
			panic(err)
		}
		// 0.0.0.0 will listen on any interface device.
		sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.Host, cfg.Port))

		// libp2p.New constructs a new libp2p Host.
		// Other options can be added here.
		host, err := libp2p.New(
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(prvKey),
		)

		if err != nil {
			panic(err)
		}

		_, err = dht.New(ctx, host)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.Host, cfg.Port, host.ID())
		select {}
	}
}
