const wsURI = "ws://localhost:3550/";

const totalTicketInput = document.getElementById("total-ticket")

function generateURL(totalTicketRequests) {
    return `http://localhost:3500/book-ticket?tickets=${totalTicketRequests}`;
}

async function handleSubmitTickets() {
    await fetch(generateURL(totalTicketInput.value), {
        method: "POST",
        mode: "no-cors",
        credentials: "same-origin",
    }).then(()=>{
        showSnackbar(`${totalTicketInput.value} Ticket Booked!!`);
    });
}

// seats status web socket
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
    let ticketCounter = 0;

    // map event data to GUI
    for(let col = 0; col < eventData[0].length; col++){
        for (let index = 0; index < eventData.length; index++) {
            if(eventData[index][col] == "x"){
                const calculateCell =`${String.fromCharCode(65 + col)}${index+1}`; 
                const currentBox = document.getElementById(calculateCell);
                currentBox.classList.add("occupied");
                ticketCounter++;
            }
        }
    }

    // Update Remaining & Occupied Seats
    const remainingTicketsElement = document.getElementById("r_count");
    const occupiedTicketsElement = document.getElementById("o_count");
    let occupiedTickets = occupiedTicketsElement.innerHTML;
    occupiedTicketsElement.innerHTML = ticketCounter;
    remainingTicketsElement.innerHTML = 120 - occupiedTickets;

}

window.addEventListener("unload", function (){
    statusWS.close();
});

function showSnackbar(text) {
    const snackbar = document.getElementById("snackbar");
    const snackbarText = document.getElementById("snackbar-text");
    snackbarText.textContent = text;
    snackbar.classList.remove("hidden");

    // Automatically hide after a timeout
    setTimeout(() => {
        snackbar.classList.add("hidden");
    }, 3000); 
}
