async function getTagList() {
    return $.get("/api/v1/tags/get").then(function (list) {
        return list
    })
}
