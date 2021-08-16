if (window.history.replaceState) {
    window.history.replaceState(null, null, window.location.href);
}

parserTags["youtube"] = {
    openTag: function(params,content) {

        return `<figure class="image"><iframe frameborder="0" style="height:max(225px,100%); width:min(100%,400px);" allowfullscreen src="https://www.youtube.com/embed/${content}">`;
    },
    closeTag: function() {
        return "</iframe></figure>";
    },
    content: function() {
        return "";
    }
}

let LeakTime = new Date();
let Link = "";
let Leak = "";
let Image = "";

const inputTime = $("input#leakTime");
const inputLink = $("input#link");
const linkInfo = $("p#linkInfo");
const submitButton = $("button#submit");
const inputLeak = $("textarea#leak");
const inputImage = $("input#image");
const preview = $("div#preview");

LeakTime.setSeconds(0);
LeakTime.setMilliseconds(0);
$("time#timeToSet").attr("datetime", LeakTime.toISOString());

$(document).ready(function () {
    setTimeVal(inputTime, LeakTime);

    inputTime.change(() => {
        preview.removeClass("is-invisible");
        LeakTime = timeChange(inputTime, $("time#timeToSet"));
        formReady();
    });

    inputLink.on("input", () => {
        preview.removeClass("is-invisible");
        Link = linkChange(inputLink, linkInfo, $("a#source"), allowedLinks);
        formReady();
    });

    inputLeak.on("input", () => {
        preview.removeClass("is-invisible");
        Leak = leakChange(inputLeak, $("p#leakPreview"),
            val=>{return BBCodeParser.process(val.replaceAll("\n", "<br>"))});
        formReady();
    })

    inputImage.on("input", () => {
        preview.removeClass("is-invisible");
        imageChange(inputImage, $("img#previewImage")).then(r => {
            Image = r
        });
        formReady();
    })
});

function formReady() {
    submitButton.prop("disabled", true);
    let linkValid =
        //{{ if ge .user.AuthLevel $linkLess }}
        true;
    //{{else}}
    Link;
    //{{end}}
    if (Leak && LeakTime && User && linkValid) {
        submitButton.prop("disabled", false);
    }
}

function post() {
    let form = $('<form></form>');

    form.attr("method", "post");
    form.addClass("is-hidden");
    form.attr("action", "/leaks/create");

    $.each({
        description: Leak,
        time: LeakTime.getTime(),
        image_url: Image,
        source_url: Link
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

