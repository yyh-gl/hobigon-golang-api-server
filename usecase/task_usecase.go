package usecase

import (
	"fmt"
	"os/exec"
)

type TaskUsecase struct {}

func (taskUsecase TaskUsecase) Notify() error {
	output, err := exec.Command("/Users/yyh-gl/.anyenv/envs/rbenv/shims/ruby", "~/workspaces/Ruby/today-tasks-notifier/main.rb").Output()
	if err != nil {
		fmt.Println("=====Error=====")
		fmt.Println(string(output))
		fmt.Println(err)
		fmt.Println("=====Error=====")
		return err
	}

	fmt.Println("=====Success=====")
	fmt.Println(string(output))
	fmt.Println("=====Success=====")

	return nil
}
