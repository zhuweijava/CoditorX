var editor = {
    channel: undefined,
    codemirror: undefined,
    currentFileName: undefined,
    init: function () {
        this._initWS();
        this._initCM();
    },
    _initCM: function () {
        // FIXME: hard coding
        //var docName = coditor.workspace + "/README.txt";
        var docName = "workspaces/admin/workspace/README.md";
        var doc = OpenDoc(docName);

        var $editor = $("#editor");
        $editor.val(doc.content);
        
        editor.currentFileName = docName;

        editor.codemirror = CodeMirror.fromTextArea($editor[0], {
            autofocus: true,
            lineNumbers: true,
            theme: "blackboard"
        });

        var mode = CodeMirror.findModeByFileName(docName);
        if (mode) {
            editor.codemirror.setOption("mode", mode.mode);
        }

        var request = newRequest();
        request.docName = docName;
        request.offset = 0;
        request.color = coditor.color;

        $.ajax({
            async: false,
            url: "/doc/setCursor",
            type: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    return false;
                }
            }
        });

        editor.codemirror.setSize('100%', $(".main").height() - $(".menu").height());
        editor.codemirror.on('changes', function (cm, changes) {
            if (changes && changes[0] && "setValue" === changes[0].origin) {
                return;
            }

            var request = newRequest();
            request.cmd = "commit";
            request.content = cm.getValue();
            request.docName = docName;
            request.user = coditor.sessionUsername;
            request.cursor = cm.getCursor();
            request.color = coditor.color;

            editor.channel.send(JSON.stringify(request));
        });
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