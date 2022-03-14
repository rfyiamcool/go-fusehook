## go-fusehook

**feature**

- control disk's block io ratelimit
- inject io error for custom path

`pkg/hookfs` copy from `https://github.com/osrg/hookfs`

### dep

the tool base on `fusermount3`.

```bash
yum -y install fuse3
```

[go-fuse mount source code](https://github.com/hanwen/go-fuse/blob/934a183ed91446d218b5471c4df9f93db039f6e1/fuse/mount_linux.go)

### Run Example

mount path, control disk's block io rate.

```
$ go build -o blkio cmd/blkio/main.go
$ ./blkio --write-rate=100kb/1s --read-rate=100kb/1s "/mnt/hookfs" "/original"
^C
```

umount target path.

```
$ fusermount -u "/mnt/hookfs"
```

### API Usage

...
