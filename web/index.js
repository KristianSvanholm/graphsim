var PIXI = require('pixi.js')
var PIXVI = require('pixi-viewport')

const run = (async () => {

    const url = "http://localhost:8080/graph"
    const response = await fetch(url)
    if (!response.ok) {
        console.log("FUCK", response.status)
    }

    const json = await response.json();
    console.log(json);

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

    console.log(window.innerWidth, window.innerHeight)
    app.stage.addChild(viewport);

    viewport.drag().pinch().wheel().decelerate();
    const graphics = new PIXI.Graphics();

    let links = json.links
    let nodes = json.nodes

    for (let i= 0; i< links.length; i++) {
        let src = nodes[links[i].Src]
        let dst = nodes[links[i].Dst]
        console.log(src, dst)
        graphics.moveTo(src.Pos.X, src.Pos.Y)
        graphics.lineTo(dst.Pos.X, dst.Pos.Y)
    }

    graphics.stroke({color: 0xffffff, pixelLine:true, width: 1})

    for( let i = 0; i< nodes.length; i++ ){
        let n = nodes[i]
        console.log(n.Pos.X,n.Pos.Y)
        graphics.circle(n.Pos.X,n.Pos.Y,5)
        graphics.fill(0xde3249)
    }

    viewport.addChild(graphics)

})

run();
