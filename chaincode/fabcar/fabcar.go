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
		MedicalRecord{Patient: "Jacob Henderson", Doctor: "MF Doom", Procedure: "Vascetomy", Cost: "400000"},
		MedicalRecord{Patient: "La'var Ball", Doctor: "Leo Muskang", Procedure: "STI Test", Cost: "00"},
		MedicalRecord{Patient: "George Stephanoppolis", Doctor: "Marley Davis", Procedure: "Immunizations", Cost: "10000"},
		MedicalRecord{Patient: "Anika Ghoshis", Doctor: "Robert DeNiro", Procedure: "12 Stitches", Cost: "103000"},
		MedicalRecord{Patient: "David Cameron", Doctor: "Rob Hood", Procedure: "Yearly Checkup", Cost: "1230400"},
		MedicalRecord{Patient: "Ryan Pagan", Doctor: "Bridget Dewey", Procedure: "Hepatitis C Vaccination", Cost: "00"},
		MedicalRecord{Patient: "Chary Adamo", Doctor: "VanBailey", Procedure: "Check Blood Pressure", Cost: "5000"},
		MedicalRecord{Patient: "Paris Hilton", Doctor: "Allen Po", Procedure: "Botox Injection 120cc", Cost: "15000000"},
		MedicalRecord{Patient: "Tata Holden", Doctor: "Nano Bot", Procedure: "Liposuction 12 lbs", Cost: "2200300"},
		MedicalRecord{Patient: "Joe Crawford", Doctor: "Brian Nina", Procedure: "Cast Broken Arm", Cost: "404400"},
		MedicalRecord{Patient:"John Smith", Doctor:"Strange", Procedure:"Brain Surgery", Cost:"50"},
		MedicalRecord{Patient:"John Jacob Jingleheimer Schmidtt", Doctor:"Grenaldi", Procedure:"Turn your head and cough", Cost:"50"},
	}

	i := 0
	for i < len(records) {
		fmt.Println("i is ", i)
		recordAsBytes, _ := json.Marshal(records[i])
		APIstub.PutState("REC"+strconv.Itoa(i), recordAsBytes)
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
	record.Patient = args[1]

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
