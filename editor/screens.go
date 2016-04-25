package editor

const globalCss = `
body {
    background-color: #ddd;
}

h1 {
    font-size: 48px
}

button {
    margin: 10px 0 0 0;
    padding: 10px;
    background-color: #22d;
    font-size: 24px;
    color: #bbb;
    text-align: center;
}

button:hover {
    background-color: #55e;
}

button:active {
    background-color: #005;
}

.main {
    padding: 20px;
}

.overview {
    padding: 20px;
}

.overview button {
    width: 220px;
    margin: 10px 10px 0 0;
}

`

const mainMenuHtml = `
<html>
    <body>
        <div class="main">
            <h1>Main Menu</h1>
            <button onclick=open>Open</button>
            <button onclick=save>Save</button>
            <button onclick=exit>Exit</button>
        </div>
    </body>
</html>
`

const overviewMenuHtml = `
<html>
    <body>
        <div class="overview">
            <button onclick=import>Import</button>
            <button onclick=newGroup>New group</button>
        </div>
    </body>
</html>
`
