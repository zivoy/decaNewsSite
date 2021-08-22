window.addEventListener('load', function () {
    prevButton.click(function () {
        goToPage(currPage - 1)
    });
    nextButton.click(function () {
        goToPage(currPage + 1)
    });
    $(".dropdown .button").click(function () {
        $(this).parents('.dropdown').toggleClass('is-active');
    });

    excludeOr.click(function () {
        sessionStorage.setItem("excludeFilterRadio", "or")
        load()
    })
    excludeAnd.click(function () {
        sessionStorage.setItem("excludeFilterRadio", "and")
        load()
    })
    includeOr.click(function () {
        sessionStorage.setItem("includeFilterRadio", "or")
        load()
    })
    includeAnd.click(function () {
        sessionStorage.setItem("includeFilterRadio", "and")
        load()
    })

    let includeUpdate = function () {
        sessionStorage.setItem("includeList", includeTagSelector.value)
        includeTags = splitList(includeTagSelector.value)
        load()
    }
    includeTagSelector.on("after.add", includeUpdate)
    includeTagSelector.on("after.remove", includeUpdate)
    let excludeUpdate = function () {
        sessionStorage.setItem("excludeList", excludeTagSelector.value)
        excludeTags = splitList(excludeTagSelector.value)
        load()
    }
    excludeTagSelector.on("after.add", excludeUpdate)
    excludeTagSelector.on("after.remove", excludeUpdate)

    populateAmountList()
    itemList.removeClass("is-hidden")
    articlesRaw = itemList.children()

    updateNumberOfPages()
    goToPage(currPage, false, true, true)
});
document.addEventListener("click", function (e) {
    let amountDropDownCLicked = false
    let filterDropDownClicked = false
    let includeFilter = false
    let includeSelector = false
    let excludeFilter = false
    let excludeSelector = false

    clickedOn(e, amountDropdown[0], function () {
        amountDropDownCLicked = true
    })

    clickedOn(e, filterDropdown[0], function () {
        filterDropDownClicked = true
    })
    clickedOn(e, includeTagSelector.input, function () {
        includeFilter = true
    })
    clickedOn(e, includeTagSelector.dropdown, function () {
        includeSelector = true
    })
    clickedOn(e, excludeTagSelector.input, function () {
        excludeFilter = true
    })
    clickedOn(e, excludeTagSelector.dropdown, function () {
        excludeSelector = true
    })

    if (!amountDropDownCLicked) {
        amountDropdown.removeClass("is-active")
    }
    if (!filterDropDownClicked) {
        let deleteButton = $(e.target).data("tag");
        if (deleteButton !== undefined)
            return
        filterDropdown.removeClass("is-active")
    } else {
        includeTagSelector.dropdown.hidden = !(includeFilter || includeSelector)
        excludeTagSelector.dropdown.hidden = !(excludeFilter || excludeSelector)
    }
});

function clickedOn(e, targetNode, callback) {
    let targetElement = e.target
    do {
        if (targetElement === targetNode) {
            callback()
            return;
        }
        // Go up the DOM
        targetElement = targetElement.parentNode;
    } while (targetElement);
}

window.onpopstate = function (e) {
    if (e.state) {
        goToPage(e.state.leak, false, true)
    }
};


const choices = [-1, 5, 10, 20, 50]


const prevButton = $("a.pagination-previous")
const nextButton = $("a.pagination-next")
const amountDropdown = $("#amountDropdown")
const filterDropdown = $("#filterItemsDropdown")
const itemList = $("#leakItems")
const pagination = $("nav.pagination")

// filter
let includeAnd = $("#includeAnd");
let includeOr = $("#includeOr");
let excludeAnd = $("#excludeAnd");
let excludeOr = $("#excludeOr");
let includeItems = $("#includeTags");
let excludeItems = $("#excludeTags");

BulmaTagsInput.attach(includeItems[0], {
    source: getTagList,
    closeDropdownOnItemSelect: true,
    selectable: false,
    tagClass: "",
    freeInput: false
})

BulmaTagsInput.attach(excludeItems[0], {
    source: getTagList,
    closeDropdownOnItemSelect: true,
    selectable: false,
    tagClass: "",
    freeInput: false
})

let includeTagSelector = includeItems[0].BulmaTagsInput()
let excludeTagSelector = excludeItems[0].BulmaTagsInput()

let perPage = parseInt(getStorageDefault("amountPerPage", choices[0]));
let numberOfPages = 0
let articles = []
let articlesRaw = {}
let _iTags = getStorageDefault("includeList", "", sessionStorage)
includeTagSelector.add(_iTags)
let includeTags = splitList(_iTags)
_iTags = getStorageDefault("excludeList", "", sessionStorage)
excludeTagSelector.add(_iTags)
let excludeTags = splitList(_iTags)

const ellipsis = `<li><span class="pagination-ellipsis">&hellip;</span></li>`;
const currentPage = `<li><a class="pagination-link is-current" aria-label="Page {p}" aria-current="page">{p}</a></li>`;
const gotoPage = `<li><a class="pagination-link" aria-label="Goto page {p}">{p}</a></li>`;

function paginate(page, high) {
    let low = 1;
    let items = [];

    function addPage(number) {
        if (number === page) {
            return currentPage.replaceAll("{p}", number);
        }
        return gotoPage.replaceAll("{p}", number);
    }

    let i
    // low side
    if (page - low < 4) {
        for (i = low; i < low + 5 && i <= high; i++)
            items.push(addPage(i));
        i--;
        if (i !== high) {
            items.push(ellipsis);
            items.push(addPage(high));
        }
    }
    // high side
    else if (high - page < 4) {
        for (i = high; i > high - 5 && i >= low; i--)
            items.unshift(addPage(i));
        i++;
        if (i !== low) {
            items.unshift(ellipsis);
            items.unshift(addPage(low));
        }
    } else {
        items.push(addPage(low));
        items.push(ellipsis);

        for (i = page - 1; i <= page + 1; i++)
            items.push(addPage(i));

        items.push(ellipsis);
        items.push(addPage(high));
    }

    $("ul.pagination-list").html(items.join("\n"));

    $("a.pagination-link").click(function (a) {
        goToPage(a.target.text)
    })
}

function goToPage(page, reload = true, inplace = false, override = false) {
    page = parseInt(page)
    page = Math.min(numberOfPages, Math.max(page, 1))

    if (page !== currPage || override)
        if (inplace)
            window.history.replaceState({leak: page}, `Leak page ${page}`, `/leaks/list/${page}`);
        else
            window.history.pushState({leak: page}, `Leak page ${page}`, `/leaks/list/${page}`);
    else if (reload)
        location.reload();
    currPage = page
    load()
}


function updateNumberOfPages(nPerPage) {
    if (nPerPage !== undefined) {
        perPage = nPerPage
        localStorage.setItem("amountPerPage", nPerPage)
    }
    if (perPage === -1)
        numberOfPages = 1
    else
        numberOfPages = Math.ceil(articlesAmount / perPage)
}

function load() {
    prevButton.addClass("is-invisible")
    nextButton.addClass("is-invisible")
    pagination.removeClass("is-hidden")
    updateNumberOfPages()

    if (numberOfPages === 1) {
        pagination.addClass("is-hidden")
    }

    paginate(currPage, numberOfPages)
    if (currPage !== 1)
        prevButton.removeClass("is-invisible")
    if (currPage !== numberOfPages)
        nextButton.removeClass("is-invisible")

    loadFilter()
    let numArticles
    if (articles.length !== articlesRaw.length)
        numArticles = `${articles.length}(${articlesRaw.length}) leak${articles.length === 1 ? '' : 's'}`
    else
        numArticles = `${articles.length} leak${articles.length === 1 ? '' : 's'}`
    $("#leakNumber").text(numArticles)

    updateDropDownText()
    loadItems()
    sessionStorage.setItem("currentPage", window.location.pathname)
    imageErr()
}

// dropdown
function populateAmountList() {
    $("#leakAmounts").empty()
    choices.forEach(function (i) {
        $('<a/>', {
            "class": `dropdown-item ${i === perPage ? "is-active" : ""}`,
            text: i === -1 ? "all" : i,
            on: {
                click: function () {
                    amountDropdown.removeClass("is-active")
                    updateNumberOfPages(i)
                    goToPage(currPage, false, true, true)
                    populateAmountList()
                }
            }
        }).appendTo("#leakAmounts");
    })
}

function updateDropDownText() {
    if (perPage === -1)
        $("#dropDownText").text("All leaks")
    else
        $("#dropDownText").text(`${perPage} leak${perPage === 1 ? "" : "s"} per page`)
}

function loadItems() {
    itemList.empty()
    if (perPage === -1)
        $.each(articles, function (_, k) {
            itemList.append(k)
        })
    else {
        let start = (currPage - 1) * perPage
        for (let i = start; i < start + perPage; i++)
            itemList.append(articles[i])
    }
}


function loadFilter() {
    if (getStorageDefault("includeFilterRadio", "or", sessionStorage) === "and") {
        includeAnd.prop('checked', true);
    } else {
        includeOr.prop('checked', true);
    }
    if (getStorageDefault("excludeFilterRadio", "or", sessionStorage) === "and") {
        excludeAnd.prop('checked', true);
    } else {
        excludeOr.prop('checked', true);
    }
    filterLeaks()
}

function filterLeaks() {
    articles = []
    articlesRaw.each(function (_, k) {
            let tags = []
            let tagA = $(k).find(".tags")
            if (tagA.length !== 0)
                tagA.find(".tag").each(function (_, e) {
                    tags[tags.length] = ($(e).text())
                })


            let include = includeTags.length === 0

            // include logic
            if (includeOr.is(':checked')) {
                for (let i in tags) {
                    let tag = tags[i]
                    if (includeTags.includes(tag)) {
                        include = true
                        break
                    }
                }
            } else {
                // and
                if (includeTags.length !== 0) {
                    include = tags.length !== 0
                    for (let i in tags) {
                        let tag = tags[i]
                        if (!includeTags.includes(tag)) {
                            include = false
                            break
                        }
                    }
                }
            }

            // exclude
            if (excludeOr.is(':checked')) {
                for (let i in tags) {
                    let tag = tags[i]
                    if (excludeTags.includes(tag)) {
                        include = false
                        break
                    }
                }
            } else {
                // and
                if (excludeTags.length !== 0) {
                    let a = tags.length !== 0
                    for (let i in tags) {
                        let tag = tags[i]
                        if (!excludeTags.includes(tag)) {
                            a = false
                        }
                    }
                    if (include && a) {
                        include = false
                    }
                }
            }

            // console.log(tags)
            if (include)
                articles[articles.length] = k
        }
    )
}

function splitList(string) {
    let items = string.split(",")
    if (items.length === 1 && items[0] === "") {
        items = []
    }
    return items
}

function clearFilters() {
    sessionStorage.setItem("excludeList", "")
    sessionStorage.setItem("includeList", "")
    includeTagSelector.removeAll()
    excludeTagSelector.removeAll()
    includeTags = []
    excludeTags = []
    sessionStorage.removeItem("includeFilterRadio")
    sessionStorage.removeItem("excludeFilterRadio")

    load()
}
