package main

import (
        "encoding/json"
        "fmt"
        "log"
        "os"

        "github.com/containernetworking/cni/pkg/invoke"
        "github.com/containernetworking/cni/pkg/skel"
        "github.com/containernetworking/cni/pkg/types"
         "github.com/containernetworking/cni/pkg/version"
)

var logger = log.New(os.Stderr, "", 0)

type NetConf struct {
        types.NetConf
        Delegate map[string]interface{} `json:"delegate"`
}

func loadNetConf(bytes []byte) (*NetConf, error) {
        n := &NetConf{}
        if err := json.Unmarshal(bytes, n); err != nil {
                return nil, fmt.Errorf("failed to load netconf: %v", err)
        }
        return n, nil
}

func isString(i interface{}) bool {
        _, ok := i.(string)
        return ok
}

func checkDelegate(netconf map[string]interface{}) error {
        if netconf["type"] == nil {
                return fmt.Errorf("delegate must have the field 'type'")
        }

        if !isString(netconf["type"]) {
                return fmt.Errorf("delegate field 'type' must be a string")
        }
        return nil
}

func delegateAdd(args string, netconf map[string]interface{}) (bool, error) {
        netconfBytes, err := json.Marshal(netconf)
        if err != nil {
                return true, fmt.Errorf("error serializing: v", err)
        }

        if os.Setenv("CNI_IFNAME", args) != nil {
                return true, fmt.Errorf("error when setting CNI_IFNAME")

        }

        result, err := invoke.DelegateAdd(netconf["type"].(string), netconfBytes)

        if err != nil {
                return true, fmt.Errorf("Error in invoking delegate add: %q - %v", netconf["type"].(string), err)

        }

        return false, result.Print()

}

func vietstack(netconf *NetConf, args string) error {
        // define result with type is error
        //var result error

        if err := checkDelegate(netconf.Delegate); err != nil {
                return fmt.Errorf("Delegate fails: %v", err)
        }

        r, err := delegateAdd(args, netconf.Delegate)
        if r != false && err != nil {
                logger.Printf("failing at: %p - %v", r, err)
                return err
        }

        return err

}

// Adding network interface to the pod
func cmdAdd(args *skel.CmdArgs) error {
        // var result error
        n, err := loadNetConf(args.StdinData)
        if err != nil {
                return err
        }

        result := vietstack(n, args.IfName)
        if result != nil {
                logger.Printf("CmdAdd handler failed: %v", result)
        }
        return result
}

func cmdDel(args *skel.CmdArgs) error {
        return nil
}

func main() {
        skel.PluginMain(cmdAdd, cmdDel, version.All)
}

