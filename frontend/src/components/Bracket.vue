<script lang="ts">

import chroma from 'chroma-js';
import {axiosApiInstance} from '../main';
import * as constants from './constants';

import { defineComponent} from 'vue';

var theBracket;

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


function calcRatio() {
      let ctx = document.createElement("canvas").getContext("2d"),
          dpr = window.devicePixelRatio || 1
      return dpr / 1;
  }


const material_font = new FontFace( 'material-icons',
  // pass the url to the file in CSS url() notation
  'url(https://fonts.gstatic.com/s/materialicons/v48/flUhRq6tzZclQEJ-Vdg-IuiaDsNcIhQ8tQ.woff2)' );
document.fonts.add( material_font ); // add it to the document's FontFaceSet



class Bracket {

    x: number;
    y: number;
    selection: Selection | null;
    orientation: number;
    data: Tournament
    tournamentId: number;
    rootNode? : Game;
    canvas? : HTMLCanvasElement;
    winnersHeaderButtons: HTMLElement | null;
    losersHeaderButtons: HTMLElement | null;

    drawcenter : number;
    ratio: number;

    Preferences: Preferences;

    width: number;
    height: number;

    isFullScreen: boolean;


    constructor(tournament, preferences) {
          this.data = tournament

          this.isFullScreen = false;
          if (preferences != null) {
              this.Preferences = new Preferences(preferences)
          } else {
              this.Preferences = DefaultPreferences()
          }
          this.x = 0;
          this.y = 0;
          this.selection = null;
          this.orientation = constants.RIGHT_TO_LEFT;
          this.canvas = document.getElementById("bracket-canvas") as HTMLCanvasElement;

          this.winnersHeaderButtons = document.getElementById("winners-header-buttons");
          this.losersHeaderButtons = document.getElementById("losers-header-buttons");

          this.drawcenter = 0;
          this.ratio = 0;

          if (tournament != null) {
             let br = tournament.bracket;
             var finalGame = new Game(this, tournament, br, br.nodes[br.rootNodeId]);
             this.rootNode = finalGame
             this.tournamentId = tournament.id
          } else {
              this.tournamentId = 0;
          }


          this.width = 0;
          this.height = 0;

          this.canvas.addEventListener("mousedown",
          this.bracketMousedownHandler.bind(this), false)
          document.addEventListener("keydown", this.keypressHandler.bind(this), false);

       
    }

    getMousePos(canvas, evt) {
        var rect = canvas.getBoundingClientRect();
        return {
            x: evt.clientX - rect.left,
            y: evt.clientY - rect.top
        };
    }

    keypressHandler (evt) {

        if (this.selection != null) {


            var hasWinner = false;
            if (this.selection.node.state.result != null) {
                hasWinner = true;
            }
            var slot = 0
            if (evt.code == 'Digit1') {
                slot = 1
            } else if (evt.code == 'Digit2') {
                slot = 2
            } else if ((evt.code == 'KeyU') && (hasWinner == true)) {
                let self = this
                var url = "/tournaments/" + self.tournamentId
                    + "/games/" + this.selection.node.Id + "/winner";

                axiosApiInstance.delete(url)
                    .then(function (response) {
                        self.update(response.data)
                    })
                    .catch(function (error) {
                        console.log(error);
                    })
                    .then(function () {
                    })
            }

            if ((slot != 0) && (hasWinner == false)) {
                var url = "/tournaments/" + this.tournamentId
                    + "/games/" + this.selection.node.Id + "/winner/"
                    + slot;
                //console.log("url = " + url);

                let self = this

                axiosApiInstance.post(url)
                    .then(function (response) {
                        var navTournamentStatus = document.getElementById("nav-tournament-state");
                        if (navTournamentStatus != null) {
                            navTournamentStatus.innerHTML = "Status: " + self.data.tournamentState;
                        }
                        self.update(response.data)
                    })
                    .catch(function (error) {
                        console.log(error);
                    })
                    .then(function () {
                    })
            } 

        }
    }

    bracketMousedownHandler (evt) {

        var canvas = document.getElementById("bracket-canvas");
        var fullscreen_div = document.getElementById("fullscreen_div");
        var mousePos = this.getMousePos(canvas, evt);
        var message = 'Mouse position: ' + mousePos.x + ',' + mousePos.y;
        if (this.rootNode != null) {
            this.selection = this.rootNode.IntersectGame(mousePos.x, mousePos.y)
        }

        if (this.selection != null) {
            //console.log("selection: " + this.selection.node.Id);
            this.render()
            return
        }

        //var message = 'Mouse position: ' + mousePos.x + ',' + mousePos.y;
        //console.log(message);

        let absMousePos = mousePos;
        absMousePos.x -= window.scrollX;
        absMousePos.y -= window.scrollY;

    }


          resize (data, rootNode, canvas) {
              //console.log("data = " + JSON.stringify(data));
              let size = this.calcSize(data.bracket, rootNode);
              canvas.width = size.width;
              canvas.height = size.height;
              canvas.style.width = size.width + "px";
              canvas.style.height = size.height + "px";
              this.width = size.width;
              this.height = size.height;
              //resizeCanvasToDisplaySize(this.canvas);

              this.ratio = calcRatio();

              console.log("*******size= " + JSON.stringify(size) + ", ratio = " + this.ratio)
          } 


        setup() {

            let canvas = this.canvas;
            if (canvas != null) {
                let ctx = canvas.getContext('2d');
                material_font.load().then( () => {
                    }).catch( console.error );

            }


        }
        render() {

            let data = this.data;
            let rootNode = this.rootNode;
            let canvas = this.canvas;

            let ctx;
            if (canvas != null) {
                ctx = canvas.getContext('2d');

                ctx.clearRect(0, 0, canvas.width, canvas.height);

                //ctx.fillStyle = chroma(constants.COLOR_DARK_BACKGROUND);
                ctx.fillStyle = chroma(this.Preferences.GetBGColor());
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                if (this.orientation == constants.RIGHT_TO_LEFT) {
                    this.renderRightToLeft(data, rootNode, canvas);
                }
            }
        };

        renderRightToLeft (data, rootNode, canvas) {

            let ctx = canvas.getContext('2d');
            var width = this.width;

            let drawcenter = this.height / 2 - 100;
            if (rootNode) {
                drawcenter = rootNode.span.upper + 22;
            }

            ctx.strokeStyle = chroma(constants.COLOR_3);
            ctx.fillStyle = chroma(constants.COLOR_3);
            ctx.textAlign = "center";
            ctx.font = "36px Arial";
            ctx.strokeText("Winners", width - 100, (drawcenter) - (this.height / 4));


            ctx.strokeText("Losers", width - 100, (drawcenter) + (this.height / 4));

            ctx.strokeStyle = chroma(this.Preferences.GetGameBorderColor());

            ctx.lineWidth = 1;
            ctx.beginPath();
            ctx.moveTo(width - (constants.NODE_WIDTH + 10), drawcenter);
            ctx.lineTo(width, drawcenter);
            ctx.closePath();
            ctx.stroke();

            ctx.strokeStyle = chroma(constants.COLOR_3);
            ctx.fillStyle = chroma(constants.COLOR_3);
            ctx.textAlign = "center";
            ctx.font = "18px Arial";
            var winner = this.getWinner();
            ctx.strokeText(winner, width - ((constants.NODE_WIDTH + 10) / 2),
                (drawcenter) - 8);

            if (this.rootNode != null) {
            this.rootNode.renderRightToLeft(ctx,
                (width) - (constants.NODE_WIDTH + (constants.NODE_WIDTH - 10)),
                (drawcenter) - (constants.NODE_HEIGHT / 2), this.selection)
            }
        }

        getWinner () {

            if (this.data != null) {
              if ((this.data.bracket.root) && (this.data.bracket.root.state.result != null)) {
                  var winnerParticipantNum = this.data.bracket.root.state.result.winningParticipant
                  if (winnerParticipantNum > 0) {
                      var p = this.data.participants[winnerParticipantNum];
                      if (p != null) {
                          return p.name;
                      }
                  }

              }
              return "<TBD>"
            }
        };

        calcSize(bracket, rootNode) {

            if (rootNode) {
                //let width = ((bracket.depth+1) * constants.NODE_WIDTH) +
                //    ((bracket.depth+1) * constants.NODE_SPACE*2);
                console.log("losersDepth = " + bracket.losersDepth);
                let losersDepth = bracket.losersDepth;
                let width = (losersDepth+1) * (constants.NODE_WIDTH + constants.NODE_SPACE) + 180
                let height = bracket.root.span.upper + bracket.root.span.lower + 100;

                return {width: width, height: height}
            }
            return {width: constants.WIDTH, height: constants.HEIGHT};
        }
        

        resizeCanvasToDisplaySize(canvas) {
            // look up the size the canvas is being displayed
            const width = canvas.clientWidth;
            const height = canvas.clientHeight;


            // If it's resolution does not match change it
            if (canvas.width !== width || canvas.height !== height) {
                canvas.width = width;
                canvas.height = height;
                return true;
            }

            return false;
        }

        update (tournament: Tournament): void {
            this.tournamentId = tournament.id
            //self.data = tournament;
            this.data = tournament
            if (this.data.bracket.rootNodeId != 0) {
                let br = this.data.bracket
                var finalGame = new Game(this, tournament, br, br.nodes[br.rootNodeId]);
                this.rootNode = finalGame;
                this.resize(tournament, finalGame, this.canvas);
                this.render()
            } else {
                console.log("No rootNodeID")
            }
        }

   }
  /*
       Leaf level games are 2nd from bottom level nodes.


       A node that has a null left and right child are leaf nodes. A Leaf node either starting position
        where is has a playerId or it represents a drop position in which case it has a gameId
       A Game object is not created for leaf nodes, instead the parent node will have either the player
       or game field set to the player or game the leaf node represents.
   */


 class Game {

     TheBracket: Bracket;
     tournament: Tournament;
     bracket: IBracket
     orientation: number;
     left: Game | null;
     right: Game | null;
     game: Node;
     participant1: number;
     participant2: number;
     dropGame1: number;
     dropGame2: number;
     x: number;
     y: number;
     span: Span;
     gridSpan: Span;
     width: number;
     selection: number;

     connectorColor: chroma.Color;
     Id: number;
     participant: number;
     dropGame: number;
     isLosersSide: boolean;
     subtype: number;

      constructor(itsBracket, Tournament, bracket, TheGame) {

          this.TheBracket = itsBracket;
          this.tournament = Tournament;
          this.bracket = bracket;
          this.game = TheGame;
          this.selection = 0;

          this.orientation = constants.RIGHT_TO_LEFT
          this.participant1 = -1;
          this.participant2 = -1;
          this.dropGame1 = 0;
          this.dropGame2 = 0;

          this.x = 0;
          this.y = 0;
          this.left = null
          this.right = null
          this.span = TheGame.span;
          this.gridSpan = TheGame.gridSpan;
          this.width = constants.NODE_WIDTH;
          this.connectorColor = chroma(itsBracket.Preferences.GetConnectorColor());

          this.Id = TheGame.id;
          this.participant = TheGame.participant;
          this.dropGame = TheGame.drop;
          this.subtype = TheGame.nodeSubType

          if (TheGame.tier == 2) {
              this.isLosersSide = true
          } else {
              this.isLosersSide = false
          }

          if ((TheGame.left.id == 0) && (TheGame.right.id == 0)) {
              let Null = { valueOf: ()=>null }
              return;
          }
          else {
          }

          if (TheGame.left.id != 0) {
              
              var node = this.bracket.nodes[TheGame.left.id];
              var typ = node.nodeType
              if (typ == 1) {
                  this.participant1 = node.participant
              } else if (typ == 2) {
                this.left = new Game(this.TheBracket, this.tournament, this.bracket, node);
                if (node != null) {
                    this.participant1 = node.participant
                }
              } else if (typ == 3) {
                this.dropGame1 = node.drop;
                this.participant1 = node.participant
                this.left = null;
              }
          }
          
          if (TheGame.right.id != 0) {
              var node = this.bracket.nodes[TheGame.right.id];
              var typ = node.nodeType
              if (typ == 1) {
                this.participant2 = node.participant
              } else if (typ == 2) {
                this.right = new Game(this.TheBracket, this.tournament, this.bracket, node);
                if (node != null) {
                    this.participant2 = node.participant
                }
              } else if (typ == 3) {
                this.dropGame2 = node.drop;
                this.participant2 = node.participant
                this.right = null;
              }
          } 
      }

    getLevel() {
        this.game.level;
    }

    playerName(num) {
        var p = this.tournament.participants[num];
        if (p != null) {
            return p.name;
        } else {
            return "BUY"
        }
    }

    highlight (ctx, reversed) {

        let width = constants.NODE_WIDTH;
        ctx.strokeStyle = chroma(constants.COLOR_22);
        ctx.lineWidth = 3;
        if (reversed) {
            ctx.strokeRect(this.x + 24, this.y, width - 26, constants.NODE_HEIGHT);
            ctx.strokeRect(this.x, this.y, width, constants.NODE_HEIGHT);
        } else {
            ctx.strokeRect(this.x, this.y, width - 26, constants.NODE_HEIGHT);
            ctx.strokeRect(this.x, this.y, width, constants.NODE_HEIGHT);
        }

        ctx.lineWidth = 1;

    }

    renderGame (ctx, selection) {

        ctx.strokeStyle = chroma("black");
        ctx.fillStyle = chroma("black");
        var darken = 0;
        if (this.game.state.result != null) {
            ctx.strokeStyle = chroma("black");
            ctx.fillStyle = chroma("black");
            darken = 2;
        }

        this.frameGame(ctx,  darken)

        var winningSlot = 0;
        if (this.game.state.result != null) {
            winningSlot = this.game.state.result.winningSlot;
        }

        this.renderPlayerIndicator(ctx, 1, winningSlot);
        this.renderPlayerIndicator(ctx, 2, winningSlot);

        ctx.strokeStyle = chroma(this.TheBracket.Preferences.GetGameFontColor());

        this.renderGameId(ctx,false);
        if (selection != null && (selection.node.Id == this.Id)) {
            this.highlight(ctx, false)
        }
    }


    renderGameId (ctx, reversed) {
        let width = constants.NODE_WIDTH;

        ctx.textAlign = "center";
        ctx.font = "15px Arial";
        ctx.fontWeight = "bold";
        if (reversed) {
            ctx.strokeText(this.Id, this.x+12, this.y + (constants.NODE_HEIGHT / 2 + 6));
        } else {
            ctx.strokeText(this.Id, this.x + width - 12, this.y + (constants.NODE_HEIGHT / 2 + 6))
        }

    }

    renderPlayerIndicator (ctx, slot, winningSlot) {

        let width = constants.NODE_WIDTH
        ctx.font = "14px Arial";
        ctx.textAlign = "center";

        let display = "TBD"
        display = this.getDisplay(slot);

        var clr = null;

        if (display.startsWith(">>") === true) {
            clr = chroma(constants.COLOR_1);
        } else {
            clr = chroma(this.TheBracket.Preferences.GetGameFontColor())
        }
        ctx.strokeStyle = clr;
        ctx.fillStyle = clr;

        if (winningSlot != 0) {
            if (winningSlot == slot) {
                ctx.font = "20px Arial";
                ctx.fontWeight = "bold";
                clr = chroma(this.TheBracket.Preferences.GetGameWinnersFontColor())
                ctx.fillStyle = clr;
                display = "*" + display ;
            } else {
                ctx.fontWeight = "lighter";
                ctx.font = "13px Arial";
                clr = chroma(this.TheBracket.Preferences.GetGameLosersFontColor())
                ctx.fillStyle = clr;
            }
            ctx.strokeStyle = chroma("black");
        }

        var yoffset = 0;
        if (slot == 1) {
            yoffset = (constants.NODE_HEIGHT / 2) - 5;
        } else {
            yoffset = (constants.NODE_HEIGHT) - 5;
        }
        ctx.strokeText(display, this.x + ((width - 24) / 2), this.y + yoffset);
        ctx.fillText(display, this.x + ((width - 24) / 2), this.y + yoffset);
    }

    frameGame (ctx, darken)  {

        let height = constants.NODE_HEIGHT;
        let width = constants.NODE_WIDTH;

        // draw the box
        ctx.textAlign = "Left";
        ctx.lineWidth = 1;
        var borderColor = chroma(this.TheBracket.Preferences.GetGameBorderColor());
        //ctx.strokeStyle = borderColor;
        ctx.strokeStyle = chroma("black");


        ctx.clearRect(this.x, this.y, width, height);
        ctx.strokeRect(this.x, this.y, width, height);
        var bkColor;
        if (this.isLosersSide) {
            bkColor = chroma(this.TheBracket.Preferences.GetGameLosersBackgroundColor());
        } else {
            bkColor = chroma(this.TheBracket.Preferences.GetGameWinnersBackgroundColor());
        }
        ctx.fillStyle = bkColor.darken(darken);
        ctx.fillRect(this.x + width - 25,  this.y+1, 24, height - 2);


        var clr1 = chroma(this.TheBracket.Preferences.GetSlot1BGColor());
        ctx.fillStyle = clr1.darken(darken);
        ctx.fillRect(this.x+1, this.y+1, width - 26,height / 2 - 2);

        var clr2 = chroma(this.TheBracket.Preferences.GetSlot2BGColor());
        ctx.fillStyle = clr2.darken(darken);
        ctx.fillRect(this.x+1, this.y + 1 + (height / 2), width - 26, height / 2 - 2);

        ctx.strokeStyle = chroma("black");

        ctx.strokeRect(this.x+width-26, this.y, 26, height);
        ctx.beginPath();
        ctx.moveTo(this.x, this.y + height/2);
        ctx.lineTo(this.x + width - 25, this.y + height / 2);
        ctx.closePath();
        ctx.stroke();

    }

    renderRightToLeft (ctx, x, y, selection) {

        this.x = x;
        this.y = y;
        this.renderGame(ctx, selection);

        let leftDepth = 0, rightDepth = 0;
        let width = constants.NODE_WIDTH;
        var space = constants.NODE_SPACE;
        if (this.left) {

            ctx.strokeStyle = this.connectorColor;

            var leftY = y - this.left.span.lower;
                drawSegment(ctx,
                x - space,
                leftY + (constants.NODE_HEIGHT / 2),
                x,
                y + (constants.NODE_HEIGHT / 4))

            leftDepth = this.left.renderRightToLeft(ctx,
                x - width - space,
                leftY, selection)
        }

        if (this.right) {


            ctx.strokeStyle = this.connectorColor;

            var rightY = y + this.right.span.upper;

            drawSegment(ctx,
                x - space,
                rightY + (constants.NODE_HEIGHT / 2),
                x,
                y + ((constants.NODE_HEIGHT / 4) * 3))

            rightDepth = this.right.renderRightToLeft(ctx,
                x - width - space, rightY, selection)

        }

        if ((leftDepth == 0) && (rightDepth == 0)) {
            return 1
        }
        return Math.max(leftDepth, rightDepth) + 1

    }

    getDisplay (slot) {

        if (slot == 1) {

            if (this.isLosersSide) {
                if (this.participant1 != 0) {
                    return this.playerName(this.participant1)
                } else if (this.dropGame1 != 0) {
                    return "L"+this.dropGame1
                }
            } else {
                if (this.participant1 != 0) {
                    return this.playerName(this.participant1)
                }
            }

            return "";
        } else {

            if (this.isLosersSide) {
                if (this.participant2 != 0) {
                    return this.playerName(this.participant2)
                } else if (this.dropGame2 != 0) {
                    return "L"+this.dropGame2
                }
            } else  {
                if (this.participant2 != 0) {
                    return this.playerName(this.participant2)
                }
            }
        }
        return ""
    }

    SetSelection (value) {
        this.selection = value
    }


    IntersectGame (x, y) {
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


 function DefaultPreferences() {

     let prefs = {
    "brackets.background-color": "#E4DEDD",
	"brackets.connector-color": "#efdf20",
	"brackets.game.border-color": "#1F1C21",
	"brackets.game.slot1.background-color": "#765350",
	"brackets.game.slot2.background-color": "#2269a0",
	"brackets.game.slot1.font-color": "#221C1B",
	"brackets.game.slot2.font-color": "#221C1B",
	"brackets.game.font-color": "#e8e9e8",
	"brackets.game.winners.font-color": "#f9eFeE",
	"brackets.game.losers.font-color": "#2249a3",
	"brackets.game.winners.background-color": "#d62720",
	"brackets.game.losers.background-color": "#2333aD"}

     let p = new Preferences(prefs)
     return p

 }

class Preferences {

    prefs: {}
    constructor (prefs) {
    	this.prefs = prefs;
    }

    GetBGColor () {
    	    return this.prefs["brackets.background-color"];
    }
    GetConnectorColor() {
            return this.prefs["brackets.connector-color"];
    }
    GetSlot1BGColor () {
            return this.prefs["brackets.game.slot1.background-color"];
    }
    GetSlot2BGColor () {
            return this.prefs["brackets.game.slot2.background-color"];
    }
    GetGameWinnersFontColor() {
            return this.prefs["brackets.game.winners.font-color"];
    }
    GetGameLosersFontColor() {
            return this.prefs["brackets.game.losers.font-color"];
    }
    GetGameBorderColor () {
            return this.prefs["brackets.game.border-color"];
    }
    GetGameWinnersBackgroundColor () {
            return this.prefs["brackets.game.winners.background-color"];
    }
    GetGameLosersBackgroundColor () {
            return this.prefs["brackets.game.losers.background-color"];
    }
    GetGameFontColor () {
            return this.prefs["brackets.game.font-color"];
    }
    Update (prefs) {
            this.prefs = prefs;
    }

}

export default defineComponent({
    name: 'Bracket',
    inject: ['router'],
    props: {
        currentTab: String,
        tournamentId: Number,

    },
    watch: {
        currentTab(val) {
            if (val == 'Bracket') {
                this.activateTab()
            }
        }

    },

    data: function() {
        return {
            preferences: null,
            bracket: null,
            isLoading: true,
            error: "",
            tournament: {id:0} as Tournament
        }
    },

    mounted: function() {
        let cmp = this
        cmp.activateTab()
    },

    updated() {

    },


    methods: {
        createBracket() {
            //console.log("createBracket")
            let cmp = this;
            let url = "/tournaments/" + cmp.tournamentId+"/generate"
            axiosApiInstance.put(url)
                .then(function (response) {
                    cmp.tournament = response.data
                    theBracket = new Bracket(cmp.tournament, cmp.preferences)
                    theBracket.setup()
                    theBracket.update(response.data)
                })
                .catch(function (error) {
                    cmp.error = error.toString()
                })
                .then(function () {
                    cmp.isLoading = false
            });
        },
        activateTab() {
            
            
            let promise1 = this.fetchPreferences()
            let promise2 = this.fetchTournament()
            var cmp = this
            Promise.all([promise1, promise2]).then(function() {

                if (cmp.tournament != null) {

                    if (cmp.tournament['tournamentState'] == 'Registration' || 
                            cmp.tournament['bracket']['root'] == null) {
                        cmp.createBracket()
                    } else {
                        theBracket = new Bracket(cmp.tournament, cmp.preferences)
                        theBracket.setup()
                        theBracket.update(cmp.tournament)
                    }
                } else {
                    console.log("No Tournament")
                }
            })

        },

        fetchPreferences() {
            let cmp = this;
            let url = "/preferences"

            return axiosApiInstance.get(url)
                .then(function (response) {
                    cmp.preferences = response.data
                })
                .catch(function (error) {
                    console.log(error);
                })
                .then(function () {
                });
        },    
        fetchTournament() {

            let cmp = this;
            let url = "/tournaments/" + cmp.tournamentId
            return axiosApiInstance.get(url)
                .then(function (response) {
                    cmp.tournament = response.data
                })
                .catch(function (error) {
                    cmp.error = error.toString()
                    console.log(error);
                })
                .then(function () {
                    cmp.isLoading = false
                });
        },
        Shuffle() {
            let cmp = this;
            let url = "/tournaments/" + cmp.tournamentId + "/randomize"
            return axiosApiInstance.put(url)
                .then(function (response) {
                    cmp.fetchTournament().then(function() {
                        cmp.createBracket()
                    })
                })
                .catch(function (error) {
                    cmp.error = error.toString()
                    console.log(error);
                })
                .then(function () {
                });


        },
        toggleFullscren (e) {
            const target = e.target.parentNode.parentNode.parentNode

            $q.fullscreen.toggle(target)
            .then(() => {
            })
            .catch((err) => {
                alert(err)
            })
        }
    }, 
    computed: {

            GamesRemaining() {
                /*
                if (this.tournament != null) {
                let val = ((2*this.tournament.participants.length)-1) - this.tournament.gamesPlayed
                val.toString()
                } 
                */
                return "?"
            }

    }

        

})

</script>


<template>

<div>

  <div id="bracket-panel" class="container">
    <div id="main-bracket-div">
      <div id="bracket-div" style="position: static">
        <canvas id="bracket-canvas" style="margin: 0px"></canvas>
        <!--
                        <input id="fullscreen_div" type="image" src="/Fullscreen.png"
                                name="fullscreen_button" alt="Submit" width="40" height="40"
                                style="pointer-events: none; position: fixed;
                                    top: 112px; left: 14px;  z-index: 1000;"/>
                                    -->
        <!--
                        <q-btn
                            style="top: 112px; left: 14px;  z-index: 1000; position: fixed"
                            color="secondary"
                            @click="toggleFullscreen"
                            :icon="$q.fullscreen.isActive ? 'fullscreen_exit' : 'fullscreen'"
                            :label="$q.fullscreen.isActive ? 'Exit Fullscreen' : 'Go Fullscreen'"
                        ></q-btn>
                        -->
      </div>
    </div>
  </div>
  </div>
</template>


<style scoped>
#bracket-panel {
  min-height: 300px;
}
</style>
