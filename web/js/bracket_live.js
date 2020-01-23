'use strict';


import {Bracket, Game, Slot} from "./bracket.js";
import * as constants from './constants.js';



function getUrlVars() {

    var vars = {};
    var parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
        vars[key] = value;
    });
    return vars;
}


window.onload = function () {

    this.bracket = new BracketLive();
//    bracket.bracketBase.loadData(canvas)


    this.canvas = document.getElementById("bracket_canvas")
    bracket.bracketBase.loadData(this.canvas)

};


class BracketLive {


      constructor () {
          this.fullscreen = document.getElementById("fullscreen");
          this.canvas = document.getElementById("bracket_canvas")
          this.bracketBase = new Bracket()
//          this.bracketBase.loadData(this.canvas)

          this.hasInput = false
          this.canvas.addEventListener("mousedown",
              this.mousedownHandler.bind(this), false)

          $("#bracket_canvas").keydown(this.keypressHandler.bind(this))
          document.addEventListener("keydown", this.keypressHandler.bind(this), false)

      }

      keypressHandler(evt) {

          var self = this
          console.log("code: " + evt.code)
          if (this.bracketBase.selection != null) {
              var slot = 0
              if (evt.code == 'Digit1') {
                  slot = 1
              } else if (evt.code == 'Digit2') {
                  slot = 2
              }

              if (slot != 0) {
                  var url = "/api/tournaments/" + window.getTournamentID()
                      + "/games/" + this.bracketBase.selection.node.Id + "/winner/"
                      + slot;
                  console.log("url = " + url)
                  $.post( url, function(data) {
                      console.log(data)
                      var finalGame = new Game(data, data.nodes[data.rootNodeId])
                      self.bracketBase.data = data
                      self.bracketBase.rootNode = finalGame;
                      self.bracketBase.render(data, canvas)
                  })
                  .fail(function() {
                      alert( "error completing game" );
                  })
              }
              else if (evt.code == 'KeyU') {
                  var url = "/api/tournaments/" + window.getTournamentID()
                      + "/games/" + this.bracketBase.selection.node.Id + "/winner";

                  $.ajax
                  ({
                      type: "DELETE",
                      //the url where you want to sent the userName and password to
                      url: url,
                      dataType: 'json',
                      async: false,
                      success: function (data) {
                          var finalGame = new Game(data, data.nodes[data.rootNodeId])
                          self.bracketBase.data = data
                          self.bracketBase.rootNode = finalGame;
                          self.bracketBase.render(data, canvas)
                      },
                      fail: function () {
                          alert( "error reseting game");
                      }
                  });


              }

          }
      }

      getMousePos(canvas, evt) {
          var rect = canvas.getBoundingClientRect();
          console.log("x: " + (evt.clientX - rect.left) + ", y: " + (evt.clientY - rect.top));
          return {
              //x: evt.clientX - rect.left,
              //y: evt.clientY - rect.top
              x: evt.clientX - rect.left,
              y: evt.clientY - rect.top
          };
      };

      render() {
          this.bracketBase.loadData(this.canvas)
      };


      mousedownHandler(evt) {

          var canvas = document.getElementById("bracket_canvas")
          var fullscreen_div = document.getElementById("fullscreen_div")
          var mousePos = this.getMousePos(canvas, evt);
          var message = 'Mouse position: ' + mousePos.x + ',' + mousePos.y;
          this.bracketBase.selection = this.bracketBase.rootNode.IntersectGame(mousePos.x, mousePos.y)

          if (this.bracketBase.selection != null) {
              console.log("selection: " + this.bracketBase.selection.node.Id);
              this.bracketBase.render(this.data, canvas);
              return
          }


//          console.log("x = " + fullscreen_div.offsetLeft + ", y = " + fullscreen_div.offsetTop);
//          console.log("height = " + fullscreen_div.height22k22 + ", width = " + fullscreen_div.width);
//          console.log("mouse.x = " + mousePos.x + ", mouse.y = " + mousePos.y);

          if ((mousePos.x >= fullscreen_div.offsetLeft &&
               mousePos.x <= fullscreen_div.offsetLeft + fullscreen_div.width) &&
              (mousePos.y >= fullscreen_div.offsetTop &&
              mousePos.y <= fullscreen_div.offsetTop + fullscreen_div.height) ) {

              var elem = document.getElementById("container-canvas");
              openFullscreen(elem);

          }


//          if mousePos.x

      }



  }


function openFullscreen(elem) {
    if (elem.requestFullscreen) {
        elem.requestFullscreen();
    } else if (elem.mozRequestFullScreen) { /* Firefox */
        elem.mozRequestFullScreen();
    } else if (elem.webkitRequestFullscreen) { /* Chrome, Safari and Opera */
        elem.webkitRequestFullscreen();
    } else if (elem.msRequestFullscreen) { /* IE/Edge */
        elem.msRequestFullscreen();
    }
}

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
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

function checkCookie() {
    var user = getCookie("username");
    if (user != "") {
        alert("Welcome again " + user);
    } else {
        user = prompt("Please enter your name:", "");f
        if (user != "" && user != null) {
            setCookie("username", user, 365);
        }
    }
}




