class Config {
    // VITE_APIHOSTADDRESS should be used as an env var at the prod server
    static backendServerAddress = import.meta.env.VITE_APIHOSTADDRESS
    static pollingInterval = 1000
    static pricePollingInterval = 10*this.pollingInterval
    static version = import.meta.env.VITE_FE_VERSION
}

export default Config