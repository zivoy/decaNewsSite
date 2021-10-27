function getStorageDefault(key, defaultVal, storage = localStorage) {
    let val = storage.getItem(key)
    if (!val) {
        storage.setItem(key, defaultVal)
        val = defaultVal
    }
    return val
}

function imageErr() {
    $('img').error(function () {
        $(this).attr('src', '/static/DecaFans-big.png').addClass('no-img');
    });
}

$(document).ready(function () {

    // Check for click events on the navbar burger icon
    $(".navbar-burger").click(function () {

        // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
        $(".navbar-burger").toggleClass("is-active");
        $(".navbar-menu").toggleClass("is-active");

    });
});

async function getTagList() {
    return $.get("/api/v1/tags/get").then(function (list) {
        return list
    })
}
