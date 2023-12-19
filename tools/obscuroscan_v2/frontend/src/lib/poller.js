class Poller {
    constructor(fetchCallback, interval) {
        this.fetchCallback = fetchCallback;
        this.interval = interval;
        this.timer = null;
    }

    // Start polling - executes the fetchCallback immediately and every interval thereafter
    start() {
        this.stop(); // Ensure previous intervals are cleared
        this.fetchCallback();
        this.timer = setInterval(async () => {
            await this.fetchCallback();
        }, this.interval);
    }

    stop() {
        if (this.timer) {
            clearInterval(this.timer);
            this.timer = null;
        }
    }
}

export default Poller