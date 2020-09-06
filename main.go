package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
 * @Description: 随机自避行走
 * @Author: Senkita
 * @Date: 2020-09-04 12:06:47
 * @LastEditors: Senkita
 * @LastEditTime: 2020-09-07 01:06:27
 */

// Node 节点属性
// * posX：x坐标、posY：y坐标、selfWeight：自身权重、through：是否途径点
type Node struct {
	posX       int
	posY       int
	selfWeight int
	through    bool
}

// Route 路线属性
// * routeArr：路线、length：路线长度
type Route struct {
	routeArr [][]int
	length   int
}

// 互斥锁、二维网格、总权重、路线数组、方向数组、新线路、线路集合
var (
	mutex        sync.Mutex
	matrix       [88][88]Node
	totalWeight  int
	routeArr     [][]int
	directionArr []string
	newRoute     Route
	route        []Route
)

func initialize(matrix *[88][88]Node) {
	/**
	 * @description: 初始化网格
	 * @param {*[88][88]Node} matrix - [二维矩阵]
	 * @return {type}
	 * @author: Senkita
	 */
	for a := 0; a < len(matrix); a++ {
		for b := 0; b < len(matrix[a]); b++ {
			matrix[a][b].posX = a
			matrix[a][b].posY = b
			matrix[a][b].through = false
		}
	}
}

func computeWeight(totalWeight *int) {
	/**
	 * @description: 计算总权重
	 * @param {*int} totalWeight - [总权重]
	 * @return {type}
	 * @author: Senkita
	 */
	mutex.Lock()
	*totalWeight = 0

	for a := 0; a < len(matrix); a++ {
		for b := 0; b < len(matrix[a]); b++ {
			*totalWeight += matrix[a][b].selfWeight
		}
	}
	mutex.Unlock()
}

func determineOutset(routeArr *[][]int) Node {
	/**
	 * @description: 确定初始点，能加权使用加权随机，无加权则使用普通随机
	 * @param {*[][]int} routeArr - [线路数组]
	 * @return {Node} matrix[a][b] - [起始点]
	 * @author: Senkita
	 */
	// 置空路线数组
	*routeArr = append([][]int{})

	// 定出加权随机数和普通随机数
	randomA := rand.Intn(len(matrix))
	randomB := rand.Intn(len(matrix[0]))

	if totalWeight != 0 {
		random := rand.Intn(totalWeight)

		for a := 0; a < len(matrix); a++ {
			for b := 0; b < len(matrix[a]); b++ {
				random -= matrix[a][b].selfWeight
				if random <= 0 {
					*routeArr = append(*routeArr, []int{matrix[a][b].posX, matrix[a][b].posY})
					return matrix[a][b]
				}
			}
		}
	}
	*routeArr = append(*routeArr, []int{matrix[randomA][randomB].posX, matrix[randomA][randomB].posY})
	return matrix[randomA][randomB]
}

func judgeDirection(node Node, directionArr *[]string) []string {
	/**
	* @description: 判断可走方向
	* @param {Node} node - 点
			  {*[]string} directionArr - 方向数组
	* @return {*[]string} directionArr - 方向数组
	* @author: Senkita
	*/
	// * 清空方向数组
	*directionArr = append([]string{})

	if node.posY-1 != -1 && matrix[node.posX][node.posY-1].through == false {
		*directionArr = append(*directionArr, "上")
	}
	if node.posY+1 != len(matrix[0]) && matrix[node.posX][node.posY+1].through == false {
		*directionArr = append(*directionArr, "下")
	}
	if node.posX-1 != -1 && matrix[node.posX-1][node.posY].through == false {
		*directionArr = append(*directionArr, "左")
	}
	if node.posX+1 != len(matrix) && matrix[node.posX+1][node.posY].through == false {
		*directionArr = append(*directionArr, "右")
	}

	return *directionArr
}

func randomDirection(node Node, directionArr []string, routeArr *[][]int, matrix *[88][88]Node) Node {
	/**
	 * @description: 普通随机确定下一步方向并行走(计入路线数组)
	 * @param {Node} node - 点
			  {[]string} directionArr - [方向数组]
			  {*[][]int} routeArr - [线路数组]
			  [*[88][88]Node] matrix - [二维矩阵]
	 * @return {Node} nextNode - [下一点]
	 * @author: Senkita
	*/

	random := rand.Intn(len(directionArr))
	direction := directionArr[random]

	var nextNode Node

	switch direction {
	case "上":
		nextNode.posX = node.posX
		nextNode.posY = node.posY - 1
	case "下":
		nextNode.posX = node.posX
		nextNode.posY = node.posY + 1
	case "左":
		nextNode.posX = node.posX - 1
		nextNode.posY = node.posY
	default:
		nextNode.posX = node.posX + 1
		nextNode.posY = node.posY
	}
	matrix[nextNode.posX][nextNode.posY].through = true
	*routeArr = append(*routeArr, []int{nextNode.posX, nextNode.posY})

	return nextNode
}

func duplicate(newRoute Route) bool {
	/**
	 * @description: 判重
	 * @param {Route} newRoute - [新路线]
	 * @return {bool}
	 * @author: Senkita
	 */
	for _, v := range route {
		if reflect.DeepEqual(newRoute.routeArr, v.routeArr) {
			return true
		}
	}
	return false
}

func cumulativeWeight(routeArr [][]int, matrix *[88][88]Node) {
	/**
	* @description: 累计权重
	* @param {[][]int} routeArr - [路线数组]
			  {*[88][88]Node} matrix - [二维矩阵]
	* @return {type}
	* @author: Senkita
	*/
	for _, v := range routeArr {
		matrix[v[0]][v[1]].selfWeight++
	}
}

func routeToFile(fileName string, routeArr [][]int) {
	/**
	* @description: 线路输出到文本
	* @param {string} fileName - [文件名]
			  {[][]int} routeArr - [路线数组]
	* @return {type}
	* @author: Senkita
	*/
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	writer := bufio.NewWriter(file)

	for _, v := range routeArr {
		content := strings.Join([]string{strconv.Itoa(v[0]), strconv.Itoa(v[1])}, ",")

		writer.WriteString(content)
		writer.WriteString("\r\n")
	}
	writer.Flush()
}

func weightToFile(fileName string, matrix [88][88]Node) {
	/**
	* @description: 权重输出到文本
	* @param {string} fileName - [文件名]
			  {[][]Node} matrix - [二维数组]
	* @return {type}
	* @author: Senkita
	*/
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	writer := bufio.NewWriter(file)

	for a := 0; a < len(matrix); a++ {
		for b := 0; b < len(matrix[a]); b++ {
			content := strings.Join([]string{strconv.Itoa(matrix[a][b].posX), strconv.Itoa(matrix[a][b].posY), strconv.Itoa(matrix[a][b].selfWeight)}, " ")

			writer.WriteString(content)
			writer.WriteString("\r\n")

		}

	}

	writer.Flush()
}

func process(runtimeStr string) {
	/**
	 * @description: 流程
	 * @param {string} runtimeStr - [启动时间戳]
	 * @return {type}
	 * @author: Senkita
	 */
	for n := 0; n < 100; {
		initialize(&matrix)

		computeWeight(&totalWeight)
		point := determineOutset(&routeArr)
	subRoute:
		for {
			directionArr = judgeDirection(point, &directionArr)

			if len(directionArr) == 0 {
				if len(routeArr) >= 30 {
					newRoute.routeArr = routeArr
					newRoute.length = len(routeArr)
					if !duplicate(newRoute) {
						route = append(route, newRoute)
						cumulativeWeight(newRoute.routeArr, &matrix)
						n++

						routeFileName := strings.Join([]string{"./data/", runtimeStr, "/", strconv.Itoa(n), ".txt"}, "")

						routeToFile(routeFileName, newRoute.routeArr)
					}
					break subRoute
				}
				break subRoute
			}
			point = randomDirection(point, directionArr, &routeArr, &matrix)
		}
	}

	weightFileName := strings.Join([]string{"./data/", runtimeStr, "/", runtimeStr, ".txt"}, "")

	weightToFile(weightFileName, matrix)
}

func main() {
	runtimeStr := strconv.FormatInt(time.Now().Unix(), 10)
	dirName := strings.Join([]string{".", "data", runtimeStr}, "/")
	err := os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}
	process(runtimeStr)
}
