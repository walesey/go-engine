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
    padding: 10px 10px 8px 20px;
    background-color: #44d;
    font-size: 24px;
    color: #bbb;
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

.overview .panel {
    margin: 10px 0 0 0;
    padding: 10px;
    background-color: #bbc;
}

.overview button {
    width: 22px;
    padding: 3px;
    background-color: #fff;
    margin: 0 5px 0 0;
}

.overview button:hover {
    background-color: #66f;
}

.overview button:active {
    background-color: #33f;
}

.fileBrowser {
    padding: 20px;
}

.fileBrowser button {
    margin: 10px 10px 0 0;
    width: 200px;
}

.fileBrowser .fileView {
    height: 600px;
    background-color: #fff;
    padding: 5px;
}

.fileBrowser .content button {
    width: 100%;
    height: 20px;
    padding: 5px 0 5px 45%;
    margin: 0;
    background-color: #0000;
    font-size: 8px;
    color: #000;
}

.fileBrowser .content button:hover {
    background-color: #00f5;
    color: #999;
}

.fileBrowser .content button:active {
    background-color: #00f9;
    color: #999;
}

.fileBrowser input {
    width: 100%;
    font-size: 16px;
    background-color: #fff;
    padding: 2px;
}

.progressBar {
    padding: 10px;
}

.progressBar .progress {
    width: 310px;
    height: 40px;
    padding: 0 0 0 3px;
    background-color: #666;
}

.progressBar .progress div {
    width: 10px;
    height: 20px;
    margin: 10px 0 10px 5px;
    float: left;
}

.nodeEditor {
    padding: 20px;
}

.nodeEditor input {
    width: 100%;
    font-size: 16px;
    background-color: #fff;
    padding: 2px;
    margin: 5px 0 0 0;

}

.nodeEditor button {
    margin: 10px 10px 0 0;
    width: 200px;
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
            <div class="panel">
                <button onclick=newGroup><img src=planetOpen></img></button>
                <button onclick=copyGroup><img src=copy></img></button>
                <button onclick=referenceGroup><img src=reference></img></button>
                <button onclick=deleteGroup><img src=trash></img></button>
                <button onclick=editGroup><img src=cog></img></button>
                <button onclick=import><img src=geometry></img></button>
                <button onclick=scale><img src=scale></img></button>
                <button onclick=translate><img src=translate></img></button>
                <button onclick=rotate><img src=rotate></img></button>
                <button onclick=reset><img src=reset></img></button>
            </div>
            <div id="overviewTree" class="panel"></div>
        </div>
    </body>
</html>
`

const nodeEditMenuHtml = `
<html>
    <body>
        <div class="nodeEditor">
            <h3>Name</h3>
            <input type="text" id=name></input>
            <h3>Classes</h3>
            <input type="text" id=class1></input>
            <input type="text" id=class2></input>
            <input type="text" id=class3></input>
            <button onclick=nodeEditorOk>Ok</button>
            <button onclick=nodeEditorCancel>Cancel</button>
        </div>
    </body>
</html>
`

const fileBrowserHtml = `
<html>
    <body>
        <div class="fileBrowser">
            <h1 id="heading">File Browser</1>
            <div class="content">
                <button onclick=fileBrowserScrollup>^</button>
                <div id="fileView" class="fileView"></div>
                <button onclick=fileBrowserScrollDown>v</button>
            </div>
            <input id="filePathInput" type="text"></input>
            <button onclick=fileBrowserOpen>Open</button>
            <button onclick=fileBrowserCancel>Cancel</button>
        </div>
    </body>
</html>
`

const progressBarHtml = `
<html>
    <body>
        <div class="progressBar">
            <div id="progress" class="progress">
                <div id="progress1"></div>
                <div id="progress2"></div>
                <div id="progress3"></div>
                <div id="progress4"></div>
                <div id="progress5"></div>
                <div id="progress6"></div>
                <div id="progress7"></div>
                <div id="progress8"></div>
                <div id="progress9"></div>
                <div id="progress10"></div>
                <div id="progress11"></div>
                <div id="progress12"></div>
                <div id="progress13"></div>
                <div id="progress14"></div>
                <div id="progress15"></div>
                <div id="progress16"></div>
                <div id="progress17"></div>
                <div id="progress18"></div>
                <div id="progress19"></div>
                <div id="progress20"></div>
            </div>
            <div id="progressBarMessage"></div>
        </div>
    </body>
</html>
`
