function _fixTime(element) {
    let rawValue = element.attr("datetime")
    let parsed = parseInt(rawValue)
    if (!isNaN(parsed) && parsed.toString() === rawValue) {
        rawValue = parsed
    }
    let date = new Date(rawValue)
    let formatted = date.toLocaleTimeString([], {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
    element.text(formatted)
}

function fixTime() {
    $("time").each(function () {
        _fixTime($(this))
    })

}

fixTime()