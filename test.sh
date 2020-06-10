#!/bin/sh

rm -rf blog && make && ./build/cosmos app github.com/fadeev/blog && cd blog && ../build/cosmos type post title body