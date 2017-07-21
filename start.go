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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//pb "github.com/hyperledger/fabric/protos/peer"
)
type User struct {
	Email string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Password string `json:"password"`
	Ccn string `json:"ccn"`
	Phone string `json:"phone"` 
}

type Project struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Postdate string `json:"postdate"`
	Enddate string `json:"enddate"`
	Minfund int `json:"minfund"`
	Maxfund int `json:"maxfund"`
	SponsorEmail string `json:"sponsoremail"`
	Status string `json:"status"`
	PledgeAmount int `json:"pledgeamount"`
}

type ProjectUpdates struct {
	ProjectName string `json:"pname"`
	Date string `json:"udate"`
	Text string `json:"utext"`
}

type ProjectLikes struct {
	ProjectName string `json:"pname"`
	UserEmail string `json:"uemail"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "init_user_old" {
		return t.init_user_old(stub, args)
	} else if function == "init_project_old" {
		return t.init_project_old(stub, args)
	} else if function == "init_user" {
		return t.init_user(stub, args)
	} else if function == "init_project" {
		return t.init_project(stub, args)
	} else if function == "init_pledge" {
		return t.init_pledge(stub, args)
	} else if function == "init_project_likes" {
		return t.init_project_likes(stub, args)
	} else if function == "init_project_updates" {
		return t.init_project_updates(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) init_project_old(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_project_old")
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	project_id := args[0]
	name := args[1]
	owner := args[2]
	moneyGoal := args[3]
	moneyDonated := "0"

	//check if owner exists
	//user, err := get_user(stub, owner)
	//if err != nil {
	//	fmt.Println("Failed to find user - " + owner)
	//	return nil, errors.New(err.Error())
	//}

	//check if ticket id already exists
	project, err := stub.GetState(project_id)
	if err == nil {
		fmt.Println("This project already exists - " + project_id)
		fmt.Println(project)
		return nil, errors.New("This ticket project exists - " + project_id)  //all stop a ticket by this id exists
	}

	//build the ticket json string manually
	str := `{
		"docType":"project", 
		"project_id": "` + project_id + `", 
		"name": "` + name + `", 
		"owner": "` + owner + `", 
"moneyGoal": "` + moneyGoal + `",
"moneyDonated": "` + moneyDonated + `"
	}`
	err = stub.PutState(project_id, []byte(str)) 	//store project with id as key
	if err != nil {
		return nil, errors.New(err.Error())
	}

	fmt.Println("- end init_project_old")
	return nil, nil
}

func (t *SimpleChaincode) init_user_old(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_user_old")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}
				  
	email := args[0]
	firstname := args[1]
	lastname := args[2]
	
	//check if user already exists
	_, err = stub.GetState(email)
	if err == nil {
		fmt.Println("This user already exists - " + email)
		return nil, errors.New("This user exists - " + email)  //all stop a ticket by this id exists
	}

	//build the user json string manually
	str := `{
		"docType":"user", 
		"email": "` + email + `", 
		"firstname": "` + firstname + `", 
		"lastname": "` + lastname + `"
	}`
	err = stub.PutState(email, []byte(str)) 	//store project with id as key
	if err != nil {
		return nil, errors.New(err.Error())
	}

	fmt.Println("- end init_user_old")
	return nil, nil
}


func (t *SimpleChaincode) init_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_user")
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}
	user := User{};
	user.Email =  args[0]
	user.Firstname = args[1]
	user.Lastname = args[2]
	user.Password =  args[3]
	user.Ccn = args[4]
	user.Phone = args[5]
	fmt.Println(user)
	

	//store user
	userAsBytes, _ := json.Marshal(user)	//convert to array of bytes
	err = stub.PutState(user.Email, userAsBytes)	  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return nil, errors.New(err.Error())
	}
	
	fmt.Println("- end init_user")
	return nil, nil	
}

func (t *SimpleChaincode) init_project(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_project")
	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}
	project := Project{};
	project.Name = args[0]
	project.Description = args[1]
	project.Postdate = args[2]
	project.Enddate = args[3]
	project.Minfund, _ = strconv.Atoi(args[4])
	project.Maxfund, _ = strconv.Atoi(args[5])
	project.SponsorEmail = args[6]
	project.Status = args[7]
	project.PledgeAmount = 0
	fmt.Println(project)
	

	//store project
	projectAsBytes, _ := json.Marshal(project)	//convert to array of bytes
	err = stub.PutState(project.Name, projectAsBytes)	  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store project")
		return nil, errors.New(err.Error())
	}
	
	fmt.Println("- end init_project")
	return nil, nil	
}

func (t *SimpleChaincode) init_project_likes(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_project_likes")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	project_likes := ProjectLikes{};
	project_likes.ProjectName = args[0]
	project_likes.UserEmail = args[1]
	fmt.Println(project_likes)
	

	//store project_likes
	projectLikesAsBytes, _ := json.Marshal(project_likes)	//convert to array of bytes
	err = stub.PutState(project_likes.ProjectName, projectLikesAsBytes)	  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store project_likes")
		return nil, errors.New(err.Error())
	}
	
	fmt.Println("- end init_project_likes")
	return nil, nil	
}

func (t *SimpleChaincode) init_project_updates(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_project_updates")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	project_updates := ProjectUpdates{};
	project_updates.ProjectName = args[0]
	project_updates.Date = args[1]
	project_updates.Text = args[2]
	fmt.Println(project_updates)
	

	//store project_updates
	projectUpdatesAsBytes, _ := json.Marshal(project_updates)	//convert to array of bytes
	err = stub.PutState(project_updates.ProjectName, projectUpdatesAsBytes)	  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store project_updates")
		return nil, errors.New(err.Error())
	}
	
	fmt.Println("- end init_project_updates")
	return nil, nil	
}

func (t *SimpleChaincode) init_pledge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Starting init_pledge")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	projectName := args[0]
	pledgeAmount := args[1]
	project, _ := read(stub, projectName)
	project.pledgeAmount += pledgeAmount
	err = stub.PutState(projectName, [1]string{projectName})
	if err != nil {
		fmt.Println("Could not store pledge")
		return nil, errors.New(err.Error())
	}
	
	fmt.Println("- end init_pledge")
	return nil, nil	
}
