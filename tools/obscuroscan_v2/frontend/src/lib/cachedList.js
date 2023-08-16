class CachedList {
    constructor() {
        this.items = [];
    }

    add(item) {
        if (!this.items.some(i => i.hash === item.hash)) {
            this.items.push(item);
        }
    }

    get() {
        return this.items.slice(-5);
    }
}

export default CachedList