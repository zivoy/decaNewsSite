if (window.history.replaceState) {
    window.history.replaceState(null, null, window.location.href);
}

parserTags["youtube"] = {
    openTag: function (params, content) {

        return `<figure class="image"><iframe frameborder="0" style="height:max(225px,100%); width:min(100%,400px);" allowfullscreen src="https://www.youtube.com/embed/${content}">`;
    },
    closeTag: function () {
        return "</iframe></figure>";
    },
    content: function () {
        return "";
    }
}

const inputTime = $("input#leakTime");
const inputLink = $("input#link");
const linkInfo = $("p#linkInfo");
const submitButton = $("button#submit");
const inputLeak = $("textarea#leak");
const inputImage = $("input#image");
const preview = $("div#preview");
const inputTitle = $("input#title");
const tagsInput = $("input#tags");

LeakTime.setSeconds(0);
LeakTime.setMilliseconds(0);
$("time#timeToSet").attr("datetime", LeakTime.toISOString());

BulmaTagsInput.attach(tagsInput[0], {
    source: async function () {
        return $.get("/api/v1/tags/get").then(function (list) {
            return list
        })
    },
    closeDropdownOnItemSelect: false,
    selectable: false,
})

$(document).ready(function () {
    setTimeVal(inputTime, LeakTime);

    inputTime.change(() => {
        LeakTime = timeChange(inputTime, $("time#timeToSet"));
        formReady();
    });

    inputLink.on("input", () => {
        Link = linkChange(inputLink, linkInfo, $("a#source"), allowedLinks);
        formReady();
    });

    inputLeak.on("input", () => {
        Leak = leakChange(inputLeak, $("p#leakPreview"),
            val => {
                return BBCodeParser.process(val.replaceAll("\n", "<br>"))
            });
        formReady();
    })

    inputImage.on("input", () => {
        Image = "";
        $("div.dropdown").addClass("is-up");
        imageChange(inputImage, $("img#previewImage")).then(r => {
            Image = r
            $("div.dropdown").removeClass("is-up")
        });
        formReady();
    })

    inputTitle.on("input", () => {
        Title = titleChange(inputTitle, $("b#leakTitle"), "DecaLeak ##########");
        formReady();
    });
});

function formReady() {
    preview.removeClass("is-invisible");
    stopLeave = true;
    submitButton.prop("disabled", true);
    if (Leak && LeakTime && User && canLinkLess()) {
        submitButton.prop("disabled", false);
    }
}

function post() {
    stopLeave = false;
    let form = $('<form></form>');

    form.attr("method", "post");
    form.addClass("is-hidden");
    form.attr("action", "/leaks/create");

    $.each({
        description: Leak,
        time: LeakTime.getTime(),
        image_url: Image,
        source_url: Link,
        title: Title,
        tags: tagsInput.val()
    }, function (key, value) {
        let field = $('<input/>');
        field.attr("type", "hidden");
        field.attr("name", key);
        field.attr("value", value);

        form.append(field);
    });

    $(document.body).append(form);
    form.submit();
}

