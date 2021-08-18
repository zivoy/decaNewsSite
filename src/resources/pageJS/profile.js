function reloadCache() {
    $("button#reload").addClass("is-loading")
    $.post(`/admin/api/clearCache/user/${viewerID}`, function () {
        location.reload()
    });
}

function togglePosting() {
    $("button#cPost").addClass("is-loading")
    $.post(`/admin/api/togglePosting/${viewerID}`, function () {
        location.reload()
    });
}

function setRank(rank) {
    $(`button#rank_${rank}`).addClass("is-loading")
    $.post(`/admin/api/updateRank/${viewerID}`, {rank: rank}, function () {
        location.reload()
    });
}

//todo move to the admin page

sessionStorage.setItem("currentPage", window.location.pathname)
