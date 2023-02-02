# redisdel

Scan for redis keys matching a given pattern and delete them.

Steps - 
1. set following in `redisdel.conf` - 
```
{
    "ServerConfigs": {
        "Address": "<redis-host>:<port>"
    },
    "Password": "<password>"
}
```

2. run
```
go build
./redisdel "key-pattern-to-delete"
```
