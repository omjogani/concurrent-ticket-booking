const wsURI = "ws://localhost:3550/";

const totalTicketInput = document.getElementById("total-ticket")
let ws;

function connectToWS(wsURI, totalTicketRequests) {
    return new WebSocket(wsURI+ `book-ticket?tickets=${totalTicketRequests}`);
}

function handleSubmitTickets() {
    ws = connectToWS(wsURI, totalTicketInput.value);
    ws.onerror = function () {
        console.log("Web Socket Connection Error!");
    }
    ws.onmessage = function(event) {
        console.log("got update");
        mapDataToUI(event);
    }
}

window.addEventListener("unload", function (){
    ws.close();
});

// ==========================
const statusWS = new WebSocket(wsURI);

// handle error on statusWS
statusWS.onerror = function() {
    console.log(`Error Sync Seats Status!`);
}

// continuous listen seats status
statusWS.onmessage = function(event){
    mapDataToUI(event);
}
    

function mapDataToUI(event){
    const eventData = JSON.parse(event.data);
    console.log(eventData);

    // map event data to GUI
    for(let col = 0; col < eventData[0].length; col++){
        for (let index = 0; index < eventData.length; index++) {
            if(eventData[index][col] == "x"){
                const calculateCell =`${String.fromCharCode(65 + col)}${index+1}`; 
                console.log(calculateCell);
                const currentBox = document.getElementById(calculateCell);
                currentBox.classList.add("occupied");
            }
        }
    }
}