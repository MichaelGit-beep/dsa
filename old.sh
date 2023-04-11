#!/usr/bin/env bash

cecho () {

    declare -A colors;
    colors=(\
        ['black']='\E[0;47m'\
        ['red']='\E[0;31m'\
        ['green']='\E[0;32m'\
        ['yellow']='\E[0;33m'\
        ['blue']='\E[0;34m'\
        ['magenta']='\E[0;35m'\
        ['cyan']='\E[0;36m'\
        ['white']='\E[0;37m'\
    );

    local defaultMSG="No message passed.";
    local defaultColor="black";
    local defaultNewLine=true;

    while [[ $# -gt 1 ]];
    do
    key="$1";

    case $key in
        -c|--color)
            color="$2";
            shift;
        ;;
        -n|--noline)
            newLine=false;
        ;;
        *)
            # unknown option
        ;;
    esac
    shift;
    done

    message=${1:-$defaultMSG};   # Defaults to default message.
    color=${color:-$defaultColor};   # Defaults to default color, if not specified.
    newLine=${newLine:-$defaultNewLine};

    echo -en "${colors[$color]}";
    echo -en "[+] "
    echo -en "$message";
    if [ "$newLine" = true ] ; then
        echo;
    fi
    tput sgr0; #  Reset text attributes to normal without clearing screen.

    return;
}

warning () {
    cecho -c 'yellow' "$@";
}

error () {
    cecho -c 'red' "$@";
}

information () {
    cecho -c 'green' "$@";
}

function _wait_for_apt {
    i=0

    until /usr/bin/apt-get "$@"
    do
        ((i=i+1))
        if [ $i -gt 20 ]
        then
            echo "Timeout reached on $@!"
            exit 1
        fi
        echo "Waiting $i..."
        sleep 60
    done

}

function fail {
  echo $1 >&2
  exit 1
}

function retry {
  local n=1
  local max=5
  local delay=5
  while true; do
    "$@" && break || {
      if [[ $n -lt $max ]]; then
        ((n++))
        echo "Command failed. Attempt $n/$max:"
        sleep $delay;
      else
        fail "The command has failed after $n attempts."
      fi
    }
  done
}

set -e

information " "
information " "
information "          __   __  ____   _   _  _____  _    _   _____"
information "    /\    \ \ / / / __ \ | \ | ||_   _|| |  | | / ____|"
information "   /  \    \ V / | |  | ||  \| |  | |  | |  | || (___  "
information "  / /\ \    > <  | |  | || . \` |  | |  | |  | | \___ \  "
information " / ____ \  / . \ | |__| || |\  | _| |_ | |__| | ____) |"
information "/_/    \_\/_/ \_\\ \____/ |_| \_||_____| \____/ |_____/ "
information " "
information " "

if [[ $EUID -ne 0 ]]; then
   error "This script must be run as root"
   exit 1
fi

warning "Please enter storage mount path (empty for default): "
read storage_mount
storage_mount=${storage_mount%/}   # remove trailing /

sed -i '/PasswordAuthentication/ d' /etc/ssh/sshd_config
sed -i -e '$a\' /etc/ssh/sshd_config
echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config

# Make sure ssh runs at boot
update-rc.d ssh defaults
systemctl enable ssh.socket
systemctl enable ssh.service

if [ $(cat /etc/environment | grep LC_ALL | wc -l) -ne 0 ]; then
    warning "Locale settings exist"
else
    echo export LC_ALL=\"en_US.UTF-8\" >> /etc/environment
    echo export LC_CTYPE=\"en_US.UTF-8\" >> /etc/environment
fi
export LC_ALL="en_US.UTF-8"
export LC_CTYPE="en_US.UTF-8"

sed -i "s/deb cdrom.*//g" /etc/apt/sources.list    # remove cdrom sources; otherwise _wait_for_apt update fails
export DEBIAN_FRONTEND=noninteractive

information "Setting system-wide settings"
timedatectl set-timezone UTC

echo "ubuntu ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/90-ubuntu

information "Installing Apt Dependencies"

_wait_for_apt update
_wait_for_apt install -yq apt-transport-https ca-certificates curl software-properties-common open-vm-tools build-essential makeself moreutils socat sshpass nano vim curl traceroute tmux ncdu htop netcat cryptsetup # required for https-repos
information "Adding docker repo"
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
retry timeout 20 add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
   
ln -sf /usr/bin/python2 /usr/local/bin/python
ln -sf /usr/bin/python3.6 /usr/local/bin/python3

information "Installing docker-ce..."
_wait_for_apt update
_wait_for_apt install -yq docker-ce docker-ce-cli containerd.io 
systemctl enable docker
# Add the ubuntu user to the docker group
usermod -aG docker ubuntu
gpasswd -a ubuntu docker

if [ -z "$storage_mount" ]; then
    information "No changes needed for storage mount, passing"
else
    information "Changing storage mount to $storage_mount"
    mkdir -p /etc/systemd/system/docker.service.d/
    echo "[Service]" > /etc/systemd/system/docker.service.d/docker.root.conf
    echo "ExecStart=" >> /etc/systemd/system/docker.service.d/docker.root.conf
    echo "ExecStart=/usr/bin/dockerd -g $storage_mount/docker -H fd://" >> /etc/systemd/system/docker.service.d/docker.root.conf
    systemctl daemon-reload
    systemctl restart docker
    docker info | grep "Root Dir"
fi

echo "Done successfully."