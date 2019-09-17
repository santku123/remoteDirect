package action 

import (
        "fmt"
        "os"
        "bytes"
        "encoding/json"
        "net/http"
        "context"
        "strings"
        "github.com/spf13/cobra"
        "opendev.org/airship/airshipctl/pkg/environment"
        redfish "github.com/nordix-airship/go-redfish/client"
)

// declare global varibles here
var actionType string // action for Command Reset
var transportType string // transport type used for communication http/https
var libType string // transport type used for communication http/https
var endpoint string // end point with port
var v1RequestFormat = new(redfishV1Format) // Request format
var client http.Client
var hostid string // instance id for the host


// Redfish Version V1 response structure defined
type redfishV1Format struct {
        OdataType         string `json:"@odata.type"`
        Name              string `json:"Name"`
        MembersOdataCount int    `json:"Members@odata.count"`
        Members           []struct {
                OdataID string `json:"@odata.id"`
        } `json:"Members"`
        OdataContext     string `json:"@odata.context"`
        OdataID          string `json:"@odata.id"`
        RedfishCopyright string `json:"@Redfish.Copyright"`
}

// Decode the JSON response from the target
func decodeJsonResponse(url string, target interface{}) error {
        r, err := client.Get(url)
        if err != nil {
                return err
        }
        defer r.Body.Close()

        return json.NewDecoder(r.Body).Decode(target)
}

// Build Action
func NewRemoteDirectTargetAction(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
        remoteTargetaction := &cobra.Command{
                Use:   "action",
                Short: "airshppctl remotedirect target action <On/Forceoff/force> endpoint <ip:port> transport <http/https> libtype <redfish/go-redfish>",
                Run: func(cmd *cobra.Command, args []string) {
                        actionName, _ := cmd.Flags().GetString("action")
                        if actionName == "" {
                                actionName = "action"
                                fmt.Println("You did not specify the action, please specify action")
                                os.Exit(1)

                        }
                        endpointName, _ := cmd.Flags().GetString("endpoint")
                        if endpointName == "" {
                                actionName = "endpoint"
                                fmt.Println("You did not specify the endpoint, please specify endpoint")
                                os.Exit(1)

                        }
                        
                        transportName, _ := cmd.Flags().GetString("transport")
                        if transportName == "" {
                                actionName = "transport"
                                fmt.Println("You did not specify the transport, please specify transport")
                                os.Exit(1)

                        }
                        libName, _ := cmd.Flags().GetString("libtype")
                        if libName == "" {
                                libName = "redfish"
                                fmt.Println("You did not specify the libType, please specify libType")
                                os.Exit(1)
                        }
                        endpoint = endpointName
                        actionType = actionName
                        transportType = transportName
                        if libName == "redfish" {
                          RunRedfish()
                          fmt.Println("1")
                        } else {
                          fmt.Println("2")
                          RunGoRedfish()
                        }     
                },
        }
        remoteTargetaction.Flags().StringP("endpoint", "T", "", "",)
        remoteTargetaction.PersistentFlags().StringVar(&actionType, "action", "On", "ForceOff or On")
        remoteTargetaction.PersistentFlags().StringVar(&transportType, "transport", "http", "https")
        remoteTargetaction.PersistentFlags().StringVar(&libType, "libtype", "redfish", "go-redfish")
        return remoteTargetaction
}

// Construct JSON request and send via HTTP request and look for failures if not not close the session
func RunGoRedfish() {
        cfg := &redfish.Configuration{
                  BasePath:      transportType+"://"+endpoint,
                  DefaultHeader: make(map[string]string),
                  UserAgent:     "go-redfish/client",
        }

        redfishApi := redfish.NewAPIClient(cfg).DefaultApi

        // Retrieve System ID(s) from avaiable systems
        s1, _, _ := redfishApi.ListSystems(context.Background())
        s := strings.Split(s1.Members[0].OdataId, "/")
        system_id := s[4] // 4th string is the systemID
        system, _, _ := redfishApi.GetSystem(context.Background(), system_id)
        fmt.Println(system)
}

func RunRedfish() {
        type Payload struct {
                ResetType string `json:"ResetType"`
        }

        data := Payload{
                ResetType: actionType,
        }
        payloadBytes, err := json.Marshal(data)
        if err != nil {
                // handle err
        }
        body := bytes.NewReader(payloadBytes)

        decodeJsonResponse(transportType+"://"+endpoint+"/redfish/v1/Systems", v1RequestFormat)
        p := &hostid
        *p = v1RequestFormat.Members[0].OdataID
        url1 := transportType+"://"+endpoint+hostid+"/Actions/ComputerSystem.Reset/"
        req, err := http.NewRequest("POST", url1, body)
        if err != nil {
                // handle err
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        if err != nil {
             fmt.Println("Error ")
             fmt.Println(err)
        }
        fmt.Println("Response" , resp )
        defer resp.Body.Close()
}


