package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxworkers = 10

var chanel = make(chan *customer, 10)
var workchan = make(chan *worker, maxworkers)

//var stop = make(chan bool)
var wg sync.WaitGroup

func main() {
	fmt.Printf("Hello, the store is open\n\n")
	wg.Add(3)
	go send_customers(30)
	go gen_workers(20)
	go receive()
	wg.Wait()
	summary()
	fmt.Println("See you tomorrow")
}

func summary() {
	fmt.Printf(" ------------  SUMMARY ------------  \n")
	n := len(workchan)
	fmt.Printf("%d workers\n", n)
	grandtotal := 0.0
	ncusto := 0
	for i := 0; i < n; i++ {
		w, more := <-workchan
		grandtotal += w.made
		ponctu := "."
		if more {
			ponctu = ","
		}
		ncusto += w.ncusto
		fmt.Printf("%s treated %d customers and made %.2f € today%s\n", w.name, w.ncusto, w.made, ponctu)

	}
	fmt.Printf("\nCongratulations, we made %.2f € and %d happy customers today guys!\n", grandtotal, ncusto)
	close(workchan)
}
func receive() {
	for {

		cu, ok := <-chanel
		if !ok {
			defer wg.Done()
			return
		}
		w := <-workchan
		w.custo = cu
		fmt.Printf("                    %s available\n", w.name)
		fmt.Printf("                    received %s, sent to %s\n", cu.name, w.name)
		wg.Add(1)
		go w.checkout()
	}
}

func send_customers(ncusto int) {

	for i := 0; i < ncusto; i++ {
		//randomly select a number of products
		npr := rand.Intn(8)
		//generate products list
		prs := []product{}
		for j := 0; j < npr; j++ {
			p := rand.Float64() * 15.0
			prs = append(prs, product{name: fmt.Sprintf("product_%d", j+1), price: p})
		}
		nc := customer{name: fmt.Sprintf("customer_%d", i+1), prods: prs}
		//time.Sleep(22 * time.Millisecond)
		fmt.Printf("%s arrived\n", nc.name)
		/*select {
		case chanel <- &nc:
		default:
			time.Sleep(100 * time.Millisecond)
		}*/
		chanel <- &nc
	}
	close(chanel)
	wg.Done()
}

func gen_workers(nworkers int) {
	defer wg.Done()
	for i := 0; i < nworkers; i++ {
		if i > maxworkers-1 {
			fmt.Println("sorry, the maximum number of workers has been reached ! ... can not accelerate")
			return
		}
		d := rand.Intn(15) * 20
		nc := worker{name: fmt.Sprintf("worker_%d", i+1), delay: d, custo: nil}
		//time.Sleep(22 * time.Millisecond)
		//fmt.Printf("     %s available\n", nc.name)
		workchan <- &nc
	}
	//close(chanel)

}

type customer struct {
	name  string
	prods []product
}
type product struct {
	name  string
	price float64
}

type worker struct {
	name   string
	delay  int
	custo  *customer
	made   float64
	ncusto int
}

func (w *worker) checkout() {
	//fmt.Println(w.custo)

	ticket := ""
	ticket = fmt.Sprintf("%s                                                                           @%s PAID TO %s:\n", ticket, w.custo.name, w.name)
	total := 0.0
	for _, p := range w.custo.prods {
		time.Sleep(time.Duration(w.delay) * time.Millisecond)
		ticket = fmt.Sprintf("%s                                                                           %s --- %.2f €\n", ticket, p.name, p.price)
		total += p.price
	}
	ticket = fmt.Sprintf("%s                                                                           total=%.2f €\n", ticket, total)
	fmt.Println(ticket)
	w.custo = nil
	w.made += total
	w.ncusto++
	workchan <- w
	wg.Done()
}
