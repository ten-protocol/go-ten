class Poller {
    constructor(fetchCallback, interval) {
        this.fetchCallback = fetchCallback;
        this.interval = interval;
        this.timer = null;
    }

    start() {
        this.stop(); // Ensure previous intervals are cleared
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