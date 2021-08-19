let stopLeave = false

window.addEventListener("beforeunload", function (e) {
    delete e["returnValue"]
    if (stopLeave === true) {
        // Cancel the event
        e.preventDefault();
        e.returnValue = "You have made changes to the page";

        return e.returnValue;
    }
});

function setTimeVal(field, date) {
    let now = new Date(date);
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    field.val(now.toISOString().slice(0, 16));
    field.attr("max", now.toISOString().slice(0, 16));
}

function timeChange(field, preview) {
    field.removeClass("is-danger");
    let time;
    try {
        time = new Date(field.val());
        preview.attr("datetime", time.toISOString());
        fixTime()
    } catch (err) {
        field.addClass("is-danger")
        time = null
        preview.text("invalid date")
    }
    return time
}

function titleChange(field, preview, def = "") {
    let val = field.val().trim();
    let prev = val
    if (val.length > 60) {
        prev = val.substring(0, 59) + "…"
    }
    if (val === "")
        prev = def
    preview.html(prev)
    return val
}

function linkChange(field, info, source, allowedLinks) {
    field.removeClass("is-danger is-success")
    info.removeClass("is-danger is-success")
    if (source !== undefined)
        source.addClass("has-text-danger")
    info.text("")
    let link = ""
    let val = field.val().trim();
    if (val === "") {
        field.addClass("is-danger")
    }

    if (allowedLinks(val).some(x => x)) {
        field.addClass("is-success")
        link = val
        if (source !== undefined) {
            source.removeClass("has-text-danger")
            source.attr("href", val)
        }
    } else {
        info.text("This does not seem to be an authorised source URL")
        info.addClass("is-danger")
        field.addClass("is-danger")
    }
    return link
}

function leakChange(field, preview, processor) {
    field.removeClass("is-danger")
    let val = field.val().trim();
    if (val === "") {
        field.addClass("is-danger")
    }
    preview.html(processor(val))
    return val
}

async function imageChange(field, previewImage) {
    let container = previewImage.parent().parent();
    let val = field.val();

    function isValidImageUrl(url, callback) {
        $('<img>', {
            src: url,
            load: function () {
                callback(true);
            },
            error: function () {
                callback(false);
            }
        });
    }

    let done = false
    let valid = false

    if (val !== "") {
        isValidImageUrl(val, function (result) {
            if (result) {
                valid = true
                field.removeClass("is-danger")
                previewImage.attr("src", val)
                container.removeClass("is-hidden")
            } else {
                previewImage.attr("src", "")
                container.addClass("is-hidden")
                field.addClass("is-danger")
                val = ""
            }
            done = true
        });
    }
    await new Promise(r => {
        let timeout = () => {
            if (done) {
                if (valid)
                    r()
            } else {
                setTimeout(timeout, 10)
            }
        }
        timeout()
    });
    return val
}

function AllowedLinks(ListOfLinks) {
    return function (val) {
        let list = [];
        for (let i of ListOfLinks) {
            list.push(val.match(new RegExp(i)));
        }
        return list;
    }
}