'use strict';

import { Bracket } from './bracket.js';
import * as constants from './constants.js';

var bracketEdit = null

window.onload = function () {
    bracketEdit = new BracketEdit()
    bracketEdit.bracketBase.loadData(this.canvas)
};

class BracketEdit {

    constructor() {
        this.canvas = document.getElementById("bracket_canvas")
        this.bracketBase = new Bracket()

        this.hasInput = false
        this.canvas.removeEventListener("mousedown",
            this.bracket.mousedownHandler)
        this.canvas.addEventListener("mousedown",
            this.mousedownHandler.bind(this), false)
        document.addEventListener("keydown",
            this.keypressHandler.bind(this), false)
        $("#bracket_canvas").keydown(this.keypressHandler.bind(this))
    }

    keypressHandler(evt) {

    }

    getMousePos(canvas, evt) {
        var rect = canvas.getBoundingClientRect();
        return {
            x: evt.clientX - rect.left,
            y: evt.clientY - rect.top
        };
    };

    mousedownHandler(evt) {

        var that = this
        var canvas = document.getElementById("bracket_canvas")
        var mousePos = this.getMousePos(canvas, evt);
        var message = 'Mouse position(editor): ' + mousePos.x + ',' + mousePos.y;
        console.log(message)
        this.bracketBase.selection = this.bracketBase.rootNode.IntersectGame(mousePos.x, mousePos.y)

        if (this.selection != null) {
            console.log("intersects with node " + this.bracketBase.selection.node.Id)
            var x = this.bracketBase.selection.node.x
            var y = this.bracketBase.selection.node.y
            if (this.bracketBase.selection.slot.slotNum == 1) {
                x = x + 10
                y = y + 5
                this.addInput(x, y)
            } else {
                x = x + 10
                y = y + 30
                this.addInput(x, y)
            }

        } else {
            console.log("No node hit")
        }

        //this.render(this.data, canvas)
    }

    addInput(x, y) {

        var input = document.createElement('input');

        input.type = 'text';
        input.style.position = 'fixed';
        input.style.left = (x - 4) + 'px';
        input.style.top = (y - 4) + 'px';

        input.onkeydown = this.handleEnter;

        document.body.appendChild(input);

        input.focus();

        this.hasInput = true;
    }

    handleEnter(e) {
        var keyCode = e.keyCode;
        if (keyCode === 13) {
            this.drawText(this.value, parseInt(this.style.left, 10), parseInt(this.style.top, 10));
            document.body.removeChild(this);
            this.hasInput = false;
        }
    }

    drawText(txt, x, y) {
        this.bracketBase.ctx.textBaseline = 'top';
        this.bracketBase.ctx.textAlign = 'left';
        this.bracketBase.ctx.font = font;
        this.bracketBase.ctx.fillText(txt, x - 4, y - 4);
    }
}




