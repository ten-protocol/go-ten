class CachedList {
    constructor() {
        // todo these should just be one storage
        this.items = [];
        this.itemsByHash = {};
    }

    add(item) {
        if (!this.items.some(i => i.hash === item.hash)) {
            this.items.push(item);
        }
    }

    get() {
        return this.items.slice(-5);
    }

    addByHash(item) {
        this.itemsByHash[item.Header.hash] = item
    }

    getByHash(hash) {
        return this.itemsByHash[hash]
    }
}

export default CachedList