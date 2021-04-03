let wallet = {
    getWallets: function() {
        // Create message
        let message = { "name": "getWallets" };

        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            if (message.payload.owners == null || message.payload.owners.length == 0) {
                index.createWalletPage(true)
            }

            // Empty list
            document.getElementById("login-owners").innerHTML = "";

            // Process login
            for (let i = 0; i < message.payload.owners.length; i++) {
                let option = document.createElement("option");
                option.value = message.payload.owners[i];
                option.text = message.payload.owners[i];
                document.getElementById("login-owners").appendChild(option)
            }
        })
    },
    walletOwner: "",
    setWalletOwner: function(owner) {
        this.walletOwner = owner
        document.getElementById("current-wallet").innerHTML = owner
    },
    walletPassphrase: "",
    getWallet: function() {
        // Create message
        let message = { name: "getWallet", payload: { owner: this.walletOwner, passphrase: this.walletPassphrase } };
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            index.hidePages();
            index.showMain();
            index.setWalletInfo(message.payload.wallet);
        })

    },
    createWallet: function() {
        let owner = document.getElementById("create-owner").value;
        let passphrase = document.getElementById("create-passphrase").value;
        let passphraseBis = document.getElementById("create-passphrase-bis").value;

        if (owner.length == 0) {
            asticode.notifier.info("Owner cannot be empty.");
            return
        }

        if (passphrase.length == 0) {
            asticode.notifier.info("Passphrase cannot be empty.");
            return;
        }

        if (passphrase != passphraseBis) {
            asticode.notifier.info("The two passphrases are not equals. Please fix it.");
            return;
        }

        // Create message
        let message = { "name": "createWallet", payload: { owner: owner, passphrase: passphrase } };
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            index.hidePages();
            index.loginWalletPage();
            asticode.notifier.info("Wallet created successfully.");
        })
    },
    addKeypair: function() {
        // Create message
        let message = { "name": "addKeypair", payload: { owner: this.walletOwner, passphrase: this.walletPassphrase } };
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            index.hidePages();
            wallet.getWallet();
            asticode.notifier.info("Keypair created successfully.");
        })
    },
    taintKeypair: function(pub) {
        // Create message
        let message = { "name": "taintkeypair", payload: { owner: this.walletOwner, passphrase: this.walletPassphrase, pub: pub } };
        // Send message
        asticode.loader.show();
        console.log(message);
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            index.hidePages();
            wallet.getWallet();
            asticode.notifier.info("Keypair created successfully.");
        })
    },
    signTransaction: function() {
        // Create message
        let signMessage = document.getElementById("sign-transaction-message").value;
        let pub = document.getElementById("sign-transaction-pk").value;
        let message = { "name": "signtransaction", payload: { owner: this.walletOwner, passphrase: this.walletPassphrase, pub: pub, message: signMessage } };
        // Send message
        asticode.loader.show();
        console.log(message);
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            index.showSignedMessage(message.payload.message);
            asticode.notifier.info("Message signed successfully.");
        })
    },
    verifyTransaction: function() {
        // Create message
        let signMessage = document.getElementById("verify-transaction-message").value;
        let signature = document.getElementById("verify-transaction-signature").value;
        let pub = document.getElementById("verify-transaction-pk").value;
        let message = { "name": "verifytransaction", payload: { owner: this.walletOwner, passphrase: this.walletPassphrase, pub: pub, message: signMessage, signature: signature } };
        // Send message
        asticode.loader.show();
        console.log(message);
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            index.showVerifiedMessage(message.payload.verified);
            asticode.notifier.info("Message verfied successfully.");
        })
    },
    checkBalance: async function() {
        let pub = document.getElementById("check-balance-pk").value;
        if (pub == null || pub == "") {
            asticode.notifier.error("Public key is required");
        }

        let response = await fetch('https://lb.testnet.vega.xyz/parties/' + pub + '/accounts');
        const accounts = await response.json(); //extract JSON from the http response

        index.showAccounts(accounts);
    },
    getAssetValue: async function(balance, asset) {
        let assetsResp = await fetch('https://lb.testnet.vega.xyz/assets');
        const assetsVal = await assetsResp.json(); //extract JSON from the http response
        for (let i = 0; i < assetsVal.assets.length; i++) {
            if (assetsVal.assets[i].id == asset) {
                let val = balance / (10 ** assetsVal.assets[i].decimals)
                return {
                    "name": assetsVal.assets[i].name,
                    "value": val,
                }
            }
        }
        return null
    },
    getMarketValue: async function(marketID) {
        let marketsResp = await fetch('https://lb.testnet.vega.xyz/markets');
        const marketsVal = await marketsResp.json(); //extract JSON from the http response
        for (let i = 0; i < marketsVal.markets.length; i++) {
            if (marketsVal.markets[i].id == marketID) {

                return marketsVal.markets[i]
            }
        }
        return null
    },
    getPositions: async function() {
        let pub = document.getElementById("check-balance-pk").value;
        if (pub == null || pub == "") {
            asticode.notifier.error("Public key is required");
        }

        let response = await fetch('https://lb.testnet.vega.xyz/parties/' + pub + '/positions');
        const accounts = await response.json(); //extract JSON from the http response

        index.showPositions(accounts);
    },
    getServiceStatus: async function() {
        url = "http://" + config.walletConfig.Host + ":" + config.walletConfig.Port + "/api/v1/status"
        fetch(url)
            .then(function(response) {
                response.json().then(function(resp) {
                    index.updateServiceStatus(resp.success);
                });
            }).catch(function() {
                index.updateServiceStatus(false)
            });
    },
    startService: function() {
        // Create message
        let message = { "name": "startService" };

        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            } else {
                asticode.notifier.info("Service started successfully");
            }

        })
    },
    stopService: function() {
        // Create message
        let message = { "name": "stopService" };

        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            console.log(message);

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            } else {
                asticode.notifier.info("Service stopped successfully");
            }

        })
    },

};