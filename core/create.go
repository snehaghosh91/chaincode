package main

import (
	"encoding/json"
	"fmt"
	_ "strconv"
	_ "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
type User struct {
	Email string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]                                   //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value))         //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

func init_user(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_user")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}
	user := User{};
	user.Email =  args[0]
	user.Firstname = args[1]
	user.Lastname = args[2]
	fmt.Println(user)
        email := args[0]
        _, err = stub.GetState(email)
	if err != nil {
		fmt.Println("This user already exists - " + email)
		return nil, shim.Error(err.Error())
	}

	//store user
	userAsBytes, _ := json.Marshal(user)	//convert to array of bytes
	err = stub.PutState(user.Email, userAsBytes)	  //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_user marble")
	return shim.Success(nil)
}
