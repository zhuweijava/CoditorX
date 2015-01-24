var editor = {
    channel: undefined,
    quill: undefined,
    init: function () {
        this._initWS();
        this._initQuill();
    },
    _initQuill: function () {
        editor.quill = new Quill('#quillEditor', {
            modules: {
                'toolbar': {container: '#quillToolbar'},
                'link-tooltip': true,
                'image-tooltip': true,
                'multi-cursor': true
            },
            theme: 'snow'
        });

        editor.quill.addModule('authorship', {
            authorId: coditor.sessionUsername,
            color: coditor.color,
            enabled: true
        });

        // FIXME: hard coding
        //var docName = coditor.workspace + "/README.txt";
        var docName = "workspaces/admin/workspace/README.txt";
        var doc = OpenDoc(docName);
        editor.quill.setText(doc.content);
        
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

        editor.quill.on('text-change', function (delta, source) {
            if (source === 'api') {
                console.log("An API call triggered this change.");
            } else if (source === 'user') {
//                console.log("A user action triggered this change.");
//                console.log(delta);

                var request = newRequest();
                request.cmd = "text-change";
                request.delta = delta;
                request.docName = docName;
                request.user = coditor.sessionUsername;
                request.color = coditor.color;

                editor.channel.send(JSON.stringify(request));
            }
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
            if (data.cmd === "text-change") {
                editor.quill.updateContents(data.output);
                
                var authorship = editor.quill.getModule('authorship')
                authorship.addAuthor(data.user, data.color);
                
                var cursorManager = editor.quill.getModule('multi-cursor');
                cursorManager.setCursor(data.user, data.offset, data.user, data.color);
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