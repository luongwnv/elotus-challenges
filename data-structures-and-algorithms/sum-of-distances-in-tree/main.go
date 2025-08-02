package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func sumOfDistancesInTree(n int, edges [][]int) []int {
	graph := make([][]int, n)
	for _, edge := range edges {
		a, b := edge[0], edge[1]
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	answer := make([]int, n)
	count := make([]int, n)
	for i := range count {
		count[i] = 1
	}

	var dfs1 func(node, parent, depth int)
	dfs1 = func(node, parent, depth int) {
		answer[0] += depth
		for _, neighbor := range graph[node] {
			if neighbor != parent {
				dfs1(neighbor, node, depth+1)
				count[node] += count[neighbor]
			}
		}
	}

	var dfs2 func(node, parent int)
	dfs2 = func(node, parent int) {
		for _, neighbor := range graph[node] {
			if neighbor != parent {
				answer[neighbor] = answer[node] - count[neighbor] + (n - count[neighbor])
				dfs2(neighbor, node)
			}
		}
	}

	if n == 1 {
		return []int{0}
	}

	dfs1(0, -1, 0)
	dfs2(0, -1)

	return answer
}

func main() {
	var n int
	fmt.Print("Enter n: ")
	_, err := fmt.Scan(&n)
	if err != nil || n < 1 || n > 30000 {
		fmt.Println("Invalid input: n must be between 1 and 30000.")
		return
	}

	fmt.Print("Enter edges (format: [[a,b],[c,d],...]): ")
	reader := bufio.NewReader(os.Stdin)
	edgesStr, _ := reader.ReadString('\n')
	edgesStr = strings.TrimSpace(edgesStr)

	var edges [][]int
	if edgesStr != "" && edgesStr != "[]" {
		edgesStr = strings.Trim(edgesStr, "[]")
		parts := strings.Split(edgesStr, "],[")
		for _, part := range parts {
			part = strings.Trim(part, "[]")
			if part != "" {
				nums := strings.Split(part, ",")
				if len(nums) == 2 {
					a, _ := strconv.Atoi(strings.TrimSpace(nums[0]))
					b, _ := strconv.Atoi(strings.TrimSpace(nums[1]))

					if a < 0 || a >= n || b < 0 || b >= n {
						fmt.Printf("Invalid input: Node values must be between 0 and %d.\n", n-1)
						return
					}

					if a == b {
						fmt.Printf("Invalid input: Self-loops are not allowed (ai != bi).\n")
						return
					}

					edges = append(edges, []int{a, b})
				}
			}
		}
	}

	if len(edges) != n-1 {
		fmt.Printf("Invalid input: Expected %d edges for n=%d, but got %d edges.\n", n-1, n, len(edges))
		return
	}

	result := sumOfDistancesInTree(n, edges)
	fmt.Printf("Output: %v\n", result)
}
