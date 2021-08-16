function loadIframe(placeholder, data) {
    data.load = function () {
        let $this = $(this)
        $this.removeClass("is-hidden")
        placeholder.remove()
        // this.style.height = this.contentWindow.document.body.offsetHeight + 'px';
    }
    data.class += " is-hidden"
    let iframe = $('<iframe>', data);
    let frameLocation = placeholder.parent()
    frameLocation.append(iframe)
}

loadIframe($("a.deca-news"),{
    src: "https://www.deca.net/news/",
    class: "deca-news",
    style: "position: static; visibility: visible; display: inline-block; width: 100%; padding: 0px; " +
        "border: none; max-width: 100%; min-width: 180px; margin-top: 0px; margin-bottom: 0px; min-height: 200px;",
    height: "4000px",
    scrolling: "no",
    frameborder: "0"
})

//maybe have this change reactively
let discord = "https://discord.com/widget?id=765917490734694412&theme=light"
if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    discord = "https://discord.com/widget?id=765917490734694412&theme=dark"
}
loadIframe($("a.deca-discord"),{
    src: discord,
    class: "deca-discord",
    sandbox:"allow-popups allow-popups-to-escape-sandbox allow-same-origin allow-scripts",
    allowtransparency: "true",
    scrolling: "no",
    width:"100%",
    frameborder: "0",
    height:"98%"
})

// loadIframe($("a.deca-linkedin"),{
//     src: "https://www.linkedin.com/company/megadodo-simulation-games/posts/",
//     class: "deca-linkedin",
//     style: "position: static; visibility: visible; display: inline-block; width: 100%; padding: 0px; " +
//         "border: none; max-width: 100%; min-width: 180px; margin-top: 0px; margin-bottom: 0px; min-height: 200px;",
//     height: "4000px",
//     scrolling: "no",
//     frameborder: "0"
// })