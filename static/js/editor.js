var editor = {
    channel: undefined,
    codemirror: undefined,
    currentFileName: undefined,
    init: function () {
        this._initWS();
    },
    _initWS: function () {
        this.channel = new ReconnectingWebSocket(coditor.conf.Channel + '/editor/ws?sid=' + coditor.sessionId);

        this.channel.onopen = function () {
            console.log('[editor onopen] connected');
        };

        this.channel.onmessage = function (e) {
            console.log('[editor onmessage]' + e.data);

            var data = JSON.parse(e.data);

            if ("changes" === data.cmd) {
                var cursor = editor.codemirror.getCursor();
                editor.codemirror.setValue(data.content);
                editor.codemirror.setCursor(cursor);
            }
        };

        this.channel.onclose = function (e) {
            console.log('[editor onclose] disconnected (' + e.code + ')');
        };

        this.channel.onerror = function (e) {
            console.log('[editor onerror] ' + JSON.parse(e));
        };
    }
};