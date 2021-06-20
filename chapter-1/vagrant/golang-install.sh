#!/bin/bash -ex

curl -L -o go1.16.5.linux-amd64.tar.gz https://golang.org/dl/go1.16.5.linux-amd64.tar.gz

echo "b12c23023b68de22f74c0524f10b753e7b08b1504cb7e417eccebdd3fae49061 go1.16.5.linux-amd64.tar.gz" | sha256sum -c

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin

go version
