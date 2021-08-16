window.addEventListener('load', function () {
    prevButton.click(function(){
        goToPage(currPage-1)
        load()
    })
    nextButton.click(function(){
        goToPage(currPage+1)
        load()
    })
    load()
});

const prevButton = $("a.pagination-previous")
const nextButton = $("a.pagination-next")

let perPage = 1;

const ellipsis = `<li><span class="pagination-ellipsis">&hellip;</span></li>`;
const currentPage = `<li><a class="pagination-link is-current" aria-label="Page {p}" aria-current="page">{p}</a></li>`;
const gotoPage = `<li><a class="pagination-link" aria-label="Goto page {p}">{p}</a></li>`;

function paginate(page, high) {
    let low = 1;
    let items = [];

    function addPage(number){
        if (number===page){
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

    $("a.pagination-link").click(function(a){
        goToPage(a.target.text)
        load()
    })
}

function goToPage(page){
    page = parseInt(page)
    if (page !== currPage)
        window.history.pushState({leak:page}, `Leak page ${page}`, `/leaks/list/${page}`);
    else
        location.reload();
    currPage = page
}

function load() {
    prevButton.addClass("is-invisible")
    nextButton.addClass("is-invisible")
    let max = Math.ceil(articlesAmount / perPage)
    page = Math.min(max,Math.max(currPage,1))
    if (page !== currPage)
        goToPage(page)

    paginate(page, max)
    if (page !== 1)
        prevButton.removeClass("is-invisible")
    if (page !== max)
        nextButton.removeClass("is-invisible")
}