package editor

const globalCss = `
body {
    background-color: #ddd;
}

.inner {
    padding: 10px;
    font-size: 28px;
}

button {
    margin: 5px 0 5px 0;
    padding: 5px;
    background-color: #22d;
    font-color: #bbb;
}

button:hover {
    background-color: #55e;
}

button:active {
    background-color: #005;
}
`

const mainMenuHtml = `
<html>
    <body>
        <div class="inner">
            <h1>Main Menu</h1>
            <button>Open</button>
            <button>Save</button>
            <button>Exit</button>
        </div>
    </body>
</html>
`
