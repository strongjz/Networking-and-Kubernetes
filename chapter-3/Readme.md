
# Chapter 3 Container Networking Intro

The following steps show how to create the networking setup.

1. Create a host with a root network namespace.
2. Create two new network namespace.
3. Create two veth pair.
4. Move one side of each veth pair into a new network namespace.
5. Address side of the veth pair inside the new network namespace.
6. Create a bridge interface.
7. Attach bridge to the host interface.
8. Attach one side of each veth pair to the bridge interface.
9. Test

```bash
      br-veth0          veth0  +-------------+
         +--------------------- + net0       |
         |       192.168.1.100 +-------------+
+--------+
|        |
| br1    | 192.168.1.10
|        |
+--------+
         |              veth1  +-------------+
         +---------------------+ net1        |
      br-veth1   192.168.1.101 +-------------+
```

## 1. Create a host with a root network namespace.

Follow the steps from Chapter 1 to start a Vagrant Host. 

Connect to the machine 

```bash
vagrant ssh
```
## 2. Create a new network namespace.

```bash
vagrant@ubuntu-xenial:~$ sudo ip netns list
vagrant@ubuntu-xenial:~$ sudo ip netns add net0
vagrant@ubuntu-xenial:~$ sudo ip netns add net1
vagrant@ubuntu-xenial:~$ sudo ip netns list
net1
net0
```

## 3. Create veth pairs.

```bash
vagrant@ubuntu-xenial:~$ sudo ip link add veth0 type veth peer name br-veth0
vagrant@ubuntu-xenial:~$ sudo ip link add veth1 type veth peer name br-veth1
```

```bash
vagrant@ubuntu-xenial:~$ ip link list veth1
7: veth1@br-veth1: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN mode DEFAULT group default qlen 1000
    link/ether 2a:92:85:81:50:50 brd ff:ff:ff:ff:ff:ff
vagrant@ubuntu-xenial:~$ ip link list veth0
5: veth0@br-veth0: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN mode DEFAULT group default qlen 1000
    link/ether 0a:a1:be:3c:89:3d brd ff:ff:ff:ff:ff:ff
```

## 4. Move one side of the veth pair into a new network namespace.

Move veth0 int net0 namespace, and veth1 into net1

```bash
vagrant@ubuntu-xenial:~$ sudo ip link set veth0 netns net0
vagrant@ubuntu-xenial:~$ sudo ip link set veth1 netns net1
```

Examine the network namespaces 

```bash
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1 ip link list veth1
7: veth1@if6: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN mode DEFAULT group default qlen 1000
    link/ether 2a:92:85:81:50:50 brd ff:ff:ff:ff:ff:ff link-netnsid 0

vagrant@ubuntu-xenial:~$ sudo ip netns exec net0 ip link list veth0
5: veth0@if4: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN mode DEFAULT group default qlen 1000
    link/ether 0a:a1:be:3c:89:3d brd ff:ff:ff:ff:ff:ff link-netnsid 0

```

## 5. Address side of the veth pair inside the new network namespace.

Address veths 

```bash
vagrant@ubuntu-xenial:~$ sudo ip netns exec net0 ip addr add "192.168.1.100/24" dev veth0
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1 ip addr add "192.168.1.101/24" dev veth1
```

Turn up the veth side 
```bash
sudo ip netns exec net0 ip link set dev veth0 up
sudo ip netns exec net1 ip link set dev veth1 up
```

```bash
vagrant@ubuntu-xenial:~$ sudo ip netns exec net0 ip link list
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN mode DEFAULT group default qlen 1
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
5: veth0@if4: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state LOWERLAYERDOWN mode DEFAULT group default qlen 1000
    link/ether 0a:a1:be:3c:89:3d brd ff:ff:ff:ff:ff:ff link-netnsid 0
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1 ip link list
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN mode DEFAULT group default qlen 1
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
7: veth1@if6: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state LOWERLAYERDOWN mode DEFAULT group default qlen 1000
    link/ether 2a:92:85:81:50:50 brd ff:ff:ff:ff:ff:ff link-netnsid 0

```

## 6. Create a bridge interface.

Create a bridge interface, turn it on, and address it. 

```bash
sudo ip link add name br1 type bridge
sudo ip link set br1 up
sudo ip addr add 192.168.1.10/24 brd + dev br1
```

## 7. Attach bridge to the host interface.

```bash
sudo ip link set enp0s8 master br1
```

Turn each side of the veth pair that will attach to the bridge, up. 
```bash
sudo ip link set br-veth0 up
sudo ip link set br-veth1 up
```

## 8. Attach one side of the veth pair to the bridge interface.

```bash
sudo ip link set br-veth0 master br1
sudo ip link set br-veth1 master br1
```

## 9. Test. 

From the host ping .10, .100 and .101

```bash
vagrant@ubuntu-xenial:~$ ping 192.168.1.10
PING 192.168.1.10 (192.168.1.10) 56(84) bytes of data.
64 bytes from 192.168.1.10: icmp_seq=1 ttl=64 time=0.011 ms
64 bytes from 192.168.1.10: icmp_seq=2 ttl=64 time=0.038 ms
64 bytes from 192.168.1.10: icmp_seq=3 ttl=64 time=0.038 ms
64 bytes from 192.168.1.10: icmp_seq=4 ttl=64 time=0.021 ms
```

```bash
vagrant@ubuntu-xenial:~$ ping 192.168.1.100
PING 192.168.1.100 (192.168.1.100) 56(84) bytes of data.
64 bytes from 192.168.1.100: icmp_seq=1 ttl=64 time=0.017 ms
64 bytes from 192.168.1.100: icmp_seq=2 ttl=64 time=0.027 ms
64 bytes from 192.168.1.100: icmp_seq=3 ttl=64 time=0.048 ms
```

```bash
vagrant@ubuntu-xenial:~$ ping 192.168.1.101
PING 192.168.1.101 (192.168.1.101) 56(84) bytes of data.
64 bytes from 192.168.1.101: icmp_seq=1 ttl=64 time=0.049 ms
64 bytes from 192.168.1.101: icmp_seq=2 ttl=64 time=0.028 ms
64 bytes from 192.168.1.101: icmp_seq=3 ttl=64 time=0.038 ms
```
Now from the perspective namespaces ping each other. 

```bash
vagrant@ubuntu-xenial:~$ sudo ip netns exec net0 ping -c 4 192.168.1.101
PING 192.168.1.101 (192.168.1.101) 56(84) bytes of data.
64 bytes from 192.168.1.101: icmp_seq=1 ttl=64 time=0.085 ms
64 bytes from 192.168.1.101: icmp_seq=2 ttl=64 time=0.054 ms
64 bytes from 192.168.1.101: icmp_seq=3 ttl=64 time=0.029 ms
64 bytes from 192.168.1.101: icmp_seq=4 ttl=64 time=0.030 ms

--- 192.168.1.101 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3000ms
rtt min/avg/max/mdev = 0.029/0.049/0.085/0.023 ms
```

```bash
vagrant@ubuntu-xenial:~$ sudo ip netns exec net1 ping -c 4 192.168.1.100
PING 192.168.1.100 (192.168.1.100) 56(84) bytes of data.
64 bytes from 192.168.1.100: icmp_seq=1 ttl=64 time=0.020 ms
64 bytes from 192.168.1.100: icmp_seq=2 ttl=64 time=0.045 ms
64 bytes from 192.168.1.100: icmp_seq=3 ttl=64 time=0.029 ms
64 bytes from 192.168.1.100: icmp_seq=4 ttl=64 time=0.034 ms

--- 192.168.1.100 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 2998ms
rtt min/avg/max/mdev = 0.020/0.032/0.045/0.009 ms
```


