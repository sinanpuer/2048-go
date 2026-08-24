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
    // The ambient wall is flavor for the waiting/menu-ish screens; hide it
    // once someone is actively playing so it doesn't compete for attention
    // with their own board.
    var wall = document.getElementById("ambient-wall");
    if (wall) wall.classList.toggle("hidden", id === "screen-play");
  }

  // ---------- ambient background: a wall of small self-playing boards,
  // mirroring the animated backdrop on the desktop app's main menu ----------

  function ambientEmptyBoard() {
    return [[0, 0, 0, 0], [0, 0, 0, 0], [0, 0, 0, 0], [0, 0, 0, 0]];
  }

  function ambientSpawnTile(b) {
    var empty = [];
    for (var r = 0; r < 4; r++) {
      for (var c = 0; c < 4; c++) {
        if (b[r][c] === 0) empty.push([r, c]);
      }
    }
    if (!empty.length) return false;
    var pos = empty[Math.floor(Math.random() * empty.length)];
    b[pos[0]][pos[1]] = Math.random() < 0.1 ? 4 : 2;
    return true;
  }

  function ambientFreshBoard() {
    var b = ambientEmptyBoard();
    ambientSpawnTile(b);
    ambientSpawnTile(b);
    return b;
  }

  function ambientCompressRow(row) {
    var values = row.filter(function (v) { return v !== 0; });
    var merged = [];
    for (var i = 0; i < values.length; i++) {
      if (i < values.length - 1 && values[i] === values[i + 1]) {
        merged.push(values[i] * 2);
        i++;
      } else {
        merged.push(values[i]);
      }
    }
    while (merged.length < 4) merged.push(0);
    var changed = false;
    for (var j = 0; j < 4; j++) {
      if (merged[j] !== row[j]) changed = true;
    }
    return { row: merged, changed: changed };
  }

  function ambientMoveLeft(b) {
    var nb = [], moved = false;
    for (var r = 0; r < 4; r++) {
      var res = ambientCompressRow(b[r]);
      nb.push(res.row);
      if (res.changed) moved = true;
    }
    return { board: nb, moved: moved };
  }

  function ambientReverseRow(row) { return row.slice().reverse(); }

  function ambientMoveRight(b) {
    var res = ambientMoveLeft(b.map(ambientReverseRow));
    return { board: res.board.map(ambientReverseRow), moved: res.moved };
  }

  function ambientTranspose(b) {
    var t = ambientEmptyBoard();
    for (var r = 0; r < 4; r++) {
      for (var c = 0; c < 4; c++) t[c][r] = b[r][c];
    }
    return t;
  }

  function ambientMoveUp(b) {
    var res = ambientMoveLeft(ambientTranspose(b));
    return { board: ambientTranspose(res.board), moved: res.moved };
  }

  function ambientMoveDown(b) {
    var res = ambientMoveRight(ambientTranspose(b));
    return { board: ambientTranspose(res.board), moved: res.moved };
  }

  var AMBIENT_MOVES = [ambientMoveUp, ambientMoveDown, ambientMoveLeft, ambientMoveRight];

  function ambientRenderBoard(entry) {
    var tiles = entry.el.children;
    var idx = 0;
    for (var r = 0; r < 4; r++) {
      for (var c = 0; c < 4; c++) {
        var v = entry.board[r][c];
        var colors = tileColors(v);
        var tile = tiles[idx++];
        tile.style.background = v === 0 ? "#3a372f" : colors[0];
        tile.style.color = colors[1];
        tile.textContent = v === 0 ? "" : String(v);
      }
    }
  }

  function ambientTick(entry) {
    var order = [0, 1, 2, 3].sort(function () { return Math.random() - 0.5; });
    var moved = false, next = null;
    for (var i = 0; i < order.length; i++) {
      var res = AMBIENT_MOVES[order[i]](entry.board);
      if (res.moved) {
        next = res.board;
        moved = true;
        break;
      }
    }
    entry.board = moved ? next : ambientFreshBoard();
    if (moved) ambientSpawnTile(entry.board);
    ambientRenderBoard(entry);
  }

  function initAmbientWall() {
    var wall = document.createElement("div");
    wall.id = "ambient-wall";
    var overlay = document.createElement("div");
    overlay.id = "ambient-overlay";
    document.body.insertBefore(overlay, document.body.firstChild);
    document.body.insertBefore(wall, document.body.firstChild);

    var boardPx = 4 * 17 + 3 * 3 + 2 * 8; // tile + gap + padding, matches CSS
    var cols = Math.max(1, Math.ceil(window.innerWidth / boardPx));
    var rows = Math.max(1, Math.ceil(window.innerHeight / boardPx));
    wall.style.gridTemplateColumns = "repeat(" + cols + ", " + boardPx + "px)";

    for (var i = 0; i < cols * rows; i++) {
      var boardEl = document.createElement("div");
      boardEl.className = "ambient-board";
      for (var t = 0; t < 16; t++) {
        boardEl.appendChild(document.createElement("div")).className = "ambient-tile";
      }
      wall.appendChild(boardEl);

      var entry = { el: boardEl, board: ambientFreshBoard() };
      ambientRenderBoard(entry);

      (function (entry) {
        var tick = function () {
          ambientTick(entry);
          setTimeout(tick, 900 + Math.random() * 600);
        };
        setTimeout(tick, Math.random() * 1200);
      })(entry);
    }
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

  function randomRoomCode() {
    var chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"; // no 0/O/1/I
    var code = "";
    for (var i = 0; i < 5; i++) code += chars[Math.floor(Math.random() * chars.length)];
    return code;
  }

  var params = new URLSearchParams(location.search);
  var roomCode = params.get("room");
  var isNewRoom = !roomCode;
  if (!roomCode) {
    roomCode = randomRoomCode();
    params.set("room", roomCode);
    history.replaceState(null, "", location.pathname + "?" + params.toString());
  }
  if (isNewRoom) {
    el("room-share").textContent = "Dein Einladungslink zum Teilen:";

    if (typeof qrcode === "function") {
      var qr = qrcode(0, "M");
      qr.addData(location.href);
      qr.make();
      var qrBox = el("room-qr");
      qrBox.innerHTML = qr.createSvgTag({ cellSize: 5, margin: 8, scalable: true });
      qrBox.classList.remove("hidden");
    }

    var copyBtn = el("copy-link-btn");
    copyBtn.classList.remove("hidden");
    copyBtn.addEventListener("click", function () {
      var reset = function () { copyBtn.textContent = "Link kopieren"; };
      var showCopied = function () {
        copyBtn.textContent = "Kopiert!";
        setTimeout(reset, 1500);
      };
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(location.href).then(showCopied, function () {
          copyBtn.textContent = "Kopieren fehlgeschlagen";
          setTimeout(reset, 1500);
        });
      } else {
        // Fallback for browsers/contexts without the async Clipboard API.
        var tmp = document.createElement("textarea");
        tmp.value = location.href;
        tmp.style.position = "fixed";
        tmp.style.opacity = "0";
        document.body.appendChild(tmp);
        tmp.select();
        try {
          document.execCommand("copy");
          showCopied();
        } catch (e) {
          copyBtn.textContent = "Kopieren fehlgeschlagen";
          setTimeout(reset, 1500);
        }
        document.body.removeChild(tmp);
      }
    });
  } else {
    el("room-share").textContent = "Raum: " + roomCode;
  }

  function connect() {
    var proto = location.protocol === "https:" ? "wss://" : "ws://";
    socket = new WebSocket(proto + location.host + "/ws?room=" + encodeURIComponent(roomCode));
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

  el("restart-btn").addEventListener("click", function () {
    send({ type: "restart" });
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

    var isHost = youId === msg.hostId;
    el("restart-btn").classList.toggle("hidden", !(isGameOver && isHost));
    el("restart-wait").classList.toggle("hidden", !(isGameOver && !isHost));
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

  initAmbientWall();
  showScreen("screen-name");
})();
