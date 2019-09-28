# md5tree

## Why?

I wrote this because I needed a tool to save large directory trees of 
md5 hashes to periodically verify the data integrity of archived files.
Backups are great, but bitrot is the essence of evil.  After searching, 
I found that nobody really seemed to have solved this problem - at least 
in a way that fits my needs - so I decided to write something and make 
it available to the world.  I hope at least one other person finds that 
it makes their lives easier.

Pull requests welcome.

## Usage

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

If `-c` is specified, the exit code of md5tree will be 2 if a hash fails.
Otherwise it will exit with 0.

## Notes

- Ignores symlinks, sockets, devices, named pipes
- Console output attempts to maintain compatibility with md5sum; please submit a github issue or pull request if you run accross an issue.

