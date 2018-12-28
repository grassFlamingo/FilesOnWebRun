/**
 * A File For FilesOnWeb
 * Created by aliy at November 23, 2018
 */

const JSON_RESPONSE_STSATE = {
    Ok: 0x01,
    Error: 0xf0
};

var WORK_DATA_ITERATOR = {
    current: 0,
    length: 0,
    listdata: null,
    clear: function(){
        this.current = 0;
        this.length = 0;
        this.listdata = null;
    },
    setList: function (workdir) {
        if (workdir == null) {
            return false;
        }
        if(workdir.State != JSON_RESPONSE_STSATE.Ok){
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
        for(; head < workdir.length; head++){
            if(!workdir[head].IsDir){
                break;
            }
        }
        var tail = head+1;
        for (; tail < workdir.length; tail++) {
            if (workdir[tail].IsDir) {
                var t = workdir[head];
                workdir[head] = workdir[tail];
                workdir[tail] = t;
                head += 1;
            }
        }
        for(; head < workdir.length; head++){
            if(workdir[head].IsImg){
                break;
            }
        }
        tail = head+1;
        for (; tail < workdir.length; tail++) {
            if (!workdir[tail].IsImg) {
                var t = workdir[head];
                workdir[head] = workdir[tail];
                workdir[tail] = t;
                head += 1;
            }
        }
        return true;
    },
    next: function () {
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
    },
    reset: function () {
        this.current = 0;
    }
};

function randomColor() {
    var raint = Math.floor(Math.random() * 0xffffff + 0x7f7f7f)
    // 0x7f7f7f
    raint = raint & 0xffffff;
    return `#${raint.toString(16)}`;
}