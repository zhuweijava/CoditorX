var session = {
    init: function () {
        this._initWS();
    },
    _initWS: function () {
        var sessionWS = new ReconnectingWebSocket(coditor.conf.Channel + '/session/ws?sid=' + coditor.sessionId);

        sessionWS.onopen = function () {
            console.log('[session onopen] connected');
        };

        sessionWS.onmessage = function (e) {
            console.log('[session onmessage]' + e.data);
        };
        sessionWS.onclose = function (e) {
            console.log('[session onclose] disconnected (' + e.code + ')');
        };
        sessionWS.onerror = function (e) {
            console.log('[session onerror] ' + JSON.parse(e));
        };
    }
};