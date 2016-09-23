var nativeSocket;

try {
    nativeSocket = new WebSocket('ws://localhost:13921/auth')
} catch (e) {
    console.log("Error: ", err);
}

var scytherReady = false;

var handshake = function() {
    nativeSocket.send({
        type: 'handshake',
        value: 'hello scyther native'
    });
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

var handleMessage = function(message) {
    switch(message.type) {
        case 'handshake':
            validateHandshake(message);
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
});

nativeSocket.addEventListener('message', function() {
    // received a message
});
