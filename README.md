## go-fusehook

**feature**

- control disk's block io rate
- inject io error for custom path
- io delay
- more 

![go-fusehook](./docs/fuse.jpg)

some code from `https://github.com/osrg/hookfs`

### dep

`go-fusehook` is based on `fusermount3` package, so you need to install `fuse3`.

```bash
yum -y install fuse3
```

[go-fuse mount source code](https://github.com/hanwen/go-fuse/blob/934a183ed91446d218b5471c4df9f93db039f6e1/fuse/mount_linux.go)

### Run Example

mount path, control disk's block io rate.

```
$ mkdir -p /mnt/hookfs-test
$ mkdir -p /data/source
$ go build -o blkio cmd/blkio/main.go
$ ./blkio --write-rate=100kb/1s --read-rate=100kb/1s "/mnt/hookfs-test" "/data/source"
^C
```

umount target path.

```
$ fusermount -u "/mnt/hookfs-test"
```

### API Usage

...
