var nativeSocket;

try {
    nativeSocket = new WebSocket('ws://localhost:13921/auth')
} catch (e) {
    console.log("Error: ", err);
}

var scytherReady = false;

var handshake = function() {
    nativeSocket.send(JSON.stringify({
        type: 'handshake',
        value: 'hello scyther native'
    }));
};

var validateHandshake = function(message) {
    if (message.type === 'handshake' && message.value === 'hello scyther web') {
        // We have a Scyther server running, yay!
        scytherReady = true;
        console.log('Scyther is ready! \\m/')
    } else {
        console.log('Not a Scyther server :(');
    }
};

var populateDiscovered = function(message) {
    message = JSON.parse(message);
    var tagTotal = "";
    var deviceData = {};
    for (var bt of message) {
        deviceData[bt.ID] = bt.ManufacturerData;
    }
    for (var bt in deviceData) {
        tagTotal += "<tr><td>" + bt + "</td><td>" + deviceData[bt] + "</td></tr>";
    }
    document.getElementById("placeholder").innerHTML = "<table class=\"table table-striped table-hover\"><tr><th>Device ID</th><th>Manufacturer Data</th></tr></thead>" + tagTotal + "</table>";
};

var handleMessage = function(message) {
    message = JSON.parse(message.data);
    // console.log(message);
    switch(message.type) {
        case 'handshake':
            validateHandshake(message);
            break;
        case 'response':
            populateDiscovered(message.value);
            break;
        default:
            console.log('Unsupported message');
    }
}

nativeSocket.addEventListener('open', function() {
    // ready to start handshake
    handshake();
});

nativeSocket.addEventListener('error', function(err) {
    // oops! something went wrong
    scytherReady = false;
    console.log("Error: ", err);
});

nativeSocket.addEventListener('close', function() {
    // server closed the connection
    scytherReady = false;
    console.log("Connection closed")
});

nativeSocket.addEventListener('message', function(msg) {
    // received a message
    handleMessage(msg);
});

var btn = document.getElementById("discover");
btn.addEventListener('click', function() {
    nativeSocket.send(JSON.stringify({type: "privilege", value: ""}));
});
