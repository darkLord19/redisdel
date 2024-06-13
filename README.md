# redisdel

Scan for redis keys matching a given pattern and delete them.

## Prepare correct `redisdel.conf`

`For single server deployment of redis`:
```
{
    "ServerConfigs": {
        "Address": "<redis-host>:<port>"
    },
    "Password": "<password>"
}
```
`For sentinel deployment of redis`:
```
{
    "Password": "<redis master node password>"
    "Username": "<redis master node username>,
    "MasterName": "<redis sentinel set master name>",
    "Addresses": "<array of addresses belonging to redis sentinel>",
    "Password": "<redis sentinel password>"
}
```

## Build
```
go build
```

## Run
```
go build
./redisdel "key-pattern-to-delete"
```
