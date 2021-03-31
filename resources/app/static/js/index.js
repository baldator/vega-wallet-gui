let index = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    licence: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
            index.checkLogin();
            config.getConfig();
            // Update service status every 5 seconds
            var intervalId = window.setInterval(function() {
                wallet.getServiceStatus();
            }, 5000);

        })
    },
    openBrowser: function(url) {
        require("electron").shell.openExternal(url);
    },
    openVegaConsole: function() {
        this.openBrowser(config.vegaConsoleUrl);
    },
    hidePages: function() {
        [].forEach.call(document.querySelectorAll('.page'), function(el) {
            el.style.display = 'none';
        });
    },
    hideRightPane: function() {
        [].forEach.call(document.querySelectorAll('.right-pane'), function(el) {
            el.style.display = 'none';
        });
    },
    hideConfigurationStatus: function() {
        [].forEach.call(document.querySelectorAll('.configuration-status'), function(el) {
            el.style.display = 'none';
        });
    },
    showMain: function() {
        document.getElementById("main").style.display = "block";
    },
    showManageKeys: function() {
        this.hideRightPane();
        document.getElementById("manage-keys").style.display = "block";
    },
    showCheckBalance: function() {
        this.hideRightPane();
        document.getElementById("check-balance").style.display = "block";
    },
    showServiceStatus: function() {
        this.hideRightPane();
        document.getElementById("service-status").style.display = "block";
    },
    showSignTransaction: function() {
        this.hideRightPane();
        document.getElementById("sign-transaction").style.display = "block";
    },
    showVerifyTransaction: function() {
        this.hideRightPane();
        document.getElementById("verify-transaction").style.display = "block";
    },

    checkLogin: function() {
        if (wallet.walletOwner == "" || wallet.walletPassphrase == "") {
            this.hidePages();
            document.getElementById("login").style.display = "block";
            wallet.getWallets();
        } else {
            this.hidePages();
            document.getElementById("main").style.display = "block";
        }
    },
    loginWalletPage: function() {
        this.hidePages();
        document.getElementById("login").style.display = "block";
    },
    loginWalletAction: function() {
        // add parameters validation
        wallet.setWalletOwner(document.getElementById("login-owners").value);
        wallet.walletPassphrase = document.getElementById("login-passphrase").value;

        wallet.getWallet()
    },
    createWalletPage: function(showModal) {
        if (showModal) {
            asticode.notifier.info("No wallet found. Do you want to create a wallet?");
        }
        this.hidePages();
        document.getElementById("create-wallet").style.display = "block";

    },
    appendKeyPair: function(keypair) {
        let div = document.createElement("div");
        div.innerHTML = `<div class="manage-key-value">Public key: ${keypair.pub}</div><div class="manage-key-value">Private key: ${keypair.priv}</div>
        <div class="manage-key-value">Metadata: ${JSON.stringify(keypair.meta)}</div><div class="manage-key-value">Tainted: ${keypair.tainted}</div><div class="manage-key-value">Algorithm: ${keypair.algo}</div>`
        console.log("tainted: " + keypair.tainted)
        console.log(typeof keypair.tainted)
        if (keypair.tainted == "false" || keypair.tainted == false) {
            div.innerHTML += `<br><div><button type="submit" onclick="wallet.taintKeypair('${keypair.pub}')">Taint</button></div>`
        }
        div.innerHTML += `<hr/>`
        document.getElementById("manage-keys-keys").appendChild(div)
    },
    setWalletInfo: function(wallet) {
        document.getElementById("manage-keys-owner-value").innerHTML = wallet.Owner;
        document.getElementById("manage-keys-keys").innerHTML = ""

        for (let i = 0; i < wallet.Keypairs.length; i++) {
            this.appendKeyPair(wallet.Keypairs[i]);
        }

        document.getElementById("sign-transaction-pk").innerHTML = '<option value="" disabled selected>Select a public key</option>'
        document.getElementById("verify-transaction-pk").innerHTML = '<option value="" disabled selected>Select a public key</option>'
        document.getElementById("check-balance-pk").innerHTML = '<option value="" disabled selected>Select a public key</option>'

        // add keypairs in select dropdown
        for (let i = 0; i < wallet.Keypairs.length; i++) {
            let option = document.createElement("option");
            option.value = wallet.Keypairs[i].pub;
            option.text = wallet.Keypairs[i].pub;
            document.getElementById("sign-transaction-pk").appendChild(option)
        }
        for (let i = 0; i < wallet.Keypairs.length; i++) {
            let option = document.createElement("option");
            option.value = wallet.Keypairs[i].pub;
            option.text = wallet.Keypairs[i].pub;
            document.getElementById("verify-transaction-pk").appendChild(option)
        }
        for (let i = 0; i < wallet.Keypairs.length; i++) {
            let option = document.createElement("option");
            option.value = wallet.Keypairs[i].pub;
            option.text = wallet.Keypairs[i].pub;
            document.getElementById("check-balance-pk").appendChild(option)
        }

    },
    showSignedMessage: function(message) {
        document.getElementById("sign-transaction-response").style.display = "block";
        document.getElementById("sign-transaction-response-value").value = message;
    },
    showVerifiedMessage: function(verfied) {
        document.getElementById("verify-transaction-response").style.display = "block";
        document.getElementById("verify-transaction-response-value").value = verfied;
    },
    showConfiguration: function(message) {
        this.hideConfigurationStatus();
        if (message.PathExists) {
            document.getElementById("service-status-conf-ok").style.display = "block";
            document.getElementById("service-status-conf-ok-path").innerHTML = message.path;
        } else {
            document.getElementById("service-status-conf-nok").style.display = "block";
        }
    },
    showPositions: async function(accounts) {
        document.getElementById("check-balance-result-data").innerHTML = "";
        document.getElementById("check-balance-result").style.display = "block";
        for (let i = 0; i < accounts.accounts.length; i++) {
            let accountReadable = await wallet.getAssetValue(accounts.accounts[i].balance, accounts.accounts[i].asset);
            let divRow = document.createElement("div");
            divRow.classList.add("divTableRow");
            let divName = document.createElement("div");
            divName.classList.add("divTableCell");
            divName.innerHTML = accountReadable.name;
            divRow.appendChild(divName);
            let divValue = document.createElement("div");
            divValue.classList.add("divTableCell");
            divValue.innerHTML = accountReadable.value;
            divRow.appendChild(divValue);

            document.getElementById("check-balance-result-data").appendChild(divRow);
        }

        if (accounts.accounts.length == 0) {
            document.getElementById("check-balance-result").style.display = "none";
            document.getElementById("check-balance-result-empty").style.display = "block";
        } else {
            document.getElementById("check-balance-result-empty").style.display = "none";
        }
    },
    updateServiceStatus: function(status) {
        [].forEach.call(document.querySelectorAll('.service-status'), function(el) {
            el.style.display = 'none';
        });
        if (status) {
            document.getElementById("service-status-ok").style.display = "block";
        } else {
            document.getElementById("service-status-nok").style.display = "block";
        }
    },
    logout: function() {
        wallet.setWalletOwner("");
        wallet.walletPassphrase = "";
        this.hidePages();
        this.checkLogin();
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return { payload: "payload" };
                case "licence":
                    index.licence(message.payload)
                    return { payload: "licence" }
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    }
};