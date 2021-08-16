function getStorageDefault(key, defaultVal, storage = localStorage) {
    let val = storage.getItem(key)
    if (!val) {
        storage.setItem(key, defaultVal)
        val = defaultVal
    }
    return val
}
