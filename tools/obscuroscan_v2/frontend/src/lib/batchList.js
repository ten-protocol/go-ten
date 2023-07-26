class BatchList {
    constructor() {
        this.items = [];
    }

    add(item) {
        if (!this.items.includes(item)) {
            this.items.push(item);
        }
    }

    get() {
        return this.items.slice(-5);
    }
}

export default BatchList