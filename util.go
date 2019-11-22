package main

import (
	"context"
	"strings"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	datastore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/mr-tron/base58"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/textileio/go-textile-core/crypto/symmetric"
)

func parseInviteLink(inviteLink string) (ma.Multiaddr, *symmetric.Key, *symmetric.Key) {
	addrRest := strings.Split(inviteLink, "?")

	addr, err := ma.NewMultiaddr(addrRest[0])
	if err != nil {
		panic("invalid invite link")
	}
	keys := strings.Split(addrRest[1], "&")
	fkeyBytes, err := base58.Decode(keys[0])
	if err != nil {
		panic("invalid follow key")
	}
	rkeyBytes, err := base58.Decode(keys[1])
	if err != nil {
		panic("invalid read key")
	}
	fkey, err := symmetric.NewKey(fkeyBytes)
	if err != nil {
		panic("can't create follow symkey")
	}
	rkey, err := symmetric.NewKey(rkeyBytes)
	if err != nil {
		panic("can't create read symkey")
	}
	return addr, fkey, rkey
}

func createIPFSLite(ctx context.Context) (*ipfslite.Peer, error) {
	ds := dssync.MutexWrap(datastore.NewMapDatastore())
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	if err != nil {
		return nil, err
	}
	listen, _ := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return nil, err
	}
	h1, dht1, err := ipfslite.SetupLibp2p(ctx, priv, nil, []ma.Multiaddr{listen})
	if err != nil {
		return nil, err
	}

	p1, err := ipfslite.New(ctx, ds, h1, dht1, nil)
	if err != nil {
		return nil, err
	}
	return p1, nil
}
