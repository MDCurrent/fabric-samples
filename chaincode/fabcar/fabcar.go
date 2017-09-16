/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the record structure, with 4 properties.  Structure tags are used by encoding/json library
type MedicalRecord struct {
	Patient   string `json:"patient"`
	Doctor    string `json:"doctor"`
	Procedure string `json:"procedure"`
	Cost      string `json:"cost"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	return shim.Success(function)
	if function == "queryRecord" {
		return s.queryRecord(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createRecord" {
		return s.createRecord(APIstub, args)
	} else if function == "queryAllRecords" {
		return s.queryAllRecords(APIstub)
	} else if function == "changeRecordPatient" {
		return s.changeRecordPatient(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recordAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(recordAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	records := []MedicalRecord{
		MedicalRecord{Patient: "Toyota", Doctor: "Prius", Procedure: "blue", Cost: "Tomoko"},
		MedicalRecord{Patient: "Ford", Doctor: "Mustang", Procedure: "red", Cost: "Brad"},
		MedicalRecord{Patient: "Hyundai", Doctor: "Tucson", Procedure: "green", Cost: "Jin Soo"},
		MedicalRecord{Patient: "Volkswagen", Doctor: "Passat", Procedure: "yellow", Cost: "Max"},
		MedicalRecord{Patient: "Tesla", Doctor: "S", Procedure: "black", Cost: "Adriana"},
		MedicalRecord{Patient: "Peugeot", Doctor: "205", Procedure: "purple", Cost: "Michel"},
		MedicalRecord{Patient: "Chery", Doctor: "S22L", Procedure: "white", Cost: "Aarav"},
		MedicalRecord{Patient: "Fiat", Doctor: "Punto", Procedure: "violet", Cost: "Pari"},
		MedicalRecord{Patient: "Tata", Doctor: "Nano", Procedure: "indigo", Cost: "Valeria"},
		MedicalRecord{Patient: "Holden", Doctor: "Barina", Procedure: "brown", Cost: "Shotaro"},
	}

	i := 0
	for i < len(records) {
		fmt.Println("i is ", i)
		carAsBytes, _ := json.Marshal(records[i])
		APIstub.PutState("REC"+strconv.Itoa(i), carAsBytes)
		fmt.Println("Added", records[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var record = MedicalRecord{Patient: args[1], Doctor: args[2], Procedure: args[3], Cost: args[4]}

	recordAsBytes, _ := json.Marshal(record)
	APIstub.PutState(args[0], recordAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllRecords(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "REC0"
	endKey := "REC999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllRecords:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeRecordPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	recordAsBytes, _ := APIstub.GetState(args[0])
	record := MedicalRecord{}

	json.Unmarshal(recordAsBytes, &record)
	record.Cost = args[1]

	recordAsBytes, _ = json.Marshal(record)
	APIstub.PutState(args[0], recordAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
