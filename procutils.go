package procutils

import (
   "errors"
//   "time"
   "os/exec"
)

const (
   MsgRunning = iota
)

type procMessage struct {
   MsgType  int
   Tag      int
}

const (
   cmdStop = iota
)

type procCmd struct {
   Cmd   int
}

type processManager struct {

   MsgQueue       chan procMessage
   procChannels   map[int](chan procCmd)

}

func NewProcessManager() processManager {

   pm := processManager{}
   pm.MsgQueue = make(chan procMessage)
   pm.procChannels = make(map[int](chan procCmd))

   return pm
}

func processRunner(tag int, command string, msgQueue chan<- procMessage, cmdQueue <-chan procCmd) {

   cmd := exec.Command(command)

   err := cmd.Start()

   if err != nil {
      msg := procMessage{}
      msg.Tag = tag
      msg.MsgType = MsgRunning
      msgQueue <- msg
   } else {
      // TODO report error
      return
   }

}

func (pm *processManager) RunProcess(tag int, cmd string, environ []string, ) error {

   if _, ok := pm.procChannels[tag]; ok {
      return errors.New("tag already in use")
   }

   procCmdQueue := make(chan procCmd)
   pm.procChannels[tag] = procCmdQueue

   go processRunner(tag, cmd, pm.MsgQueue, procCmdQueue)

   return nil
}

func (pm *processManager) StopProcess(tag int) error {

   if _, ok := pm.procChannels[tag]; !ok {
      return errors.New("unknown tag")
   }

   cmd := procCmd{}
   cmd.Cmd = cmdStop
   pm.procChannels[tag] <- cmd

   return nil
}

