# Syncer
Syncer is a tool to sync a struct between two programs over the network using GRPC. 
## WARNING
Do not use this project directly at the moment, the api is changing and all the needed tests have not been created.

## Endpoint
Each struct that you want to sync needs an endpoint. Best way to sync multiple structs is to use a struct to combine the
two for syncing. Struct tag `extractor:"-"` will block the field from being synced. All types besides interfaces can be
created on the other side. As long as interfaces exist on the other side the data within will be updated but unless it
is defined beforehand then the injector has no knowledge of the base type and cannot create it. 

```go
package main

import (
	"github.com/kjbreil/syncer/endpoint"
	"github.com/kjbreil/syncer/endpoint/settings"
)

package main

import (
    "github.com/kjbreil/syncer/endpoint"
    "github.com/kjbreil/syncer/endpoint/settings"
    "net"
)

type dataOne struct {
    String string
    wontSync string
    WontSync string  `extractor:"-"`
}

type dataTwo struct {
    Int int
}


type synced struct {
    DataOne *dataOne
    DataTwo *dataTwo
}

func main() {
    s := synced{
        DataOne: &dataOne{
            String: "String1",
            wontSync: "String2",
            WontSync: "String3",
        },
        DataTwo: &dataTwo{
            Int: 1,
        },
    }
    settings := settings.Settings{
        Port:       45012,
        Peers:       []net.TCPAddr{{
            IP:   net.ParseIP("10.0.2.2"),
            Port: 45012,
        }},
        AutoUpdate: true,
    }
    ep, err := endpoint.New(&s, &settings)
    if err!= nil {
        panic(err)
    }
    ep.Run(false)
    
    ...
}



```