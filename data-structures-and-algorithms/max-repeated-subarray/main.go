package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func findLength(nums1 []int, nums2 []int) int {
	m, n := len(nums1), len(nums2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	maxLen := 0

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if nums1[i-1] == nums2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
				}
			}
		}
	}
	return maxLen
}

func main() {
	fmt.Print("Enter nums1 (format: [1,1,1,1,1]): ")
	reader := bufio.NewReader(os.Stdin)
	nums1Str, _ := reader.ReadString('\n')
	nums1Str = strings.TrimSpace(nums1Str)

	var nums1 []int
	if nums1Str == "" || nums1Str == "[]" {
		fmt.Println("Invalid input: nums1 length must be between 1 and 1000.")
		return
	}

	nums1Str = strings.Trim(nums1Str, "[]")
	parts := strings.Split(nums1Str, ",")

	if len(parts) < 1 || len(parts) > 1000 {
		fmt.Println("Invalid input: nums1 length must be between 1 and 1000.")
		return
	}

	for _, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Println("Invalid input: nums1 contains non-integer values.")
			return
		}
		if num < 0 || num > 100 {
			fmt.Println("Invalid input: nums1 values must be between 0 and 100.")
			return
		}
		nums1 = append(nums1, num)
	}

	fmt.Print("Enter nums2 (format: [1,1,1,1,1]): ")
	nums2Str, _ := reader.ReadString('\n')
	nums2Str = strings.TrimSpace(nums2Str)

	var nums2 []int
	if nums2Str == "" || nums2Str == "[]" {
		fmt.Println("Invalid input: nums2 length must be between 1 and 1000.")
		return
	}

	nums2Str = strings.Trim(nums2Str, "[]")
	parts = strings.Split(nums2Str, ",")

	if len(parts) < 1 || len(parts) > 1000 {
		fmt.Println("Invalid input: nums2 length must be between 1 and 1000.")
		return
	}

	for _, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			fmt.Println("Invalid input: nums2 contains non-integer values.")
			return
		}
		if num < 0 || num > 100 {
			fmt.Println("Invalid input: nums2 values must be between 0 and 100.")
			return
		}
		nums2 = append(nums2, num)
	}

	result := findLength(nums1, nums2)
	fmt.Println(result)
}
