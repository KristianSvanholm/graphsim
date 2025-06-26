var PIXI = require('pixi.js')
var PIXVI = require('pixi-viewport')

let itterations = 125
let start = performance.now()

var itt = document.getElementById("itt");
var search = document.getElementById("search");
var duration = document.getElementById("duration")
itt.oninput = function(){
    itterations = this.value
}
search.onclick = function(){
    const canvases = document.body.querySelectorAll('canvas');
    canvases.forEach(canvas => canvas.remove());

    start = performance.now();
    run(itterations)
    
}

const run = (async () => {

    const url = "http://localhost:8080/graph?itt="+itterations
    const response = await fetch(url)
    if (!response.ok) {
        console.log("Ooops", response.status)
    }

    const json = await response.json();

    // Create the application helper and add its render target to the page
    const app = new PIXI.Application();
    await app.init({ antialias: true, resizeTo: window })
    document.body.appendChild(app.canvas);

    const viewport = new PIXVI.Viewport({
        screenWidth: window.innerWidth,
        screenHeight: window.innerHeight,
        worldwidth: 1920*10,
        worldheight: 1080*10,
        events: app.renderer.events,
    })

    app.stage.addChild(viewport);

    viewport.drag().pinch().wheel().decelerate();
    const graphics = new PIXI.Graphics();

    let links = json.links
    let nodes = json.nodes
    let quads = json.quads

    for( let i = 0; i< quads.length; i++ ){
        let q = quads[i]
        // Left
        graphics.moveTo(q.Pos.X, q.Pos.Y)
        graphics.lineTo(q.Pos.X, q.Pos.Y + q.Size)
       
        // Right
        graphics.moveTo(q.Pos.X + q.Size, q.Pos.Y)
        graphics.lineTo(q.Pos.X + q.Size, q.Pos.Y + q.Size)

        // Top
        graphics.moveTo(q.Pos.X, q.Pos.Y)
        graphics.lineTo(q.Pos.X + q.Size, q.Pos.Y)

        // Bottom
        graphics.moveTo(q.Pos.X, q.Pos.Y + q.Size)
        graphics.lineTo(q.Pos.X + q.Size, q.Pos.Y + q.Size)

    }

    for (let i= 0; i< links.length; i++) {
        let src = nodes[links[i].Src]
        let dst = nodes[links[i].Dst]
        graphics.moveTo(src.Pos.X, src.Pos.Y)
        graphics.lineTo(dst.Pos.X, dst.Pos.Y)
    }

    graphics.stroke({color: 0xffffff, pixelLine:true, width: 1})

    for( let i = 0; i< nodes.length; i++ ){
        let n = nodes[i]
        graphics.circle(n.Pos.X,n.Pos.Y,5)
        graphics.fill(0xde3249)
    }

    viewport.addChild(graphics)

    const end = performance.now();
    const dur = end - start;
    console.log(duration)
    duration.textContent = dur+"ms"
})

run();
