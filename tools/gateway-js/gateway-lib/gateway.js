const axios = require("axios");

const obscuroGatewayVersion = "v1";
const pathJoin = obscuroGatewayVersion + "/join/";
const pathAuthenticate = obscuroGatewayVersion + "/authenticate/";

class Gateway {
    constructor(httpURL, wsURL, provider) {
        this.httpURL = httpURL;
        this.wsURL = wsURL;
        this.userId = '';
        this.provider = provider;
    }

    async join() {
        try {
            const response = await axios.get(`${this.httpURL}${pathJoin}`);
            if (response.status !== 200) {
                throw new Error(`Failed to get userID. Status code: ${response.status}`);
            }
            // todo make further checks on the data
            this.userId = response.data;
        } catch (error) {
            throw new Error(`Failed to get userID. ${error}`);
        }
    }

    async registerAccount(address) {
        const message = `Register ${this.userID} for ${address.toLowerCase()}`;
        let signature = ""

        try {
                signature = await this.provider.request({
                    method: "personal_sign",
                    params: [message, address]
                })
            } catch (err) {
                throw new Error(`Failed to sign message. ${err}`);
            }

        // todo make further checks on the data
        if (signature === -1) {
            return "Signing failed"
        }

        try {
            const authenticateUserURL = pathAuthenticate+"?token="+this.userId
            const authenticateFields = {"signature": signature, "message": message}
            const authenticateResp = await axios.post(
                authenticateUserURL,
                authenticateFields,
                {
                    headers: {
                        "Accept": "application/json",
                        "Content-Type": "application/json"
                    }
                }
            );
            // todo make further checks on the data
            return await authenticateResp.text()
        } catch (error) {
            throw new Error(`Failed to register account. ${error}`);
        }
    }

    userID() {
        return this.userId;
    }

    http() {
        return `${this.httpURL}/v1/?token=${this.userId}`;
    }

    ws() {
        return `${this.wsURL}/v1/?token=${this.userId}`;
    }
}

module.exports = Gateway;