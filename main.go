package main

import "fmt"

func main() {
	//pids, err := GetPortPID(9999)
	//if err != nil {
	//	panic(err)
	//}

	//pids := []int{79422}
	//
	//for _, pid := range pids {
	//	process, err := os.FindProcess(pid)
	//	if err != nil {
	//		panic(err)
	//	}
	//	err = process.Signal(syscall.Signal(0))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = process.Signal(syscall.SIGTERM)
	//	if err != nil {
	//		panic(err)
	//	}
	//	err = WaitProcessDone(pid)
	//	fmt.Printf("process.Wait, err:%s\n", err)
	//}

	c := make(chan int, 1)

	c <- 1
	close(c)

	for i := range c {
		fmt.Println(i)
	}
}

//func WaitProcessDone(pid int) error {
//	for {
//		process, err := os.FindProcess(pid)
//		if err != nil {
//			return err
//		}
//		// check if process done
//		err = process.Signal(syscall.Signal(0))
//		if err != nil {
//			if errors.Is(err, os.ErrProcessDone) {
//				fmt.Println("process done")
//			} else {
//				fmt.Printf("WaitProcessDone|err:%s\n", err)
//			}
//			return nil
//		}
//	}
//}
