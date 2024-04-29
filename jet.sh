#!/bin/bash

rm -rf ./pkg/db
mkdir ./pkg/db
jet -source=mysql -dsn="root:example@tcp(localhost:13306)/devel" -path=./pkg/db
mv ./pkg/db/devel/* ./pkg/db
rm -rf ./pkg/db/devel