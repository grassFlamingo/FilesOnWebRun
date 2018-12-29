/**
 * A File For FilesOnWeb
 * Created by aliy at November 23, 2018
 */

class FilesIterator {
    constructor(workdir) {
        return this.setup(workdir);
    }

    clear() {
        this.current = 0;
        this.length = 0;
        this.listdata = null;
    }


    setup(workdir) {
        if (workdir == null || workdir.State != JSON_RESPONSE_STSATE.Ok) {
            this.clear()
            return false;
        }
        workdir = workdir.Data;
        this.length = workdir.length;
        this.listdata = workdir;
        this.current = 0;
        if (this.length < 1) {
            return;
        }
        var head = 0;
        for (; head < workdir.length; head++) {
            if (!workdir[head].IsDir) {
                break;
            }
        }
        var tail = head + 1;
        for (; tail < workdir.length; tail++) {
            if (workdir[tail].IsDir) {
                var t = workdir[head];
                workdir[head] = workdir[tail];
                workdir[tail] = t;
                head += 1;
            }
        }
        for (; head < workdir.length; head++) {
            if (workdir[head].IsImg) {
                break;
            }
        }
        tail = head + 1;
        for (; tail < workdir.length; tail++) {
            if (!workdir[tail].IsImg) {
                var t = workdir[head];
                workdir[head] = workdir[tail];
                workdir[tail] = t;
                head += 1;
            }
        }
        return true;
    }

    next() {
        if (this.current < this.length) {
            var dat = this.listdata[this.current];
            var ind = this.current++;
            return {
                index: ind,
                data: dat,
            };
        } else {
            return {
                index: -1,
                data: null,
            }
        }
    }

    getSubIterator(pattern) {
        var out = FilesIterator(null);
        out.listdata = new Array();
        var oldList = this.listdata;
        for (var i in oldList) {
            if (pattern.test(oldList[i])) {
                out.listdata.push(oldList[i]);
            }
        }
        out.length = out.listdata.length;
        out.current = 0;
    }

    reset() {
        this.current = 0;
    }
}

var WORK_DATA_ITERATOR = new FilesIterator();

function randomColor() {
    var raint = Math.floor(Math.random() * 0xffffff + 0x7f7f7f)
    // 0x7f7f7f
    raint = raint & 0xffffff;
    return `#${raint.toString(16)}`;
}

const JSON_RESPONSE_STSATE = {
    Ok: 0x0100,
    Error: 0x0200,
    BAD_CMD: 0x0201,
    BAD_DIR: 0x0202,
    BAD_OPEN: 0x0203,
    toString: function (state) {
        var ans = "Status OK";
        switch (state) {
            case this.Ok:
                ans = "Status OK";
                break;
            case this.Error:
                ans = "State Error";
                break;
            case this.BAD_CMD:
                ans = "Bad Command";
                break;
            case this.BAD_DIR:
                ans = "Bad Directory";
                break;
            case this.BAD_OPEN:
                ans = "Can't Open";
                break;
        }
        return ans;
    }
};

const FileTypeUtil = {
    fileTypeMap: {
        "zip": "achive",
        "rar": "achive",
        "7z": "achive",
        "gz": "achive",
        "pdf": "pdf",
        "txt": "txt",
        "md": "markdown",
        "markdown": "markdown",
        "doc": "word",
        "docx": "word",
        "xls": "excel",
        "xlsx": "excel",
        "ppt": "powerpont",
        "pptx": "powerpont",
        "jpg": "image",
        "png": "image",
        "gif": "image",
        "bmp": "image",
        "svg": "image",
        "jpeg": "image",
        "unknown": "unknown",
    },
    folderIcon: "/img/folder.svg",
    getSuffix: function(filename){
        var suffix = "unknown";
        for (var i = filename.length - 1; i >= 0; i++) {
            if (filename[i] == '.') {
                suffix = filename.substring(i);
                suffix = suffix.toLowerCase();
                break;
            }
        }
        return suffix;
    },
    getFileIcon: function (filename) {
        var suffix = this.getSuffix(filename)
        console.log(suffix);
        return `/img/${thid.fileTypeMap[suffix]}`;
    },
}

const OSUtil = {
    unpackServerResponse: function (data, errorfun) {
        if (data.State == JSON_RESPONSE_STSATE.Ok) {
            return data.Data;
        }
        errorfun = errorfun != null ? errorfun : function (msg) {
            console.log(`Unpack Server Response Fail with message ${msg}`)
        }
        errorfun(JSON_RESPONSE_STSATE.toString(data.State));
    },


}