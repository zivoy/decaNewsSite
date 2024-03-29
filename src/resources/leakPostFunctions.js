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

function setTimeVal(field, now) {
    field.value(`${now.getFullYear()}/${now.getMonth()+1}/${now.getDate()} ${now.getHours()}:${now.getMinutes()}`)
}

function calenderOptions(leakTime) {
    return {
        color: "info",
        showHeader: false,
        showFooter: false,
        displayMode: "inline",//"dialog",
        // enableYearSwitch: false,
        startDate: leakTime,
        startTime: leakTime,
        maxDate: new Date(),
        dateFormat: 'yyyy/MM/dd',
        displayYearsCount: 5,
        minuteSteps: 1,
    };
}

function timeChange(field, preview) {
    let time;
    try {
        time = new Date(field.value());
        if (preview !== undefined)
            preview.attr("datetime", time.toISOString());
        fixTime()
    } catch (err) {
        time = null
        if (preview !== undefined)
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
    if (preview !== undefined)
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
    if (preview !== undefined)
        preview.html(processor(val))
    return val
}

async function imageChange(field, previewImage) {
    let container = previewImage.parent().parent();
    let val = field.val();

    function isValidImageUrl(url, callback) {
        if (val === "") {
            callback(true);
            return;
        }
        if (!url.includes("://")) {
            callback(false);
            return;
        }

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

    isValidImageUrl(val, function (result) {
        container.addClass("is-hidden")

        if (result) {
            valid = true
            field.removeClass("is-danger")
            if (previewImage !== undefined)
                previewImage.attr("src", val)
            if (val !== "") {
                container.removeClass("is-hidden")
            }
        } else {
            if (previewImage !== undefined)
                previewImage.attr("src", "")
            container.addClass("is-hidden")
            val = ""
        }
        done = true
    });
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