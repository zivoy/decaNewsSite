window.addEventListener('load', function () {
    prevButton.click(function () {
        goToPage(currPage - 1)
    });
    nextButton.click(function () {
        goToPage(currPage + 1)
    });
    $("#dropButton").click(function () {
        dropDown.toggleClass("is-active")
    });

    populateAmountList()
    itemList.removeClass("is-hidden")
    articles = itemList.children()

    updateNumberOfPages()
    goToPage(currPage, false, true, true)
});

window.onpopstate = function (e) {
    if (e.state) {
        console.log(e)
        goToPage(e.state.leak, false, true)
    }
};


const choices = [-1, 5, 10, 20, 50]


const prevButton = $("a.pagination-previous")
const nextButton = $("a.pagination-next")
const dropDown = $("#amountDropdown")
const itemList = $("#leakItems")
const pagination = $("nav.pagination")

let perPage = parseInt(getStorageDefault("amountPerPage", choices[0]));
let numberOfPages = 0
let articles = {}

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
                    dropDown.removeClass("is-active")
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
        articles.each(function (_, k) {
            itemList.append(k)
        })
    else {
        let start = (currPage - 1) * perPage
        for (let i = start; i < start + perPage; i++)
            itemList.append(articles[i])
    }
}
