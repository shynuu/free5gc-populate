# Populate free5GC DB with CLI

Populate the DB of free5GC using a CLI (no need for webconsole anymore), useful for container based deployments or non-UI deployments.

## How to use it:

- Clone this repository
- Compile the tool: `go build -o populate -x main.go`
- Run free5GC
- Set the desired parameter in the [config](/config.yaml) file, namely: the mongo URL of the free5gc DB, the MCC, the MNC, the Permanent Subscription Key, the Operator Code (OP or OPC), the Sequence number (SQN), the Authentication Management Field, The slices parameters and the IMSI of the UE which will register on the network
- Run the populate tool: `./populate --config <configuration file>`
- Connect to the mongo DB and check that the following collections were properly modified: `policyData.ues.amData`, `policyData.ues.smData`, `policyData.ues.flowRule`, `subscriptionData.authenticationData.authenticationSubscription`, `subscriptionData.provisionedData.amData`, `subscriptionData.provisionedData.smData`, `subscriptionData.provisionedData.smfSelectionSubscriptionData`

## Configuration file

```yaml
# Mongo Configuration
mongo:
  name: free5gc
  url: mongodb://localhost:27017

# Mobile Country Code value of HPLMN
mcc: "208"
# Mobile Network Code value of HPLMN (2 or 3 digits)
mnc: "93"
# Permanent subscription key
key: 8baf473f2f8fd09487cccbd7097c6862
# Operator code (OP or OPC) of the UE
op: 8e27b6af0e692e750f32667a3b14605d
# Sequence number
sqn: 16f3b3f70fc2
# Authentication Management Field
amf: 8000

# List all the slices
slices:
  - sst: 01
    sd: 010203
    varqi: 9
    dnn: internet
  - sst: 01
    sd: 112233
    varqi: 9
    dnn: internet2
    # PDU Sessions Type for this slice
    # possible values are: "IPV4", "IPV6", "IPV4V6", "UNSTRUCTURED", "ETHERNET"
    # when omitted, results in "IPV4"
    pdu-session-type: "IPV4V6" 

# All IMSI of UEs
imsi:
  - imsi-208930000000001
  - imsi-208930000000002
  - imsi-208930000000003
```

Populate the DB of free5GC using a CLI, useful for container based deployments - Check https://github.com/shynuu/slice-aware-ntn for full description

