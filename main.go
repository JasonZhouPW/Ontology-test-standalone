package main

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"

	"flag"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"io/ioutil"
	"math"
	"time"
)

var (
	walletFile string
	pwd        string
	acct2      string
	codefile   string
	skipinit   bool
)

func init() {
	flag.StringVar(&walletFile, "wallet", "./wallet.dat", "wallet file")
	flag.StringVar(&pwd, "pwd", "", "account password")
	flag.StringVar(&codefile, "code", "./oep4.avm", "oep4 file")
	flag.StringVar(&acct2, "acct2", "", "account2 base58 address")
	flag.BoolVar(&skipinit, "skipinit", false, "skip init method")

	flag.Parse()
}

func main() {

	fmt.Printf("use wallet file:%s\n", walletFile)
	if len(pwd) == 0 {
		fmt.Println("pwd is empty")
		return
	}
	fmt.Printf("use code file:%s\n", codefile)
	code, err := ioutil.ReadFile(codefile)
	if err != nil {
		fmt.Printf("load code file error:%s\n", err.Error())
		return
	}
	fmt.Printf("acct2:%s\n", acct2)
	if len(acct2) == 0 {
		fmt.Printf("acct2 is empty \n")
		return
	}

	fmt.Println("starting test OEP4 ")
	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress("http://polaris2.ont.io:20336")
	wallet, err := ontSdk.OpenWallet(walletFile)
	if err != nil {
		fmt.Printf("Error!!:wallet open failed\n")
		return
	}

	//we get the first account of the wallet by your password
	signer, err := wallet.GetDefaultAccount([]byte(pwd))
	if err != nil {
		fmt.Printf("Error!!:password error\n")
		return
	}
	account2, err := wallet.GetAccountByAddress(acct2, []byte(pwd))
	if err != nil {
		fmt.Printf("Error!!:GetContractAddressByaddress error\n")
		return
	}

	codeHash := string(code)
	fmt.Printf("codeHash is: %s\n", codeHash)
	fmt.Printf("codeHash len is :%d\n", len(codeHash))
	contractAddress, _ := utils.GetContractAddress(codeHash)

	fmt.Printf("CodeAddress :%s\n", contractAddress.ToHexString())

	//set timeout
	timeoutSec := 60 * time.Second

	gasprice := uint64(500)
	deploygaslimit := uint64(2000000000000)

	txHash, err := ontSdk.NeoVM.DeployNeoVMSmartContract(gasprice, deploygaslimit,
		signer,
		true,
		codeHash,
		"onta",
		"1.0",
		"author",
		"abc@ont.io",
		"Ont Anniversary")

	if err != nil {
		fmt.Println("deploy contract error: %s\n", err.Error())
	}
	_, err = ontSdk.WaitForGenerateBlock(timeoutSec)
	if err != nil {
		fmt.Printf("WaitForGenerateBlock error\n")
		return
	}
	fmt.Printf("the deploy contract txhash is %x\n", txHash)
	fmt.Println("deploy contract succeed!")

	fmt.Println("========== start invoke contract test ==========")
	invokeGasLimit := uint64(20000)

	if !skipinit {
		fmt.Println("test init function ")

		txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice,
			invokeGasLimit,
			signer,
			contractAddress,
			[]interface{}{"init", []interface{}{}})

		if err != nil {
			fmt.Printf("invoke init error: %s\n", err)
			return
		}

		//WaitForGenerateBlock
		_, err = ontSdk.WaitForGenerateBlock(30*time.Second, 1)
		if err != nil {
			fmt.Printf("invoke init WaitForGenerateBlock error: %s\n", err)
			return
		}

		//GetEventLog, to check the result of invoke
		events, err := ontSdk.GetSmartContractEvent(txHash.ToHexString())
		if err != nil {
			fmt.Printf("invoke invoke GetSmartContractEvent error:%s\n", err)
			return
		}
		if events.State == 0 {
			fmt.Println("invoke invoke failed invoked exec state return 0")
			return
		}

		fmt.Printf("events.Notify:%v", events.Notify)
		for _, notify := range events.Notify {
			fmt.Printf("%+v\n", notify)
		}

		fmt.Println("========== test init function ==========")
	}

	fmt.Println("========== testing totalSupply ==========")
	obj, err := ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"totalSupply", []interface{}{}})

	totalSupply, err := obj.Result.ToInteger()
	if err != nil {
		fmt.Printf("Get totalSupply error:%s", err)

		return
	}

	fmt.Printf("total supply is %d\n", totalSupply)
	fmt.Println("========== testing totalSupply end ==========")

	fmt.Println("========== testing name ==========")

	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"name", []interface{}{}})

	name, err := obj.Result.ToString()
	if err != nil {
		fmt.Printf("Get name error:%s", err)
		return
	}

	fmt.Printf("name is %s\n", name)
	fmt.Println("========== testing name end========== ")
	fmt.Println("========== testing symbol ==========")

	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"symbol", []interface{}{}})

	symbol, err := obj.Result.ToString()
	if err != nil {
		fmt.Println("Get symbol error:%s", err)
		return
	}

	fmt.Printf("symbol is %s\n", symbol)
	fmt.Println("========== testing symbol end ==========")

	fmt.Println("========== testing decimals ==========")
	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"decimals", []interface{}{}})

	decimal, err := obj.Result.ToInteger()
	if err != nil {
		fmt.Printf("Get decimals error:%s", err)
		return
	}

	fmt.Printf("decimals is %d\n", decimal)
	fmt.Println("========== testing decimals end ==========")

	fmt.Println("========== testing balanceOf ==========")
	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"balanceOf", []interface{}{signer.Address[:]}})
	if err != nil {
		fmt.Printf("Get balanceOf error:%s", err)

		return
	}

	balance, err := obj.Result.ToInteger()
	if err != nil {
		fmt.Printf("Get balanceOf error:%s", err)
		return
	}

	//
	fmt.Printf("balance of %s is %d\n", signer.Address.ToBase58(), balance)
	fmt.Println("========== testing balanceOf  end ==========")

	fmt.Println("========== testing transfer ==========")

	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		signer,
		contractAddress,
		[]interface{}{"transfer", []interface{}{signer.Address[:], account2.Address[:], 10000000000000}})
	if err != nil {
		fmt.Printf("invoke transfer error: %s\n", err)
		return
	}

	//WaitForGenerateBlock
	_, err = ontSdk.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		fmt.Printf("invoke transfer error: %s\n", err)
		return
	}

	//GetEventLog, to check the result of invoke
	events, err := ontSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("invoke transfererror:%s", err)
		return
	}
	if events.State == 0 {
		fmt.Printf("invoke transfer exec state return 0")
		return
	}
	notify := events.Notify[0]
	invokeState := notify.States.([]interface{})

	method, _ := common.HexToBytes(invokeState[0].(string))
	addFromTmp, _ := common.HexToBytes(invokeState[1].(string))
	addFrom, _ := common.AddressParseFromBytes(addFromTmp)

	addToTmp, _ := common.HexToBytes(invokeState[2].(string))
	addTo, _ := common.AddressParseFromBytes(addToTmp)
	tmp, _ := common.HexToBytes(invokeState[3].(string))
	amount := common.BigIntFromNeoBytes(tmp)
	fmt.Printf("states[method:%s,from:%s,to:%s,value:%d]\n", method, addFrom.ToBase58(), addTo.ToBase58(), amount.Int64())
	fmt.Println("========== testing transfer end ==========")

	fmt.Println("========== testing balanceOf ==========")
	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"balanceOf", []interface{}{signer.Address[:]}})
	if err != nil {
		fmt.Printf("Get balanceOf error:%s", err)

		return
	}

	balance, err = obj.Result.ToInteger()
	if err != nil {
		fmt.Printf("Get balanceOf error:%s", err)
		return
	}

	//
	fmt.Printf("balance of %s is %d\n", signer.Address.ToBase58(), balance)
	fmt.Println("========== testing balanceOf end ==========")

	fmt.Println("========== testing balanceOf ==========")
	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"balanceOf", []interface{}{account2.Address[:]}})
	if err != nil {
		fmt.Printf("Get balanceOf error:%s", err)

		return
	}

	balance, err = obj.Result.ToInteger()
	if err != nil {
		fmt.Printf("Get balanceOf error:%s", err)
		return
	}

	//
	fmt.Printf("balance of %s is %d\n", account2.Address.ToBase58(), balance)
	fmt.Println("========== testing balanceOf end ==========")

	fmt.Println("========== testing approve ==========")
	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		signer,
		contractAddress,
		[]interface{}{"approve", []interface{}{signer.Address[:], account2.Address[:], 60000000000}})
	if err != nil {
		fmt.Printf("invoke approve error: %s\n", err)
	}

	//WaitForGenerateBlock
	_, err = ontSdk.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		fmt.Printf("invoke approve WaitForGenerateBlock error: %s", err)
		return
	}

	//GetEventLog, to check the result of invoke
	events, err = ontSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("invoke approve GetSmartContractEvent error:%s", err)
		return
	}
	if events.State == 0 {
		fmt.Printf("invoke approve failed  exec state return 0")
		return
	}
	notify = events.Notify[0]
	invokeState = notify.States.([]interface{})

	method, _ = common.HexToBytes(invokeState[0].(string))
	addFromTmp, _ = common.HexToBytes(invokeState[1].(string))
	addFrom, _ = common.AddressParseFromBytes(addFromTmp)

	addToTmp, _ = common.HexToBytes(invokeState[2].(string))
	addTo, _ = common.AddressParseFromBytes(addToTmp)
	tmp, _ = common.HexToBytes(invokeState[3].(string))
	amount = common.BigIntFromNeoBytes(tmp)
	fmt.Printf("states[method:%s,owner:%s,spender:%s,amount:%d]\n", method, addFrom.ToBase58(), addTo.ToBase58(), amount.Int64())

	fmt.Println("========== testing approve end ==========")

	fmt.Println("========== testing allowance signer ==========")

	obj, err = ontSdk.NeoVM.PreExecInvokeNeoVMContract(contractAddress, []interface{}{"allowance", []interface{}{signer.Address[:], account2.Address[:]}})
	if err != nil {
		fmt.Printf("get allowance NewNeoVMSInvokeTransaction error:%s\n", err)

		return
	}

	allowance, err := obj.Result.ToInteger()
	if err != nil {
		fmt.Printf("get allowance PrepareInvokeContract error:%s", err)
		return
	}

	fmt.Printf("allowance is %d\n", allowance)
	fmt.Println("========== testing allowance signer end ==========")
	fmt.Println("========== testing transferFrom ==========")

	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		account2,
		contractAddress,
		[]interface{}{"transferFrom", []interface{}{account2.Address[:], signer.Address[:], account2.Address[:], 30000000000}})
	if err != nil {
		fmt.Printf("invoke transferFrom InvokeNeoVMSmartContract error: %s\n", err)
	}

	//WaitForGenerateBlock
	_, err = ontSdk.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		fmt.Printf("invoke transferFrom WaitForGenerateBlock error: %s\n", err)
		return
	}

	//GetEventLog, to check the result of invoke
	events, err = ontSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("invoke transferFrom GetSmartContractEvent error:%s", err)
		return
	}
	if events.State == 0 {
		fmt.Printf("invoke transferFrom failed invoked exec state return 0")
		return
	}
	notify = events.Notify[0]
	invokeState = notify.States.([]interface{})

	method, _ = common.HexToBytes(invokeState[0].(string))
	addFromTmp, _ = common.HexToBytes(invokeState[1].(string))
	addFrom, _ = common.AddressParseFromBytes(addFromTmp)

	addToTmp, _ = common.HexToBytes(invokeState[2].(string))
	addTo, _ = common.AddressParseFromBytes(addToTmp)
	tmp, _ = common.HexToBytes(invokeState[3].(string))
	amount = common.BigIntFromNeoBytes(tmp)
	fmt.Printf("states[method:%s,from:%s,to:%s,value:%d]\n", method, addFrom.ToBase58(), addTo.ToBase58(), amount.Int64())

	fmt.Println("========== testing transferFrom  end ==========")

	fmt.Println("========== functional testing end ==========")
	fmt.Println("========== start exception testing ==========")
	fmt.Println("========== testing transfer -1 ==========")

	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		signer,
		contractAddress,
		[]interface{}{"transfer", []interface{}{signer.Address[:], account2.Address[:], -1}})
	if err == nil {
		fmt.Println("Testing Error !!!: transfer -1 should failed!")
		return
	}

	fmt.Println("========== testing transfer end ==========")

	fmt.Println("========== testing approve -1 ==========")
	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		signer,
		contractAddress,
		[]interface{}{"approve", []interface{}{signer.Address[:], account2.Address[:], -1}})
	if err == nil {
		fmt.Println("Testing Error !!!: approve -1 should failed!")
		return
	}

	fmt.Println("========== testing approve end ==========")

	fmt.Println("========== testing transferFrom -1 ==========")

	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		account2,
		contractAddress,
		[]interface{}{"transferFrom", []interface{}{account2.Address[:], signer.Address[:], account2.Address[:], -1}})
	if err == nil {
		fmt.Println("Testing Error !!!: transferFrom -1 should failed!")
		return
	}

	fmt.Println("========== testing transferFrom  end ==========")

	fmt.Println("========== testing transfer MaxInt64 ==========")

	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		signer,
		contractAddress,
		[]interface{}{"transfer", []interface{}{signer.Address[:], account2.Address[:], math.MaxInt64}})
	if err == nil {
		fmt.Println("Testing Error !!!: transfer MaxInt64 should failed!")
		return
	}

	fmt.Println("========== testing transfer end ==========")
	fmt.Println("========== testing approve math.MaxInt64 ==========")
	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		signer,
		contractAddress,
		[]interface{}{"approve", []interface{}{signer.Address[:], account2.Address[:], math.MaxInt64}})
	if err == nil {
		fmt.Println("Testing Error !!!: approve MaxInt64 should failed!")
		return
	}

	fmt.Println("========== testing approve end ==========")

	fmt.Println("========== testing transferFrom math.MaxInt64 ==========")

	txHash, err = ontSdk.NeoVM.InvokeNeoVMContract(gasprice, invokeGasLimit,
		account2,
		contractAddress,
		[]interface{}{"transferFrom", []interface{}{account2.Address[:], signer.Address[:], account2.Address[:], math.MaxInt64}})
	if err == nil {
		fmt.Printf("invoke transferFrom InvokeNeoVMSmartContract error: %s\n", err)
	}

	fmt.Println("========== testing transferFrom end ==========")

	fmt.Println("========== all tests passed!! ==========")
}
