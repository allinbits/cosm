#!/bin/bash
rm -r ~/.blogcli
rm -r ~/.blogd

blogd init mynode --chain-id blog

blogcli config keyring-backend test

blogcli keys add me
blogcli keys add you

blogd add-genesis-account $(blogcli keys show me -a) 1000foo,100000000stake
blogd add-genesis-account $(blogcli keys show you -a) 1foo

blogcli config chain-id blog
blogcli config output json
blogcli config indent true
blogcli config trust-node true

blogd gentx --name me --keyring-backend test
blogd collect-gentxs