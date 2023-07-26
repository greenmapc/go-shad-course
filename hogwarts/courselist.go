//go:build !solution

package hogwarts

import "fmt"

const (
	GrayVertex  = '1'
	BlackVertex = '2'
)

func GetCourseList(prereqs map[string][]string) []string {
	colors := make(map[string]int)
	parents := []string{}

	for subject := range prereqs {
		_, exists := colors[subject]
		if !exists {
			fmt.Println("starting for ", subject)
			topologicalSort(prereqs, subject, colors, &parents)
		}
	}

	return parents
}

func topologicalSort(graph map[string][]string, v string, colors map[string]int, parents *[]string) {
	colors[v] = GrayVertex

	for _, next := range graph[v] {
		color, exists := colors[next]

		if exists && color == GrayVertex {
			panic("Cycle in the graph!")
		}

		if !exists {
			topologicalSort(graph, next, colors, parents)
		}
	}
	*parents = append(*parents, v)
	colors[v] = BlackVertex
}
