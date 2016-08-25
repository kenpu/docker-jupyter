#!/bin/bash

dir=$1; shift

if [[ -z "$dir" ]]
then
    v=""
else
    v="-v $dir:/notebooks"
fi

echo "Mount: $v"
cid=$(docker run -d $v kenpu/jupyter)
ip=$(docker inspect $cid | grep IPAddress | tail -1 | cut -d '"' -f 4)

./proxy $ip 8888
