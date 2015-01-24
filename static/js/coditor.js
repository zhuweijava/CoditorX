var coditor = {
    conf: undefined,
    sessionId: undefined,
    sessionUsername: undefined,
    color: undefined,
    workspace: undefined,
    i18n: undefined,
    init: function () {
        // 点击隐藏弹出层
        $("body").bind("mouseup", function (event) {
            $(".frame").hide();
        });
        
        // 禁止鼠标右键菜单
        document.oncontextmenu = function () {
            return false;
        };
        
        this.conf = conf;
        this.sessionId = sessionId;
        this.sessionUsername = sessionUsername;
        this.color = color;
        this.i18n = i18n;
        this.workspace = workspace;
    },
    getClassBySuffix: function (suffix) {
        var iconSkin = "ico-file-other";
        switch (suffix) {
            case "html":
            case "htm":
                iconSkin = "ico-file-html";
                break;
            case "go":
                iconSkin = "ico-file-go";
                break;
            case "css":
                iconSkin = "ico-file-css";
                break;
            case "txt":
                iconSkin = "ico-file-text";
                break;
            case "sql":
                iconSkin = "ico-file-sql";
                break;
            case "properties":
                iconSkin = "ico-file-pro";
                break;
            case "md":
                iconSkin = "ico-file-md";
                break;
            case "js", "json":
                iconSkin = "ico-file-js";
                break;
            case "xml":
                iconSkin = "ico-file-xml";
                break;
            case "jpg":
            case "jpeg":
            case "bmp":
            case "gif":
            case "png":
            case "svg":
            case "ico":
                iconSkin = "ico-file-img";
                break;
        }

        return iconSkin;
    }
};

$(document).ready(function () {
    menu.init();
    side.init();
    coditor.init();
    session.init();
    notification.init();
    editor.init();
});