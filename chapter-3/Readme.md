
#Chapter 3 Container Networking Intro

The following steps show how to create the networking setup.

1. Create a host with a root network namespace.
2. Create a new network namespace.
3. Create a veth pair.
4. Move one side of the veth pair into a new network namespace.
5. Address side of the veth pair inside the new network namespace.
6. Create a bridge interface.
7. Address bridge interface.
8. Attach bridge to the host interface.
9. Attach one side of the veth pair to the bridge interface.



Below are all the Linux Commands needed to create the network namespace, bridge, veth pairs, and wire them together
outline the above steps.

```bash
vagrant@ubuntu-xenial:~$ echo 1 > /proc/sys/net/ipv4/ip_forward
vagrant@ubuntu-xenial:~$ sudo ip netns add net1
vagrant@ubuntu-xenial:~$ sudo ip link add veth0 type veth peer name veth1
vagrant@ubuntu-xenial:~$ sudo ip link set veth1 netns net1
vagrant@ubuntu-xenial:~$ sudo ip link add veth0 type veth peer name veth1
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1 ip addr add 192.168.1.101/24 dev veth1
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1 ip link set dev veth1 up
vagrant@ubuntu-xenial:~$ sudo ip link add br0 type bridge
vagrant@ubuntu-xenial:~$ sudo ip link set dev br0 up
vagrant@ubuntu-xenial:~$ sudo ip link set enp0s3 master br0
vagrant@ubuntu-xenial:~$ sudo ip link set veth0 master br0
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1  ip route add default via 192.168.1.100
```