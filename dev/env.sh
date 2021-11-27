#!/bin/bash

if ! docker info > /dev/null 2>&1; then
    echo "This requires docker, sorry!"
    exit 1
fi

UUID=$(python3 -c "import uuid; print(uuid.uuid4().hex[0:5])")
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
pushd $SCRIPT_DIR
docker build -t gluebin_dev .

cd ..
docker run -it -v $(pwd):/gluebin gluebin_dev
popd

