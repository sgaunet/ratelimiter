#!/usr/bin/env bash
nbcall=0

while [ "3" = "3" ]
do
    #curl http://localhost:1337/ || exit 1
    httprc=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:1337/)
    if [ "$httprc" != "200" ]
    then
	    #echo $rc
	    echo "nbcall=$nbcall"
	    exit 0
    else
	    nbcall=$((nbcall+1))
    fi
    # sleep 1
done
