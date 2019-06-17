# ontology-test-standalone
ontology neovm smart contract test project

## Run the test
```go run main.go <params>``` or 
```go build``` and run ```./ontology-test-standalone <params>```

## Usage
run parameters:

```--walletfile```  : wallet.dat file , default: ./wallet.dat

```--pwd```  : password of the account (should be same for the wallet file)

```--acct2``` : for the transfer / transferFrom / approve method testing, base58 format

```--skipinit``` : skip the init method ,default true

```--code```  : avm code file , should contains the hex format strings, default:./oep4.avm

