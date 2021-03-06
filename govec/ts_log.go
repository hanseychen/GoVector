package govec

import (
	"bytes"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"os"
	"strconv"
	"time"
)

type GoTSLog struct {
	GoLog
}

func InitGoVectorTimeStamp(processid string, logfilename string) *GoTSLog {
	gv := InitGoVector(processid, logfilename)
	gvts := &GoTSLog{*gv}
	gvts.SetLogFunc(gvts.LogThis)
	// Need to do a reinitialization to have a timestamped initializaiton event
	gvts.reinitialize()
	return gvts
}

func (gv *GoTSLog) truncate() bool {
	complete := true
	file, err := os.OpenFile(gv.logfile, os.O_TRUNC | os.O_WRONLY, 0600)
	if err != nil {
		complete = false
	}
	defer file.Close()
	return complete
}

func (gv *GoTSLog) reinitialize() {
	complete := gv.truncate()
	if !complete {
		gv.logger.Println("Something went wrong during re-initialization")
	}
	vc, err := vclock.FromBytes(gv.currentVC)
	if err != nil {
		gv.logger.Println(err.Error())
	}
	ok := gv.LogThis( "Initialization Complete", gv.pid , vc.ReturnVCString())
	if !ok {
		gv.logger.Println("Something went wrong during re-initialization")
	}
}

func (gv *GoTSLog) LogThis(Message string, ProcessID string, VCString string) bool {
	complete := true
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	buffer.WriteString(" ")
	buffer.WriteString(ProcessID)
	buffer.WriteString(" ")
	buffer.WriteString(VCString)
	buffer.WriteString("\n")
	buffer.WriteString(Message)
	buffer.WriteString("\n")
	output := buffer.String()

	gv.output += output
	if !gv.buffered {
		complete = gv.Flush()
	}

	if gv.printonscreen == true {
		gv.logger.Println(output)
	}
	return complete
}
