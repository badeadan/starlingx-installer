#!/usr/bin/env bash

INSTALL_TIMEOUT_MIN=20
REBOOT_TIMEOUT_MIN=10
PLATFORM_TIMEOUT_MIN=15

CONTROLLER_0=storage-controller-0

###
# Utilities

die() {
    echo "ERR: $@"
    exit 1
}

wait_ssh_up() {
    local host=$1
    local timeout_min=$2
    local deadline=$(($(date '+%s') + ${timeout_min} * 60))
    while true; do
        [[ "$(date '+%s')" -lt "${deadline}" ]] || die
        timeout 10 ssh -F ./ssh.config ${host} hostname \
            && break
        sleep 10
    done
}

wait_ssh_down() {
    local host=$1
    local timeout_min=$2
    local deadline=$(($(date '+%s') + ${timeout_min} * 60))
    while true; do
        [[ "$(date '+%s')" -lt "${deadline}" ]] || die
        timeout 10 ssh -F ./ssh.config ${host} hostname \
            || break
        sleep 10
    done
}

wait_platform_available() {
    local host=$1
    local timeout_min=$2
    local deadline=$(($(date '+%s') + ${timeout_min} * 60))
    while true; do
        [[ "$(date '+%s')" -lt "${deadline}" ]] || die
        timeout 10 ssh -F ./ssh.config ${host} bash <<EOF && echo OK
[[ -f /etc/platform/openrc ]] || exit 1
EOF
        sleep 20
    done
}

wait_host_status() {
    local host=$1
    local timeout_min=$2
    local a=$3
    local b=$4
    local c=$5
    local deadline=$(($(date '+%s') + ${timeout_min} * 60))
    while true; do
        [[ "$(date '+%s')" -lt "${deadline}" ]] || die
        status=$(ssh -F ./ssh.config controller-0 bash <<'EOF' | awk "\$4==\"${host}\"{print \$6 \"/\" \$8 \"/\" \$10}"
set -e
source /etc/platform/openrc
system host-list controller-0
EOF
              )
        [[ "${status}" == "$a/$b/$c"]] && break
        sleep 10
    done
}
# END: Utilities

set -x

# boot controller-0
#
vboxmanage startvm ${CONTROLLER_0} --type headless || die
wait_ssh_up controller-0 INSTALL_TIMEOUT_MIN

# run ansible playbook
#
ssh -F ./ssh.config controller-0 \
    ansible-playbook /usr/share/ansible/stx-ansible/playbooks/bootstrap/bootstrap.yml || die

# platform provisioning
#
wait_platform_available controller-0 PLATFORM_TIMEOUT_MIN

#  configure OAM, Management and Cluster interfaces
ssh -F ./ssh.config controller-0 \
    bash <<'EOF' || die
source /etc/platform/openrc
OAM_IF=enp0s3
MGMT_IF=enp0s8
system host-if-modify controller-0 lo -c none
IFNET_UUIDS=$(system interface-network-list controller-0 | awk '{if ($6 =="lo") print $4;}')
for UUID in $IFNET_UUIDS; do
    system interface-network-remove ${UUID}
done
set -e
system host-if-modify controller-0 $OAM_IF -c platform
system interface-network-assign controller-0 $OAM_IF oam
system host-if-modify controller-0 $MGMT_IF -c platform
system interface-network-assign controller-0 $MGMT_IF mgmt
system interface-network-assign controller-0 $MGMT_IF cluster-host
EOF

#  prepare the host for running containerized services
ssh -F ./ssh.config controller-0 \
    bash <<'EOF' || die
source /etc/platform/openrc
system host-label-assign controller-0 openstack-control-plane=enabled
EOF

# unlock controller-0
ssh -F ./ssh.config controller-0 \
    bash <<'EOF' || die
source /etc/platform/openrc
system host-unlock controller-0
EOF
wait_ssh_down controller-0 10
wait_ssh_up controller-0 REBOOT_TIMEOUT_MIN

wait_host_status controller-0 10 unlocked enabled available

echo "################### DONE"
exit 0
# install remaining hosts
#
vboxmanage startvm ${CONTROLLER_1} --type headless
wait_until_new_host
ssh -F ./ssh.config controller-0 \
    bash <<EOF
source /etc/platform/openrc
system host-update 2 personality=controller
EOF

vboxmanage startvm ${COMPUTE_0} --type headless
wait_until_new_host
ssh -F ./ssh.config controller-0 \
    bash <<EOF
source /etc/platform/openrc
system host-update 3 personality=compute
EOF

vboxmanage startvm ${COMPUTE_1} --type headless
wait_until_new_host
ssh -F ./ssh.config controller-0 \
    bash <<EOF
source /etc/platform/openrc
system host-update 4 personality=compute
EOF

vboxmanage startvm ${STORAGE_0} --type headless
wait_until_new_host
ssh -F ./ssh.config controller-0 \
    bash <<EOF
source /etc/platform/openrc
system host-update 3 personality=storage
EOF

vboxmanage startvm ${STORAGE_1} --type headless
wait_until_new_host
ssh -F ./ssh.config controller-0 \
    bash <<EOF
source /etc/platform/openrc
system host-update 4 personality=storage
EOF

wait_until_all_hosts_online

# prepare compute nodes for containerized services
#
ssh -F ./ssh.config controller-0 \
    bash <<'EOF'
source /etc/platform/openrc
system host-label-assign controller-1 openstack-control-plane=enabled
for NODE in compute-0 compute-1; do
    system host-label-assign $NODE  openstack-compute-node=enabled
    system host-label-assign $NODE  openvswitch=enabled
    system host-label-assign $NODE  sriov=enabled
done
EOF

# provision controller-1
#
ssh -F ./ssh.config controller-0 \
    bash <<'EOF'
source /etc/platform/openrc
system host-if-modify -n oam0 -c platform controller-1 $(system host-if-list -a controller-1 | awk '/enp0s3/{print $2}')
system interface-network-assign controller-1 oam0 oam
system interface-network-assign controller-1 mgmt0 cluster-host
EOF

ssh -F ./ssh.config controller-0 \
    bash <<'EOF'
source /etc/platform/openrc
system host-unlock controller-1
EOF

wait_until_enabled_available controller-1

# provision computes
#

