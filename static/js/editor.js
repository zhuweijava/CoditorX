var editor = {
    channel: undefined,
    codemirror: undefined,
    init: function () {
        this._initWS();
        this._initCM();
    },
    _initCM: function () {
        // FIXME: hard coding
        //var docName = coditor.workspace + "/README.txt";
        var docName = "workspaces/admin/workspace/README.txt";
        var doc = OpenDoc(docName);

        var textArea = document.getElementById("editor");
        textArea.value = doc.content;

        editor.codemirror = CodeMirror.fromTextArea(textArea, {
            autofocus: true
        });
    },
    _initWS: function () {
        this.channel = new ReconnectingWebSocket(coditor.conf.Channel + '/editor/ws?sid=' + coditor.sessionId);

        this.channel.onopen = function () {
            console.log('[editor onopen] connected');
        };

        this.channel.onmessage = function (e) {
            console.log('[editor onmessage]' + e.data);
        };

        this.channel.onclose = function (e) {
            console.log('[editor onclose] disconnected (' + e.code + ')');
        };

        this.channel.onerror = function (e) {
            console.log('[editor onerror] ' + JSON.parse(e));
        };
    }
};