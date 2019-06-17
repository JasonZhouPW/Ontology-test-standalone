package main

import (
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"math"
	"time"
	"flag"
	"io/ioutil"
)

var(
	walletFile  string
	pwd string
	acct2 string
	codefile string
)

func init(){
	flag.StringVar(&walletFile,"wallet","./wallet.dat","wallet file")
	flag.StringVar(&pwd,"pwd","","account password")
	flag.StringVar(&codefile,"code","./oep4.avm","oep4 file")
	flag.StringVar(&acct2,"acct2","","account2 base58 address")

	flag.Parse()
}

func main() {

	fmt.Printf("use wallet file:%s\n",walletFile)
	if len(pwd) == 0{
		fmt.Println("pwd is empty")
		return
	}
	fmt.Printf("use code file:%s\n",codefile)
	code, err := ioutil.ReadFile(codefile)
	if err != nil {
		fmt.Printf("load code file error:%s\n",err.Error())
		return
	}
	fmt.Printf("acct2:%s\n",acct2)
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
	//AS3SCXw8GKTEeXpdwVw7EcC4rqSebFYpfb
	account2, err := wallet.GetAccountByAddress(acct2, []byte(pwd))
	if err != nil {
		fmt.Printf("Error!!:GetContractAddressByaddress error\n")
		return
	}

	skipinit := true
	codeHash := string(code)
	//codeHash := "56c56b681953797374656d2e53746f726167652e476574436f6e746578746a00527ac414e98f4998d837fcdd44a50561f7f32140c7c6c2606a51527ac40800008a5d784563016a52527ac401016a53527ac401026a54527ac401036a55527ac46c5fc56b6a00527ac46a51527ac46a52527ac46a51c304696e69747d9c7c75640f006a00c365a4016c75666203006a51c3046e616d657d9c7c75640f006a00c3653b026c75666203006a51c30673796d626f6c7d9c7c75640f006a00c36539026c75666203006a51c308646563696d616c737d9c7c75640f006a00c3652a026c75666203006a51c30b746f74616c537570706c797d9c7c75640f006a00c36514026c75666203006a51c30962616c616e63654f667d9c7c756414006a52c300c36a00c36518026c75666203006a51c3087472616e736665727d9c7c75641e006a52c352c36a52c351c36a52c300c36a00c3651c026c75666203006a51c30d7472616e736665724d756c74697d9c7c756412006a52c36a00c3654b036c75666203006a51c30c7472616e7366657246726f6d7d9c7c756423006a52c353c36a52c352c36a52c351c36a52c300c36a00c36525046c75666203006a51c307617070726f76657d9c7c75641e006a52c352c36a52c351c36a52c300c36a00c36548036c75666203006a51c309616c6c6f77616e63657d9c7c756419006a52c351c36a52c300c36a00c365b7056c7566620300006c756653c56b6a00527ac46a51527ac46a52527ac4616c756653c56b6a00527ac46a00c355c36a00c300c3681253797374656d2e53746f726167652e476574640a00006c75666282006a00c352c36a00c355c36a00c300c3681253797374656d2e53746f726167652e5075746a00c352c36a00c353c36a00c351c37e6a00c300c3681253797374656d2e53746f726167652e507574087472616e73666572006a00c351c36a00c352c354c176c9681553797374656d2e52756e74696d652e4e6f74696679516c75666c756652c56b6a00527ac40f4f4e5420416e6e69766572736172796c756652c56b6a00527ac4044f4e54416c756652c56b6a00527ac4586c756652c56b6a00527ac46a00c355c36a00c300c3681253797374656d2e53746f726167652e4765746c756653c56b6a00527ac46a51527ac46a00c353c36a51c37e6a00c300c3681253797374656d2e53746f726167652e4765746c756658c56b6a00527ac46a51527ac46a52527ac46a53527ac46a51c3681b53797374656d2e52756e74696d652e436865636b5769746e65737376640c00756a53c3007da07c75f16a00c353c36a51c37e6a54527ac46a54c36a00c300c3681253797374656d2e53746f726167652e4765746a55527ac46a53c36a55c37da17c75f16a53c36a55c37d9c7c756425006a54c36a00c300c3681553797374656d2e53746f726167652e44656c657465622b006a53c36a55c36a00c36503046a54c36a00c300c3681253797374656d2e53746f726167652e5075746a00c353c36a52c37e6a56527ac46a53c36a56c36a00c300c3681253797374656d2e53746f726167652e4765746a00c36585036a56c36a00c300c3681253797374656d2e53746f726167652e507574087472616e736665726a51c36a52c36a53c354c176c9681553797374656d2e52756e74696d652e4e6f74696679516c756657c56b6a00527ac46a51527ac4006a52527ac46a51c36a53527ac46a53c3c06a54527ac46a52c36a54c39f6432006a53c36a52c3c36a55527ac46a52c351936a52527ac46a55c352c36a55c351c36a55c300c36a00c36554fef162caff516c756655c56b6a00527ac46a51527ac46a52527ac46a53527ac46a51c3681b53797374656d2e52756e74696d652e436865636b5769746e657373f16a53c36a51c36a00c365d6fd7da17c7576640c00756a53c3007da07c75f16a53c36a00c354c36a51c37e6a52c37e6a00c300c3681253797374656d2e53746f726167652e50757408617070726f76616c6a51c36a52c36a53c354c176c9681553797374656d2e52756e74696d652e4e6f74696679516c75665bc56b6a00527ac46a51527ac46a52527ac46a53527ac46a54527ac46a51c3681b53797374656d2e52756e74696d652e436865636b5769746e657373f16a00c353c36a52c37e6a55527ac46a55c36a00c300c3681253797374656d2e53746f726167652e4765746a56527ac46a54c36a56c37da17c7576640c00756a54c3007da07c75f16a00c354c36a52c37e6a51c37e6a57527ac46a57c36a00c300c3681253797374656d2e53746f726167652e4765746a58527ac46a00c353c36a53c37e6a59527ac46a54c36a58c37da17c75f16a54c36a58c37d9c7c75644d006a57c36a00c300c3681553797374656d2e53746f726167652e44656c6574656a54c36a56c36a00c3654e016a55c36a00c300c3681253797374656d2e53746f726167652e5075746253006a54c36a58c36a00c36523016a57c36a00c300c3681253797374656d2e53746f726167652e5075746a54c36a56c36a00c365fb006a55c36a00c300c3681253797374656d2e53746f726167652e5075746a54c36a59c36a00c300c3681253797374656d2e53746f726167652e4765746a00c3658b006a59c36a00c300c3681253797374656d2e53746f726167652e507574087472616e736665726a52c36a53c36a54c354c176c9681553797374656d2e52756e74696d652e4e6f74696679516c756654c56b6a00527ac46a51527ac46a52527ac46a00c354c36a51c37e6a52c37e6a00c300c3681253797374656d2e53746f726167652e4765746c756655c56b6a00527ac46a51527ac46a52527ac46a51c36a52c3936a53527ac46a53c36a51c37da27c75f16a53c36c756654c56b6a00527ac46a51527ac46a52527ac46a51c36a52c37da27c75f16a51c36a52c3946c7566"
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
		return	}


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
