package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Define the User structure, with 8 attributes.
type User struct {
	UID         string `json:"uid"`
	UNID        string `json:"unid"`
	UPubKey     string `json:"upubkey"`
	UserLevel   string `json:"userlevel"`
	ASLevel     string `json:"aslevel"`
	UserZone    string `json:"userzone"`
	validity    string `json:"validity"`
	UTrustLevel int    `json:"utrustlevel`
	UStatus     string `json:"status"`
}

// Define the Device structure, with 9 attributes.
type Device struct {
	DID         string `json:"did"`
	DNID        string `json:"dnid"`
	DPubKey     string `json:"dpubkey"`
	DType       string `json:""dtype`
	SLevel      string `json:"slevel"`
	Dzone       string `json:"Dzone"`
	TATimeStart int    `json:"tatimestart"`
	TATimeEnd   int    `json:"tatimeend"`
	DTrustLevel int    `json:"dtrustlevel"`
	DStatus     string `json:"dstatus"`
}

//Define the Request Structure, With 6 Properties
type Request struct {
	RType       string `json:"type"`
	ActionType  string `json:"actiontype"`
	RequesterID string `json:"rid"`
	DeviceID    string `json:"did"`
	Time        int    `json:"time"`
	Permission  string `json:"permission"`
}

// define 3 counter to keep track of the user, device, and request number.
type EntryCounter struct {
	UserCount    int `json:"usercount"`
	DeviceCount  int `json:"devicecount"`
	RequestCount int `json:"requestcount"`
}

type QueryResultU struct {
	Key    string `json:"Key"`
	Record *User
}

type QueryResultD struct {
	Key    string `json:"Key"`
	Record *Device
}

type QueryResultR struct {
	Key    string `json:"Key"`
	Record *Request
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []User{
		User{UID: "User0", UNID: "127.2.2.0", UPubKey: "sjhohsov23hjv23", UserLevel: "Admin", ASLevel: "High", UserZone: "A", validity: "N/A", UTrustLevel: 100, UStatus: "Active"},
		User{UID: "User1", UNID: "127.8.8.0", UPubKey: "00fnjmso3hiol2p", UserLevel: "Admin", ASLevel: "High", UserZone: "B", validity: "N/A", UTrustLevel: 100, UStatus: "Active"},
	}

	i := 0
	for i < len(users) {
		userAsBytes, _ := json.Marshal(users[i])
		err := ctx.GetStub().PutState("User"+strconv.Itoa(i), userAsBytes)
		fmt.Println("Added", users[i])
		i = i + 1

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// addUser fuction Register user.
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, id string, uid string, unid string, upubkey string, userlevel string, aslevel string, userzone string, validity string, utrustlevel int, status string) error {

	b := adminOrNot(id)

	if b != true {
		return fmt.Errorf("Only the Administrator can add a User")
	}
	var user = User{UID: uid, UNID: unid, UPubKey: upubkey, UserLevel: userlevel, ASLevel: aslevel, UserZone: userzone, validity: validity, UTrustLevel: utrustlevel, UStatus: status}

	/////wrinting to ledger//////////////
	userAsBytes, _ := json.Marshal(user)
	return ctx.GetStub().PutState(uid, userAsBytes)
	/////////////////
}

func (s *SmartContract) AddDevice(ctx contractapi.TransactionContextInterface, id string, did string, dnid string, dpubkey string, dtype string, slevel string, dzone string, tatimestart int, tatimeend int, dtrustlevel int, dstatus string) error {

	b := adminOrNot(id)

	if b != true {
		return fmt.Errorf("Only the Administrator can add a Device")
	}
	var device = Device{DID: did, DNID: dnid, DPubKey: dpubkey, DType: dtype, SLevel: slevel, Dzone: dzone, TATimeStart: tatimestart, TATimeEnd: tatimeend, DTrustLevel: dtrustlevel, DStatus: dstatus}

	deviceAsBytes, _ := json.Marshal(device)
	return ctx.GetStub().PutState(did, deviceAsBytes)

}

////Update user function is resopnsible for update user's attributes.

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, uid string, userlevel string, aslevel string, userzone string, validity string) error {

	b := adminOrNot(id)

	if b != true {
		return fmt.Errorf("Only the Administrator can Update a User")
	}
	user, err := s.QueryUser(ctx, uid)

	if err != nil {
		return err
	}

	user.UserLevel = userlevel
	user.ASLevel = aslevel
	user.UserZone = userzone
	user.validity = validity

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(uid, userAsBytes)
}

//Update Device function is resopnsible for update device's attributes.

func (s *SmartContract) UpdateDevice(ctx contractapi.TransactionContextInterface, id string, did string, slevel string, dzone string, tatimestart int, tatimeend int) error {

	b := adminOrNot(id)

	if b != true {
		return fmt.Errorf("Only the Administrator can Update a Device")
	}
	device, err := s.QueryDevice(ctx, did)

	if err != nil {
		return err
	}
	device.SLevel = slevel
	device.Dzone = dzone
	device.TATimeStart = tatimestart
	device.TATimeEnd = tatimeend

	deviceAsBytes, _ := json.Marshal(device)

	return ctx.GetStub().PutState(did, deviceAsBytes)
}

////delete user update the user status field to deactive which refer the user is deleted from the system

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string, uid string, userstatus string) error {

	b := adminOrNot(id)

	if b != true {
		return fmt.Errorf("Only the Administrator can Delete a User")
	}
	user, err := s.QueryUser(ctx, uid)

	if err != nil {
		return err
	}

	user.UStatus = userstatus

	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(uid, userAsBytes)
}

///delete function update the status field to deactive which refer the device is deleted from the system

func (s *SmartContract) DeleteDevice(ctx contractapi.TransactionContextInterface, id string, did string, dstatus string) error {

	b := adminOrNot(id)

	if b != true {
		return fmt.Errorf("Only the Administrator can Delete a Device")
	}
	device, err := s.QueryDevice(ctx, did)

	if err != nil {
		return err
	}
	device.DStatus = dstatus

	deviceAsBytes, _ := json.Marshal(device)

	return ctx.GetStub().PutState(did, deviceAsBytes)
}

//accessRequest function is the access control policy
func (s *SmartContract) AccessRequestVerifier(ctx contractapi.TransactionContextInterface, args0 string, args1 string, args2 string, args3 string, args4 string, args5 string) error {

	if args1 == "U2D" {
		userAsBytes, _ := ctx.GetStub().GetState(args3)
		user := User{}
		json.Unmarshal(userAsBytes, &user)

		deviceAsBytes, _ := ctx.GetStub().GetState(args4)
		device := Device{}
		json.Unmarshal(deviceAsBytes, &device)

		if args2 == "Read" && (device.DType == "Sensor" || device.DType == "Both") {
			if user.UserLevel == "Admin" && user.UNID == device.DNID {

				var tm, _ = strconv.Atoi(args5)
				var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
				requestAsBytes, _ := json.Marshal(request)
				ctx.GetStub().PutState(args0, requestAsBytes)

				return nil
			} else if user.UserLevel == "Guest" && user.validity == "Not valid" {

				var tm, _ = strconv.Atoi(args5)
				var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
				requestAsBytes, _ := json.Marshal(request)
				ctx.GetStub().PutState(args0, requestAsBytes)

				return nil
			} else {
				if user.UserZone == device.Dzone {

					var tm, _ = strconv.Atoi(args5)
					var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
					requestAsBytes, _ := json.Marshal(request)
					ctx.GetStub().PutState(args0, requestAsBytes)

					return nil
				} else {
					var time, _ = strconv.Atoi(args5)
					if (device.TATimeStart <= time) && (time <= device.TATimeEnd) {

						var tm, _ = strconv.Atoi(args4)
						var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
						requestAsBytes, _ := json.Marshal(request)
						ctx.GetStub().PutState(args0, requestAsBytes)

						return nil
					} else {

						var tm, _ = strconv.Atoi(args5)
						var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
						requestAsBytes, _ := json.Marshal(request)
						ctx.GetStub().PutState(args0, requestAsBytes)

						return nil
					}
				}

			}

		} else if args2 == "Action" && (device.DType == "Actuator" || device.DType == "Both") {
			if user.UserLevel == "Admin" && user.UNID == device.DNID {

				var tm, _ = strconv.Atoi(args5)
				var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
				requestAsBytes, _ := json.Marshal(request)
				ctx.GetStub().PutState(args0, requestAsBytes)

				return nil
			} else if user.UserLevel == "Guest" && user.validity == "Not valid" {
				var tm, _ = strconv.Atoi(args5)
				var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
				requestAsBytes, _ := json.Marshal(request)
				ctx.GetStub().PutState(args0, requestAsBytes)

				return nil
			} else {
				if user.UserZone == device.Dzone {

					var tm, _ = strconv.Atoi(args5)
					var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
					requestAsBytes, _ := json.Marshal(request)
					ctx.GetStub().PutState(args0, requestAsBytes)

					return nil
				} else {
					var time, _ = strconv.Atoi(args5)
					if (device.TATimeStart <= time) && (time <= device.TATimeEnd) {
						if device.SLevel == user.ASLevel {

							var tm, _ = strconv.Atoi(args5)
							var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
							requestAsBytes, _ := json.Marshal(request)
							ctx.GetStub().PutState(args0, requestAsBytes)

							return nil
						}
					} else {

						var tm, _ = strconv.Atoi(args5)
						var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
						requestAsBytes, _ := json.Marshal(request)
						ctx.GetStub().PutState(args0, requestAsBytes)

						return nil
					}
				}

			}

		} else {

			var tm, _ = strconv.Atoi(args5)
			var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
			requestAsBytes, _ := json.Marshal(request)
			ctx.GetStub().PutState(args0, requestAsBytes)

			return nil
		}

	} else if args1 == "D2D" {
		rdeviceAsBytes, _ := ctx.GetStub().GetState(args2)
		rdevice := Device{}
		json.Unmarshal(rdeviceAsBytes, &rdevice)

		deviceAsBytes, _ := ctx.GetStub().GetState(args3)
		device := Device{}
		json.Unmarshal(deviceAsBytes, &device)

		if (args2 == "Read" && (device.DType == "Sensor" || device.DType == "Both")) || (args2 == "Action" && (device.DType == "Actuator" || device.DType == "Both")) {
			if rdevice.DNID == device.DNID {
				if rdevice.Dzone == device.Dzone {
					var tm, _ = strconv.Atoi(args5)
					var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
					requestAsBytes, _ := json.Marshal(request)
					ctx.GetStub().PutState(args0, requestAsBytes)

					return nil

				} else {
					var time, _ = strconv.Atoi(args5)
					if (device.TATimeStart <= time) && (time <= device.TATimeEnd) {
						if rdevice.SLevel == device.SLevel {

							var tm, _ = strconv.Atoi(args5)
							var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "ALLOW"}
							requestAsBytes, _ := json.Marshal(request)
							ctx.GetStub().PutState(args0, requestAsBytes)

							return nil

						}
					} else {

						var tm, _ = strconv.Atoi(args5)
						var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
						requestAsBytes, _ := json.Marshal(request)
						ctx.GetStub().PutState(args0, requestAsBytes)

						return nil

					}
				}
			} else {

				var tm, _ = strconv.Atoi(args5)
				var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
				requestAsBytes, _ := json.Marshal(request)
				ctx.GetStub().PutState(args0, requestAsBytes)

				return nil
			}
		} else {

			var tm, _ = strconv.Atoi(args4)
			var request = Request{RType: args1, ActionType: args2, RequesterID: args3, DeviceID: args4, Time: tm, Permission: "DENY"}
			requestAsBytes, _ := json.Marshal(request)
			ctx.GetStub().PutState(args0, requestAsBytes)

			return nil
		}

	}
	return fmt.Errorf("Error!!!! Unsupported Request")
}

//end  accessRequest function

//update Trust level function
func (s *SmartContract) TrustLevelUpdater(ctx contractapi.TransactionContextInterface, args0 string, args1 string, args2 string) error {

	requestAsBytes, _ := ctx.GetStub().GetState(args0)
	request := Request{}
	json.Unmarshal(requestAsBytes, &request)
	var RID = request.RequesterID
	var DID = request.DeviceID

	userAsBytes, _ := ctx.GetStub().GetState(RID)
	user := User{}
	json.Unmarshal(userAsBytes, &user)
	var um = user.UTrustLevel

	deviceAsBytes, _ := ctx.GetStub().GetState(DID)
	device := Device{}
	json.Unmarshal(deviceAsBytes, &device)
	var dm = device.DTrustLevel
	if args1 == "Satisfactory" && um < 100 {
		um = um + 1
		user.UTrustLevel = um
		userAsBytes, _ = json.Marshal(user)
		ctx.GetStub().PutState(RID, userAsBytes)
	} else {
		um = um - 1

		/*if um <= 0 {
			var ip[] string = {RID}
			s.deleteUser(APIstub, ip)
		}*/
	}
	if args2 == "Satisfactory" && dm < 100 {
		dm = dm + 1
		device.DTrustLevel = dm
		deviceAsBytes, _ = json.Marshal(device)
		ctx.GetStub().PutState(args1, deviceAsBytes)
	} else {
		dm = dm - 1
		device.DTrustLevel = dm
		deviceAsBytes, _ = json.Marshal(device)
		ctx.GetStub().PutState(args1, deviceAsBytes)
		/*if dm <= 0 {
			var ip[] string = {DID}
			s.deleteDevice(APIstub, ip)

		}*/
	}

	userAsBytes, _ = json.Marshal(user)
	ctx.GetStub().PutState(args1, userAsBytes)

	return nil
}

//

//query function

//query user
func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, uid string) (*User, error) {

	userAsBytes, err := ctx.GetStub().GetState(uid)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", uid)
	}

	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

//query Device
func (s *SmartContract) QueryDevice(ctx contractapi.TransactionContextInterface, did string) (*Device, error) {

	deviceAsBytes, err := ctx.GetStub().GetState(did)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if deviceAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", did)
	}

	device := new(Device)
	_ = json.Unmarshal(deviceAsBytes, device)

	return device, nil
}

//query Request
func (s *SmartContract) QueryAccessRequest(ctx contractapi.TransactionContextInterface, rid string) (*Request, error) {

	requestAsBytes, err := ctx.GetStub().GetState(rid)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if requestAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", rid)
	}

	request := new(Request)
	_ = json.Unmarshal(requestAsBytes, request)

	return request, nil
}

// QueryAllUser returns all User found in world state
func (s *SmartContract) QueryAllUser(ctx contractapi.TransactionContextInterface) ([]QueryResultU, error) {
	startKey := "User0"
	endKey := "User99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultU{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(User)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResultU := QueryResultU{Key: queryResponse.Key, Record: user}
		results = append(results, queryResultU)
	}

	return results, nil
}

// QueryAllDevice returns all devices found in world state
func (s *SmartContract) QueryAllDevice(ctx contractapi.TransactionContextInterface) ([]QueryResultD, error) {
	startKey := "Device0"
	endKey := "Device99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultD{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		device := new(Device)
		_ = json.Unmarshal(queryResponse.Value, device)

		queryResultD := QueryResultD{Key: queryResponse.Key, Record: device}
		results = append(results, queryResultD)
	}

	return results, nil
}

// QueryAllAccessRequest returns all Request found in world state
func (s *SmartContract) QueryAllAccessRequest(ctx contractapi.TransactionContextInterface) ([]QueryResultR, error) {
	startKey := "Request0"
	endKey := "Request99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResultR{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		request := new(Request)
		_ = json.Unmarshal(queryResponse.Value, request)

		queryResultR := QueryResultR{Key: queryResponse.Key, Record: request}
		results = append(results, queryResultR)
	}

	return results, nil
}

func adminOrNot(args string) bool {

	/*userAsBytes, _ := APIstub.GetState(args)
	user := User{}
	json.Unmarshal(userAsBytes, &user)
	if user.UserLevel == "Admin" {
		return true
	} else {
		return false
	}*/
	return true
}

func (s *SmartContract) QueryPermission(ctx contractapi.TransactionContextInterface, rid string) (string, error) {

	requestAsBytes, err := ctx.GetStub().GetState(rid)
	var rse = "DENY"
	if err != nil {
		return rse, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if requestAsBytes == nil {
		return rse, fmt.Errorf("%s does not exist", rid)
	}

	request := Request{}

	_ = json.Unmarshal(requestAsBytes, &request)
	var p = request.Permission
	/*request := new(Request)
	_ = json.Unmarshal(requestAsBytes, request)
	*/
	return p, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
