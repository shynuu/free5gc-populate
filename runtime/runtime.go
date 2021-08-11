// Package runtime includes all the logic of the tool
package runtime

import (
	"encoding/json"
	"fmt"

	"github.com/free5gc/MongoDBLibrary"
	"github.com/free5gc/openapi/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	authSubsDataColl = "subscriptionData.authenticationData.authenticationSubscription"
	amDataColl       = "subscriptionData.provisionedData.amData"
	smDataColl       = "subscriptionData.provisionedData.smData"
	smfSelDataColl   = "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	amPolicyDataColl = "policyData.ues.amData"
	smPolicyDataColl = "policyData.ues.smData"
	flowRuleDataColl = "policyData.ues.flowRule"
)

func toBsonM(data interface{}) (ret bson.M) {
	tmp, _ := json.Marshal(data)
	json.Unmarshal(tmp, &ret)
	return
}

func generateSubs(ueID string, servingPlmnID string, slices []Slice) *SubsData {
	authSubsData := models.AuthenticationSubscription{
		AuthenticationManagementField: PopulateConfig.AMF,
		AuthenticationMethod:          "5G_AKA", // "5G_AKA", "EAP_AKA_PRIME"
		Milenage: &models.Milenage{
			Op: &models.Op{
				EncryptionAlgorithm: 0,
				EncryptionKey:       0,
				OpValue:             PopulateConfig.OP, // Required
			},
		},
		Opc: &models.Opc{
			EncryptionAlgorithm: 0,
			EncryptionKey:       0,
			OpcValue:            PopulateConfig.OP, // Required
		},
		PermanentKey: &models.PermanentKey{
			EncryptionAlgorithm: 0,
			EncryptionKey:       0,
			PermanentKeyValue:   PopulateConfig.Key, // Required
		},
		SequenceNumber: PopulateConfig.SQN, // Required
	}

	var sliceArray = make([]models.Snssai, len(slices))
	for k, slice := range slices {
		sliceArray[k] = models.Snssai{
			Sd:  slice.Sd,
			Sst: slice.Sst,
		}
	}

	amDataData := models.AccessAndMobilitySubscriptionData{
		Gpsis: []string{
			"msisdn-0900000000",
		},
		Nssai: &models.Nssai{
			DefaultSingleNssais: sliceArray,
			SingleNssais:        sliceArray,
		},
		SubscribedUeAmbr: &models.AmbrRm{
			Downlink: "500 Mbps",
			Uplink:   "500 Mbps",
		},
	}

	var smDataData = make([]models.SessionManagementSubscriptionData, len(slices))
	for k, slice := range slices {
		smDataData[k] = models.SessionManagementSubscriptionData{
			SingleNssai: &models.Snssai{
				Sst: slice.Sst,
				Sd:  slice.Sd,
			},
			DnnConfigurations: map[string]models.DnnConfiguration{
				slice.Dnn: {
					PduSessionTypes: &models.PduSessionTypes{
						DefaultSessionType:  models.PduSessionType_IPV4,
						AllowedSessionTypes: []models.PduSessionType{models.PduSessionType_IPV4},
					},
					SscModes: &models.SscModes{
						DefaultSscMode:  models.SscMode__1,
						AllowedSscModes: []models.SscMode{models.SscMode__1},
					},
					SessionAmbr: &models.Ambr{
						Downlink: "500 Mbps",
						Uplink:   "500 Mbps",
					},
					Var5gQosProfile: &models.SubscribedDefaultQos{
						Var5qi: int32(slice.VarQI),
						Arp: &models.Arp{
							PriorityLevel: 8,
						},
						PriorityLevel: 8,
					},
				},
			},
		}
	}

	var smfSel = make(map[string]models.SnssaiInfo)
	for _, slice := range slices {
		snssai := fmt.Sprintf("%02d%s", slice.Sst, slice.Sd)
		smfSel[snssai] = models.SnssaiInfo{
			DnnInfos: []models.DnnInfo{
				{
					Dnn: slice.Dnn,
				},
			},
		}
	}

	smfSelData := models.SmfSelectionSubscriptionData{
		SubscribedSnssaiInfos: smfSel,
	}

	amPolicyData := models.AmPolicyData{
		SubscCats: []string{
			"free5gc",
		},
	}

	var smPol = make(map[string]models.SmPolicySnssaiData)
	for _, slice := range slices {
		snssai := fmt.Sprintf("%02d%s", slice.Sst, slice.Sd)
		smPol[snssai] = models.SmPolicySnssaiData{
			Snssai: &models.Snssai{
				Sd:  slice.Sd,
				Sst: slice.Sst,
			},
			SmPolicyDnnData: map[string]models.SmPolicyDnnData{
				slice.Dnn: {
					Dnn: slice.Dnn,
				},
			},
		}
	}

	smPolicyData := models.SmPolicyData{
		SmPolicySnssaiData: smPol,
	}

	return &SubsData{
		PlmnID:                            servingPlmnID,
		UeId:                              ueID,
		AuthenticationSubscription:        authSubsData,
		AccessAndMobilitySubscriptionData: amDataData,
		SessionManagementSubscriptionData: smDataData,
		SmfSelectionSubscriptionData:      smfSelData,
		AmPolicyData:                      amPolicyData,
		SmPolicyData:                      smPolicyData,
	}
}

func InsertSubscriber(ueId string, servingPlmnId string, subsData SubsData) {

	filterUeIDOnly := bson.M{"ueId": ueId}
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	authSubsBsonM := toBsonM(subsData.AuthenticationSubscription)
	authSubsBsonM["ueId"] = ueId
	amDataBsonM := toBsonM(subsData.AccessAndMobilitySubscriptionData)
	amDataBsonM["ueId"] = ueId
	amDataBsonM["servingPlmnId"] = servingPlmnId

	smDatasBsonA := make([]interface{}, 0, len(subsData.SessionManagementSubscriptionData))
	for _, smSubsData := range subsData.SessionManagementSubscriptionData {
		smDataBsonM := toBsonM(smSubsData)
		smDataBsonM["ueId"] = ueId
		smDataBsonM["servingPlmnId"] = servingPlmnId
		smDatasBsonA = append(smDatasBsonA, smDataBsonM)
	}

	smfSelSubsBsonM := toBsonM(subsData.SmfSelectionSubscriptionData)
	smfSelSubsBsonM["ueId"] = ueId
	smfSelSubsBsonM["servingPlmnId"] = servingPlmnId
	amPolicyDataBsonM := toBsonM(subsData.AmPolicyData)
	amPolicyDataBsonM["ueId"] = ueId
	smPolicyDataBsonM := toBsonM(subsData.SmPolicyData)
	smPolicyDataBsonM["ueId"] = ueId

	// flowRulesBsonA := make([]interface{}, 0, len(subsData.FlowRules))
	// for _, flowRule := range subsData.FlowRules {
	// 	flowRuleBsonM := toBsonM(flowRule)
	// 	flowRuleBsonM["ueId"] = ueId
	// 	flowRuleBsonM["servingPlmnId"] = servingPlmnId
	// 	flowRulesBsonA = append(flowRulesBsonA, flowRuleBsonM)
	// }

	MongoDBLibrary.RestfulAPIPost(authSubsDataColl, filterUeIDOnly, authSubsBsonM)
	MongoDBLibrary.RestfulAPIPost(amDataColl, filter, amDataBsonM)
	MongoDBLibrary.RestfulAPIPostMany(smDataColl, filter, smDatasBsonA)
	MongoDBLibrary.RestfulAPIPost(smfSelDataColl, filter, smfSelSubsBsonM)
	MongoDBLibrary.RestfulAPIPost(amPolicyDataColl, filterUeIDOnly, amPolicyDataBsonM)
	MongoDBLibrary.RestfulAPIPost(smPolicyDataColl, filterUeIDOnly, smPolicyDataBsonM)
	// MongoDBLibrary.RestfulAPIPostMany(flowRuleDataColl, filter, flowRulesBsonA)

}
