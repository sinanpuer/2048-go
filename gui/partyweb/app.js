(function () {
  "use strict";

  var TILE_COLORS = {
    0: ["#cdc1b4", "#cdc1b4"],
    2: ["#eee4da", "#776e65"],
    4: ["#cfc7bb", "#776e65"],
    8: ["#f2b179", "#ffffff"],
    16: ["#f59563", "#ffffff"],
    32: ["#f67c5f", "#ffffff"],
    64: ["#f65e3b", "#ffffff"],
    128: ["#edcf72", "#ffffff"],
    256: ["#edcc61", "#ffffff"],
    512: ["#edc850", "#ffffff"],
    1024: ["#edc53f", "#ffffff"],
    2048: ["#edc22e", "#ffffff"],
  };

  function tileColors(v) {
    return TILE_COLORS[v] || ["#3c3a32", "#edc22e"];
  }

  function el(id) { return document.getElementById(id); }

  function showScreen(id) {
    ["screen-name", "screen-lobby", "screen-countdown", "screen-play", "screen-overview"].forEach(function (s) {
      el(s).classList.toggle("hidden", s !== id);
    });
  }

  function renderBoard(container, board, big) {
    container.innerHTML = "";
    container.className = "board " + (big ? "board-big" : "board-mini");
    for (var r = 0; r < 4; r++) {
      for (var c = 0; c < 4; c++) {
        var v = board[r][c];
        var colors = tileColors(v);
        var tile = document.createElement("div");
        tile.className = "tile";
        tile.style.background = colors[0];
        tile.style.color = colors[1];
        tile.textContent = v === 0 ? "" : String(v);
        container.appendChild(tile);
      }
    }
  }

  var socket = null;
  var youId = null;
  var joined = false;

  function connect() {
    var proto = location.protocol === "https:" ? "wss://" : "ws://";
    socket = new WebSocket(proto + location.host + "/ws");
    socket.onmessage = function (ev) {
      var msg = JSON.parse(ev.data);
      if (msg.type === "state") handleState(msg);
    };
    socket.onclose = function () {
      if (joined) {
        el("name-error").textContent = "Verbindung verloren. Seite neu laden, um erneut beizutreten.";
      }
    };
  }

  function send(obj) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(obj));
    }
  }

  el("join-btn").addEventListener("click", function () {
    if (joined) return;
    joined = true;
    var name = el("name-input").value.trim();
    connect();
    socket.addEventListener("open", function () {
      send({ type: "join", name: name });
    });
  });

  el("start-btn").addEventListener("click", function () {
    send({ type: "start" });
  });

  function handleState(msg) {
    youId = msg.youId;

    if (msg.error) {
      el("name-error").textContent = msg.error;
      return;
    }
    el("name-error").textContent = "";

    var me = null;
    for (var i = 0; i < msg.players.length; i++) {
      if (msg.players[i].id === youId) me = msg.players[i];
    }

    if (msg.phase === "lobby") {
      showScreen("screen-lobby");
      renderLobby(msg);
    } else if (msg.phase === "countdown") {
      showScreen("screen-countdown");
      el("countdown-number").textContent = msg.countdown;
    } else if (msg.phase === "playing") {
      if (me && me.alive) {
        showScreen("screen-play");
        renderPlay(msg, me);
      } else {
        showScreen("screen-overview");
        renderOverview(msg, false);
      }
    } else if (msg.phase === "gameover") {
      showScreen("screen-overview");
      renderOverview(msg, true);
    }
  }

  function renderLobby(msg) {
    var list = el("lobby-list");
    list.innerHTML = "";
    msg.players.forEach(function (p) {
      var li = document.createElement("li");
      var nameSpan = document.createElement("span");
      nameSpan.textContent = p.name;
      li.appendChild(nameSpan);
      if (p.id === msg.hostId) {
        var tag = document.createElement("span");
        tag.className = "host-tag";
        tag.textContent = "HOST";
        li.appendChild(tag);
      }
      list.appendChild(li);
    });
    el("lobby-count").textContent = msg.players.length + " / 9 Spieler";

    var isHost = youId === msg.hostId;
    var canStart = isHost && msg.players.length >= 2;
    el("start-btn").classList.toggle("hidden", !isHost);
    el("start-btn").disabled = !canStart;
    el("lobby-wait").classList.toggle("hidden", isHost);
  }

  function renderPlay(msg, me) {
    el("play-status").innerHTML =
      "<span>Punkte: <b>" + me.score + "</b></span>" +
      (me.combo > 0 ? "<span>Combo x" + me.combo + "</span>" : "<span></span>");
    renderBoard(el("own-board"), me.board, true);

    var sidebar = el("sidebar");
    sidebar.innerHTML = "";
    msg.players.forEach(function (p) {
      if (p.id === youId) return;
      var row = document.createElement("div");
      row.className = "sidebar-row" + (p.alive ? "" : " dead");
      var dot = document.createElement("span");
      dot.className = "dot";
      row.appendChild(dot);
      var name = document.createElement("span");
      name.className = "sidebar-name";
      name.textContent = p.name;
      row.appendChild(name);
      var score = document.createElement("span");
      score.textContent = p.score;
      row.appendChild(score);
      sidebar.appendChild(row);
    });
  }

  function renderOverview(msg, isGameOver) {
    el("overview-title").textContent = isGameOver ? "Spiel vorbei" : "Du bist raus - Zuschauer-Uebersicht";
    var grid = el("overview-grid");
    grid.innerHTML = "";
    msg.players.forEach(function (p) {
      var card = document.createElement("div");
      card.className = "overview-card" + (p.alive ? "" : " dead") + (p.id === msg.winnerId ? " winner" : "");

      var name = document.createElement("div");
      name.className = "overview-name";
      name.textContent = p.name + (p.id === youId ? " (Du)" : "");
      card.appendChild(name);

      var boardEl = document.createElement("div");
      renderBoard(boardEl, p.board, false);
      card.appendChild(boardEl);

      var score = document.createElement("div");
      score.className = "overview-score";
      score.textContent = "Punkte: " + p.score;
      card.appendChild(score);

      if (isGameOver && p.id === msg.winnerId) {
        var banner = document.createElement("div");
        banner.className = "victory-banner";
        banner.textContent = "Gewonnen!";
        card.appendChild(banner);
      }

      grid.appendChild(card);
    });
  }

  function move(dir) {
    send({ type: "move", dir: dir });
  }

  el("dpad-up").addEventListener("click", function () { move("up"); });
  el("dpad-down").addEventListener("click", function () { move("down"); });
  el("dpad-left").addEventListener("click", function () { move("left"); });
  el("dpad-right").addEventListener("click", function () { move("right"); });

  document.addEventListener("keydown", function (ev) {
    if (el("screen-play").classList.contains("hidden")) return;
    switch (ev.key) {
      case "ArrowUp": move("up"); ev.preventDefault(); break;
      case "ArrowDown": move("down"); ev.preventDefault(); break;
      case "ArrowLeft": move("left"); ev.preventDefault(); break;
      case "ArrowRight": move("right"); ev.preventDefault(); break;
    }
  });

  showScreen("screen-name");
})();
