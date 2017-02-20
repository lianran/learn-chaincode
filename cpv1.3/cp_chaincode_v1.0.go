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
    "time"
    "fmt"
    "strings"
    "strconv"
    "encoding/json"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

// myChaincode 
type myChaincode struct {
}
//sp
var sp = "\n"

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
    err := shim.Start(new(myChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

// Init resets all the things
func (t *myChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
    return shim.Success(nil)
}

// Invoke is our entry point to invoke a chaincode function
func (t *myChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    switch function {

    case "create":
        return t.create(stub, args)

    case "transfer":
        return t.transfer(stub, args)

    case "createbx":
        return t.createbx(stub, args)

    case "confirmbx":
        return t.confirmbx(stub, args)

    case "rejectvx":
        return t.rejectbx(stub, args)

    case "myHistory":
        return t.myHistory(stub, args)

    case "getCountofInvoice":
        return t.getCountofInvoice(stub, args)

    case "getInvoice":
        return t.getInvoice(stub, args)

    case "getMetadata":
        return t.getMetadata(stub, args)

    case "getReimburseInfo":
        return t.getReimburseInfo(stub, args)

    case "getbx":
        return t.getbx(stub,args)

    default:
        return shim.Error("Unsupported operation")
    }
}

func (t *myChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 4{
        return shim.Error("create operation must include at last four arguments, a uuid , a from , a to and timestamp")
    }
    // get the args
    uuid := args[0]
    fromid := args[1]
    toid := args[2]
    timestamp := args[3]
    metadata := args[4]
    history := fromid
    owner := fromid
    status := "0"
    createtm := timestamp
    submittm := "0"
    confirmtm := "0"
    bxuuid := "0"

    //TODO: need some check for fromid and data
    //check fromid and toid
    if fromid == toid {
        return shim.Error("create operation failed, fromid is same with toid")
    }
    //do some check for the timestamp
    ts := time.Now().Unix() 
    tm, err := strconv.ParseInt(timestamp, 10, 64) 
    if err != nil {
        return shim.Error("bad format of the timestamp")
    }
    if tm - ts > 3600 || ts - tm > 3600 {
        return shim.Error("the timestamp is bad one !")
    }

    
    //check for existence of the bill
    oldvalue, err := stub.GetState(uuid)
    if err != nil {
        return shim.Error("create operation failed. Error accessing state(check the existence of bill): " + err.Error())
    }
    if oldvalue != nil {
        return shim.Error("existed bill!")
    } 

    key := uuid
    value := fromid + sp + toid + sp + history + sp + owner + sp + status + sp  + createtm + sp + submittm + sp + confirmtm + sp + bxuuid
    fmt.Printf("value is %s", value)

    err = stub.PutState(key, []byte(value))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return shim.Error("create operation failed. Error updating state: " + err.Error())
    }
    //store the metadata
    key = uuid + sp + "md"
    value = metadata
    fmt.Printf("value is %s", value)

    err = stub.PutState(key, []byte(value))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return shim.Error("store the metadat operation failed. Error updating state: " + err.Error())
    }

    //store the from and to 
    key = fromid + sp + timestamp + sp + uuid
    err = stub.PutState(key, []byte(timestamp))
    if err != nil {
        fmt.Printf("Error putting state for fromid : %s", err)
        return shim.Error("store the from operation failed. Error updating state: " + err.Error())
    }
    key = toid + sp + timestamp + sp + uuid
    err = stub.PutState(key, []byte(timestamp))
    if err != nil {
        fmt.Printf("Error putting state for toid : %s", err)
        return shim.Error("store the to operation failed. Error updating state: " + err.Error())
    }
    return shim.Success(nil)
}

func (t *myChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 4{
            return shim.Error("transfer operation must include at last there arguments, a uuid , a owner , a toid and timestamp")
        }
    //get the args
    key := args[0]
    uuid := key
    _owner := args[1]
    _toid := args[2]
    timestamp := args[3]


    //get the  info of uuid
    value, err := stub.GetState(key)
    if err != nil {
        return shim.Error("get operation failed. Error accessing state: " +  err.Error())
    }
    if value == nil {
        return shim.Error("this bill does not exist")
    }
    listValue := strings.Split(string(value), sp)
    fromid := listValue[0]
    toid := listValue[1]
    history := listValue[2]
    owner := listValue[3]
    status := listValue[4]
    createtm := listValue[5]
    submittm := listValue[6]
    confirmtm := listValue[7]
    bxuuid := listValue[8]
    //ToDo: some check for the owner?
    // if the person don't own it, he can transfer this bill
    if _owner != owner {
        return shim.Error("don't have the right to transfer the bill")
        //return nil, errors.New("don't have the right to transfer")
    }
    //if the owner is toid, it cann't be transfer any more
    if owner == toid {
        return shim.Error("cann't transfer bill now")
    }
    if status == "2" {
        return shim.Error("this bill has been submited adn you can't transfer it any more!")
    }
    if status == "3" {
        return shim.Error("this bill has been reimbursed!")
    }
    if status == "0"{
        status = "1"
    }

    //do some check for the timestamp
    ts := time.Now().Unix() 
    tm, err := strconv.ParseInt(timestamp, 10, 64) 
    if err != nil {
        return shim.Error("bad format of the timestamp")
    }
    if tm - ts > 3600 || ts - tm > 3600 {
        return shim.Error("the timestamp is bad one !")
    }

    history = history + "," + _toid
    owner = _toid
    newvalue := fromid + sp + toid + sp + history + sp + owner + sp + status + sp  + createtm + sp + submittm + sp + confirmtm + sp + bxuuid
    fmt.Printf("the old value is: %s", value)
    fmt.Printf("the new value is: %s", newvalue)
    err = stub.PutState(key, []byte(newvalue))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return shim.Error("transfer operation failed. Error updating state: " + err.Error())
    }
    //ToDo: some check for the state of puting 
    // add two sp have no reasons:)
    key = owner + sp + sp + uuid
    err = stub.PutState(key, []byte(timestamp))
    if err != nil {
        fmt.Printf("Error putting state for owner : %s", err)
        return shim.Error("transfer Operation failed. Error updating state: " + err.Error())
    }
    return shim.Success(nil)
}

func (t *myChaincode) createbx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 7{
        return shim.Error("createbx operation must include at last three arguments, a bxuuid, an array of uuid,  a owner , a toid, timestamp , data for reimbursement and a endtm")
    }
    //get the args
    bxuuid := args[0]
    uuids := args[1]
    owner := args[2]
    toid := args[3]
    timestamp := args[4]
    bxinfo := args[5]
    endtm := args[6]

    //do some check for the args
    //timestamp
    ts := time.Now().Unix() 
    tm, err := strconv.ParseInt(timestamp, 10, 64) 
    if err != nil {
        return shim.Error("bad format of the timestamp")
    }
    if tm - ts > 3600 || ts - tm > 3600 {
        return shim.Error("the timestamp is bad one !")
    }
    //endtm
    tm, err = strconv.ParseInt(timestamp, 10, 64) 
    if err != nil {
        return shim.Error("bad format of the endtm")
    }
    listuuid := strings.Split(uuids, ",")

    //update each uuid
    for _, uuid := range listuuid {
        _, err := t.submit(stub, uuid, owner, toid, timestamp, bxuuid)
        if err != nil{
            return shim.Error(err.Error())
        }
    }
    //store the bxuuid
    value := uuids + sp + owner + sp + toid  + sp + endtm + sp + bxinfo
    err = stub.PutState(bxuuid, []byte(value))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return shim.Error("createbx operation failed. Error updating state: " + err.Error())
    }
    return shim.Success(nil)
}


func (t *myChaincode) submit(stub shim.ChaincodeStubInterface, uuid string, _owner string, _toid string, timestamp string, bxuuid string) ([]byte, error){
    key := uuid
    //get the  info of uuid
    value, err := stub.GetState(key)
    if err != nil {
        return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
    }
    if value == nil {
        return nil, fmt.Errorf("this bill does not exist")
    }
    listValue := strings.Split(string(value), sp)
    fromid := listValue[0]
    toid := listValue[1]
    history := listValue[2]
    owner := listValue[3]
    status := listValue[4]
    createtm := listValue[5]
    //update the submittm
    submittm := timestamp
    confirmtm := listValue[7]

    // if the person don't own it, he can transfer this bill
    if _owner != owner {
        return []byte("don't have the right to submit the bill"), errors.New("don't have the right to submit")
        //return nil, errors.New("don't have the right to transfer")
    }
    if _toid != toid{
        return []byte("bad toid"), errors.New("bad toid")
    }
    if status == "2" {
        return []byte("this bill has been submited adn you can't transfer it any more!"), errors.New("this bill has been submited adn you can't transfer it any more!")
    }
    if status == "3" {
        return []byte("this bill has been reimbursed!"), errors.New("this bill has been reimbursed!")
    }
    if status == "1" || status == "0" {
        status = "2"
    }

    //update the uuid info
    newvalue := fromid + sp + toid + sp + history + sp + owner + sp + status + sp  + createtm + sp + submittm + sp + confirmtm + sp + bxuuid
    fmt.Printf("the old value is: %s", value)
    fmt.Printf("the new value is: %s", newvalue)
    err = stub.PutState(key, []byte(newvalue))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return nil, fmt.Errorf("submit operation failed. Error updating state: %s", err)
    }
    return nil, nil
}

func (t *myChaincode) confirmbx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 3 {
        return shim.Error("confirmbx operation must include at last there arguments, a bxuuid , a toid and timestamp")
    }
    //get the args
    bxuuid := args[0]
    _toid := args[1]
    timestamp := args[2]

    //do some check for the timestamp
    ts := time.Now().Unix() 
    tm, err := strconv.ParseInt(timestamp, 10, 64) 
    if err != nil {
        return shim.Error("bad format of the timestamp")
    }
    if tm - ts > 3600 || ts - tm > 3600 {
        return shim.Error("the timestamp is bad one !")
    }

    //get the info of bxuuid
    value, err := stub.GetState(bxuuid)
    listvalue := strings.SplitN(string(value), sp, 4)
    uuids := listvalue[0]
    //owner := listvalue[1]
    toid := listvalue[2]
    //endtm := listvalue[3]
    //check for the endtm
    //intendtm, err := strconv.ParseInt(endtm, 10, 64)
    //if ts > intendtm {
        //ToDo: reject it ?
        //return nil, fmt.Errorf("out of the bx time")
    //}
    if _toid != toid {
        return shim.Error(" don't hvae the right to confirmbx!")
    }

    //update each uuid
    listuuid := strings.Split(uuids, ",")
    for _, uuid := range listuuid {
        _, err := t.confirm(stub, uuid, toid, timestamp)
        if err !=nil {
            return shim.Error(err.Error())
        }
    }
    return shim.Success(nil)
}


func (t *myChaincode) confirm(stub shim.ChaincodeStubInterface, uuid string, _toid string, timestamp string) ([]byte, error){
    
    key := uuid
    //get the  info of uuid
    value, err := stub.GetState(key)
    if err != nil {
        return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
    }
    if value == nil {
        return nil, fmt.Errorf("this bill does not exist")
    }
    listValue := strings.Split(string(value), sp)
    fromid := listValue[0]
    toid := listValue[1]
    //update the history
    history := listValue[2] + "," + toid
    //update the owner
    owner := toid
    status := listValue[4]
    createtm := listValue[5]
    submittm := listValue[6]
    //update the confirmtm
    confirmtm := timestamp
    bxuuid := listValue[8]

    // if the person is not the toid
    if _toid != toid {
        return []byte("don't have the right to cnfirm the bill"), errors.New("don't have the right to confirm")
        //return nil, errors.New("don' t have the right to transfer")
    }
    if status == "1" || status == "0" {
        return []byte("this bill has not been submited "), errors.New("this bill has  not been submited")
    }
    if status == "3" {
        return []byte("this bill has been reimbursed!"), errors.New("this bill has been reimbursed!")
    }
    if status == "2" {
        status = "3"
    }

    //update the uuid info
    newvalue := fromid + sp + toid + sp + history + sp + owner + sp + status + sp  + createtm + sp + submittm + sp + confirmtm + sp + bxuuid
    fmt.Printf("the old value is: %s", value)
    fmt.Printf("the new value is: %s", newvalue)
    err = stub.PutState(key, []byte(newvalue))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return nil, fmt.Errorf("confirm operation failed. Error updating state: %s", err)
    }
    return nil, nil
}
func (t *myChaincode) rejectbx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 4 {
        return shim.Error("reject operation must include at last four arguments, a uuid , a toid , a reason and  timestamp")
    }
    //get the args
    bxuuid := args[0]
    _toid := args[1]
    reason := args[2]
    timestamp := args[3]

    //do some check for the timestamp
    ts := time.Now().Unix() 
    tm, err := strconv.ParseInt(timestamp, 10, 64) 
    if err != nil {
        return shim.Error("bad format of the timestamp")
    }
    if tm - ts > 3600 || ts - tm > 3600 {
        return shim.Error("the timestamp is bad one !")
    }

    //get the info of bxuuid
    value, err := stub.GetState(bxuuid)
    listvalue := strings.SplitN(string(value), sp, 4)
    uuids := listvalue[0]
    owner := listvalue[1]
    toid := listvalue[2]
    endtm := listvalue[3]
    bxinfo := listvalue[4]
    //check for the endtm
    // intendtm, err := strconv.ParseInt(endtm, 10, 64)
    // if ts > intendtm {
    //     //ToDo: reject it ?
    //     return nil, fmt.Errorf("out of the bx time")
    // }
    if _toid != toid {
        return shim.Error(" don't hvae the right to confirmbx!")
    }
    //update each uuid
    listuuid := strings.Split(uuids, ",")
    for _, uuid := range listuuid {
        _, err := t.reject(stub, uuid, toid, timestamp)
        if err !=nil {
            return shim.Error(err.Error())
        }
    }
    //update the uuid infor
    newValue := uuids + sp + owner + sp + toid  + sp + endtm + sp + reason + sp + bxinfo 
    err = stub.PutState(bxuuid, []byte(newValue))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return shim.Error("rejectbx operation failed. Error updating state: " + err.Error())
    }
    return shim.Success(nil)
}

func (t *myChaincode) reject(stub shim.ChaincodeStubInterface, uuid string, _toid string, timestamp string) ([]byte, error){
    key := uuid
    //get the  info of uuid
    value, err := stub.GetState(key)
    if err != nil {
        return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
    }
    if value == nil {
        return nil, fmt.Errorf("this bill does not exist")
    }
    listValue := strings.Split(string(value), sp)
    fromid := listValue[0]
    toid := listValue[1]
    history := listValue[2]
    owner := listValue[3]
    status := listValue[4]
    createtm := listValue[5]
    submittm := listValue[6]
    confirmtm := timestamp
    bxuuid := listValue[8]

    // if the person is not the toid
    if _toid != toid {
        return []byte("don't have the right to reject the bill"), errors.New("don't have the right to reject")
        //return nil, errors.New("don't have the right to transfer")
    }
    if status == "1" || status == "0" {
        return []byte("this bill has not been submited "), errors.New("this bill has  not been submited ")
    }
    if status == "3" {
        return []byte("this bill has been reimbursed!"), errors.New("this bill has been reimbursed!")
    }
    if status == "2" {
        status = "1"
    }

    //update the uuid info
    newvalue := fromid + sp + toid + sp + history + sp + owner + sp + status + sp  + createtm + sp + submittm + sp + confirmtm + sp + bxuuid
    fmt.Printf("the old value is: %s", value)
    fmt.Printf("the new value is: %s", newvalue)
    err = stub.PutState(key, []byte(newvalue))
    if err != nil {
        fmt.Printf("Error putting state %s", err)
        return nil, fmt.Errorf("reject operation failed. Error updating state: %s", err)
    }
    return nil, nil

}

func (t *myChaincode) myHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 2{
        return shim.Error("myhistory operation must include at last two arguments, owner, flag(and time s)")
    }
    owner := args[0]
    business := args[1]
    //Todo: some check for the owner?

    //get the timestamp
    ts := time.Now().Unix() 
    //timestamp := strconv.FormatInt(ts, 10)

    tm := int64(3600)
    var err error
    if len(args) >= 3{
        tm, err = strconv.ParseInt(args[2], 10, 64)
    }
    if err != nil {
        return shim.Error("getnumofbills failed. Bad format of the time: " + err.Error())
    }
    starttime := strconv.FormatInt(ts-tm, 10)
    endtime := strconv.FormatInt(ts,10)

    //check is a user or a business
    bus := true
    fmt.Printf("flag is %s", business)
    if business != "1"{
        bus = false
    }
    fmt.Println(bus)
    keysIter, err := stub.GetStateByRange("just find nothin", "just find nothin")
    if bus {
        keysIter, err = stub.GetStateByRange(owner + sp + starttime, owner + sp + endtime)
    } else {
        keysIter, err = stub.GetStateByRange(owner + sp + sp + "0", owner + sp + sp + "z")
    }

    
    if err != nil {
        return shim.Error("getnumofbills failed. Error accessing state: " + err.Error())
    }
    defer keysIter.Close()

    var keys []string
    for keysIter.HasNext() {
        key, _, iterErr := keysIter.Next()
        if iterErr != nil {
            return shim.Error("getnumofbills operation failed. Error accessing state: " + err.Error())
        }
        keys = append(keys, key)
    }
    
    jsonKeys, err := json.Marshal(keys)
    if err != nil {
        return shim.Error("keys operation failed. Error marshaling JSON: " + err.Error())
    }

    return shim.Success(jsonKeys)
}

func (t *myChaincode) getCountofInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 1{
        return shim.Error("getnumsofbills operation must include at last one argument, owner(and time s)")
    }
    owner := args[0]
    //Todo: some check for the owner?

    //get the timestamp
    ts := time.Now().Unix() 

    tm := int64(3600)
    var err error
    if len(args) >= 2{
        tm, err = strconv.ParseInt(args[1], 10, 64)
    }
    if err != nil {
        return shim.Error("getnumofbills failed. Bad format of the time: " + err.Error())
    }
    starttime := strconv.FormatInt(ts-tm, 10)
    endtime := strconv.FormatInt(ts,10)

    keysIter, err := stub.GetStateByRange(owner + sp + starttime, owner + sp + endtime)
    if err != nil {
        return shim.Error("getnumofbills failed. Error accessing state: " + err.Error())
    }
    defer keysIter.Close()

    cnt := int64(0)

    for keysIter.HasNext() {
        _, _, iterErr := keysIter.Next()
        if iterErr != nil {
            return shim.Error("getnumofbills operation failed. Error accessing state: " + err.Error())
        }
        cnt = cnt + 1
    }
    return shim.Success([]byte(strconv.FormatInt(cnt,10)))
}

func (t *myChaincode) getInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 2{
        return shim.Error("getbill operation must include at last two arguments, uuid and owner")
    }
    uuid := args[0]
    _owner := args[1]

    //ToDo: some checks?
    key := uuid
    value, err := stub.GetState(key)
    if err != nil {
        return shim.Error("get operation failed. Error accessing state: " +  err.Error())
    }
    if value == nil {
        return shim.Success([]byte("don't have this bill"))
    }
    listValue := strings.Split(string(value), sp)
    // check the ownership
    //fromid := listValue[0]
    toid := listValue[1]
    history := listValue[2]
    //owner := listValue[3]
    listHistory := strings.Split(history,",")
    flag := false
    for _, val := range listHistory{
        if _owner == val{
            flag = true
        }
    }
    if flag || _owner == toid {
        return shim.Success(value)
    }
    return shim.Success([]byte("you don't have the right to get this bill"))
}

func (t *myChaincode) getMetadata(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) < 2{
        return shim.Error("getmetadata operation must include at last two arguments, uuid and owner")
    }
    uuid := args[0]
    _owner := args[1]

    //ToDo: some checks?
    key := uuid
    value, err := stub.GetState(key)
    if err != nil {
        return shim.Error("get operation failed. Error accessing state: " + err.Error())
    }
    if value == nil {
        return shim.Success([]byte("don't have this bill"))
    }
    listValue := strings.Split(string(value), sp)
    // check the ownership
    //fromid := listValue[0]
    toid := listValue[1]
    history := listValue[2]
    //owner := listValue[3]
    listHistory := strings.Split(history,",")
    flag := false
    for _, val := range listHistory{
        if _owner == val{
            flag = true
        }
    }
    if flag != true && _owner == toid {
        return shim.Error("you don't have the right to get the matedate of this bill")
    }
    //get the metadata
    key = uuid + sp + "md"
    value, err = stub.GetState(key)
    if err != nil {
        return shim.Error("get operation failed. Error accessing state: " + err.Error())
    }
    return shim.Success(value)
}

func (t *myChaincode)getReimburseInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 2{
        return shim.Error("getmetadata operation must include at last two arguments, uuid and ownerid")
    }
    uuid := args[0]
    _owner := args[1]

    //get the info uuid
    key := uuid
    value, err := stub.GetState(key)
    if err != nil{
        return shim.Error("get operation failed. Error accessing state: " + err.Error())
    }
    if value == nil {
        return shim.Success([]byte("don't have this bill"))
    }
    listValue := strings.Split(string(value), sp)
    // get the bxid
    bxuuid := listValue[8]
    status := listValue[4]
    if status == "0" || status == "1"{
        return shim.Success([]byte("this invoice have not been reimbursed yet"))
    }
    //get the reimbuseinfo
    key = bxuuid
    value, err = stub.GetState(key)
    if err != nil {
        return shim.Error("get operation failed. Error accessing state: " + err.Error())
    }
    //check the right 
    listValue = strings.Split(string(value), sp)
    ownerid := listValue[1]
    toid := listValue[2]
    if ownerid != _owner && toid != _owner {
        return shim.Error("you don't have the right to get the reimbuseinfo")
    }
    return shim.Success(value)
}

func (t *myChaincode)getbx(stub shim.ChaincodeStubInterface, args []string) pb.Response{
    if len(args) < 2{
        return shim.Error("getmetadata operation must include at last two arguments, bxuuid and ownerid")
    }
    bxuuid := args[0]
    _owner := args[1]

    //get the info of the bxuuid
    key := bxuuid
    value, err := stub.GetState(key)
    if err != nil{
        return shim.Error("get operation failed. Error accessing state: " + err.Error())
    }
    if value == nil {
        return shim.Success([]byte("don't have this info"))
    }
    listValue := strings.Split(string(value), sp)
    // check the ownership
    owner := listValue[1]
    toid := listValue[2]
    if _owner != owner && _owner != toid{
        return shim.Error("you don't have the right to get this info")
    }
    return shim.Success(value)
}