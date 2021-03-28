let config = {
    getConfig: function() {
        // Create message
        let message = { "name": "getconfig" };

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

            index.showConfiguration(message.payload);
        })
    },
    initConfig: function(force = false) {
        // Create message
        let message = { "name": "initconfig", payload: { force: force, genRsaKey: false } };

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

            asticode.notifier.info("Configuration file created successfully.");
            config.getConif();
        })
    },
}