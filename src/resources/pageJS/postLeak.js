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

const inputLink = $("input#link");
const linkInfo = $("p#linkInfo");
const submitButton = $("button#submit");
const inputLeak = $("textarea#leak");
const inputImage = $("input#image");
const preview = $("div#preview");
const inputTitle = $("input#title");
const tagsInput = $("input#tags");
const tagList = $("div#tagsList")

LeakTime.setSeconds(0);
LeakTime.setMilliseconds(0);
$("time#timeToSet").attr("datetime", LeakTime.toISOString());

BulmaTagsInput.attach(tagsInput[0], {
    source: getTagList,
    closeDropdownOnItemSelect: false,
    selectable: false,
    tagClass: "",
})

const tagSelector = tagsInput[0].BulmaTagsInput()
let calendar

$(document).ready(function () {
    calendar = bulmaCalendar.attach('[type="datetime-local"]', calenderOptions(LeakTime))[0];

    calendar.on('select', () => {
        calendar.value()
        LeakTime = timeChange(calendar, $("time#timeToSet"));
        formReady();
        // calendar.refresh()
    });
    // calendar.on('hide', (e,s) => {
    //     calendar.save()
    // });

    LeakTime = timeChange(calendar);

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

    function tagUpdate() {
        tagList.empty()
        let tags = tagSelector.items
        if (tags.length > 0) {
            tagList.removeClass("is-hidden");
        } else {
            tagList.addClass("is-hidden");
        }
        for (let i in tags) {
            let tag = tags[i]
            $("<span>", {
                class: "tag " + "TODOColor",
            }).append(tag).appendTo(tagList)
        }
        formReady();
    }

    tagSelector.on("after.add", tagUpdate)
    tagSelector.on("after.remove", tagUpdate)
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

