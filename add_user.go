package main

import (
	"encoding/json"
	"fmt"
	_ "strconv"
	_ "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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

	var user User
	user.ObjectType = "user"
	user.Email = args[0]
	user.Firstname = args[1]
  user.Lastname = args[2]
	fmt.Println(user)

	//store user
	userAsBytes, _ := json.Marshal(user)	//convert to array of bytes
	err = stub.PutState(user.Email, userAsBytes)	  //store owner by its Email
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_user")
	return shim.Success(nil)
}
