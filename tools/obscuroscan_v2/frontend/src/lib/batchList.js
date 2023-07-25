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
        const startIdx = this.items.length - 5 > 0 ? this.items.length - 5 : 0;
        return this.items.slice(startIdx);
    }
}

export default BatchList