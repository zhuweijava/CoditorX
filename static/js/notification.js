var notification = {
    init: function () {      
        this._initWS();
    },
    _initWS: function () {
        var notificationWS = new ReconnectingWebSocket(coditor.conf.Channel + '/notification/ws?sid=' + coditor.sessionId);

        notificationWS.onopen = function () {
            console.log('[notification onopen] connected');
        };

        notificationWS.onmessage = function (e) {
            console.log('[notification onmessage]' + e.data);
        };

        notificationWS.onclose = function (e) {
            console.log('[notification onclose] disconnected (' + e.code + ')');
        };

        notificationWS.onerror = function (e) {
            console.log('[notification onerror] ' + JSON.parse(e));
        };
    }
};