function Doc(name, content, version, debug) {
    this.name = name;
    this.content = content;
    this.version = version;
    this.dmp = new diff_match_patch();
    if (debug) {
        this.debug = true;
    } else {
        this.debug = false;
    }
    return this;
}

function OpenDoc(fileName) {
    var doc = undefined;

    var request = newRequest();
    request.fileName = fileName;

    $.ajax({
        async: false,
        url: "/doc/open",
        type: "POST",
        data: JSON.stringify(request),
        success: function (data) {
            if (data.succ) {
                var d = data.doc;
                doc = new Doc(fileName, d.content, d.version);
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            // TODO
        }
    });

    return doc;
}

// commit doc to server.
Doc.prototype.commit = function () {
    var doc = this;
    var result = undefined;
    var request = newRequest();
    var file = {
        name: this.name,
        version: this.version,
        content: this.content
    };
    request.file = file;

    $.ajax({
        async: false,
        url: "/doc/commit",
        type: "POST",
        data: JSON.stringify(request),
        success: function (data) {
            result = data;
            if (data.succ) {
                doc.version = data.output.version;
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            // TODO
        }
    });
    return result;
}

// pull from server.
Doc.prototype.pull = function () {
    var doc = this;
    var result = undefined;
    var request = newRequest();
    var file = {
        name: doc.name,
        version: doc.version
    }
    request.file = file;

    $.ajax({
        async: false,
        url: "/doc/fetch",
        type: "POST",
        data: JSON.stringify(request),
        success: function (data) {
            result = data;
            if (data.succ) {
                var length = data.patchss.length
                for (var i = 0; i < length; i++) {
                    var patchsStr = data.patchss[i];
                    var patches = doc.dmp.patch_fromText(patchsStr);
                    var outputs = doc.dmp.patch_apply(patches, doc.content);
                    var result = outputs[1];
                    console.log(patches);
                    for (var i = 0; i < result.length; i++) {
                        if (!result[i]) {
                            console.log("result:" + result);
                        }
                    }
                    doc.content = outputs[0];
                }
                doc.version = data.version;
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            // TODO
        }
    });
    return result;
}

// set content.
Doc.prototype.setContent = function (content) {
    this.content = content;
}