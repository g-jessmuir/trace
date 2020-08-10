$(function () {
    var ws;

    if (window.WebSocket === undefined) {
        $("#container").append("Your browser does not support WebSockets");
        return;
    } else {
        ws = initWS();
    }

    function initWS() {
        var socket = new WebSocket("ws://localhost:8080/ws"),
            container = $("#container")
        socket.onmessage = function (e) {
            var msg = JSON.parse(e.data);
            var val = msg.data;
            container.html("")
            if ($("#rendering").length == 0) {
                container.append(`<img id="rendering"/>`);
            }
            $("#rendering").attr("src", `data:image/png;base64,${val}`);
            if (msg.status == "done") {
                container.append("<p>Finished!</p>");
            }
        }
        return socket;
    }

    $("#startBtn").click(function (e) {
        e.preventDefault();
        ws.send(JSON.stringify({
            Seed: parseInt($("#seedField").val()),
            Samples: parseInt($("#samplesField").val()),
            Threads: parseInt($("#threadsField").val())
        }));
    });
});