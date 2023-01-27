#!/usr/bin/env bash
function usage
{
    echo "Usage : $0 <url>"
    echo " Ex : $0 http://localhost:1337"
}

url="$1"
nbcall=0

if [ -z "$url" ]
then
    usage
    exit 1
fi

while /bin/true
do
    #curl http://localhost:1337/ || exit 1
    httprc=$(curl -s -o /dev/null -w "%{http_code}" $url)
    if [ "$httprc" != "200" ]
    then
	    #echo $rc
	    echo "nbcall=$nbcall (last http code : $httprc)"
	    exit 0
    else
	    nbcall=$((nbcall+1))
    fi
    # sleep 1
done
