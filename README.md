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


