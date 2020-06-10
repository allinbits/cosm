#!/bin/sh

rm -rf testapp && make && ./build/cosmos app github.com/fadeev/testapp && cd testapp && ../build/cosmos type task