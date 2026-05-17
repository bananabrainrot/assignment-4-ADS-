package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type Edge struct {
	To     string
	Weight int
}

type Graph struct {
	adj map[string][]Edge
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[string][]Edge)}
}

func (g *Graph) AddVertex(v string) {
	if _, ok := g.adj[v]; !ok {
		g.adj[v] = []Edge{}
	}
}

func (g *Graph) AddEdge(u, v string, weight int) {
	g.AddVertex(u)
	g.AddVertex(v)
	g.adj[u] = append(g.adj[u], Edge{To: v, Weight: weight})
	g.adj[v] = append(g.adj[v], Edge{To: u, Weight: weight})
}

func (g *Graph) SortedVertices() []string {
	vertices := make([]string, 0, len(g.adj))
	for v := range g.adj {
		vertices = append(vertices, v)
	}
	sort.Strings(vertices)
	return vertices
}

func (g *Graph) SortedNeighbors(v string) []Edge {
	edges := append([]Edge(nil), g.adj[v]...)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].To < edges[j].To
	})
	return edges
}

func (g *Graph) PrintAdjacencyList() {
	fmt.Println("Adjacency List:")
	for _, v := range g.SortedVertices() {
		neighbors := g.SortedNeighbors(v)
		fmt.Printf("%s: ", v)
		for i, edge := range neighbors {
			fmt.Printf("%s(%d)", edge.To, edge.Weight)
			if i < len(neighbors)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g *Graph) DFS(start string) []string {
	visited := make(map[string]bool)
	order := []string{}
	var dfsFunc func(string)
	dfsFunc = func(node string) {
		visited[node] = true
		order = append(order, node)
		for _, edge := range g.SortedNeighbors(node) {
			if !visited[edge.To] {
				dfsFunc(edge.To)
			}
		}
	}
	if _, ok := g.adj[start]; ok {
		dfsFunc(start)
	}
	return order
}

func (g *Graph) BFS(start string) []string {
	visited := make(map[string]bool)
	order := []string{}
	queue := []string{start}
	if _, ok := g.adj[start]; !ok {
		return order
	}
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)
		for _, edge := range g.SortedNeighbors(node) {
			if !visited[edge.To] {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}
	return order
}

func (g *Graph) Dijkstra(source string) (map[string]int, map[string]string) {
	dist := make(map[string]int)
	prev := make(map[string]string)
	for vertex := range g.adj {
		dist[vertex] = int(^uint(0) >> 1) // max int
	}
	if _, ok := g.adj[source]; !ok {
		return dist, prev
	}
	dist[source] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{vertex: source, priority: 0})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		u := item.vertex
		d := item.priority
		if d > dist[u] {
			continue
		}
		for _, edge := range g.adj[u] {
			alt := dist[u] + edge.Weight
			if alt < dist[edge.To] {
				dist[edge.To] = alt
				prev[edge.To] = u
				heap.Push(pq, &Item{vertex: edge.To, priority: alt})
			}
		}
	}
	return dist, prev
}

func ReconstructPath(prev map[string]string, dest string) []string {
	path := []string{}
	current := dest
	for current != "" {
		path = append([]string{current}, path...)
		current = prev[current]
	}
	return path
}

type Item struct {
	vertex   string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}

func main() {
	g := NewGraph()
	vertices := []string{"A", "B", "C", "D", "E", "F"}
	for _, v := range vertices {
		g.AddVertex(v)
	}

	g.AddEdge("A", "D", 12)
	g.AddEdge("E", "B", 1)
	g.AddEdge("F", "A", 5)
	g.AddEdge("B", "F", 1)
	g.AddEdge("A", "E", 9)
	g.AddEdge("C", "B", 1)

	g.PrintAdjacencyList()

	dfsOrder := g.DFS("C")
	fmt.Printf("DFS starting from C: %v\n", dfsOrder)

	bfsOrder := g.BFS("C")
	fmt.Printf("BFS starting from C: %v\n", bfsOrder)

	dist, prev := g.Dijkstra("A")
	fmt.Println("\nDijkstra shortest paths from A:")
	for _, v := range g.SortedVertices() {
		if dist[v] == int(^uint(0)>>1) {
			fmt.Printf("%s: unreachable\n", v)
			continue
		}
		path := ReconstructPath(prev, v)
		if len(path) == 0 {
			path = []string{v}
		}
		fmt.Printf("%s: distance=%d, path=%v\n", v, dist[v], path)
	}
}
