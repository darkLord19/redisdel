# redisdel

Scan keys matching a given pattern and delete them.

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
go run main.go "key-pattern-to-delete"
```