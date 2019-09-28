# md5tree

Usage:

To scan the current directory:

```
md5tree .
```


To save the hashes:

```
md5tree -l ~/data/hashes.ldb .
```


To see if the hashes have changed since the last scan:

```
md5tree -l ~/data/hashes.ldb -c .
```

You could also use the environment variable MD5TREE_DATADIR to specify where
the leveldb data directory is located.

Specifying a different location with `-l` will overried the `MD5TREE_DATADIR`.

```
export MD5TREE_DATADIR=${HOME}/data/md5tree.ldb
md5tree -c .
```

## Notes

- Ignores symlinks, sockets, devices, named pipes
- Console output attempts to maintain compatibility with md5sum; please submit a github issue or pull request if you run accross an issue.

