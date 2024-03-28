document.addEventListener("DOMContentLoaded", (event) => {
    // init tinymce
    /* tinymce.init({
        selector: '#tinymce',
        plugins: 'autolink lists link image charmap preview anchor pagebreak',
        toolbar_mode: 'floating',
    }) */
    // init tippy
    tippy(".avatar")
})

/* Progress bar */
//Source: https://alligator.io/js/progress-bar-javascript-css-variables/
var h = document.documentElement,
    b = document.body,
    st = "scrollTop",
    sh = "scrollHeight",
    progress = document.querySelector("#progress"),
    scroll;
var scrollpos = window.scrollY;
var header = document.getElementById("header");

document.addEventListener("scroll", function () {
    /*Refresh scroll % width*/
    scroll = ((h[st] || b[st]) / ((h[sh] || b[sh]) - h.clientHeight)) * 100;
    progress?.style.setProperty("--scroll", scroll + "%");

    /*Apply classes for slide in bar*/
    scrollpos = window.scrollY;

    if (scrollpos > 100) {
        header?.classList.remove("hidden");
        header?.classList.remove("fadeOutUp");
        header?.classList.add("slideInDown");
    } else {
        header?.classList.remove("slideInDown");
        header?.classList.add("fadeOutUp");
        header?.classList.add("hidden");
    }
})

// scroll to top
const t = document.querySelector(".js-scroll-top");
if (t) {
    t.onclick = () => {
        window.scrollTo({ top: 0, behavior: "smooth" });
    };
    const e = document.querySelector(".scroll-top path"),
        o = e.getTotalLength();
    (e.style.transition = e.style.WebkitTransition = "none"),
        (e.style.strokeDasharray = `${o} ${o}`),
        (e.style.strokeDashoffset = o),
        e.getBoundingClientRect(),
        (e.style.transition = e.style.WebkitTransition =
            "stroke-dashoffset 10ms linear");
    const n = function () {
        const t =
            window.scrollY ||
            window.scrollTopBtn ||
            document.documentElement.scrollTopBtn,
            n = Math.max(
                document.body.scrollHeight,
                document.documentElement.scrollHeight,
                document.body.offsetHeight,
                document.documentElement.offsetHeight,
                document.body.clientHeight,
                document.documentElement.clientHeight
            ),
            s = Math.max(
                document.documentElement.clientHeight,
                window.innerHeight || 0
            );
        var l = o - (t * o) / (n - s);
        e.style.strokeDashoffset = l;
    };
    n();
    const s = 100;
    window.addEventListener(
        "scroll",
        function (e) {
            n();
            (window.scrollY ||
                window.scrollTopBtn ||
                document.getElementsByTagName("html")[0].scrollTopBtn) > s
                ? t.classList.add("is-active")
                : t.classList.remove("is-active");
        },
        !1
    );
}
