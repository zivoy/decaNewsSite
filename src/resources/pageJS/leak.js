function openModal() {
    $("#archive-modal").addClass("is-active");
}

function closeModal() {
    $("#archive-modal").removeClass("is-active");
}

function archivePost() {
    $("button#archive").addClass("is-loading")
    $.post(`/api/v1/archive/${id}`, function () {
        // maybe add a message ¯\_(ツ)_/¯
        window.location.href = "/leaks";
    }).fail(function () {
        $("button#archive").removeClass("is-loading")
        alert("Error archiving leak");
    });
}


let submitButton;
let editing = false;

function startEdit() {
    $("<link/>", {
        rel: "stylesheet",
        type: "text/css",
        href: "/static/css/poster.css"
    }).appendTo("head");

    $("#leakBody").html($(`
            <fieldset>
                <div class="field">
                    <label class="label">Leak Title</label>
                    <div class="control">
                        <input id="title" name="title" class="input" type="text" placeholder="${dTitle}"
                               value="${title}">
                    </div>
                </div>
            
                <div class="field">
                    <label class="label">Leak Description</label>
                    <div class="control">
                        <textarea id="leak" name="leak" class="textarea"
                                  placeholder="Description of leak, be as descriptive as you want" required>${leak}</textarea>
                    </div>
                    <p class="help" id="bbcode"></p>
                </div>

                <div class="field">
                    <label class="label">Source Link</label>
                    <div class="control">
                        <input id="link" name="link" class="input" type="url" placeholder="Link to the leak in the chat"
                               value="${link}" required>
                    </div>
                    <p class="help" id="linkInfo"></p>
                </div>

                <div class="field">
                    <label class="label">Image</label>
                    <div class="control">
                        <input class="input" type="url" id="image" name="image" placeholder="Optional link to image"
                            value="${image}">
                    </div>
                </div>

                <div class="field">
                    <label class="label">Time</label>
                    <div class="control">
                        <input class="input" type="datetime-local" id="leakTime" name="leakTime" required>
                    </div>
                </div>
                
                <div class="field">
                    <label class="label">Tags</label>
                    <div class="control">
                        <input class="input" type="text" id="tags" data-max-chars="20" placeholder="Choose Tags" 
                        value="${tags}">
                    </div>
                    <p class="help">Tags can be at most 20 characters long</p>
                </div>

            </fieldset>
        `));
    $("p#bbcode").html("supports <a href='https://www.bbcode.org/reference.php' target='_blank' rel='noopener'>bbcode</a>")
    $("#leakImage").remove();
    $("#controlButtons").html(`
            <button class="button is-warning" id="cancelEdit" onclick="location.reload();">Cancel</button>
            <button class="button is-primary" id="saveEdit" onclick="UpdateLeak()">Save</button>
        `);

    editing = true;

    const inputTitle = $("input#title");
    const inputLink = $("input#link");
    const linkInfo = $("p#linkInfo");
    const inputLeak = $("textarea#leak");
    const inputImage = $("input#image");
    const tagsInput = $("input#tags");

    $.getScript("https://cdn.jsdelivr.net/npm/@creativebulma/bulma-tagsinput/dist/js/bulma-tagsinput.min.js",
        function () {
            BulmaTagsInput.attach(tagsInput[0], {
                source: getTagList,
                closeDropdownOnItemSelect: false,
                selectable: false,
                tagClass: "",
            })
        });

    $.getScript("https://cdn.jsdelivr.net/npm/bulma-calendar/dist/js/bulma-calendar.min.js",
        function () {
            let calendar = bulmaCalendar.attach('[type="datetime-local"]', calenderOptions(time))[0];
            setTimeVal(calendar, time)

            calendar.on('select', () => {
                calendar.value()
                time = timeChange(calendar, $("time#timeToSet"));
                formReady();
            });
        });


    submitButton = $("button#saveEdit");

    inputTitle.on("input", () => {
        title = titleChange(inputTitle, $("b#leakTitle"), dTitle);
        formReady();
    });

    inputLink.on("input", () => {
        link = linkChange(inputLink, linkInfo, undefined, allowedLinks);
        formReady();
    });

    inputLeak.on("input", () => {
        leak = leakChange(inputLeak, $("p#leakPreview"), val => {
            return val.replaceAll("\n", "<br>")
        });
        formReady();
    })

    inputImage.on("input", () => {
        imageChange(inputImage, $("img#previewImage")).then(r => {
            image = r
        });
        formReady();
    })

    function formReady() {
        stopLeave = true;
        submitButton.prop("disabled", true);
        // maybe add a check for source url :/
        if (leak && time) {
            submitButton.prop("disabled", false);
        }
    }
}

function UpdateLeak() {
    if (editing) {
        stopLeave = false;
        submitButton.addClass("is-loading")
        $.ajax({
            url: `/leaks/update/${id}`,
            type: "post",
            data: {
                description: leak,
                time: time.getTime(),
                image_url: image,
                source_url: link,
                title: title,
                tags: $("input#tags").val()
            },
            success: function () {
                submitButton.removeClass("is-loading")
                location.reload()
            }
        });
    }
}

$("#backButton").click(function () {
    window.location.href = getStorageDefault("lastPage", "/leaks", sessionStorage);
});
