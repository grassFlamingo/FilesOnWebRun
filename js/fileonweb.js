/**
 * A File For FilesOnWeb
 * Created by aliy at November 23, 2018
 */

class FilesKeeper {
    constructor(dirList) {
        this.dirFolders = [];
        this.dirFiles = [];
        this.setup(dirList);
    }

    clear() {
        this.current = 0;
        this.length = 0;
        this.listdata = null;
    }

    isAllImages(){
        if(this.dirFolders.length > 0){
            return false;
        }
        for(var i in this.dirFiles){
            if(FileTypeUtil.getFileType(this.dirFiles[i].Name) != "image"){
                return false;
            }
        }
        return true;
    }

    setup(dirList) {
        if (dirList == undefined || dirList == null) {
            return;
        }
        this.length = dirList.length;
        this.listdata = dirList;
        this.current = 0;
        if (this.length < 1) {
            return;
        }
        for (var i in dirList) {
            var titem = dirList[i];
            titem.ModeTime = Date.parse(titem.ModeTime)
            if (titem.IsDir) {
                this.dirFolders.push(titem);
            } else {
                this.dirFiles.push(titem);
            }
        }
    }

    sortAtoZ() {
        var sortmeth = function (a, b) {
            return a.Name > b.Name;
        }
        this.dirFolders.sort(sortmeth);
        this.dirFiles.sort(sortmeth);
    }


    sortZtoA() {
        var sortmeth = function (a, b) {
            return a.Name < b.Name;
        }
        this.dirFolders.sort(sortmeth);
        this.dirFiles.sort(sortmeth);
    }

    sortFirstModified() {
        var sortmeth = function (a, b) {
            return a.ModeTime > b.ModeTime;
        }
        this.dirFolders.sort(sortmeth);
        this.dirFiles.sort(sortmeth);
    }

    sortLastModified() {
        var sortmeth = function (a, b) {
            return a.ModeTime < b.ModeTime;
        }
        this.dirFolders.sort(sortmeth);
        this.dirFiles.sort(sortmeth);
    }

    getIterator(pattern, flags) {
        if(flags == undefined){
            flags = "gi"; // global ignore case
        }
        return new FilesKeeperIterator(this.dirFolders, this.dirFiles, pattern, flags);
    }
}

class FilesKeeperIterator {
    constructor(folderList, fileList, pattern, flags) {
        this.thelist = folderList.concat(fileList);
        this.index = 0;
        this.length = this.thelist.length;
        if (pattern == undefined) {
            this.filter = function (item) {
                return true;
            }
        } else {
            this.pattern = new RegExp(pattern, flags)
            this.filter = function (item) {
                return this.pattern.test(item.Name)
            }
        }
    }

    hasNext() {
        return this.index < this.length;
    }

    next() {
        var i = this.index;
        for(; i < this.length; i++){
            if(this.filter(this.thelist[i])){
                this.index = i+1;
                return {
                    index: this.index,
                    data: this.thelist[i],
                }
            }
        }
        return null;
    }
}

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
        "folder": "folder",
        "unknow": "unknow",
    },
    folderIcon: "/img/folder.svg",
    getSuffix: function (filename) {
        var suffix = "unknow";
        if (filename[filename.length - 1] == "/") {
            return "folder";
        }
        for (var i = filename.length - 1; i >= 0 && filename[i] != '/'; i--) {
            if (filename[i] == '.') {
                suffix = filename.substring(i+1);
                suffix = suffix.toLowerCase();
                break;
            }
        }
        return suffix;
    },
    getFileType: function (filename) {
        var suffix = this.getSuffix(filename);
        var type = this.fileTypeMap[suffix];
        if(type == undefined){
            type = "unknow";
        }
        return type
    },
    getFileIcon: function (filename) {
        return `/img/${this.getFileType(filename)}.svg`;
    },
}

const OSUtil = {
    unpackServerResponse: function (data, errorfun) {
        if (data.State == JSON_RESPONSE_STSATE.Ok) {
            return data.Data;
        }
        errorfun = errorfun != undefined ? errorfun : function (msg) {
            console.log(`Unpack Server Response Fail with message ${msg}`)
        }
        errorfun(JSON_RESPONSE_STSATE.toString(data.State));
    }
}