/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	
	"encoding/json"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("alcsLogger")

// SimpleChaincode example simple Chaincode implementation
type ALCSTicketsChaincode struct {
	
}

type ticketDetail struct {
	tkt_id 				string `json:"tkt_id"`
	occ_tkt_num 		string `json:"occ_tkt_num"` 		//onecallcenter_ticketnum
	excavator_name 		string `json:"excavator_name"`		
	excavator_ph_num 	string `json:"excavator_ph_num"`
	excavator_email 	string `json:"excavator_email"`
	work_for			string `json:"work_for"`
	dig_date			string `json:"dig_date"`
	street_address		string `json:"street_address"`
	city				string `json:"city"`
	state				string `json:"state"`
	mem_code			string `json:"mem_code"`
	type_of_tkt			string `json:"type_of_tkt"`
	tkt_priority		string `json:"tkt_priority"`
	hipr_info			string `json:"hipr_info"`
	latitude			string `json:"latitude"`
	longitude			string `json:"longitude"`
	tkt_remarks			string `json:"tkt_remarks"`
	explosives			bool   `json:"explosives"`
	action_code			string `json:"action_code"`
	vendor_remarks		string `json:"vendor_remarks"`
	ack				    bool   `json:"ack"`
}

//chaincode Init function
func (t *ALCSTicketsChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	return nil, nil
}
 
//chaincode Query function 
func (t *ALCSTicketsChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//chaincode Invoke function 
func (t *ALCSTicketsChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "UpdateVendorDetails" {
		return t.UpdateVendorDetails(stub, args)
	} else if function == "UpdateATTTktDetails" {
		return t.UpdateATTTktDetails(stub, args)
	} else if function == "UpdateTKTAck" {
		return t.UpdateTKTAck(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func (t *ALCSTicketsChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var tkt_id, jsonvalue string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	tkt_id = args[0] //rename for funsies
	jsonvalue = args[1]
	err = stub.PutState(tkt_id, []byte(jsonvalue)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *ALCSTicketsChaincode) UpdateVendorDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//check args 1. tkt_id, 2. action_code, 3. vendor_remarks
	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. tkt_id, action_code and remarks")
	}
	
	var tkt_id, action_code, vendor_remarks string
	tkt_id = args[0]
	action_code = args[1]
	vendor_remarks = args[2]
	
	laBytes, err := stub.GetState(tkt_id)
	if err != nil {
		logger.Error("Could not fetch loan application from ledger", err)
		return nil, err
	}
	
	var tktDetail ticketDetail
	err = json.Unmarshal(laBytes, &tktDetail)
	tktDetail.action_code = action_code
	tktDetail.vendor_remarks = vendor_remarks
	
	laBytes, err = json.Marshal(&tktDetail)
	if err != nil {
		logger.Error("Could not marshal loan application post update", err)
		return nil, err
	}

	err = stub.PutState(tkt_id, laBytes)
	if err != nil {
		logger.Error("Could not save loan application post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'vendorUpdate', description:" + tkt_id + "' Successfully updated action_code and remarks'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated Vendor data")
	return nil, nil
}

func (t *ALCSTicketsChaincode) UpdateATTTktDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//check args 1. tkt_id, 2. hipr_info, 3. latitude, 4. longitude
	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. tkt_id, hipr_info, latitude, longitude")
	}
	
	var tkt_id, hipr_info, latitude, longitude string
	tkt_id = args[0]
	hipr_info = args[1]
	latitude = args[2]
	longitude = args[3]
	
	laBytes, err := stub.GetState(tkt_id)
	if err != nil {
		logger.Error("Could not fetch loan application from ledger", err)
		return nil, err
	}
	
	var tktDetail ticketDetail
	err = json.Unmarshal(laBytes, &tktDetail)
	tktDetail.hipr_info = hipr_info
	tktDetail.latitude = latitude
	tktDetail.longitude = longitude
	
	laBytes, err = json.Marshal(&tktDetail)
	if err != nil {
		logger.Error("Could not marshal loan application post update", err)
		return nil, err
	}

	err = stub.PutState(tkt_id, laBytes)
	if err != nil {
		logger.Error("Could not save loan application post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'ATTTicketUpdate', description:" + tkt_id + "' Successfully updated hipr_info, latitude and longitude'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated ATT ticket data")
	return nil, nil
}

func (t *ALCSTicketsChaincode) UpdateTKTAck(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//check args 1. tkt_id, 2. ack
	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. tkt_id, acknowledgement")
	}
	
	var tkt_id string
	var ack bool
	tkt_id = args[0]
	ack, err := strconv.ParseBool(args[1])
	//ack = bool(args[1])

	
	laBytes, err := stub.GetState(tkt_id)
	if err != nil {
		logger.Error("Could not fetch loan application from ledger", err)
		return nil, err
	}
	
	var tktDetail ticketDetail
	err = json.Unmarshal(laBytes, &tktDetail)
	tktDetail.ack = ack
	
	laBytes, err = json.Marshal(&tktDetail)
	if err != nil {
		logger.Error("Could not marshal loan application post update", err)
		return nil, err
	}

	err = stub.PutState(tkt_id, laBytes)
	if err != nil {
		logger.Error("Could not save loan application post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'ATTUpdateACK', description:" + tkt_id + "' Successfully updated acknowledgement'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated ATT ticket acknowledgement")
	return nil, nil
}

// read - query function to read key/value pair
func (t *ALCSTicketsChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
 
func main() {
    err := shim.Start(new(ALCSTicketsChaincode))
    if err != nil {
        fmt.Println("Could not start ALCSTicketsChaincode")
    } else {
        fmt.Println("ALCSTicketsChaincode successfully started")
    }
 
}