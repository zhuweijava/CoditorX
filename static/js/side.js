/*
 * Copyright (c) 2015, B3log
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

var side = {
    init: function () {
        new Tabs({
            id: ".side"
        });

        this._initFileList();
        this._initShareList();
        this._initDialog();
    },
    _mockData: function (count) {
        var data = [],
                types = ['html', 'htm', 'sql', 'properties', 'md', 'json', 'xml',
                    'jpg', 'jpeg', 'bmp', 'gif', 'png', 'svg', 'ico', 'go', 'css', 'js', 'txt', ''];

        for (var i = 0, max = count; i < max; i++) {
            var type = types[Math.ceil(Math.random() * 18)];
            data.push({
                'isShare': false,
                'id': i,
                'name': 'fileName' + i + '.' + type,
                'type': type
            });
        }

        return data;
    },
    _initFileList: function () {
        var request = newRequest();

        $.ajax({
            url: "/files",
            type: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    return false;
                }

                var filesHTML = '<ul class="list">',
                        $files = $("#files");

                for (var i = 0, max = data.files.length; i < max; i++) {
                    filesHTML += '<li data-share="' + data.files[i].isShare
                            + '"><span class="ico-file ' + coditor.getClassBySuffix(data.files[i].type)
                            + '"></span> ' + data.files[i].name + '<li>';
                }
                $files.html(filesHTML + '<ul>');

                $files.find("li").mouseup(function (event) {
                    event.stopPropagation();

                    if (event.button === 0) { // 左键
                        $files.next().hide();
                        return;
                    }

                    // event.button === 2 右键
                    $files.next().show().css({
                        "left": "38px",
                        "top": (event.target.offsetTop - $files.parent().scrollTop() + 22) + "px"
                    });
                    if ($(this).data('share')) {
                        $files.next().find('.share').hide();
                        $files.next().find('.unshare').show();
                    } else {
                        $files.next().find('.share').show();
                        $files.next().find('.unshare').hide();
                    }

                    $files.find("li").removeClass("current");
                    $(this).addClass('current');

                    return;
                });
            }
        });
    },
    _initShareList: function () {
        var request = newRequest();

        $.ajax({
            url: "/shares",
            type: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    return false;
                }

                var filesHTML = '<ul class="list">',
                        $shareFiles = $("#shareFiles");

                for (var i = 0, max = data.shares.length; i < max; i++) {
                    var shareFile = data.shares[i];
                    var index = shareFile.docName.lastIndexOf("\\.");
                    var fileType = "";
                    if (index > 0) {
                        fileType = shareFile.docName.sub(index + 1);
                    }
                    filesHTML += '<li><span class="ico-file ' + coditor.getClassBySuffix(fileType)
                            + '"></span> ' + '/' + shareFile.owner + '/' + shareFile.docName + '<li>';
                }
                $shareFiles.html(filesHTML + '<ul>');

                $shareFiles.find("li").mouseup(function (event) {
                    event.stopPropagation();

                    if (event.button === 0) { // 左键
                        $shareFiles.next().hide();
                        return;
                    }

                    // event.button === 2 右键
                    $shareFiles.next().show().css({
                        "left": "38px",
                        "top": (event.target.offsetTop - $shareFiles.parent().scrollTop() + 22) + "px"
                    });
                    return;
                });
            }
        });
    },
    _initDialog: function () {
        $("#dialogAlert").dialog({
            "modal": true,
            "height": 36,
            "width": 260,
            "title": '提示',
            "hiddenOk": true,
            "cancelText": '确认',
            "afterOpen": function (msg) {
                $("#dialogAlert").html(msg);
            }
        });

        $(".dialog-prompt > input").keyup(function (event) {
            var $okBtn = $(this).closest(".dialog-main").find(".dialog-footer > button:eq(0)");
            if (event.which === 13 && !$okBtn.prop("disabled")) {
                $okBtn.click();
            }

            if ($.trim($(this).val()) === "") {
                $okBtn.prop("disabled", true);
            } else {
                $okBtn.prop("disabled", false);
            }
        });

        $("#dialogUnshareConfirm").dialog({
            "modal": true,
            "height": 36,
            "width": 260,
            "title": 'Unshare',
            "okText": 'Unshare',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogUnshareConfirm > b").html('"FileName"');
            },
            "ok": function () {
                var request = newWideRequest();
                request.path = wide.curNode.path;

                $.ajax({
                    type: 'POST',
                    url: config.context + '/file/remove',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (!data.succ) {
                            $("#dialogRemoveConfirm").dialog("close");
                            bottomGroup.tabs.setCurrent("notification");
                            windows.flowBottom();
                            $(".bottom-window-group .notification").focus();
                            return false;
                        }

                        $("#dialogRemoveConfirm").dialog("close");
                        tree.fileTree.removeNode(wide.curNode);

                        if (!tree.isDir()) {
                            // 是文件的话，查看 editor 中是否被打开，如打开则移除
                            for (var i = 0, ii = editors.data.length; i < ii; i++) {
                                if (editors.data[i].id === wide.curNode.tId) {
                                    $(".edit-panel .tabs > div[data-index=" + wide.curNode.tId + "]").find(".ico-close").click();
                                    break;
                                }
                            }
                        } else {
                            for (var i = 0, ii = editors.data.length; i < ii; i++) {
                                if (tree.isParents(editors.data[i].id, wide.curNode.tId)) {
                                    $(".edit-panel .tabs > div[data-index=" + editors.data[i].id + "]").find(".ico-close").click();
                                    i--;
                                    ii--;
                                }
                            }
                        }
                    }
                });
            }
        });

        $("#dialogRemoveConfirm").dialog({
            "modal": true,
            "height": 36,
            "width": 260,
            "title": 'Delete',
            "okText": 'Delete',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogRemoveConfirm > b").html('"FileName"');
            },
            "ok": function () {
                var request = newWideRequest();
                request.path = wide.curNode.path;

                $.ajax({
                    type: 'POST',
                    url: config.context + '/file/remove',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (!data.succ) {
                            $("#dialogRemoveConfirm").dialog("close");
                            bottomGroup.tabs.setCurrent("notification");
                            windows.flowBottom();
                            $(".bottom-window-group .notification").focus();
                            return false;
                        }

                        $("#dialogRemoveConfirm").dialog("close");
                        tree.fileTree.removeNode(wide.curNode);

                        if (!tree.isDir()) {
                            // 是文件的话，查看 editor 中是否被打开，如打开则移除
                            for (var i = 0, ii = editors.data.length; i < ii; i++) {
                                if (editors.data[i].id === wide.curNode.tId) {
                                    $(".edit-panel .tabs > div[data-index=" + wide.curNode.tId + "]").find(".ico-close").click();
                                    break;
                                }
                            }
                        } else {
                            for (var i = 0, ii = editors.data.length; i < ii; i++) {
                                if (tree.isParents(editors.data[i].id, wide.curNode.tId)) {
                                    $(".edit-panel .tabs > div[data-index=" + editors.data[i].id + "]").find(".ico-close").click();
                                    i--;
                                    ii--;
                                }
                            }
                        }
                    }
                });
            }
        });

        $("#dialogShare").load('/share', function () {
            $("#dialogShare").dialog({
                "modal": true,
                "height": 190,
                "width": 600,
                "title": 'Share',
                "afterOpen": function () {
                    $("#dialogShare .fileName").val(coditor.workspace + "/" + $.trim($("#files li.current").text()));
                    $("#dialogShare").find('input[type=checkbox]').prop('checked', false);
                    $("#dialogShare").find('.viewers').show();
                },
                "ok": function(){
                    var fileName = $("#dialogShare .fileName").val();
                    var editors = $("#dialogShare .editors").val();
                    var isPublic = 0;
                    if($("#dialogShare .isPublic").attr("checked")==true){
                        isPublic = 1;
                    };
                    var viewers = $("#dialogShare .viewers").val();
                    var request = newRequest();
                    request["fileName"] = fileName;
                    request["editors"] = editors;
                    request["isPublic"] = isPublic;
                    request["viewers"] = viewers;
                    $.ajax({
                        type: 'POST',
                        url: '/share',
                        data: JSON.stringify(request),
                        async: false,
                        dataType: "json",
                        success: function (data) {

                        }
                    })
                    return true
                }
            });

            $("#dialogShare").find('input[type=checkbox]').click(function () {
                if ($(this).prop('checked')) {
                    $("#dialogShare").find('.viewers').hide();
                } else {
                    $("#dialogShare").find('viewers').show();
                }
            });
        });

        $("#dialogRenamePrompt").dialog({
            "modal": true,
            "height": 52,
            "width": 260,
            "title": 'Rename',
            "okText": 'Rename',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogRenamePrompt").closest(".dialog-main").find(".dialog-footer > button:eq(0)").prop("disabled", true);
                $("#dialogRenamePrompt > input").val('FileName').select().focus();
            },
            "ok": function () {
                var name = $("#dialogRenamePrompt > input").val(),
                        request = newWideRequest();

                request.oldPath = wide.curNode.path;

                request.newPath = wide.curNode.path.substring(0,
                        wide.curNode.path.lastIndexOf(config.pathSeparator))
                        + config.pathSeparator + name;

                $.ajax({
                    type: 'POST',
                    url: config.context + '/file/rename',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (!data.succ) {
                            $("#dialogRenamePrompt").dialog("close");
                            bottomGroup.tabs.setCurrent("notification");
                            windows.flowBottom();
                            $(".bottom-window-group .notification").focus();
                            return false;
                        }

                        $("#dialogRenamePrompt").dialog("close");

                        // update tree node
                        var suffixIndex = name.lastIndexOf('.'),
                                iconSkin = wide.getClassBySuffix(name.substr(suffixIndex + 1));
                        wide.curNode.name = name;
                        wide.curNode.title = request.newPath;
                        wide.curNode.path = request.newPath;
                        wide.curNode.iconSkin = iconSkin;
                        tree.fileTree.updateNode(wide.curNode);

                        // update open editor tab name
                        for (var i = 0, ii = editors.data.length; i < ii; i++) {
                            if (wide.curNode.tId === editors.data[i].id) {
                                var info = CodeMirror.findModeByExtension(name.substr(suffixIndex + 1));
                                if (info) {
                                    editors.data[i].editor.setOption("mode", info.mime);
                                }

                                var $currentSpan = $(".edit-panel .tabs > div[data-index=" + wide.curNode.tId + "] > span:eq(0)");
                                $currentSpan.attr("title", request.newPath);
                                $currentSpan.html('<span class="' + iconSkin + 'ico"></span>' + wide.curNode.name);
                                break;
                            }
                        }
                    }
                });
            }
        });

        $("#dialogNewFilePrompt").dialog({
            "modal": true,
            "height": 52,
            "width": 260,
            "title": 'Create File',
            "okText": 'Create',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogNewFilePrompt > input").val('').focus();
                $("#dialogNewFilePrompt").closest(".dialog-main").find(".dialog-footer > button:eq(0)").prop("disabled", true);
            },
            "ok": function () {
                var request = newRequest(),
                        name = $("#dialogNewFilePrompt > input").val();

                request.name = name;

                var isOk = false;
                $.ajax({
                    async: false,
                    type: 'POST',
                    url: '/file/new',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (data.succ) {
                            side._initFileList();
                            isOk = true;
                        }
                    }
                });
                return isOk;
            }
        });
    },
    new : function () {
        $("#dialogNewFilePrompt").dialog("open");
    },
    remove: function () {
        $("#dialogRemoveConfirm").dialog("open");
    },
    share: function () {
        $("#dialogShare").dialog('open');
    },
    unshare: function () {
        $("#dialogUnshareConfirm").dialog("open");
    },
    rename: function () {
        $("#dialogRenamePrompt").dialog('open');
    },
    open: function () {
        OpenDoc(coditor.workspace + $("#files li.current").text());
    }
};