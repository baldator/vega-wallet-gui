package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"google.golang.org/grpc"
)

type createWalletRequest struct {
	Owner      string `json:"owner"`
	Passphrase string `json:"passphrase"`
}

type createWalletResponse struct {
	Wallet wallet.Wallet `json:"wallet"`
}

type getWalletRequest struct {
	Owner      string `json:"owner"`
	Passphrase string `json:"passphrase"`
}

type getWalletResponse struct {
	Wallet wallet.Wallet `json:"wallet"`
}

type signTransactionRequest struct {
	Owner      string `json:"owner"`
	Passphrase string `json:"passphrase"`
	Pub        string `json:"pub"`
	Message    string `json:"message"`
}

type signTransactionResponse struct {
	Message string `json:"message"`
}

type taintKeypairRequest struct {
	Owner      string `json:"owner"`
	Passphrase string `json:"passphrase"`
	Pub        string `json:"pub"`
}

type taintKeypairResponse struct {
	Wallet wallet.Wallet `json:"wallet"`
}

type verifyTransactionRequest struct {
	Owner      string `json:"owner"`
	Passphrase string `json:"passphrase"`
	Pub        string `json:"pub"`
	Message    string `json:"message"`
	Signature  string `json:"signature"`
}

type verifyTransactionResponse struct {
	Verified bool `json:"verified"`
}

type checkBalanceRequest struct {
	Owner string `json: "owner"`
}

type checkBalanceResponse struct {
	Owner string `json: "owner"`
}

type Owners struct {
	Owners []string `json:"owners"`
	Path   string   `json:"path"`
}

func getWallets() (Owners, error) {
	var owners Owners
	var ownersSlice []string
	if ok, err := fsutil.PathExists(walletPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return owners, fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(walletPath); err != nil {
			return owners, fmt.Errorf("error creating root directory: %v", err)
		}
	}

	files, err := ioutil.ReadDir(walletPath)
	if err != nil {
		return owners, fmt.Errorf("error reading root directory path: %v", err)
	}

	for _, f := range files {
		ownersSlice = append(ownersSlice, f.Name())
	}
	owners.Owners = ownersSlice
	owners.Path = walletPath

	return owners, nil
}

// Generate Wallet
func genWallet(walletOwner string, passphrase string, metas string) (wallet.Wallet, error) {
	var emptyWallet wallet.Wallet
	if len(walletOwner) <= 0 {
		return emptyWallet, errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return emptyWallet, errors.New("wallet passphrase is required and cannot be empty")
	}

	if ok, err := fsutil.PathExists(walletPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return emptyWallet, fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(walletPath); err != nil {
			return emptyWallet, fmt.Errorf("error creating root directory: %v", err)
		}
	}

	if err := wallet.EnsureBaseFolder(walletPath); err != nil {
		return emptyWallet, fmt.Errorf("unable to initialization root folder: %v", err)
	}

	wal, err := wallet.Read(walletPath, walletOwner, passphrase)
	if err != nil {
		if err != wallet.ErrWalletDoesNotExists {
			// this an invalid key, returning error
			return emptyWallet, fmt.Errorf("unable to decrypt wallet: %v", err)
		}
		// wallet do not exit, let's try to create it
		wal, err = wallet.Create(walletPath, walletOwner, passphrase)
		if err != nil {
			return emptyWallet, fmt.Errorf("unable to create wallet: %v", err)
		}
	}

	// at this point we have a valid wallet
	// let's generate the keypair
	// defaulting to ed25519 for now
	algo := crypto.NewEd25519()
	kp, err := wallet.GenKeypair(algo.Name())
	if err != nil {
		return emptyWallet, fmt.Errorf("unable to generate new key pair: %v", err)
	}

	if len(metas) > 0 {
		// expect ; separated metas
		metasSplit := strings.Split(metas, ";")
		for _, v := range metasSplit {
			metaVal := strings.Split(v, ":")
			if len(metaVal) != 2 {
				return emptyWallet, fmt.Errorf("invalid meta format")
			}
			kp.Meta = append(kp.Meta, wallet.Meta{Key: metaVal[0], Value: metaVal[1]})
		}
	}

	// the user did not specify any metas
	// we'll create a default one for them
	if len(kp.Meta) <= 0 {
		kp.Meta = append(
			kp.Meta,
			wallet.Meta{
				Key:   "name",
				Value: fmt.Sprintf("%v's key %v", walletOwner, len(wal.Keypairs)+1),
			},
		)
	}

	// now updating the wallet and saving it
	walletOutput, err := wallet.AddKeypair(kp, walletPath, walletOwner, passphrase)
	if err != nil {
		return emptyWallet, fmt.Errorf("unable to add keypair to wallet: %v", err)
	}

	return *walletOutput, nil
}

func createWallet(owner string, passphrase string) (createWalletResponse, error) {
	var walletResp createWalletResponse

	if len(owner) <= 0 {
		return walletResp, errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return walletResp, errors.New("wallet passphrase is required")
	}

	if ok, err := fsutil.PathExists(rootPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return walletResp, fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(rootPath); err != nil {
			return walletResp, fmt.Errorf("error creating root directory: %v", err)
		}
	}

	if err := wallet.EnsureBaseFolder(rootPath); err != nil {
		return walletResp, fmt.Errorf("unable to initialization root folder: %v", err)
	}

	wal, err := wallet.Read(rootPath, owner, passphrase)
	if err != nil {
		if err != wallet.ErrWalletDoesNotExists {
			// this an invalid key, returning error
			return walletResp, fmt.Errorf("unable to decrypt wallet: %v", err)
		}
		// wallet do not exit, let's try to create it
		wal, err = wallet.Create(rootPath, owner, passphrase)
		if err != nil {
			return walletResp, fmt.Errorf("unable to create wallet: %v", err)
		}
	}

	// at this point we have a valid wallet
	// let's generate the keypair
	// defaulting to ed25519 for now
	algo := crypto.NewEd25519()
	kp, err := wallet.GenKeypair(algo.Name())
	if err != nil {
		return walletResp, fmt.Errorf("unable to generate new key pair: %v", err)
	}

	// the user did not specify any metas
	// we'll create a default one for them
	if len(kp.Meta) <= 0 {
		kp.Meta = append(
			kp.Meta,
			wallet.Meta{
				Key:   "name",
				Value: fmt.Sprintf("%v's key %v", owner, len(wal.Keypairs)+1),
			},
		)
	}

	// now updating the wallet and saving it
	_, err = wallet.AddKeypair(kp, rootPath, owner, passphrase)
	if err != nil {
		return walletResp, fmt.Errorf("unable to add keypair to wallet: %v", err)
	}

	walletTmp, _ := wallet.Read(rootPath, owner, passphrase)
	walletResp.Wallet = *walletTmp

	return walletResp, nil
}

func getWallet(owner string, passphrase string) (getWalletResponse, error) {
	var walletResp getWalletResponse

	if len(owner) <= 0 {
		return walletResp, errors.New("wallet name is required")
	}

	if ok, err := fsutil.PathExists(walletPath); !ok {
		return walletResp, fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, owner, passphrase)
	if err != nil {
		return walletResp, fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	// print the new keys for user info
	walletResp.Wallet = *wal

	return walletResp, nil
}

func signTransaction(owner string, passphrase string, pub string, message string) (signTransactionResponse, error) {
	var walletResp signTransactionResponse

	if len(owner) <= 0 {
		return walletResp, errors.New("wallet name is required")
	}

	if len(pub) <= 0 {
		return walletResp, errors.New("pubkey is required")
	}
	if len(message) <= 0 {
		return walletResp, errors.New("data is required")
	}

	if ok, err := fsutil.PathExists(walletPath); !ok {
		return walletResp, fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, owner, passphrase)
	if err != nil {
		return walletResp, fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	dataBuf, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return walletResp, fmt.Errorf("invalid base64 encoded data: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == pub {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return walletResp, fmt.Errorf("unknown public key: %v", pub)
	}
	if kp.Tainted {
		return walletResp, fmt.Errorf("key is tainted: %v", pub)
	}

	alg, err := crypto.NewSignatureAlgorithm(crypto.Ed25519)
	if err != nil {
		return walletResp, fmt.Errorf("unable to instanciate signature algorithm: %v", err)
	}
	sig, err := wallet.Sign(alg, kp, dataBuf)
	if err != nil {
		return walletResp, fmt.Errorf("unable to sign: %v", err)
	}
	fmt.Printf("%v\n", base64.StdEncoding.EncodeToString(sig))
	walletResp.Message = base64.StdEncoding.EncodeToString(sig)

	return walletResp, nil
}

func taintKeypair(owner string, passphrase string, pub string) (taintKeypairResponse, error) {
	var walletResp taintKeypairResponse

	if len(owner) <= 0 {
		return walletResp, errors.New("wallet name is required")
	}

	if ok, err := fsutil.PathExists(walletPath); !ok {
		return walletResp, fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, owner, passphrase)
	if err != nil {
		return walletResp, fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == pub {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return walletResp, fmt.Errorf("unknown public key: %s", pub)
	}

	if kp.Tainted {
		return walletResp, fmt.Errorf("key %s is already tainted", pub)
	}

	kp.Tainted = true

	wal, err = wallet.Write(wal, rootPath, owner, passphrase)
	if err != nil {
		return walletResp, err
	}

	// print the new keys for user info
	walletResp.Wallet = *wal

	return walletResp, nil
}

func verifyTransaction(owner string, passphrase string, pub string, message string, sig string) (verifyTransactionResponse, error) {
	var walletResp verifyTransactionResponse

	if len(owner) <= 0 {
		return walletResp, errors.New("wallet name is required")
	}

	if len(pub) <= 0 {
		return walletResp, errors.New("pubkey is required")
	}
	if len(message) <= 0 {
		return walletResp, errors.New("message is required")
	}
	if len(sig) <= 0 {
		return walletResp, errors.New("data is required")
	}

	if ok, err := fsutil.PathExists(walletPath); !ok {
		return walletResp, fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, owner, passphrase)
	if err != nil {
		return walletResp, fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	dataBuf, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return walletResp, fmt.Errorf("invalid base64 encoded data: %v", err)
	}
	sigBuf, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return walletResp, fmt.Errorf("invalid base64 encoded data: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == pub {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return walletResp, fmt.Errorf("unknown public key: %v", pub)
	}

	alg, err := crypto.NewSignatureAlgorithm(crypto.Ed25519)
	if err != nil {
		return walletResp, fmt.Errorf("unable to instanciate signature algorithm: %v", err)
	}
	verified, err := wallet.Verify(alg, kp, dataBuf, sigBuf)
	if err != nil {
		return walletResp, fmt.Errorf("unable to verify: %v", err)
	}
	walletResp.Verified = verified

	return walletResp, nil
}

func checkBalance(owner string, nodeURLGrpc string) (checkBalanceResponse, error) {
	var balance checkBalanceResponse

	if len(owner) <= 0 {
		return balance, errors.New("wallet name is required")
	}

	if len(nodeURLGrpc) == 0 {
		return balance, errors.New("gRPC URL is required")
	}

	conn, err := grpc.Dial(nodeURLGrpc, grpc.WithInsecure())
	if err != nil {
		return balance, err
	}
	defer conn.Close()

	return balance, nil
}
