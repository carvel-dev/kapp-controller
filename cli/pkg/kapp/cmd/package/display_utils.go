// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"fmt"
// 	"time"

// 	"github.com/briandowns/spinner"

// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
// )

// // DisplayProgress creates an spinner instance; keeps receiving the progress messages in the channel and displays those using the spinner until an error occurs
// func DisplayProgress(initialMsg string, pp *tkgpackagedatamodel.PackageProgress) error {
// 	var (
// 		currMsg string
// 		s       *spinner.Spinner
// 		err     error
// 	)

// 	newSpinner := func() (*spinner.Spinner, error) {
// 		s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
// 		if err := s.Color("bgBlack", "bold", "fgWhite"); err != nil {
// 			return nil, err
// 		}
// 		return s, nil
// 	}
// 	if s, err = newSpinner(); err != nil {
// 		return err
// 	}

// 	writeProgress := func(s *spinner.Spinner, msg string) error {
// 		s.Stop()
// 		if s, err = newSpinner(); err != nil {
// 			return err
// 		}
// 		log.Infof("\n")
// 		s.Suffix = fmt.Sprintf(" %s", msg)
// 		s.Start()
// 		return nil
// 	}

// 	s.Suffix = fmt.Sprintf(" %s", initialMsg)
// 	s.Start()

// 	defer func() {
// 		s.Stop()
// 	}()
// 	for {
// 		select {
// 		case err := <-pp.Err:
// 			s.FinalMSG = "\n\n"
// 			return err
// 		case msg := <-pp.ProgressMsg:
// 			if msg != currMsg {
// 				if err := writeProgress(s, msg); err != nil {
// 					return err
// 				}
// 				currMsg = msg
// 			}
// 		case <-pp.Done:
// 			for msg := range pp.ProgressMsg {
// 				if msg == currMsg {
// 					continue
// 				}
// 				if err := writeProgress(s, msg); err != nil {
// 					return err
// 				}
// 				currMsg = msg
// 			}
// 			return nil
// 		}
// 	}
// }
