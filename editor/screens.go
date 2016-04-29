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

.overview button {
    width: 220px;
    margin: 10px 10px 0 0;
}

.tree {
    margin: 10px 0 0 0;
	padding: 10px 10px 0 10px;
    background-color: #bbc;
}

.tree .treeItem {
	margin: 0 0 10px 0;
}

.tree .lvl2 { margin: 0 0 10px 10px; }
.tree .lvl3 { margin: 0 0 10px 20px; }
.tree .lvl4 { margin: 0 0 10px 30px; }
.tree .lvl5 { margin: 0 0 10px 40px; }
.tree .lvl6 { margin: 0 0 10px 50px; }
.tree .lvl7 { margin: 0 0 10px 60px; }
.tree .lvl8 { margin: 0 0 10px 70px; }
.tree .lvl9 { margin: 0 0 10px 80px; }
.tree .lvl10 { margin: 0 0 10px 90px; }

.tree .treeItem div {
    height: 20px;
    margin: 0 5px 0 5px;
}

.tree .icon {
    background-color: #333;
    width: 10px;
}

.tree .label {
	width: 150px;
    font-size: 16px;
    padding: 1px;
}

.tree .delete {
	width: 25px;
    padding: 0 0 2px 2px;
}

.tree .delete:hover {
    color: #f00;
}

.tree .delete:active {
        color: #fbb;
}

.tree .closed .icon {
    background-color: #777;
}

.tree .closed .label {
    color: #999;
}


.progressBar {
	padding: 10px;
}

.progress {
	width: 310px;
    height: 40px;
	padding: 0 0 0 3px;
    background-color: #666;
}

.progress div {
	width: 10px;
	height: 20px;
	margin: 10px 0 10px 5px;
    float: left;
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
            <div class="tree"></div>
        </div>
    </body>
</html>
`

const progressBar = `
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
        </div>
    </body>
</html>
`
