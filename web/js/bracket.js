'use strict';


/*
const constants.NODE_WIDTH=120;
const constants.NODE_HEIGHT=40;
const constants.HEIGHT = 1400;
const constants.WIDTH = 1200;

const constants.SIZE_BASE = 6;

 */

import * as constants from './constants.js';

window.getTournamentID = function() {

    var idElem = document.getElementById('TournamentID');
    var id = idElem.getAttribute('value');
    return id
}
/*
window.getTournamentID = function() {
    var pathArray = window.location.pathname.split('/');
    let id = pathArray[pathArray.length - 1]
    console.log("id from url = " + id)
    return id
};
 */


//var number = getUrlVars()["x"];

function getUrlVars() {
    var vars = {};
    var parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
        vars[key] = value;
    });
    return vars;
}

window.onload = function () {

    var tId = getUrlVars()["id"];

    /*
    if (tId == null || (tId != "")) {
        var url = "/api/tournaments";
        $.post(url,
            function (data) {
                tId = data.id;
            });
    }
     */

  //  var myScroll = new IScroll('#container-canvas',{
   //         mouseWheel: true,
   //         scrollbars: true});


};


export class Bracket {

      constructor () {
          this.x = 0;
          this.y = 0;
          this.selection = null;
          this.orientation = constants.RIGHT_TO_LEFT;
          this.data = null;
          this.rootNode = null;
          this.canvas = document.getElementById("bracket_canvas");
          this.drawcenter = 0;

          // Adjust resolution
          this.ratio = calcRatio();
          //let w = this.canvas.width;
          //let h = this.canvas.height;
          //this.canvas.width = this.canvas.width * ratio;
          //this.canvas.height = this.canvas.height * ratio;
          //this.canvas.style.width = w + "px";
          //this.canvas.style.height = h + "px";
          //this.canvas.getContext("2d").setTransform(ratio, 0, 0, ratio, 0, 0);
          this.ctx = this.canvas.getContext('2d');

          //this.ctx.setTransform(ratio, 0, 0, ratio, 0, 0);

          this.resize = function(data) {
              let size = calcSize(data);
              this.width = size.width;
              this.height = size.height;

              console.log("calculated size: " + size.height + " x " + size.width);
              this.canvas.width = this.width;
              this.canvas.height = this.height;
              this.canvas.style.width = this.width + "px";
              this.canvas.style.height = this.height + "px";
          };

          this.render = function (data, canvas) {

              //this.resize(data);

              if (this.rootNode) {
                  this.drawcenter = this.rootNode.span.upper + 22;
              }
              this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

              if (this.orientation == constants.RIGHT_TO_LEFT) {
                  this.renderRightToLeft(data, canvas)
              } else if (this.orientation == constants.CENTERED) {
                  this.renderCentered(data, canvas)
              }
          };

          this.renderRightToLeft = function (data, canvas) {

              var width = this.width;

              // Render winners line, and winner name (if the tourney is complete)
              this.ctx.strokeStyle = chroma("black")
              this.ctx.lineWidth = 1
              this.ctx.beginPath();
              this.ctx.moveTo(width - (constants.NODE_WIDTH + 10), this.drawcenter);
              this.ctx.lineTo(width, this.drawcenter);
              this.ctx.closePath();
              this.ctx.stroke();

              this.ctx.strokeStyle = chroma(constants.COLOR_3);
              this.ctx.fillStyle = chroma(constants.COLOR_3);
              this.ctx.textAlign = "center"
              this.ctx.font = "18px Comic Sans MS";
              var winner = this.getWinner();
              this.ctx.strokeText(winner, width - ((constants.NODE_WIDTH + 10) / 2),
                  (this.drawcenter) - 8);

              this.rootNode.renderRightToLeft(this.ctx,
                  (width) - (constants.NODE_WIDTH + (constants.NODE_WIDTH - 10)),
                  (this.drawcenter) - (constants.NODE_HEIGHT / 2),
                  0, this.data.degree, this.selection)
          };

          this.getWinner = function () {

              console.log("this.rootNode:" + this.data.root)
              console.log("this.rootNode.state.result: " + this.data.root.state.result)
              if ((this.data.root) && (this.data.root.state.result != null)) {
                  var winnerId = this.data.root.state.result.winningPlayer
                  if (winnerId > 0) {
                      var p = this.data.players[winnerId];
                      if (p != null) {
                          return p.name;
                      }
                  }

              }
              return "<TBD>"
          };

          this.renderCentered = function (data, canvas) {

              // Render winners line, and winner name (if the tourney is complete)

              this.rootNode.left.renderCentered(this.ctx, (constants.WIDTH)  - 10,
                  (constants.HEIGHT / 2) - (constants.NODE_HEIGHT / 2),
                  1, this.data.degree, this.selection)

              this.rootNode.right.renderCentered(this.ctx,  10,
                  (constants.HEIGHT / 2) - (constants.NODE_HEIGHT / 2),
                  1, this.data.degree, this.selection)
          };

          this.loadData = function (canvas) {
              var self = this
//              var pathArray = window.location.pathname.split('/');
//              let id = pathArray[pathArray.length - 1]
//              console.log("id from url = " + id)

              let id = window.getTournamentID();

              //self.tournamentId = this.getTournamentId()
              $.getJSON("/api/tournaments/" + id, function (data) { // parenthesis are not necessary
                  console.log(data);
                  self.data = data;

                  self.resize(data);
                  /*
                  var size;
                  if (data) {
                      size = calcSize(data);
                  }
                  console.log("width: " + size.w);
                  console.log("height: " + size.h);
                  self.canvas.width = size.w;
                  self.canvas.height = size.h;
                   */

                  //self.winner = data.winner;
                  if (self.rootNode == null) {
                      var finalGame = new Game(data, data.nodes[data.rootNodeId]);
                      self.rootNode = finalGame
                  }
                  //self.resize()
                  self.render(data, canvas)
              })
          };

          this.computeWidth = function () {

          };

          this.computeHeight = function () {

          }

      }
}

  export class Slot {
      constructor(SlotData, slotNum) {
          this.playerNumber = SlotData.playerNumber
          this.winnerFromGame = SlotData.fromGameId
          this.dropdown = SlotData.dropdown
          this.playerName = SlotData.playerName
          this.slotNum = slotNum
      }


  }

  /*
       Leaf level games are 2nd from bottom level nodes.


       A node that has a null left and right child are leaf nodes. A Leaf node either starting position
        where is has a playerId or it represents a drop position in which case it has a gameId
       A Game object is not created for leaf nodes, instead the parent node will have either the player
       or game field set to the player or game the leaf node represents.
   */
  export class Game {

      constructor(Tournament, TheGame) {
          this.tournament = Tournament
          this.game = TheGame;
          this.type = TheGame.Type
          this.isLeaf = false;

          this.orientation = constants.RIGHT_TO_LEFT

          this.left = null;
          this.right = null;

          this.player1 = 0;
          this.player2 = 0;
          this.dropGame1 = 0;
          this.dropGame2 = 0;

          this.x = null;
          this.y = null;
          this.span = TheGame.span;
          this.width = constants.NODE_WIDTH;
///          this.space = 22

          this.upperBGColor = chroma(constants.COLOR_15);
          this.lowerBGColor = chroma(constants.COLOR_16);
          this.gameIdBGColor = chroma(constants.COLOR_13);
          this.losersGameIdBGColor = chroma(constants.COLOR_3);
//          this.losersGameIdBGColor = chroma("darkOrange")
          this.pendingGameIdBGColor = chroma(constants.COLOR_1);
          this.Id = TheGame.id;
///          this.player1 = TheGame.player1;
          //         this.player2 = TheGame.player2;
          this.player = TheGame.player;
          this.dropGame = TheGame.drop;

          if (TheGame.tier == 2) {
              this.isLosersSide = true
          } else {
              this.isLosersSide = false
          }

          /*
          if (BracketNode.slot1 != null) {
              this.slot1 = new Slot(BracketNode.slot1, 1)
          }

          if (BracketNode.slot2 != null) {
              this.slot2 = new Slot(BracketNode.slot2, 2)
          }
           */

          if ((TheGame.left.id == 0) && (TheGame.right.id == 0)) {
              this.isLeaf = true;
              let Null = { valueOf: ()=>null }
              return Null;
          }

          if ((TheGame.left.id != 0) && (TheGame.left.kind == 1)) {
              var node = this.tournament.nodes[TheGame.left.id];
              this.left = new Game(this.tournament, node)
              if (node != null) {
                  this.player1 = this.left.player;
              }

          } else if ((TheGame.left.id != 0) && (TheGame.left.kind == 2 || TheGame.left.kind == 3)) {
              var g = this.tournament.nodes[TheGame.left.id]
              this.player1 = g.player;
              this.dropGame1 = g.drop;
              this.left = null;
          }

          if ((TheGame.right.id != 0) && (TheGame.right.kind == 1)) {
              var node = this.tournament.nodes[TheGame.right.id];
              this.right = new Game(this.tournament, node)
              if (node != null) {
                  this.player2 = this.right.player;
              }
          } else if ((TheGame.right.id != 0) && (TheGame.right.kind == 2 || TheGame.right.kind == 3)) {
              var g = this.tournament.nodes[TheGame.right.id]
              this.player2 = g.player;
              this.dropGame2= g.drop;
              this.right = null;
          }

          /*
          if (TheGame.left.id != 0) {

              var node = this.tournament.nodes[TheGame.left.id];
              if (node) {
                  this.left = new Game(this.tournament, node)
                  if (this.left.valueOf() == null) {
                      var g = this.tournament.nodes[TheGame.left.id]
                      this.player1 = g.player;
                      this.dropGame1 = g.drop;
                      this.left = null;
                  } else if (this.left.player != 0) {
                      this.player1 = this.left.player;
                  }
              } else {
                  console.log("Error, bad node id")
              }
          }

          if (TheGame.right.id != 0) {
              var node = this.tournament.nodes[TheGame.right.id];
              if (node) {
                 this.right = new Game(this.tournament, node)
                  if (this.right.valueOf() == null) {
                      var g = this.tournament.nodes[TheGame.right.id]
                      this.player2 = g.player;
                      this.dropGame2 = g.drop;
                      this.right = null;
                  } else if (this.right.player != 0) {
                      this.player2 = this.right.player;
                  }
              } else {
                  console.log("Error, bad node id")
              }
          }
           */

          this.log = function () {
              console.log("node id: " + this.Id + ", type: " + this.type)
              console.log("   left: " + this.left);
              console.log("   right: " + this.right);
              console.log("   dropGame1: " + this.dropGame1);
              console.log("   dropGame2: " + this.dropGame2);
          }

          this.playerName = function (num) {
              var p = this.tournament.players[num];
              if (p != null) {
                  return p.name;
              } else {
                  return "BUY"
              }
          };

          this.frameBox = function (ctx, x, y, width, alpha, color1, color2, reversed) {

              // draw the box
              ctx.textAlign = "Left";
              ctx.lineWidth = 1;

              ctx.clearRect(x, y, width, constants.NODE_HEIGHT);
              ctx.strokeRect(x, y, width, constants.NODE_HEIGHT);

              ctx.fillStyle = color1.alpha(alpha);
              if (reversed) {
                  ctx.strokeRect(x + 24, y, width - 24, constants.NODE_HEIGHT);
                  ctx.fillRect(x + 25, y + 1, width - 26, constants.NODE_HEIGHT / 2 - 2);
                  ctx.fillStyle = color2.alpha(alpha);
                  ctx.fillRect(x + 25, y + 1 + (constants.NODE_HEIGHT / 2), width - 26, constants.NODE_HEIGHT / 2 - 2);
              } else {
                  ctx.strokeRect(x, y, width - 24, constants.NODE_HEIGHT);
                  ctx.fillRect(x + 1, y + 1, width - 26, constants.NODE_HEIGHT / 2 - 2);
                  ctx.fillStyle = color2.alpha(alpha);
                  ctx.fillRect(x + 1, y + 1 + (constants.NODE_HEIGHT / 2), width - 26, constants.NODE_HEIGHT / 2 - 2);
              }


          };

          this.fillBox = function (ctx, x, y, width, alpha, color, reversed) {

              ctx.fillStyle = color.alpha(alpha);
              if (reversed) {
                  ctx.fillRect(x + 1, y + 1, 22, constants.NODE_HEIGHT - 2);
                  ctx.beginPath();
                  ctx.moveTo(x + 24, y + (constants.NODE_HEIGHT / 2));
                  ctx.lineTo(x + width, y + (constants.NODE_HEIGHT / 2));
              } else {
                  ctx.fillRect(x + width - 23, y + 1, 22, constants.NODE_HEIGHT - 2);
                  ctx.beginPath();
                  ctx.moveTo(x, y + (constants.NODE_HEIGHT / 2));
                  ctx.lineTo(x + (width) - 24, y + (constants.NODE_HEIGHT / 2));
              }

              ctx.closePath();
              ctx.stroke();
              ctx.lineWidth = 1;

          };

          this.highlight = function (ctx, x, y, width, reversed) {

              ctx.strokeStyle = chroma("yellow");
              ctx.lineWidth = 1;
              if (reversed) {
                  ctx.strokeRect(x + 24, y, width - 24, constants.NODE_HEIGHT);
                  ctx.strokeRect(x, y, width, constants.NODE_HEIGHT);
              } else {
                  ctx.strokeRect(x, y, width - 24, constants.NODE_HEIGHT);
                  ctx.strokeRect(x, y, width, constants.NODE_HEIGHT);
              }


          };


          this.calcTextColor = function (slot) {
              return chroma(constants.COLOR_1);
          };

          this.drawPlayerIndicator = function (ctx, x, y, width, slot, reversed, winningSlot) {

              ctx.font = "14px Comic Sans MS";
              ctx.textAlign = "center";

              let display = "TBD";
              display = this.getDisplay(slot);

              var clr;

              if (display.startsWith(">>") === true) {
                  clr = chroma(constants.COLOR_1);
              } else {
                  clr = chroma(constants.COLOR_28);
              }

              ctx.strokeStyle = clr;
              if (clr != null) {

                  if (winningSlot != 0) {
                      if (winningSlot == slot) {
                          ctx.font = "16px Comic Sans MS";
                          clr = chroma(constants.COLOR_9)
                      } else {
                          ctx.font = "13px Comic Sans MS";
                          clr = clr.alpha(0.3)
                      }
                      ctx.strokeStyle = clr;
                  }

//                  if (this.left != null && this.left.ResolvedBuy()) {
//                      display = display + " *";
//                  }

                  var yoffset = 0;
                  if (slot == 1) {
                      yoffset = (constants.NODE_HEIGHT / 2) - 5;
                  } else {
                      yoffset = (constants.NODE_HEIGHT) - 5;
                  }

                  ctx.strokeText(display, x + ((width - 24) / 2), y + yoffset);
              }
          };

          this.drawID = function (ctx, x, y, width, reversed) {

              ctx.textAlign = "center"
              ctx.font = "16px Comic Sans MS";
              if (reversed) {
                  ctx.strokeText(this.Id, x + 12, y + (constants.NODE_HEIGHT / 2 + 6));
              } else {
                  ctx.strokeText(this.Id, x + width - 12, y + (constants.NODE_HEIGHT / 2 + 6))
              }
          };

          this.renderCentered = function (ctx, x, y, level, degree, selection) {

              var width = this.width;
              if (this.isLosersSide) {
                  width = width * 0.7
              }

              var space = width * .2

              ctx.strokeStyle = chroma("black");
              var alpha = 1.0;
              var brighten = 1;
              if (this.game.result != null) {
                  ctx.strokeStyle = chroma("black").brighten(3);
                  alpha = 0.3;
                  brighten = 3;
              }

              this.x = x;
              this.y = y;
              console.log("node: x: " + this.x + ", y: " + this.y);

              this.frameBox(ctx, x, y, width, alpha, this.upperBGColor, this.lowerBGColor, this.isLosersSide)

              var clr
              if (this.player1 == 0 || this.player2 == 0) {
                  clr = this.pendingGameIdBGColor
              } else {
                  clr = this.gameIdBGColor.alpha(alpha)
              }

              this.fillBox(ctx, x, y, width, alpha, clr, this.isLosersSide)

              // highlight if selected
              if (selection != null && (selection.node.Id == this.Id)) {
                  this.highlight(ctx, x, y, width, this.isLosersSide)
              }
              this.drawPlayerIndicator(ctx, x, y, width, 1, this.isLosersSide, 0);
              this.drawPlayerIndicator(ctx, x, y, width, 2, this.isLosersSide,  0);

              ctx.strokeStyle = chroma("black").alpha(alpha)

              if (this.game.result != null) {
                  ctx.strokeStyle = chroma("black").brighten(brighten)
              }

              if (this.Id != 0) {
                  this.drawID(ctx, x, y, width, this.isLosersSide)
              }


              if (this.left) {

                  ctx.strokeStyle = chroma("black").brighten(brighten)
                  var leftY = y - this.left.span.lower;
                  if (this.isLosersSide) {
                      var leftX = x + width + space
                      drawSegment(ctx,
                          x + width,
                          y + (constants.NODE_HEIGHT / 4),
                          leftX,
                          leftY + (constants.NODE_HEIGHT / 2))

                      this.left.renderCentered(ctx, leftX, leftY, level + 1, degree, selection)
                  } else {
                      drawSegment(ctx,
                          x - space,
                          leftY + (constants.NODE_HEIGHT / 2),
                          x,
                          y + (constants.NODE_HEIGHT / 4))

                      this.left.renderCentered(ctx,
                          x - width - space,
                          leftY,
                          level + 1, degree, selection)
                  }


              }

              if (this.right) {
                  ctx.strokeStyle = chroma("black").brighten(brighten)

                  var rightX = 0;
                  var rightY = y + this.right.span.upper;

                  if (this.isLosersSide) {

                      rightX = x + width
                      drawSegment(ctx,
                          rightX,
                          y + ((constants.NODE_HEIGHT / 4) * 3),
                          rightX + space,
                          rightY + (constants.NODE_HEIGHT / 2))

                      this.right.renderCentered(ctx,
                          rightX+space, rightY, level + 1, degree, selection)
                  } else {
                      drawSegment(ctx,
                          x - space,
                          rightY + (constants.NODE_HEIGHT / 2),
                          x,
                          y + ((constants.NODE_HEIGHT / 4) * 3))

                      this.right.renderCentered(ctx,
                          x - width - space, rightY, level + 1, degree, selection)

                  }

              }
      };

      this.renderRightToLeft = function (ctx, x, y, level, degree, selection) {

//          this.log();
          var width = this.width

          if (this.isLoserSide != true) {
              width = width * 0.7;
          }


          var space = constants.NODE_SPACE;

          ctx.strokeStyle = chroma("black");
          var alpha = 1.0;
          var brighten = 1;
          if (this.game.state.result != null) {
              ctx.strokeStyle = chroma("black").brighten(3);
              alpha = 0.3;
              brighten = 3;
          }

          this.x = x;
          this.y = y;
          console.log("node: " + this.Id + " - x: " + this.x + ", y: " + this.y);

          this.frameBox(ctx, x, y, width, alpha, this.upperBGColor, this.lowerBGColor, false);

          var clr
          if (this.isLosersSide) {
              clr = this.losersGameIdBGColor;
          } else {
              clr = this.gameIdBGColor
          }

          this.fillBox(ctx, x, y, width, alpha, clr, false)

          if (selection != null && (selection.node.Id == this.Id)) {
              this.highlight(ctx, x, y, width, false)
          }

          var winningSlot = 0;
          if (this.game.state.result != null) {
              winningSlot = this.game.state.result.winningSlot;
          }

          this.drawPlayerIndicator(ctx, x, y,width, 1, false, winningSlot);
          this.drawPlayerIndicator(ctx, x, y,width, 2, false, winningSlot);

          ctx.strokeStyle = chroma("black").alpha(alpha)

          if (this.game.state.result != null) {
              ctx.strokeStyle = chroma("black").brighten(brighten)
          }

          this.drawID(ctx, x, y, width, false);

          if (this.left) {
//              console.log("Id = " + this.Id);
//              console.log("left = " + this.left);
//              console.log("left.span = " + this.left.span);

              ctx.strokeStyle = chroma("black").brighten(brighten)
              var leftY = y - this.left.span.lower;
                  drawSegment(ctx,
                      x - space,
                      leftY + (constants.NODE_HEIGHT / 2),
                      x,
                      y + (constants.NODE_HEIGHT / 4))

                  this.left.renderRightToLeft(ctx,
                      x - width - space,
                      leftY,
                      level + 1, degree, selection)
          }

          if (this.right) {
              ctx.strokeStyle = chroma("black").brighten(brighten);

//              console.log("Id = " + this.Id);
//              console.log("right = " + this.right);
//              console.log("right.span = " + this.right.span);
              var rightY = y + this.right.span.upper;

              drawSegment(ctx,
                  x - space,
                  rightY + (constants.NODE_HEIGHT / 2),
                  x,
                  y + ((constants.NODE_HEIGHT / 4) * 3))

              this.right.renderRightToLeft(ctx,
                  x - width - space, rightY, level + 1, degree, selection)

          }
      };


          this.getDisplay = function (side) {
              if (side == 1) {

                  if (this.player1 != 0) {
                      return this.playerName(this.player1)
                  }

                  if (this.dropGame1 != 0) {
                      var prefix = ">> L";
                      return prefix + this.dropGame1.toString();
                  }

                  return "";
//                  var pos = "x: " + this.x + ", y: " + this.y;
//                  return pos
              } else {
                  if (this.player2 != 0) {
                      return this.playerName(this.player2)
                  }
                  if (this.dropGame2 != 0) {
                      var prefix = ">> L";
                      return prefix + this.dropGame2.toString();
                  }
              }
              return ""
          };

          this.SetSelected = function (value) {
              this.selected = value
          };


          this.IntersectGame = function (x, y) {

              if ((x > this.x && x < this.x + this.width) &&
                  (y > this.y && y < this.y + constants.NODE_HEIGHT)) {
                  return {node: this}
              }

              if (this.right != null) {
                  var result = this.right.IntersectGame(x, y)
                  if (result != null) {
                      return result
                  }
              }
              if (this.left != null) {
                  var result = this.left.IntersectGame(x, y)
                  if (result != null) {
                      return result
                  }
              }
              return null;

          }

      }
  }

function calcSize(data) {

      let width = ((data.depth-1) * constants.NODE_WIDTH) +
                  ((data.depth-2) * constants.NODE_SPACE);
      let height = data.root.span.upper + data.root.span.lower + 100;

      return {width: width, height: height}
}
  /*
function calcSize(data) {
    var h = data.root.span.upper + data.root.span.lower + 100;
    var w = 1000;

    var len = Object.keys(data.players).length;

    if (len >= 32) {
       // w += 600;
        w += 300 + (10 * len);
        //h += 800;
    } else if (len >= 16) {
        w += 300 + (10 * len);
        //h += 200;
    }

    return {w, h}
}
*/

function calcRatio() {
      let ctx = document.createElement("canvas").getContext("2d"),
          dpr = window.devicePixelRatio || 1,
          bsr = ctx.webkitBackingStorePixelRatio ||
              ctx.mozBackingStorePixelRatio ||
              ctx.msBackingStorePixelRatio ||
              ctx.oBackingStorePixelRatio ||
              ctx.backingStorePixelRatio || 1;
      return dpr / bsr;
  }


  function drawSegment(ctx, srcX, srcY, dstX, dstY) {

          ctx.beginPath()
          ctx.moveTo(srcX, srcY)
          ctx.lineTo(srcX + ((dstX - srcX) / 2), srcY)
          ctx.closePath()
          ctx.stroke()

          ctx.beginPath()
          ctx.moveTo(srcX + ((dstX - srcX) / 2), srcY)
          ctx.lineTo(srcX + ((dstX - srcX) / 2), dstY)
          ctx.closePath()
          ctx.stroke()

          ctx.beginPath()
          ctx.moveTo(srcX + ((dstX - srcX) / 2), dstY)
          ctx.lineTo(dstX, dstY)
          ctx.closePath()
          ctx.stroke()


  }

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}


$('#home_icon').on('click', function () {
    var page = "/home";
    window.location.assign(page)
});


