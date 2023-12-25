const wsURI = "ws://localhost:3550/book-ticket";

const totalTicketInput = document.getElementById("total-ticket")
let ws;

function connectToWS(wsURI, totalTicketRequests) {
    return new WebSocket(wsURI+ `?tickets=${totalTicketRequests}`);
}

function handleSubmitTickets() {
    ws = connectToWS(wsURI, totalTicketInput.value);
    ws.onerror = function () {
        console.log("Web Socket Connection Error!");
    }
}

window.addEventListener("unload", function (){
    ws.close();
});