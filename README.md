# README
A prometheus exporter to show version part of filename.

```
[0]$ ./rpm-exporter -name glibc -name openssl
2019/07/11 20:55:52 glibc: 2.26-32.amzn2.0.1.x86_64
2019/07/11 20:55:52 openssl: 1.0.2k-16.amzn2.1.1.x86_64

[1]$ curl -s localhost:9872/metrics | grep ^rpm
rpm_info{rpm_name="glibc",version="2.26-32.amzn2.0.1.x86_64"} 1
rpm_info{rpm_name="openssl",version="1.0.2k-16.amzn2.1.1.x86_64"} 1
```