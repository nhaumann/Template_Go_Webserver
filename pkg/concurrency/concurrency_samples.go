package concurrency_samples

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Point struct {
	x float64
	y float64
	z float64
}

const (
	SIZE                       = 10000000
	NEAREST_NEIGHBORS_DISTANCE = 10
)

func Setup() {

	fmt.Println("Generating points using wait groups...")
	gp_wg_start := time.Now()
	generatePointsWithWaitGroup(SIZE)
	gp_wg_elapsed := time.Since(gp_wg_start)
	fmt.Println("generatePointsWithWaitGroup took", gp_wg_elapsed)
	fmt.Println("===========================================")

	fmt.Println("Generating points using wait groups and pointer receiver...")
	gp_wg_pr_start := time.Now()
	generatePointsWithWaitGroupAndPointerReceiver(SIZE)
	gp_wg_pr_elapsed := time.Since(gp_wg_pr_start)
	fmt.Println("generatePointsWithWaitGroupAndPointerReceiver took", gp_wg_pr_elapsed)
	fmt.Println("===========================================")

	fmt.Println("Generating points using channel...")
	gp_ch_start := time.Now()
	generatePointsWithChannel(SIZE)
	gp_ch_elapsed := time.Since(gp_ch_start)
	fmt.Println("generatePointsWithChannel took", gp_ch_elapsed)

	fmt.Println("===========================================")

	fmt.Println("Generating points using channel and pointer receiver...")
	gp_ch_pr_start := time.Now()
	generatePointsWithChannelAndPointerReceiver(SIZE)
	gp_ch_pr_elapsed := time.Since(gp_ch_pr_start)
	fmt.Println("generatePointsWithChannelAndPointerReceiver took", gp_ch_pr_elapsed)

	fmt.Println("===========================================")

	fmt.Println("Generating points...")
	gp_start := time.Now()
	generatePoints(SIZE)
	gp_elapsed := time.Since(gp_start)
	fmt.Println("generatePoints took", gp_elapsed)

	fmt.Println("===========================================")

	fmt.Println("Generating points using pointer receiver...")
	gp_pr_start := time.Now()
	generatePointsWithPointerReceiver(SIZE)
	gp_pr_elapsed := time.Since(gp_pr_start)
	fmt.Println("generatePointsUsingPointerReceiver took", gp_pr_elapsed)

	fmt.Println("===========================================")

	ex_lr_wg_start := time.Now()
	executeLongRunningTasksUsingWaitGroups(10)
	ex_lr_wg_elapsed := time.Since(ex_lr_wg_start)
	fmt.Println("executeLongRunningTasksUsingWaitGroups took", ex_lr_wg_elapsed)

	fmt.Println("===========================================")

	ex_lr_ch_start := time.Now()
	executeLongRunningTasksUsingChannels(10)
	ex_lr_ch_elapsed := time.Since(ex_lr_ch_start)
	fmt.Println("executeLongRunningTasksUsingChannel took", ex_lr_ch_elapsed)

	fmt.Println("===========================================")

	ex_lr_start := time.Now()
	executeLongRunningTasks(10)
	ex_lr_elapsed := time.Since(ex_lr_start)
	fmt.Println("executeLongRunningTasks took", ex_lr_elapsed)

}

// GeneratePoint generates a single point with random x, y, and z values
func generatePoint() Point {
	return Point{
		x: rand.Float64()*100 + 1,
		y: rand.Float64()*100 + 1,
		z: rand.Float64()*100 + 1,
	}
}

func generatePointUsingPointerReceiver(p *Point) {
	p = &Point{
		x: rand.Float64()*100 + 1,
		y: rand.Float64()*100 + 1,
		z: rand.Float64()*100 + 1,
	}
}

func generatePointsWithWaitGroup(size int) []Point {

	var points []Point = make([]Point, size)
	var wg sync.WaitGroup
	wg.Add(size)

	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			points[i] = generatePoint()
		}(i)
	}
	wg.Wait()
	return points
}

func generatePointsWithWaitGroupAndPointerReceiver(size int) []Point {

	var points []Point = make([]Point, size)
	var wg sync.WaitGroup
	wg.Add(size)

	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			generatePointUsingPointerReceiver(&points[i])
		}(i)
	}
	wg.Wait()
	return points
}

func generatePointsWithChannel(size int) []Point {

	var points []Point = make([]Point, size)
	ch := make(chan Point, size)

	for i := 0; i < size; i++ {
		go func() {
			ch <- generatePoint()
		}()
	}

	for i := 0; i < size; i++ {
		points[i] = <-ch
	}

	return points
}

func generatePointsWithChannelAndPointerReceiver(size int) {

	var points []Point = make([]Point, size)
	ch := make(chan Point, size)

	for i := 0; i < size; i++ {
		go func() {
			var p Point
			generatePointUsingPointerReceiver(&p)
			ch <- p
		}()
	}

	for i := 0; i < size; i++ {
		points[i] = <-ch
	}

}

func generatePoints(size int) []Point {
	var points []Point = make([]Point, size)
	for i := 0; i < size; i++ {
		points[i] = generatePoint()
	}
	return points
}

func generatePointsWithPointerReceiver(size int) []Point {
	var points []Point = make([]Point, size)

	for i := 0; i < size; i++ {
		generatePointUsingPointerReceiver(&points[i])
	}
	return points
}

func executeLongRunningTasksUsingWaitGroups(iterations int) {

	var wg sync.WaitGroup
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			longRunningTask()
			wg.Done()
		}()
	}

	wg.Wait()

}

func executeLongRunningTasksUsingChannels(iterations int) {
	ch := make(chan int, iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			longRunningTask()
			ch <- 1
		}()
	}

	for i := 0; i < iterations; i++ {
		<-ch
	}
	close(ch)
}

func executeLongRunningTasks(iterations int) {
	for i := 0; i < iterations; i++ {
		longRunningTask()
	}
}

func longRunningTask() {
	time.Sleep(1 * time.Second)
}
