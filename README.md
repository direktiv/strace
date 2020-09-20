# strace

[![Build Status](https://travis-ci.org/vorteil/strace.svg?branch=master)](https://travis-ci.org/vorteil/strace) [![Maintainability](https://api.codeclimate.com/v1/badges/28304b790bcb6763c9c2/maintainability)](https://codeclimate.com/github/vorteil/strace/maintainability) <a href="https://codeclimate.com/github/vorteil/strace/test_coverage"><img src="https://api.codeclimate.com/v1/badges/28304b790bcb6763c9c2/test_coverage" /></a> [![Go Report Card](https://goreportcard.com/badge/github.com/vorteil/strace)](https://goreportcard.com/report/github.com/vorteil/strace)

strace is a simple strace-ing tool for [vorteil.io](http://www.vorteil.io).

To enable strace on an app on [vorteil.io](http://www.vorteil.io) micro vm it needs to be enabled in the program configuration:

```toml
[[program]]
  binary = "/app"
  args = "-arg1 -arg2"
  strace = true
```

This configuration enables trace on an application and prints it to stdout:

```sh
brk(0) = 94819228717056
arch_prctl(0x3001, 140730692922624, 0x7fb0d9c47230) = -1 invalid argument
access("/etc/ld.so.preload", 4) = -1 no such file or directory
openat(4294967196, "/etc/ld.so.cache", 0x80000, 0) = 3
fstat(3, 0x7ffe6af5df00) = 0
mmap(0x0, 105369, 1, 0x2, 3, 0) = 140397544280064
close(3) = 0
openat(4294967196, "/lib/x86_64-linux-gnu/libselinux.so.1", 0x80000, 0) = 3
read(3, "ELF", 832) = 832
fstat(3, 0x7ffe6af5df50) = 0
mmap(0x0, 8192, 3, 0x22, 4294967295, 0) = 140397544271872
mmap(0x0, 174600, 1, 0x802, 3, 0) = 140397544095744
mprotect(140397544120320, 135168, 0) = 0
mmap(0x7fb0d9bea000, 102400, 5, 0x812, 3, 24576) = 140397544120320
mmap(0x7fb0d9c03000, 28672, 1, 0x812, 3, 126976) = 140397544222720
mmap(0x7fb0d9c0b000, 8192, 3, 0x812, 3, 155648) = 140397544255488
mmap(0x7fb0d9c0d000, 6664, 3, 0x32, 4294967295, 0) = 140397544263680
close(3) = 0

```
