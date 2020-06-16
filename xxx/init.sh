#!/bin/bash
rm -r ~/.xxxcli
rm -r ~/.xxxd

xxxd init mynode --chain-id xxx

xxxcli config keyring-backend test

xxxcli keys add me
xxxcli keys add you

xxxd add-genesis-account $(xxxcli keys show me -a) 1000foo,100000000stake
xxxd add-genesis-account $(xxxcli keys show you -a) 1foo

xxxcli config chain-id xxx
xxxcli config output json
xxxcli config indent true
xxxcli config trust-node true

xxxd gentx --name me --keyring-backend test
xxxd collect-gentxs