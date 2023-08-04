class Config {
    // VITE_APIHOSTADDRESS should be used as an env var at the prod server
    static backendServerAddress = import.meta.env.VITE_APIHOSTADDRESS
    static pollingInterval = 5000
    static pricePollingInterval = 10*this.pollingInterval
}

export default Config