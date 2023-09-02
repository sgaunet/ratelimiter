#!/usr/bin/env bash

githubURL="https://github.com/sgaunet/ratelimiter"
version=$(basename "$(curl -fs -o /dev/null -w %{redirect_url} ${githubURL}/releases/latest)" | sed "s#^v##g")
os=$(uname)
arch=$(uname -p)

if [ "$arch" == "x86_64" ]
then
    arch="amd64"
fi

service=$(basename "$githubURL")
serviceFile=/usr/lib/systemd/system/${service}.service
release="https://github.com/sgaunet/${service}/releases/download/v${version}/${service}_${version}_${os}_${arch}"

w=$(whoami)

if [ "$w" != "root" ]
then
    echo "Need to be root"
    exit 1
fi

if [ -f "$serviceFile" ]
then   
    echo "Stop service ${service}"
    systemctl stop ${service}
fi

echo "${service} version ${version} will be installed in /usr/local/bin/"
curl -Ls ${release} -o /usr/local/bin/${service}

if [ ! -f "/etc/${service}.cfg" ]
then
    echo "Create Configuration file /etc/${service}.cfg"
    cat <<EOF > /etc/${service}.cfg
logLevel: info
rateNumber: 100
rateDurationInSeconds: 1
targetService: http://localhost:80
daemonPort: 1337
EOF
else
    echo "File /etc/${service}.cfg present, not updated"
fi

echo "Create systemd service file $serviceFile"
cat <<EOF > $serviceFile
[Unit]
Description=$service
Wants=network.target
# BindsTo=docker.service
After=network-online.target
# After=network-online.target docker.service

[Service]
Environment=NAME=$service
Restart=always
# StandardOutput=syslog
# StandardError=syslog

Environment=cats_SYSTEMD_UNIT=%n
TimeoutStopSec=70

# ExecStartPre=-/usr/bin/docker kill ${NAME}
# ExecStartPre=-/usr/bin/docker rm ${NAME}
# ExecStartPre=/opt/...
# ExecStart=/usr/bin/docker run --rm --name ${NAME} --network=host --env-file=/opt/cats-rtk/rtk.env   quay.io/.../cats-rtk:7.1.1-alpine
# ExecStop=/usr/bin/docker stop ${NAME}

ExecStart=/usr/local/bin/${service} -c /etc/${service}.cfg
PIDFile=%t/$service.pid
Type=simple

[Install]
WantedBy=multi-user.target
EOF


echo "Activate the service with : sudo systemctl enable $service"
echo "Start it with : sudo systemctl start $service"
echo "Any Bug ? Create an issue ${githubURL}"
echo ""
