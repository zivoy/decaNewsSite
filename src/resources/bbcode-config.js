/* 
 * Javascript BBCode Parser Config Options
 * @author Philip Nicolcev
 * @license MIT License
 */

var parserColors = ['gray', 'silver', 'white', 'yellow', 'orange', 'red', 'fuchsia', 'blue', 'green', 'black', '#cd38d9'];

var parserTags = {
    'b': {
        openTag: function () {
            return '<b>';
        },
        closeTag: function () {
            return '</b>';
        }
    },
    'code': {
        openTag: function () {
            return '<code>';
        },
        closeTag: function () {
            return '</code>';
        },
        noParse: true
    },
    'color': {
        openTag: function (params) {
            var colorCode = params.substr(1) || "inherit";
            BBCodeParser.regExpAllowedColors.lastIndex = 0;
            BBCodeParser.regExpValidHexColors.lastIndex = 0;
            if (!BBCodeParser.regExpAllowedColors.test(colorCode)) {
                if (!BBCodeParser.regExpValidHexColors.test(colorCode)) {
                    colorCode = "inherit";
                } else {
                    if (colorCode.substr(0, 1) !== "#") {
                        colorCode = "#" + colorCode;
                    }
                }
            }

            return '<span style="color:' + colorCode + '">';
        },
        closeTag: function () {
            return '</span>';
        }
    },
    'i': {
        openTag: function () {
            return '<i>';
        },
        closeTag: function () {
            return '</i>';
        }
    },
    'img': {
        openTag: function (params, content) {

            var myUrl = content;

            BBCodeParser.urlPattern.lastIndex = 0;
            if (!BBCodeParser.urlPattern.test(myUrl)) {
                myUrl = "";
            }

            return '<img class="bbCodeImage" alt="bbcode image preview" src="' + myUrl + '">';
        },
        closeTag: function () {
            return '';
        },
        content: function () {
            return '';
        }
    },
    'list': {
        openTag: function () {
            return '<ul>';
        },
        closeTag: function () {
            return '</ul>';
        },
        restrictChildrenTo: ["*", "li"]
    },
    'noparse': {
        openTag: function () {
            return '';
        },
        closeTag: function () {
            return '';
        },
        noParse: true
    },
    'quote': {
        openTag: function () {
            return '<q>';
        },
        closeTag: function () {
            return '</q>';
        }
    },
    's': {
        openTag: function () {
            return '<s>';
        },
        closeTag: function () {
            return '</s>';
        }
    },
    'size': {
        openTag: function (params) {
            var mySize = parseInt(params.substr(1), 10) || 0;
            if (mySize < 10 || mySize > 20) {
                mySize = 'inherit';
            } else {
                mySize = mySize + 'px';
            }
            return '<span style="font-size:' + mySize + '">';
        },
        closeTag: function () {
            return '</span>';
        }
    },
    'u': {
        openTag: function () {
            return '<span style="text-decoration:underline">';
        },
        closeTag: function () {
            return '</span>';
        }
    },
    'url': {
        openTag: function (params, content) {

            var myUrl;
            console.log(params.substr(1))

            if (!params) {
                myUrl = content.replace(/<.*?>/g, "");
            } else {
                myUrl = params;
            }

            BBCodeParser.urlPattern.lastIndex = 0;
            if (!BBCodeParser.urlPattern.test(myUrl)) {
                myUrl = "#";
            }

            return '<a href="' + myUrl + '">';
        },
        closeTag: function () {
            return '</a>';
        }
    }
};
