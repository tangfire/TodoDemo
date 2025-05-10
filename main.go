package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
	"strconv"
	"time"
)

const TASKS_FILE_NAME = "tasks.json"

type todoTask struct {
	ID        int       `json:"id"`         // 任务ID
	TaskTitle string    `json:"task_title"` // 任务内容
	CreatedAt time.Time `json:"created_at"` // 创建时间
	IsDone    bool      `json:"is_done"`    // 是否完成
	IsDeleted bool      `json:"is_deleted"` // 是否删除
}

type todoList []todoTask

var addTaskTitle string // 添加任务的内容
var doneTaskID int      // 完成任务的ID
var deleteTaskID int    // 删除任务的ID
var islistTask bool     // 是否查看未完成任务
var islistAllTask bool  // 是否查看所有任务

func init() {
	flag.StringVar(&addTaskTitle, "add", "", "添加代办任务")
	flag.IntVar(&doneTaskID, "done", 0, "完成待办任务")
	flag.IntVar(&deleteTaskID, "del", 0, "删除待办任务")
	flag.BoolVar(&islistTask, "list", false, "查看未完成任务")
	flag.BoolVar(&islistAllTask, "all", false, "查看所有任务")
	flag.Parse()

}

// 加载json文件数据
func loadData() (todoList, error) {
	var list todoList
	// 以只读的方式打开文件 如果不存在则创建 权限0777
	dataFile, err := os.OpenFile(TASKS_FILE_NAME, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return list, fmt.Errorf("打开文件失败:%s", err)
	}

	defer dataFile.Close()
	// 读取文件数据
	data, err := io.ReadAll(dataFile)
	if err != nil {
		return list, fmt.Errorf("读取文件失败:%s", err)
	}
	// 解析json数据
	if len(data) != 0 {
		err = json.Unmarshal(data, &list)
		if err != nil {
			return list, fmt.Errorf("解析json数据失败:%s", err)
		}
	}
	return list, nil

}

// 保存数据到json文件
func saveData(list todoList) error {
	// 以只读的方式打开文件 打开时清空文件 如果不存在则创建 权限0777
	dataFile, err := os.OpenFile(TASKS_FILE_NAME, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("打开文件失败:%s", err)
	}
	defer dataFile.Close()

	// 将任务列表转换为json数据
	data, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("转换json数据失败:%s", err)
	}
	// 将json数据写入文件
	_, err = dataFile.Write(data)
	if err != nil {
		return fmt.Errorf("写入文件失败:%s", err)
	}

	return nil
}

// 添加任务
func addTask(list todoList) {
	task := todoTask{
		ID:        len(list) + 1,
		TaskTitle: addTaskTitle,
		CreatedAt: time.Now(),
		IsDone:    false,
	}
	list = append(list, task)
	err := saveData(list)
	if err != nil {
		fmt.Println("error: %s\n", err)
	} else {
		fmt.Println("success 成功添加任务")
	}

}

// 完成任务
func doneTask(list todoList) {
	task := &list[doneTaskID-1]
	task.IsDone = true
	err := saveData(list)
	if err != nil {
		fmt.Println("error: %s\n", err)
	} else {
		fmt.Println("success 成功删除任务")
	}

}

// 删除任务
func deleteTask(list todoList) {
	task := &list[deleteTaskID-1]
	task.IsDeleted = true
	err := saveData(list)
	if err != nil {
		fmt.Println("error: %s\n", err)
	} else {
		fmt.Println("success 成功删除任务")
	}
}

// 将任务列表转换为字符串切片
func listToString(list todoList) [][]string {
	var result [][]string
	for _, task := range list {
		result = append(result, []string{
			strconv.Itoa(task.ID),
			task.TaskTitle,
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			strconv.FormatBool(task.IsDone),
			strconv.FormatBool(task.IsDeleted),
		})
	}
	return result
}

// 打印任务列表
func printTask(isAll bool, list todoList) {
	data := listToString(list)
	table := tablewriter.NewWriter(os.Stdout)
	if isAll {
		table.SetHeader([]string{"ID", "待办任务", "创建时间", "是否完成"})
		for _, v := range data {
			if v[4] == "false" {
				table.Append(v[:4])
			}
		}
	} else {
		table.SetHeader([]string{"ID", "待办任务", "创建时间"})
		for _, v := range data {
			if v[3] == "false" && v[4] == "false" {
				table.Append(v[:3])
			}
		}
	}
	table.Render()
}

func main() {
	list, err := loadData()
	if err != nil {
		fmt.Println("error:%s", err)
		return
	}
	switch {
	case islistTask:
		printTask(false, list)
	case islistAllTask:
		printTask(true, list)
	case doneTaskID != 0:
		doneTask(list)
	case addTaskTitle != "":
		addTask(list)
	case deleteTaskID != 0:
		deleteTask(list)
	}

}
