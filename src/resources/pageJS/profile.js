function reloadCache() {
    $("button#reload").addClass("is-loading")
    $.post("/admin/api/clearCache/user/{{$viewing.UID}}", function () {
        location.reload()
    });
}

function togglePosting() {
    $("button#cPost").addClass("is-loading")
    $.post("/admin/api/togglePosting/{{$viewing.UID}}", function () {
        location.reload()
    });
}

function setRank(rank) {
    $(`button#rank_${rank}`).addClass("is-loading")
    $.post("/admin/api/updateRank/{{$viewing.UID}}", {rank: rank}, function () {
        location.reload()
    });
}

//todo move to the admin page

sessionStorage.setItem("currentPage", window.location.pathname)
